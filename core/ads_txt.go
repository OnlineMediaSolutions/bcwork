package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	adstxt "github.com/m6yf/bcwork/modules/ads_txt"
	"github.com/m6yf/bcwork/modules/compass"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/m6yf/bcwork/modules/logger"
	"github.com/spf13/viper"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

type AdsTxtService struct {
	historyModule       history.HistoryModule
	compassModule       compass.CompassModule
	adstxtModule        adstxt.AdsTxtLinesCreater
	lowPerformanceCache *LowPerformanceCache
}

func NewAdsTxtService(
	ctx context.Context,
	historyModule history.HistoryModule,
	compassModule compass.CompassModule,
	adstxtModule adstxt.AdsTxtLinesCreater,
) *AdsTxtService {
	service := &AdsTxtService{
		historyModule:       historyModule,
		compassModule:       compassModule,
		adstxtModule:        adstxtModule,
		lowPerformanceCache: &LowPerformanceCache{},
	}

	// updating ads.txt metadata every N minutes
	go func(ctx context.Context) {
		var (
			start                  time.Time
			defaultMinutesToUpdate time.Duration = 60 * time.Minute
		)

		minutesToUpdate := viper.GetDuration(config.AdsTxtMetadataUpdateRateKey) * time.Minute
		if minutesToUpdate == 0 {
			minutesToUpdate = defaultMinutesToUpdate
		}

		ticker := time.NewTicker(minutesToUpdate)

		for {
			select {
			case <-ticker.C:
				logger.Logger(ctx).Info().Msg("start updating ads txt metadata")
				start = time.Now()

				data, err := service.GetGroupByDPAdsTxtTable(ctx, nil)
				if err != nil {
					logger.Logger(ctx).Err(err).Msg("cannot get group by dp table to update ads txt metadata")
					continue
				}

				logger.Logger(ctx).Info().Msg("sending data to update ads txt metadata")

				err = service.adstxtModule.UpdateAdsTxtMetadata(ctx, data)
				if err != nil {
					logger.Logger(ctx).Err(err).Msg("cannot update ads txt metadata")
					continue
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
			logger.Logger(ctx).Info().Msgf("ads txt metadata successfully updated in %v", time.Since(start).String())
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

// TODO: pagination
type AdsTxtGetOptions struct {
	Filter     AdsTxtGetFilter        `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
}

type AdsTxtGetFilter struct {
	// TODO: filters
	PublisherID               string `json:"publisher_id"`
	Domain                    string `json:"domain"`
	DomainStatus              string `json:"domain_status"`
	DemandPartnerNameExtended string `json:"demand_partner_name_extended"`
	AdsTxtLine                string `json:"ads_txt_line"`
	// dpc.id as demand_partner_connection_id,
	// dpc."media_type",
	// d.manager_id as demand_manager_id,
	// at2.demand_status,
	// at2.status,
	// d.is_approval_needed,
	// dpc.is_required_for_ads_txt as is_required,
	// at2.last_scanned_at,
	// at2.error_message
	DemandPartnerId       filter.StringArrayFilter `json:"demand_partner_id,omitempty"`
	DemandPartnerName     filter.StringArrayFilter `json:"demand_partner_name,omitempty"`
	IsDemandPartnerActive *filter.BoolFilter       `json:"is_demand_partner_active,omitempty"`
}

func (filter *AdsTxtGetFilter) QueryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)
	if filter == nil {
		return mods
	}

	if len(filter.DemandPartnerId) > 0 {
		mods = append(mods, filter.DemandPartnerId.AndIn(models.DpoColumns.DemandPartnerID))
	}

	if len(filter.DemandPartnerName) > 0 {
		mods = append(mods, filter.DemandPartnerName.AndIn(models.DpoColumns.DemandPartnerName))
	}

	if filter.IsDemandPartnerActive != nil {
		mods = append(mods, filter.IsDemandPartnerActive.Where(models.DpoColumns.Active))
	}

	return mods
}

func (a *AdsTxtService) GetMainAdsTxtTable(ctx context.Context, ops *AdsTxtGetOptions) ([]*dto.AdsTxt, error) {
	query := fmt.Sprintf(`
		with main_table as (
			select 
				dense_rank() over (order by t.id) as cursor_id,
				t.*,
				p."name" as publisher_name,
				p.account_manager_id,
				p.campaign_manager_id
			from (
				%v
				union 
				%v
				union 
				%v
			) as t
			join publisher p on p.publisher_id = t.publisher_id
		)
		select 
			*
		from main_table
		where true %v
		;
	`,
		adsTxtDemandPartnerConnectionBaseQuery,
		fmt.Sprintf(adsTxtdemandPartnerChildrenBaseQuery, "d.demand_partner_name"),
		fmt.Sprintf(adsTxtSeatOwnersBaseQuery, "''", "''", "0", "null", "true", dynamicPublisherIDPlaceholder, ""),
		ops.Pagination.GetWhereClause("cursor_id"), // TODO: to constants
	)

	mainTable, err := a.getAdsTxtTableByQueryWithUsersFullNames(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve main table: %w", err)
	}

	return mainTable, nil
}

func (a *AdsTxtService) GetGroupByDPAdsTxtTable(ctx context.Context, ops *AdsTxtGetOptions) (map[string]*dto.AdsTxtGroupedByDPData, error) {
	query := fmt.Sprintf(`
		with group_by_dp_table as (
			select
				dense_rank() over (order by t.publisher_id, t."domain", t.demand_partner_name, t.demand_partner_connection_id) as cursor_id,
				t.*,
				p."name" as publisher_name,
				p.account_manager_id,
				p.campaign_manager_id,
				sum(case 
					when t.status = 'added' then 1
					else 0
				end) over (partition by t.publisher_id, t."domain", t.demand_partner_name, t.demand_partner_connection_id) as added, 
				count(t.status) over (partition by t.publisher_id, t."domain", t.demand_partner_name, t.demand_partner_connection_id) as total,
				bool_and(case 
					when t.status = 'added' AND t.is_required and t.demand_status = 'approved' then true
					when not t.is_required then true
					else false
				end) over (partition by t.publisher_id, t."domain", t.demand_partner_name, t.demand_partner_connection_id) as is_ready_to_go_live
			from (
				%v
				where dpc.is_required_for_ads_txt
				union 
				%v
				union all
				%v
				join demand_partner_connection dpc on d.demand_partner_id = dpc.demand_partner_id 
			) as t
			join publisher p on t.publisher_id = p.publisher_id 
			where t.is_demand_partner_active
			order by t.publisher_id, t."domain", t.demand_partner_name, t.demand_partner_connection_id, t.demand_partner_name_extended
		)
		select
			*
		from group_by_dp_table
		where true %v
		;
	`,
		adsTxtDemandPartnerConnectionBaseQuery,
		fmt.Sprintf(adsTxtdemandPartnerChildrenBaseQuery, "d.demand_partner_name"),
		fmt.Sprintf(adsTxtSeatOwnersBaseQuery, "d.demand_partner_id", "d.demand_partner_name", "dpc.id", "dpc.media_type", "d.active", dynamicPublisherIDPlaceholder, "join dpo d on d.seat_owner_id = so.id"),
		ops.Pagination.GetWhereClause("cursor_id"),
	)

	var rawTable []*dto.AdsTxt
	err := queries.Raw(query).Bind(ctx, bcdb.DB(), &rawTable)
	if err != nil {
		return nil, err
	}

	usersMap, err := getUsersMap(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	demandPartnersWithConnectionsMap := make(map[string]struct{})
	groupByDpTable := make(map[string]*dto.AdsTxtGroupedByDPData)
	for _, row := range rawTable {
		demandPartnersWithConnectionsMap[fmt.Sprintf("%v:%v", row.DemandPartnerName, row.DemandPartnerConnectionID.Int)] = struct{}{}

		name := fmt.Sprintf("%v:%v:%v:%v", row.PublisherID, row.Domain, row.DemandPartnerName, row.DemandPartnerConnectionID.Int)

		row.AccountManagerFullName = usersMap[row.AccountManagerID.String]
		row.CampaignManagerFullName = usersMap[row.CampaignManagerID.String]
		row.DemandManagerFullName = usersMap[row.DemandManagerID.String]

		dpData, ok := groupByDpTable[name]
		if !ok {
			dpData = &dto.AdsTxtGroupedByDPData{
				Parent:   row,
				Children: []*dto.AdsTxt{row},
			}
			groupByDpTable[name] = dpData
		} else {
			dpData.ProcessParentRow(row)
			dpData.Children = append(dpData.Children, row)
		}
	}

	// mirroring
	mirroredDomains, err := models.PublisherDomains(
		models.PublisherDomainWhere.MirrorPublisherID.IsNotNull(),
	).
		All(ctx, bcdb.DB())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get domains with mirror feature: %w", err)
	}

	for _, mirroredDomain := range mirroredDomains {
		for demandPartnerWithConnection := range demandPartnersWithConnectionsMap {
			sourceKey := fmt.Sprintf("%v:%v:%v", mirroredDomain.PublisherID, mirroredDomain.Domain, demandPartnerWithConnection)
			mirroredKey := fmt.Sprintf("%v:%v:%v", mirroredDomain.MirrorPublisherID.String, mirroredDomain.Domain, demandPartnerWithConnection)

			sourceValue, ok := groupByDpTable[sourceKey]
			if !ok {
				logger.Logger(ctx).Warn().Str("source_key", sourceKey).Msg("cannot get source data for mirroring update")
				continue
			}

			mirroredValue, ok := groupByDpTable[mirroredKey]
			if !ok {
				logger.Logger(ctx).Warn().Str("mirrored_key", mirroredKey).Msg("cannot get mirrored data for mirroring update")
				continue
			}

			sourceValue.Parent.Mirror(mirroredValue.Parent)

			if len(sourceValue.Children) != len(mirroredValue.Children) {
				logger.Logger(ctx).Warn().
					Int("len_source_children", len(sourceValue.Children)).
					Int("len_mirrored_children", len(mirroredValue.Children)).
					Msg("source and mirrored children have different length")
				continue
			}

			for i := range sourceValue.Children {
				if sourceValue.Children[i].DemandPartnerNameExtended != mirroredValue.Children[i].DemandPartnerNameExtended {
					logger.Logger(ctx).Warn().
						Str("source_child_demand_partner_name", sourceValue.Children[i].DemandPartnerNameExtended).
						Str("mirrored_child_demand_partner_name", mirroredValue.Children[i].DemandPartnerNameExtended).
						Msg("source and mirrored children have different demand partner names")
					continue
				}
				sourceValue.Children[i].Mirror(mirroredValue.Children[i])
			}
		}
	}

	return groupByDpTable, nil
}

func (a *AdsTxtService) GetAMAdsTxtTable(ctx context.Context, ops *AdsTxtGetOptions) ([]*dto.AdsTxt, error) {
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
				p.campaign_manager_id
			from (
				%v
				union 
				%v
				union all
				%v
			) as t
			join publisher p on p.publisher_id = t.publisher_id
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
		amTable, err = a.getAdsTxtTableByQueryWithUsersFullNames(ctx, query)
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

func (a *AdsTxtService) GetCMAdsTxtTable(ctx context.Context, ops *AdsTxtGetOptions) ([]*dto.AdsTxt, error) {
	query := fmt.Sprintf(`
		with cm as (
			select 
				t.*,
				p."name" as publisher_name,
				p.account_manager_id,
				p.campaign_manager_id,
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

	cmTable, err := a.getAdsTxtTableByQueryWithUsersFullNames(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve cm table: %w", err)
	}

	return cmTable, nil
}

func (a *AdsTxtService) GetMBAdsTxtTable(ctx context.Context, ops *AdsTxtGetOptions) ([]*dto.AdsTxt, error) {
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
			-- children
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
	mod, err := models.AdsTXTS(models.AdsTXTWhere.ID.EQ(data.ID)).One(ctx, bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed to get ads txt line with id [%v] to update: %w", data.ID, err)
	}

	columns := make([]string, 0, 3)
	columns = append(columns, models.AdsTXTColumns.UpdatedAt)
	mod.UpdatedAt = null.TimeFrom(time.Now())

	if mod.DomainStatus != data.DomainStatus {
		mod.DomainStatus = data.DomainStatus
		columns = append(columns, models.AdsTXTColumns.DomainStatus)
	}

	if mod.DemandStatus != data.DemandStatus {
		mod.DemandStatus = data.DemandStatus
		columns = append(columns, models.AdsTXTColumns.DemandStatus)
	}

	if len(columns) == 1 {
		return errors.New("there are no new values to update ads txt line")
	}

	_, err = mod.Update(ctx, bcdb.DB(), boil.Whitelist(columns...))
	if err != nil {
		return fmt.Errorf("failed to update ads txt line: %w", err)
	}

	return nil
}

func (a *AdsTxtService) getAdsTxtTableByQueryWithUsersFullNames(ctx context.Context, query string) ([]*dto.AdsTxt, error) {
	var table []*dto.AdsTxt
	err := queries.Raw(query).Bind(ctx, bcdb.DB(), &table)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve ads txt lines: %w", err)
	}

	usersMap, err := getUsersMap(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	for i := range table {
		table[i].AccountManagerFullName = usersMap[table[i].AccountManagerID.String]
		table[i].CampaignManagerFullName = usersMap[table[i].CampaignManagerID.String]
		table[i].DemandManagerFullName = usersMap[table[i].DemandManagerID.String]
	}

	return table, nil
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

func getUsersMap(ctx context.Context) (map[string]string, error) {
	users, err := models.Users().All(ctx, bcdb.DB())
	if err != nil {
		return nil, err
	}

	usersMap := make(map[string]string, len(users))
	for _, user := range users {
		userID := strconv.Itoa(user.ID)
		usersMap[userID] = user.FirstName + " " + user.LastName
	}

	return usersMap, nil
}
