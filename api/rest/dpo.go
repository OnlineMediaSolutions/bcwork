package rest

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

// DemandPartnerOptimizationSetHandler Update demand partner optimization rule for a publisher.
// @Description Update demand partner optimization rule for a publisher.
// @Tags DPO
// @Accept json
// @Produce json
// @Param options body dto.DPORuleUpdateRequest true "Demand Partner Optimization update rule"
// @Success 200 {object} core.DemandPartnerOptimizationUpdateResponse
// @Security ApiKeyAuth
// @Router /dpo/set [post]
func (o *OMSNewPlatform) DemandPartnerOptimizationSetHandler(c *fiber.Ctx) error {
	data := &dto.DPORuleUpdateRequest{}
	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to parse metadata update payload", err)
	}

	ruleID, err := o.dpoService.SetDPORule(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to set dpo rule", err)
	}

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
func (o *OMSNewPlatform) DemandPartnerOptimizationGetHandler(c *fiber.Ctx) error {
	data := &core.DPOFactorOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error when parsing request body for /dpo/get", err)
	}

	pubs, err := o.dpoService.GetJoinedDPORule(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve DPO data", err)
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
func (o *OMSNewPlatform) DemandPartnerOptimizationDeleteHandler(c *fiber.Ctx) error {
	var dpoRules []string
	if err := c.BodyParser(&dpoRules); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to parse array of dpo rules to delete", err)
	}

	err := o.dpoService.DeleteDPORule(c.Context(), dpoRules)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Error in delete query execution", err)
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
func (o *OMSNewPlatform) DemandPartnerOptimizationUpdateHandler(c *fiber.Ctx) error {
	ruleId := c.Query("rid")
	factorStr := c.Query("factor")
	factor, err := strconv.ParseFloat(factorStr, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to parse factor", err)
	}

	err = o.dpoService.UpdateDPORule(c.Context(), ruleId, factor)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update dpo rule", err)
	}

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
