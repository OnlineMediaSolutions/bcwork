package rest

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"net/http"
	"strconv"
	"time"
)

// DemandReportGetRequest contains filter parameters for retrieving events
type DemandPartnerOptimizationUpdateRequest struct {
	DemandPartner string `json:"demand_partner_id"`
	Publisher     string `json:"publisher"`
	Domain        string `json:"domain,omitempty"`
	Country       string `json:"country,omitempty"`
	OS            string `json:"os,omitempty"`
	DeviceType    string `json:"device_type,omitempty"`
	PlacementType string `json:"placement_type,omitempty"`
	Browser       string `json:"browser,omitempty"`

	Factor float64 `json:"factor"`
}

// DemandPartnerOptimizationUpdateRespose
type DemandPartnerOptimizationUpdateResponse struct {
	// in: body
	Status string `json:"status"`
	RuleID string `json:"rule_id"`
}

var dpo_query = `SELECT dpo.*, dpo_rule.*, publisher.name
FROM dpo
JOIN dpo_rule on dpo.demand_partner_id = dpo_rule.demand_partner_id
JOIN publisher on dpo_rule.publisher = publisher.publisher_id`

var dpo_where_query = ` WHERE dpo.demand_partner_id = '%s'`

// DemandPartnerOptimizationSetHandler Update demand partner optimization rule for a publisher.
// @Description Update demand partner optimization rule for a publisher.
// @Tags DPO
// @Accept json
// @Produce json
// @Param options body DemandPartnerOptimizationUpdateRequest true "Demand Partner Optimization update rule"
// @Success 200 {object} DemandPartnerOptimizationUpdateResponse
// @Security ApiKeyAuth
// @Router /dpo/set [post]
func DemandPartnerOptimizationSetHandler(c *fiber.Ctx) error {

	data := &DemandPartnerOptimizationUpdateRequest{}
	if err := c.BodyParser(&data); err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to parse metadata update payload")

		return c.SendStatus(http.StatusBadRequest)
	}

	if data.DemandPartner == "" {
		c.SendString("'demand_partner_id' is mandatory")
		return c.SendStatus(http.StatusBadRequest)
	}

	//_ := strconv.FormatFloat(data.Factor, 'f', 2, 64)
	if data.Factor < 0 || data.Factor > 100 {
		c.SendString("'factor' should be a positive number  <= 100")
		return c.SendStatus(http.StatusBadRequest)
	}

	dpoRule := core.DemandPartnerOptimizationRule{
		DemandPartner: data.DemandPartner,
		Publisher:     data.Publisher,
		Domain:        data.Domain,
		Country:       data.Country,
		OS:            data.OS,
		DeviceType:    data.DeviceType,
		PlacementType: data.PlacementType,
		Browser:       data.Browser,
		Factor:        data.Factor,
	}

	ruleID, err := dpoRule.Save(c.Context())
	if err != nil {
		return err
	}

	go func() {
		err := core.SendToRT(context.Background(), data.DemandPartner)
		if err != nil {
			log.Error().Err(err).Msg("failed to update RT metadata for dpo")
		}
	}()

	return c.JSON(DemandPartnerOptimizationUpdateResponse{
		Status: "ok",
		RuleID: ruleID,
	})
}

// DemandPartnerOptimizationGetHandler Get demand partner optimization rules for publisher.
// @Description Get demand partner optimization rules for publisher.
// @Tags DPO
// @Param dpid query string true "demand partner ID"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /dpo/get [get]
func DemandPartnerOptimizationGetHandler(c *fiber.Ctx) error {

	dpid := c.Query("dpid")

	c.Set("Content-Type", "application/json")

	var results []JoinedDpo
	query := buildQuery(dpid)
	err := queries.Raw(query).Bind(c.Context(), bcdb.DB(), &results)

	if err != nil {
		return errors.Wrapf(err, "Failed to fetch dpo data")
	}

	return c.JSON(results)
}

func buildQuery(dpid string) string {
	if dpid == "" {
		return dpo_query
	} else {
		return fmt.Sprintf(dpo_query+dpo_where_query, dpid)
	}
}

type JoinedDpo struct {
	DemandPartnerID string      `json:"demand_partner_id"`
	IsInclude       bool        `json:"is_include"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       *time.Time  `json:"updated_at"`
	RuleId          string      `json:"rule_id"`
	Publisher       string      `json:"publisher"`
	Domain          string      `json:"domain"`
	Country         string      `json:"country"`
	Browser         null.String `json:"browser"`
	OS              null.String `json:"os,omitempty"`
	DeviceType      string      `json:"device_type"`
	PlacementType   null.String `json:"placement_type"`
	Factor          float64     `json:"factor"`
	Name            string      `json:"name"`
}

// DemandPartnerOptimizationGetHandler Delete demand partner optimization rule for publisher.
// @Description Delete demand partner optimization rule for publisher.
// @Tags dpo
// @Param rid query string true "rule ID"
// @Produce json
// @Security ApiKeyAuth
// @Router /dpo/delete [delete]
func DemandPartnerOptimizationDeleteHandler(c *fiber.Ctx) error {

	ruleId := c.Query("rid")
	if ruleId == "" {
		c.SendString("'rid' (rule id_ is mandatory")
		return c.SendStatus(http.StatusBadRequest)
	}

	c.Set("Content-Type", "application/json")

	rule, err := models.DpoRules(models.DpoRuleWhere.RuleID.EQ(ruleId)).One(c.Context(), bcdb.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to fetch dpo rule")
	}

	deleted, err := rule.Delete(c.Context(), bcdb.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to delete dpo rule")
	}

	if deleted > 0 {
		go func() {
			err := core.SendToRT(context.Background(), rule.DemandPartnerID)
			if err != nil {
				log.Error().Err(err).Msg("failed to update RT metadata for dpo")
			}
		}()
	}

	c.Set("Content-Type", "application/json")

	return c.JSON(map[string]interface{}{"status": "ok"})

}

// DemandPartnerOptimizationUpdateHandler Update demand partner optimization rule by rule id.
// @Description Update demand partner optimization rule by rule id..
// @Tags dpo
// @Param rid query string true "rule ID"
// @Param factor query int true "factor (0-100)"
// @Produce json
// @Security ApiKeyAuth
// @Router /dpo/update [get]
func DemandPartnerOptimizationUpdateHandler(c *fiber.Ctx) error {

	ruleId := c.Query("rid")
	if ruleId == "" {
		c.SendString("'rid' (rule id_ is mandatory")
		return c.SendStatus(http.StatusBadRequest)
	}

	factorStr := c.Query("factor")
	if factorStr == "" {
		c.SendString("'factor' (factor is mandatory (0-100)")
		return c.SendStatus(http.StatusBadRequest)
	}

	factor, err := strconv.ParseFloat(factorStr, 64)
	if err != nil {
		c.SendString("'factor' should be numeric (0-100)")
		return c.SendStatus(http.StatusBadRequest)
	}
	if factor > 100 || factor < 0 {
		c.SendString("'factor' should be numeric (0-100)")
		return c.SendStatus(http.StatusBadRequest)
	}

	c.Set("Content-Type", "application/json")

	rule, err := models.DpoRules(models.DpoRuleWhere.RuleID.EQ(ruleId)).One(c.Context(), bcdb.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to fetch dpo rule")
	}

	rule.Factor = factor

	updated, err := rule.Update(c.Context(), bcdb.DB(), boil.Whitelist(models.DpoRuleColumns.Factor))
	if err != nil {
		return errors.Wrapf(err, "failed to delete dpo rule")
	}

	if updated > 0 {
		go func() {
			err := core.SendToRT(context.Background(), rule.DemandPartnerID)
			if err != nil {
				log.Error().Err(err).Msg("failed to update RT metadata for dpo")
			}
		}()
	}

	c.Set("Content-Type", "application/json")

	return c.JSON(map[string]interface{}{"status": "ok"})

}

var htmlDemandPartnerOptimization = `
<html>
<head>
     <link href="https://unpkg.com/tailwindcss@^1.0/dist/tailwind.min.css" rel="stylesheet">
</head>
<body>
<div class="md:flex justify-center md:items-center">
   <div class="mt-1 flex md:mt-0 md:ml-4">
    <img class="filter invert h-40 w-40" src="https://onlinemediasolutions.com/wp-content/themes/brightcom/assets/images/oms-logo.svg" alt="">
  </div>
<div class="min-w-0">
    <h2 class="p-3 text-2xl font-bold leading-7 text-purple-600 sm:text-3xl sm:truncate">
      Current Publisher Factors 
    </h2>
  </div>
 
</div>


<div class="flex flex-col">
  <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
    <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
      <div class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" class="font-bold px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                Key
              </th>
               <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                  Factor
               </th>
               <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                  Create At
               </th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-900 uppercase tracking-wider">
                  Committed
               </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
              {{data}}
          </tbody>
        </table>
      </div>
    </div>
  </div>
</div>
</body>
</html>`

var rowDemandPartnerOptimization = `<tr>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     %s
                 </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     %s
                  </td>
                   <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     %s
                 </td>
                 <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                     %d
                 </td>
                        
            </tr>`
