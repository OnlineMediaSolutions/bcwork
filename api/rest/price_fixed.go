package rest

import (
	"bytes"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasttemplate"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"golang.org/x/text/message"
)

// DemandReportGetRequest contains filter parameters for retrieving events
type FixedPriceUpdateRequest struct {
	Publisher string  `json:"publisher"`
	Domain    string  `json:"domain"`
	Price     float64 `json:"price"`
	Mobile    bool    `json:"mobile"`
}

// FixedPriceUpdateRespose
type FixedPriceUpdateRespose struct {
	// in: body
	Status string `json:"status"`
}

// FixedPricePostHandler Update javascript tag guaranteed price
// @Description Update javascript tag guaranteed price
// @Tags MetaData
// @Accept json
// @Produce json
// @Param options body FixedPriceUpdateRequest true "FixedPrice update Options"
// @Success 200 {object} FixedPriceUpdateRespose
// @Security ApiKeyAuth
// @Router /price/fixed [post]
func FixedPricePostHandler(c *fiber.Ctx) error {
	data := &FixedPriceUpdateRequest{}
	if err := c.BodyParser(&data); err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to parse metadata update payload")

		return c.SendStatus(http.StatusBadRequest)
	}

	if data.Publisher == "" {
		c.SendString("'publisher' is mandatory")
		return c.SendStatus(http.StatusBadRequest)
	}

	if data.Domain == "" {
		c.SendString("'domain' is mandatory")
		return c.SendStatus(http.StatusBadRequest)
	}

	price := strconv.FormatFloat(data.Price, 'f', 2, 64)
	//log.Info().Interface("update", data).Msg("metadata update parsed")
	mod := models.MetadataQueue{
		Key:           "price:fixed:" + data.Publisher + ":" + data.Domain,
		TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
		Value:         []byte(price),
	}

	if data.Mobile {
		mod.Key = "mobile:" + mod.Key
	}

	//log.Info().Interface("update", data).Msg("metadata update parsed")

	err := mod.Insert(c.Context(), bcdb.DB(), boil.Infer())
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to insert metadata update to queue")
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.JSON(FixedPriceUpdateRespose{
		Status: "ok",
	})
}

// FixedPriceGetHandler Get fixed price rates
// @Description Get fixed price rates
// @Tags MetaData
// @Accept json
// @Produce html
// @Security ApiKeyAuth
// @Router /price/fixed [get]
func FixedPriceGetAllHandler(c *fiber.Ctx) error {
	query := `select metadata_queue.*
from metadata_queue,(select key,max(created_at) created_at from metadata_queue where key like '%price:fixed%' group by key) last
where last.created_at=metadata_queue.created_at and last.key=metadata_queue.key order by metadata_queue.key`

	records := models.MetadataQueueSlice{}
	err := queries.Raw(query).Bind(c.Context(), bcdb.DB(), &records)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch all price factors", err)
	}

	if c.Query("format") == "html" {
		c.Set("Content-Type", "text/html")
		b := bytes.Buffer{}
		p := message.NewPrinter(message.MatchLanguage("en"))

		for _, rec := range records {
			b.WriteString(p.Sprintf(rowFixedPrice, rec.Key, rec.Value, rec.CreatedAt.Format("2006-01-02 15:04"), rec.CommitedInstances))
		}
		t := fasttemplate.New(htmlFixedPrice, "{{", "}}")
		s := t.ExecuteString(map[string]interface{}{
			"data": b.String(),
		})

		return c.SendString(s)
	} else {
		c.Set("Content-Type", "application/json")
		return c.JSON(records)
	}
}

var htmlFixedPrice = `
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
      Current Publisher Fixed Price Rates 
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
                  Price
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

var rowFixedPrice = `<tr>
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
