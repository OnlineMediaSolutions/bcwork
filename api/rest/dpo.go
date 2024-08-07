package rest

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/rs/zerolog/log"
	"strconv"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

// DemandPartnerOptimizationGetHandler Get demand partner optimization rules for publisher.
// @Description Get demand partner optimization rules for publisher.
// @Tags DPO
// @Param options body core.DPOGetOptions true "options"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /dp/get [post]
func DemandPartnerGetHandler(c *fiber.Ctx) error {

	data := &core.DPOGetOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body parsing error")
	}

	pubs, err := core.GetDpos(c.Context(), data)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve DPO'/s")
	}
	return c.JSON(pubs)
}

// DemandPartnerOptimizationSetHandler Update demand partner optimization rule for a publisher.
// @Description Update demand partner optimization rule for a publisher.
// @Tags DPO
// @Accept json
// @Produce json
// @Param options body core.DemandPartnerOptimizationUpdateRequest true "Demand Partner Optimization update rule"
// @Success 200 {object} core.DemandPartnerOptimizationUpdateResponse
// @Security ApiKeyAuth
// @Router /dpo/set [post]
func DemandPartnerOptimizationSetHandler(c *fiber.Ctx) error {

	data := &core.DemandPartnerOptimizationUpdateRequest{}
	err := c.BodyParser(&data)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to parse metadata update payload")
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
			log.Error().Err(err).Msg("Failed to update RT metadata for dpo")
		}
	}()

	return utils.DpoSuccessResponse(c, fiber.StatusOK, ruleID, "Dpo successfully added")
}

// DemandPartnerOptimizationGetHandler Get demand partner optimization rules for publisher.
// @Description Get demand partner optimization rules for publisher.
// @Tags DPO
// @Param options body core.DPOFactorOptions true "options"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /dpo/get [post]
func DemandPartnerOptimizationGetHandler(c *fiber.Ctx) error {

	data := &core.DPOFactorOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error when parsing request body for /dpo/get")
	}
	pubs, err := core.GetJoinedDPORule(c.Context(), data)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to retrieve DPO data, %s", err))
	}
	return c.JSON(pubs)
}

// DemandPartnerOptimizationGetHandler Delete demand partner optimization rule for publisher.
// @Description Delete demand partner optimization rule for publisher.
// @Tags DPO
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param options body []string true "options"
// @Router /dpo/delete [delete]
func DemandPartnerOptimizationDeleteHandler(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/json")
	var dpoRules []string
	if err := c.BodyParser(&dpoRules); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("Failed to parse array of dpo rules to delete, %s", err))
	}
	deleteQuery := core.CreateDeleteQuery(dpoRules)

	_, err := queries.Raw(deleteQuery).Exec(bcdb.DB())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("%s", err.Error()))
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "DPO rules were deleted")
}

// DemandPartnerOptimizationUpdateHandler Update demand partner optimization rule by rule id.
// @Description Update demand partner optimization rule by rule id..
// @Tags DPO
// @Param rid query string true "rule ID"
// @Param factor query int true "factor (0-100)"
// @Produce json
// @Security ApiKeyAuth
// @Router /dpo/update [get]
func DemandPartnerOptimizationUpdateHandler(c *fiber.Ctx) error {

	ruleId := c.Query("rid")
	factorStr := c.Query("factor")
	factor, err := strconv.ParseFloat(factorStr, 64)
	c.Set("Content-Type", "application/json")

	rule, err := models.DpoRules(models.DpoRuleWhere.RuleID.EQ(ruleId)).One(c.Context(), bcdb.DB())

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("Failed to delete dpo rule, %s", err))
	}

	rule.Factor = factor
	rule.Active = true
	updated, err := rule.Update(c.Context(), bcdb.DB(), boil.Whitelist(models.DpoRuleColumns.Factor, models.DpoRuleColumns.Active))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("Failed to delete dpo rule, %s", err))
	}

	if updated > 0 {
		go func() {
			err := core.SendToRT(context.Background(), rule.DemandPartnerID)
			if err != nil {
				log.Error().Err(err).Msg("Failed to update RT metadata for dpo")
			}
		}()
	}

	c.Set("Content-Type", "application/json")
	return utils.SuccessResponse(c, fiber.StatusOK, "Ok")
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
