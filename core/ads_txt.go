package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	adstxt "github.com/m6yf/bcwork/modules/ads_txt"
	"github.com/m6yf/bcwork/modules/compass"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/m6yf/bcwork/modules/logger"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

type AdsTxtService struct {
	historyModule       history.HistoryModule
	compassModule       compass.CompassModule
	adstxtModule        adstxt.AdsTxtManager
	lowPerformanceCache *LowPerformanceCache
}

func NewAdsTxtService(
	ctx context.Context,
	historyModule history.HistoryModule,
	compassModule compass.CompassModule,
	adstxtModule adstxt.AdsTxtManager,
) *AdsTxtService {
	service := &AdsTxtService{
		historyModule:       historyModule,
		compassModule:       compassModule,
		adstxtModule:        adstxtModule,
		lowPerformanceCache: &LowPerformanceCache{},
	}

	go func(ctx context.Context) {
		var defaultMinutesToUpdateMetadata time.Duration = 60 * time.Minute

		minutesToUpdateMetadata := viper.GetDuration(config.AdsTxtMetadataUpdateRateKey) * time.Minute
		if minutesToUpdateMetadata == 0 {
			minutesToUpdateMetadata = defaultMinutesToUpdateMetadata
		}

		ticker := time.NewTicker(minutesToUpdateMetadata)

		for {
			select {
			case <-ticker.C:
				err := service.updateAdsTxtMetadata(ctx)
				if err != nil {
					logger.Logger(ctx).Err(err).Msg("cannot update ads txt metadata")
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}(ctx)

	return service
}

type LowPerformanceCache struct {
	partitions     string
	lowPerformance map[string]bool
	revenuePerDP   map[string]float64
	totalRevenue   map[string]float64
}

type ReportResult struct {
	Data struct {
		Result []struct {
			Domain        string  `json:"Domain"`
			DemandPartner string  `json:"DemandPartner"`
			Revenue       float64 `json:"Revenue"`
		} `json:"result"`
	} `json:"data"`
}

func (a *AdsTxtService) GetMainAdsTxtTable(ctx context.Context, ops *AdsTxtGetMainOptions) (*dto.AdsTxtResponse, error) {
	const cteName = "main_table"

	cteMods := ops.Filter.queryModMain()
	total, err := models.AdsTXTMainViews(cteMods...).Count(ctx, bcdb.DB())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve count from main table: %w", err)
	}

	cteMods = append(cteMods, qm.Select(
		getCursorIDExpression(ops.Order, cursorIDOrderByDefault),
		`*`,
	)).
		Order(ops.Order, nil, cursorIDColumnName)
	raw, args := queries.BuildQuery(models.AdsTXTMainViews(cteMods...).Query)

	qmods := []qm.QueryMod{
		qm.With(fmt.Sprintf("%v as (%v)", cteName, strings.ReplaceAll(raw, ";", "")), args...),
		qm.Select(`*`),
		qm.From(cteName),
	}
	qmods = append(qmods, ops.Pagination.DoV2(cursorIDColumnName, len(args))...)

	var mainTable []*dto.AdsTxt
	err = models.NewQuery(qmods...).Bind(ctx, bcdb.DB(), &mainTable)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve main table: %w", err)
	}

	for _, row := range mainTable {
		// statuses mapping
		row.Status = dto.StatusMap[row.Status]
		row.DomainStatus = dto.DomainStatusMap[row.DomainStatus]
		row.DemandStatus = dto.DPStatusMap[row.DemandStatus]
	}

	return &dto.AdsTxtResponse{
		Data:  mainTable,
		Total: total,
	}, nil
}

func (a *AdsTxtService) GetGroupByDPAdsTxtTable(ctx context.Context, ops *AdsTxtGetGroupByDPOptions) (*dto.AdsTxtGroupByDPResponse, error) {
	const cteName = "group_by_dp"

	total, err := getGroupByDPTotal(ctx, ops)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve count from group by dp table: %w", err)
	}

	cteMods := ops.Filter.queryModGroupByDP()
	cteMods = append(cteMods, qm.Select(
		getCursorIDExpression(ops.Order, groupByDPIDColumnName),
		`*`,
	))
	raw, args := queries.BuildQuery(models.AdsTXTGroupByDPViews(cteMods...).Query)

	qmods := []qm.QueryMod{
		qm.With(fmt.Sprintf("%v as (%v)", cteName, strings.ReplaceAll(raw, ";", "")), args...),
		qm.Select(`*`),
		qm.From(cteName),
	}
	qmods = append(qmods, ops.Pagination.DoV2(cursorIDColumnName, len(args))...)

	var rawTable []*dto.AdsTxt
	err = models.NewQuery(qmods...).Bind(ctx, bcdb.DB(), &rawTable)
	if err != nil {
		return nil, err
	}

	groupByDpTableMap := make(map[string]*dto.AdsTxtGroupedByDP)
	for _, row := range rawTable {
		name := fmt.Sprintf("%v:%v:%v:%v", row.PublisherID, row.Domain, row.DemandPartnerID, row.DemandPartnerConnectionID.Int)
		// statuses mapping
		row.Status = dto.StatusMap[row.Status]
		row.DomainStatus = dto.DomainStatusMap[row.DomainStatus]
		row.DemandStatus = dto.DPStatusMap[row.DemandStatus]

		dpData, ok := groupByDpTableMap[name]
		if !ok {
			groupByDPData := &dto.AdsTxtGroupedByDP{AdsTxt: &dto.AdsTxt{}}
			groupByDPData.FromAdsTxt(row)
			groupByDPData.GroupedLines = append(groupByDPData.GroupedLines, row)

			dpData = groupByDPData
			groupByDpTableMap[name] = dpData
		} else {
			dpData.ProcessParentRow(row)
			dpData.GroupedLines = append(dpData.GroupedLines, row)
		}
	}

	groupByDpTable := make([]*dto.AdsTxtGroupedByDP, 0, len(groupByDpTableMap))
	for _, row := range groupByDpTableMap {
		groupByDpTable = append(groupByDpTable, row)
	}

	response := &dto.AdsTxtGroupByDPResponse{
		Data:  groupByDpTable,
		Total: total,
		Order: ops.Order,
	}

	sort.Sort(response)

	return response, nil
}

func (a *AdsTxtService) GetAMAdsTxtTable(ctx context.Context, ops *AdsTxtGetBaseOptions) ([]*dto.AdsTxt, error) {
	type seatOwnersWithActiveDP struct {
		SeatOwnerName string
		HasActiveDP   bool
	}

	var seatOwners []*seatOwnersWithActiveDP
	err := queries.Raw(`
		select
			so.seat_owner_name,
			bool_or(d.active) as has_active_dp 
		from seat_owner so 
		join dpo d on d.seat_owner_id = so.id
		group by so.seat_owner_name
	`).Bind(ctx, bcdb.DB(), &seatOwners)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve seat owner map: %w", err)
	}

	seatOwnerMap := make(map[string]bool, len(seatOwners))
	for _, seatOwner := range seatOwners {
		seatOwnerMap[seatOwner.SeatOwnerName] = seatOwner.HasActiveDP
	}

	query := fmt.Sprintf(`
		with am as (
			select 
				t.*,
				p."name" as publisher_name,
				p.account_manager_id,
				p.campaign_manager_id,
				u1.first_name || ' ' || u1.last_name as account_manager_full_name,
				u2.first_name || ' ' || u2.last_name as campaign_manager_full_name,
				u3.first_name || ' ' || u3.last_name as demand_manager_full_name
			from (
				%v
				union 
				%v
				union all
				%v
			) as t
			join publisher p on p.publisher_id = t.publisher_id
			left join "user" u1 on u1.id::varchar = p.account_manager_id
			left join "user" u2 on u2.id::varchar = p.campaign_manager_id
			left join "user" u3 on u3.id = t.demand_manager_id
			where t.status != 'not_scanned' and t.domain_status != 'paused'
		)
		select 
			*
		from am a
		order by a.id;
	`,
		adsTxtDemandPartnerConnectionBaseQuery,
		fmt.Sprintf(adsTxtdemandPartnerChildrenBaseQuery, "dpc.dp_child_name"),
		fmt.Sprintf(adsTxtSeatOwnersBaseQuery, "d.demand_partner_id", "so.seat_owner_name", "null", "null", "d.active", dynamicPublisherIDPlaceholder, "join dpo d on d.seat_owner_id = so.id"),
	)

	var (
		errGroup errgroup.Group
		amTable  []*dto.AdsTxt
	)

	errGroup.Go(func() error {
		err := queries.Raw(query).Bind(ctx, bcdb.DB(), &amTable)
		if err != nil {
			return fmt.Errorf("failed to retrieve am table: %w", err)
		}

		return nil
	})

	errGroup.Go(func() error {
		err = a.updatePerfomanceData()
		if err != nil {
			return fmt.Errorf("failed to update perfomance data: %w", err)
		}

		return nil
	})

	err = errGroup.Wait()
	if err != nil {
		return nil, err
	}

	for _, row := range amTable {
		key := fmt.Sprintf("%s_%s_%s", strings.ToLower(row.Domain), strings.ToLower(row.DemandPartnerName), "nb")
		action := dto.AdsTxtActionNone
		hasActiveDP, isSeatOwner := seatOwnerMap[row.DemandPartnerName]

		if slices.Contains(
			[]string{
				dto.DPStatusRejected,
				dto.DPStatusRejectedTQ,
				dto.DPStatusWillNotBeSent,
				dto.DPStatusDisabledSPO,
				dto.DPStatusApprovedPaused,
				dto.DPStatusDisabledNoImps,
			},
			row.DemandStatus,
		) || (!isSeatOwner && !row.IsDemandPartnerActive) {
			action = dto.AdsTxtActionNone
			if row.Status == dto.AdsTxtStatusAdded {
				action = dto.AdsTxtActionRemove
			}
		} else if row.Status == dto.AdsTxtStatusAdded {
			action = dto.AdsTxtActionKeep
			if row.ErrorMessage.Valid {
				action = dto.AdsTxtActionFix
			}
		} else {
			action = dto.AdsTxtActionAdd
		}

		if !isSeatOwner &&
			row.Status == dto.AdsTxtStatusAdded &&
			!slices.Contains([]string{dto.AdsTxtActionFix, dto.AdsTxtActionRemove}, action) &&
			a.lowPerformanceCache.lowPerformance[key] {
			action = dto.AdsTxtActionLowPerfomance
		}

		if isSeatOwner && !hasActiveDP {
			if action == dto.AdsTxtActionAdd {
				action = dto.AdsTxtActionNone
			} else if slices.Contains([]string{dto.AdsTxtActionKeep, dto.AdsTxtActionFix}, action) {
				action = dto.AdsTxtActionRemove
			}
		}

		row.Action = action
	}

	return amTable, nil
}

func (a *AdsTxtService) GetCMAdsTxtTable(ctx context.Context, ops *AdsTxtGetBaseOptions) ([]*dto.AdsTxt, error) {
	query := fmt.Sprintf(`
		with cm as (
			select 
				t.*,
				p."name" as publisher_name,
				p.account_manager_id,
				p.campaign_manager_id,
				u1.first_name || ' ' || u1.last_name as account_manager_full_name,
				u2.first_name || ' ' || u2.last_name as campaign_manager_full_name,
				u3.first_name || ' ' || u3.last_name as demand_manager_full_name,
				case
					when t.is_approval_needed and t.is_required then 1
					when t.is_approval_needed and not t.is_required then 2
					when not t.is_approval_needed and t.is_required then 3
					else 0
				end as approval_group
			from (
				%v
				union 
				%v
			) as t
			join publisher p on p.publisher_id = t.publisher_id
			left join "user" u1 on u1.id::varchar = p.account_manager_id
			left join "user" u2 on u2.id::varchar = p.campaign_manager_id
			left join "user" u3 on u3.id = t.demand_manager_id
			where t.is_demand_partner_active and t.demand_status in ('pending', 'not_sent')
		)
		select 
			*
		from cm c
		where 
			(c.approval_group = 1 and c.domain_status in ('new', 'active') and c.status = 'added')
			or (c.approval_group = 2 and c.domain_status in ('new', 'active') and c.status in ('added', 'no', 'not_scanned'))
			or (c.approval_group = 3 and c.domain_status = 'active' and c.status = 'added')
		order by c.id;
	`,
		adsTxtDemandPartnerConnectionBaseQuery,
		fmt.Sprintf(adsTxtdemandPartnerChildrenBaseQuery, "d.demand_partner_name"),
	)

	var cmTable []*dto.AdsTxt
	err := queries.Raw(query).Bind(ctx, bcdb.DB(), &cmTable)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve cm table: %w", err)
	}

	return cmTable, nil
}

func (a *AdsTxtService) GetMBAdsTxtTable(ctx context.Context, ops *AdsTxtGetBaseOptions) ([]*dto.AdsTxt, error) {
	query := `
		select 
			demand_partner_name_extended,
			seat_owner_name,
			score,
			ads_txt_line
		from (
			-- seat_owners
			select 
				so.seat_owner_name || ' - Direct' as demand_partner_name_extended,
				so.seat_owner_name,
				case 
					when so.seat_owner_name in ('OMS', 'Brightcom') then 0
					else min(score)
				end as score,
				so.seat_owner_domain || ', ' || 
					replace(so.publisher_account, '%s', 'XXXXX') || ', ' || 
					'DIRECT' || 
				case 
					when so.certification_authority_id is not null 
					then ', ' || so.certification_authority_id 
				else '' 
				end as ads_txt_line,
				true as active,
				true as is_seat_owner
			from seat_owner so
			join dpo d on so.id = d.seat_owner_id
			where d.active
			group by so.seat_owner_name, so.seat_owner_domain, so.publisher_account, so.certification_authority_id
			union
			-- demand_partners
			select 
				d.demand_partner_name || ' - ' || d.demand_partner_name as demand_partner_name_extended,
				coalesce(so.seat_owner_name, d.demand_partner_name),
				d.score,
				dpc.dp_domain || ', ' || 
					dpc.publisher_account || ', ' || 
					case 
						when dpc.is_direct then 'DIRECT' 
						else 'RESELLER' 
					end || 
					case 
						when dpc.certification_authority_id is not null 
						then ', ' || dpc.certification_authority_id 
					else '' 
				end as ads_txt_line,
				d.active,
				false as is_seat_owner
			from dpo d 
			join demand_partner_connection dpc on d.demand_partner_id = dpc.demand_partner_id
			left join seat_owner so on d.seat_owner_id = so.id
			where dpc.is_required_for_ads_txt
			union
			-- GroupedLines
			select 
				d.demand_partner_name || ' - ' || dpc.dp_child_name as demand_partner_name_extended,
				coalesce(so.seat_owner_name, dpc.dp_child_name),
				d.score,
				dpc.dp_domain || ', ' || 
					dpc.publisher_account || ', ' || 
					case 
						when dpc.is_direct then 'DIRECT' 
						else 'RESELLER' 
					end || 
					case 
						when dpc.certification_authority_id is not null 
						then ', ' || dpc.certification_authority_id 
					else '' 
				end as ads_txt_line,
				true as active,
				false as is_seat_owner
			from dpo d 
			join demand_partner_connection dpc2 on d.demand_partner_id = dpc2.demand_partner_id
			join demand_partner_child dpc on dpc2.id = dpc.dp_connection_id 
			left join seat_owner so on d.seat_owner_id = so.id
			where dpc.is_required_for_ads_txt
		) as t
		where active
		order by t.score, t.is_seat_owner, t.demand_partner_name_extended;
	`

	var mbTable []*dto.AdsTxt
	err := queries.Raw(query).Bind(ctx, bcdb.DB(), &mbTable)
	if err != nil {
		return nil, err
	}

	for i := range mbTable {
		mbTable[i].ID = i + 1
	}

	return mbTable, nil
}

func (a *AdsTxtService) UpdateAdsTxt(ctx context.Context, data *dto.AdsTxtUpdateRequest) error {
	if data.DemandPartnerID != nil {
		var demandPartnerConnections models.DemandPartnerConnectionSlice
		err := models.DemandPartnerConnections(
			qm.Select(models.DemandPartnerConnectionColumns.ID),
			models.DemandPartnerConnectionWhere.DemandPartnerID.EQ(*data.DemandPartnerID),
		).Bind(ctx, bcdb.DB(), &demandPartnerConnections)
		if err != nil {
			return fmt.Errorf("failed to get demand partner connection ids while updating demand status for ads txt lines: %w", err)
		}

		var demandPartnerConnectionIDs []int
		for _, row := range demandPartnerConnections {
			demandPartnerConnectionIDs = append(demandPartnerConnectionIDs, row.ID)
		}

		var demandPartnerChildren models.DemandPartnerChildSlice
		err = models.DemandPartnerChildren(
			qm.Select(models.DemandPartnerChildColumns.ID),
			models.DemandPartnerChildWhere.DPConnectionID.IN(demandPartnerConnectionIDs),
		).Bind(ctx, bcdb.DB(), &demandPartnerChildren)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to get demand partner child ids while updating demand status for ads txt lines: %w", err)
		}

		var demandPartnerChildIDs []int
		for _, row := range demandPartnerChildren {
			demandPartnerChildIDs = append(demandPartnerChildIDs, row.ID)
		}

		mods := qmods.QueryModsSlice{
			models.AdsTXTWhere.DemandPartnerConnectionID.IN(demandPartnerConnectionIDs),
			models.AdsTXTWhere.Domain.IN(data.Domain),
		}
		if len(demandPartnerChildIDs) > 0 {
			mods = append(mods,
				models.AdsTXTWhere.DemandPartnerChildID.IN(demandPartnerChildIDs),
			)
		}

		_, err = models.AdsTXTS(mods...).UpdateAll(ctx, bcdb.DB(), models.M{models.AdsTXTColumns.DemandStatus: *data.DemandStatus})
		if err != nil {
			return fmt.Errorf("failed to update demand status for ads txt lines: %w", err)
		}
	} else {
		// if there is no demand partner id, updating domain status
		_, err := models.AdsTXTS(
			models.AdsTXTWhere.Domain.IN(data.Domain),
		).
			UpdateAll(ctx, bcdb.DB(), models.M{models.AdsTXTColumns.DomainStatus: *data.DomainStatus})
		if err != nil {
			return fmt.Errorf("failed to update domain status for ads txt lines: %w", err)
		}
	}

	go a.adstxtModule.UpdateAdsTxtMaterializedViews(ctx)

	return nil
}

func (a *AdsTxtService) updatePerfomanceData() error {
	const (
		compassKey          = "compass"
		bothKey             = "both"
		newBidderKey        = "nb"
		newBidderReportPath = "/report-dashboard/report-new-bidder"
		compassReportPath   = "/report-dashboard/report-query" //nolint:gosec
	)

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth-1, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := time.Date(currentYear, currentMonth-1, 1, 23, 59, 59, 0, currentLocation).AddDate(0, 1, -1)

	lpPartitions := fmt.Sprintf(
		"%s-%s",
		firstOfMonth.Format("2006-01-02"),
		lastOfMonth.Format("2006-01-02"),
	)

	if a.lowPerformanceCache.partitions != lpPartitions {
		a.lowPerformanceCache.partitions = lpPartitions
		a.lowPerformanceCache.lowPerformance = nil
		a.lowPerformanceCache.revenuePerDP = nil
		a.lowPerformanceCache.totalRevenue = nil
	}

	if a.lowPerformanceCache.lowPerformance == nil {
		body := buildGetLowPerfomanceRequestBody(firstOfMonth, lastOfMonth)
		var (
			compassResult, nbResult ReportResult
			errGroup                errgroup.Group
		)

		errGroup.Go(func() error {
			data, err := a.compassModule.Request(compassReportPath, http.MethodPost, body, true)
			if err != nil {
				return err
			}

			err = json.Unmarshal(data, &compassResult)
			if err != nil {
				return err
			}

			return nil
		})

		errGroup.Go(func() error {
			data, err := a.compassModule.Request(newBidderReportPath, http.MethodPost, body, true)
			if err != nil {
				return err
			}

			err = json.Unmarshal(data, &nbResult)
			if err != nil {
				return err
			}

			return nil
		})

		err := errGroup.Wait()
		if err != nil {
			return err
		}

		totalDomainsRevenue := make(map[string]float64)
		totalDomainsRevenuePerDP := make(map[string]float64)
		totalRevenuePerDP := make(map[string]float64)

		// for _, row := range compassResult.Data.Result {
		// 	domain := strings.ToLower(row.Domain)
		// 	demandPartner := strings.ToLower(row.DemandPartner)
		// 	ckey := fmt.Sprintf("%s_%s_%s", domain, demandPartner, compassKey)
		// 	bkey := fmt.Sprintf("%s_%s_%s", domain, demandPartner, bothKey)

		// 	totalDomainsRevenue[domain] += row.Revenue

		// 	totalDomainsRevenuePerDP[ckey] += row.Revenue
		// 	totalDomainsRevenuePerDP[bkey] += row.Revenue

		// 	totalRevenuePerDP[fmt.Sprintf("%s_%s", demandPartner, compassKey)] += row.Revenue
		// 	totalRevenuePerDP[fmt.Sprintf("%s_%s", demandPartner, bothKey)] += row.Revenue
		// }

		for _, row := range nbResult.Data.Result {
			domain := strings.ToLower(row.Domain)
			demandPartner := strings.ToLower(row.DemandPartner)
			nbkey := fmt.Sprintf("%s_%s_%s", domain, demandPartner, newBidderKey)
			// bkey := fmt.Sprintf("%s_%s_%s", domain, demandPartner, bothKey)

			totalDomainsRevenue[domain] += row.Revenue

			totalDomainsRevenuePerDP[nbkey] += row.Revenue
			// totalDomainsRevenuePerDP[bkey] += row.Revenue

			totalRevenuePerDP[fmt.Sprintf("%s_%s", demandPartner, newBidderKey)] += row.Revenue
			// totalRevenuePerDP[fmt.Sprintf("%s_%s", demandPartner, bothKey)] += row.Revenue
		}

		a.lowPerformanceCache.totalRevenue = map[string]float64{
			compassKey:   0,
			bothKey:      0,
			newBidderKey: 0,
		}

		for key, value := range totalRevenuePerDP {
			if strings.HasSuffix(key, bothKey) {
				a.lowPerformanceCache.totalRevenue[bothKey] += value
			} else if strings.HasSuffix(key, compassKey) {
				a.lowPerformanceCache.totalRevenue[compassKey] += value
			} else {
				a.lowPerformanceCache.totalRevenue[newBidderKey] += value
			}
		}

		lowPerformance := make(map[string]bool)
		for key, value := range totalDomainsRevenuePerDP {
			domain := strings.Split(key, "_")[0]
			lowPerformance[key] = (value/totalDomainsRevenue[domain])*100 < 1
		}

		a.lowPerformanceCache.lowPerformance = lowPerformance
		a.lowPerformanceCache.revenuePerDP = totalRevenuePerDP
	}

	return nil
}

func (a *AdsTxtService) updateAdsTxtMetadata(ctx context.Context) error {
	logger.Logger(ctx).Info().Msg("start updating ads txt metadata")
	start := time.Now()

	data, err := a.GetGroupByDPAdsTxtTable(ctx, &AdsTxtGetGroupByDPOptions{})
	if err != nil {
		return fmt.Errorf("cannot get group by dp table to update ads txt metadata: %w", err)
	}

	logger.Logger(ctx).Info().Msg("sending data to update ads txt metadata")

	err = a.adstxtModule.UpdateAdsTxtMetadata(ctx, data)
	if err != nil {
		return fmt.Errorf("cannot update ads txt metadata: %w", err)
	}

	logger.Logger(ctx).Info().Msgf("ads txt metadata successfully updated in %v", time.Since(start).String())

	return nil
}

func buildGetLowPerfomanceRequestBody(firstOfMonth, lastOfMonth time.Time) []byte {
	return []byte(fmt.Sprintf(`
		{
			"data": {
				"date": {
					"range": [
						"%v",
						"%v"
					],
					"interval": "month"
				},
				"dimensions": [
					"Domain", "DemandPartner"
				],
				"metrics": [
					"Revenue"
				]
			}
		}`,
		firstOfMonth.Format("2006-01-02 00:00:00"),
		lastOfMonth.Format("2006-01-02 15:04:05"),
	))
}

func getGroupByDPTotal(ctx context.Context, ops *AdsTxtGetGroupByDPOptions) (int64, error) {
	type amountOfRows struct {
		Value int64 `boil:"total"`
	}

	totalMods := ops.Filter.queryModGroupByDP().Add(
		qm.Select(fmt.Sprintf("count(distinct %v) as total", groupByDPIDColumnName)),
	)
	var total amountOfRows
	err := models.AdsTXTGroupByDPViews(totalMods...).Bind(ctx, bcdb.DB(), &total)
	if err != nil {
		return 0, err
	}

	return total.Value, nil
}

func getCursorIDExpression(orderOps order.Sort, defaultColumnName string) string {
	const cursorIDExpression = "dense_rank() over (order by %s)"

	if len(orderOps) == 0 {
		return fmt.Sprintf(`%v as %v`, fmt.Sprintf(cursorIDExpression, defaultColumnName), cursorIDColumnName)
	}

	cursorIDOrderBy := make([]string, 0, len(orderOps)+1)
	for _, orderField := range orderOps {
		order := "ASC"
		if orderField.Desc {
			order = "DESC"
		}
		cursorIDOrderBy = append(cursorIDOrderBy, fmt.Sprintf("%v %v", orderField.Name, order))
	}
	cursorIDOrderBy = append(cursorIDOrderBy, fmt.Sprintf("%v ASC", defaultColumnName))

	return fmt.Sprintf(`%v as %v`, fmt.Sprintf(cursorIDExpression, strings.Join(cursorIDOrderBy, ",")), cursorIDColumnName)
}
