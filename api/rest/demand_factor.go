package rest

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/valyala/fasttemplate"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/text/message"
)

type DemandFactorUpdateRequest struct {
	DemandPartner string  `json:"demand_partner"`
	Factor        float64 `json:"factor"`
}

// DemandFactorUpdateRespose
type DemandFactorUpdateRespose struct {
	// in: body
	Status string `json:"status"`
}

// StagingTestHandler Update new bidder demand factor for demand partner.
// @Description Update new bidder demand factor for demand partner.
// @Tags Staging Test
// @Accept json
// @Produce json
// @Success 200 {object}  utils.BaseResponse
// @Security ApiKeyAuth
// @Router /hello/world/get [get]
func StagingTestHandler(c *fiber.Ctx) error {

	return utils.SuccessResponse(c, fiber.StatusOK, "Hello staging env")

}

// DemandFactorPostHandler Update new bidder demand factor for demand partner.
// @Description Update new bidder demand factor for demand partner.
// @Tags MetaData
// @Accept json
// @Produce json
// @Param options body DemandFactorUpdateRequest true "DemandFactor update Options"
// @Success 200 {object} DemandFactorUpdateRespose
// @Security ApiKeyAuth
// @Router /demand/factor [post]
func DemandFactorPostHandler(c *fiber.Ctx) error {
	data := &DemandFactorUpdateRequest{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "failed to parse metadata update payload", err)
	}

	if data.DemandPartner == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "", errors.New("'demand_partner' is mandatory"))
	}

	factor := strconv.FormatFloat(data.Factor, 'f', 2, 64)
	if data.Factor > 1 || data.Factor < 0 {
		return utils.ErrorResponse(c, http.StatusBadRequest, "", errors.New("'factor' can only be between 0 and 1"))
	}
	//log.Info().Interface("update", data).Msg("metadata update parsed")
	mod := models.MetadataQueue{
		Key:           "demand:factor:" + data.DemandPartner,
		TransactionID: bcguid.NewFromf(data.DemandPartner, time.Now()),
		Value:         []byte(factor),
	}

	//log.Info().Interface("update", data).Msg("metadata update parsed")

	err := mod.Insert(c.Context(), bcdb.DB(), boil.Infer())
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "failed to insert metadata update to queue", err)
	}

	return c.JSON(DemandFactorUpdateRespose{
		Status: "ok",
	})
}

func DemandFactorSetHandler(c *fiber.Ctx) error {
	demand := c.Query("demand")
	if demand == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "", errors.New("'demand' is mandatory"))
	}

	factor := c.Query("factor")
	if factor == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "", errors.New("'factor' is mandatory"))
	}

	parseFactor, err := strconv.ParseFloat(factor, 64)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "", errors.New("failed to parse factor"))
	}
	if parseFactor > 1 || parseFactor < 0 {
		return utils.ErrorResponse(c, http.StatusBadRequest, "", errors.New("'factor' can only be between 0 and 1"))
	}

	mod := models.MetadataQueue{
		Key:           "demand:factor:" + demand,
		TransactionID: bcguid.NewFromf(demand, time.Now()),
		Value:         []byte(factor),
	}

	//log.Info().Interface("update", data).Msg("metadata update parsed")

	err = mod.Insert(c.Context(), bcdb.DB(), boil.Infer())
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "failed to insert metadata update to queue", err)
	}

	return c.SendStatus(http.StatusOK)
}

func DemandFactorGetHandler(c *fiber.Ctx) error {
	demand := c.Query("demand")
	if demand == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "", errors.New("'demand' is mandatory"))
	}

	key := "demand:factor:" + demand

	//log.Info().Interface("update", data).Msg("metadata update parsed")

	meta, err := models.MetadataQueues(models.MetadataQueueWhere.Key.EQ(key), qm.OrderBy("created_by desc")).One(c.Context(), bcdb.DB())
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "", fmt.Errorf("failed to fetch %s", key))
	}

	c.Set("Content-Type", "application/json")

	return c.JSON(meta)
}

func DemandFactorGetAllHandler(c *fiber.Ctx) error {
	query := `select metadata_queue.*
from metadata_queue,(select key,max(created_at) created_at from metadata_queue where key like '%demand:factor%' group by key) last
where last.created_at=metadata_queue.created_at and last.key=metadata_queue.key order by metadata_queue.key`

	records := models.MetadataQueueSlice{}
	err := queries.Raw(query).Bind(c.Context(), bcdb.DB(), &records)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch all demand factors", err)
	}

	if c.Query("format") == "html" {
		c.Set("Content-Type", "text/html")
		b := bytes.Buffer{}
		p := message.NewPrinter(message.MatchLanguage("en"))

		for _, rec := range records {
			b.WriteString(p.Sprintf(rowDemandFactor, rec.Key, rec.Value, rec.CreatedAt.Format("2006-01-02 15:04"), rec.CommitedInstances))
		}
		t := fasttemplate.New(htmlDemandFactor, "{{", "}}")
		s := t.ExecuteString(map[string]interface{}{
			"data": b.String(),
		})

		return c.SendString(s)
	} else {
		c.Set("Content-Type", "application/json")
		return c.JSON(records)
	}
}

var htmlDemandFactor = `
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

var rowDemandFactor = `<tr>
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
