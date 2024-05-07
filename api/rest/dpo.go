package rest

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
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

// DemandPartnerOptimizationSetHandler Update demand partner optimization rule for a publisher.
// @Description Update demand partner optimization rule for a publisher.
// @Tags dpo
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
// @Tags dpo
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /dpo/get [get]
func DemandPartnerOptimizationGetHandler(c *fiber.Ctx) error {

	dpid := c.Query("dpid")
	if dpid == "" {
		c.SendString("'dpid' is mandatory")
		return c.SendStatus(http.StatusBadRequest)
	}

	c.Set("Content-Type", "application/json")

	dpo, err := models.Dpos(models.DpoWhere.DemandPartnerID.EQ(dpid)).One(c.Context(), bcdb.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to fetch dpo data")
	}

	return c.JSON(dpo)
}

// DemandPartnerOptimizationGetHandler Delete demand partner optimization rule for publisher.
// @Description Delete demand partner optimization rule for publisher.
// @Tags dpo
// @Produce json
// @Security ApiKeyAuth
// @Router /dpo/delete [delete]
func DemandPartnerOptimizationDeleteHandler(c *fiber.Ctx) error {

	publisher := c.Query("publisher")
	if publisher == "" {
		c.SendString("'publisher' is mandatory")
		return c.SendStatus(http.StatusBadRequest)
	}

	c.Set("Content-Type", "application/json")

	return c.JSON(map[string]interface{}{})

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
