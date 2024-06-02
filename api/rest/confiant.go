package rest

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"net/http"
	"time"
)

type ConfiantUpdateRequest struct {
	Publisher string  `json:"publisher_id"`
	Domain    string  `json:"domain"`
	Hash      string  `json:"confiant_key"`
	Rate      float64 `json:"rate"`
}

// ConfiantPostHandler Update and enable Confiant setup
// @Description Update and enable Confiant setup (publisher is mandatory, domain is optional)
// @Tags metadata
// @Accept json
// @Produce json
// @Param options body ConfiantUpdateRequest true "Confiant update Options"
// @Success 200 {object} SendStatus
// @Security ApiKeyAuth
// @Router /confiant [post]
func ConfiantPostHandler(c *fiber.Ctx) error {

	data := &ConfiantUpdateRequest{}
	if err := c.BodyParser(&data); err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to parse metadata update payload")
		return c.SendStatus(http.StatusBadRequest)
	}

	if data.Publisher == "" {
		c.SendString("'publisher' is mandatory")
		return c.SendStatus(http.StatusBadRequest)
	}

	errMessage := updateMetaDataQueue(c, data)
	if len(errMessage) != 0 {
		return c.Status(http.StatusInternalServerError).JSON(Response{Status: "error", Message: errMessage})
	}

	err := updateConfiant(c, data)
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("Failed to update Confiant table with the following")
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).JSON(Response{Status: "ok", Message: "Confiant table was successfully updated"})
}

func updateMetaDataQueue(c *fiber.Ctx, data *ConfiantUpdateRequest) string {

	val, err := json.Marshal(data.Hash)
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("Failed to parse hash value")
		return "Failed to parse hash value"
	}

	mod := models.MetadataQueue{
		Key:           "confiant:" + data.Publisher,
		TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
		Value:         val,
	}

	if data.Domain != "" {
		mod.Key = mod.Key + ":" + data.Domain
	}

	err = mod.Insert(c.Context(), bcdb.DB(), boil.Infer())
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("Failed to insert metadata update to queue")
		return "failed to insert metadata update to queue"
	}
	return ""
}

func updateConfiant(c *fiber.Ctx, data *ConfiantUpdateRequest) error {

	modConf := models.Confiant{
		PublisherID: data.Publisher,
		ConfiantKey: data.Hash,
		Rate:        data.Rate,
		Domain:      data.Domain,
	}

	return modConf.Upsert(c.Context(), bcdb.DB(), true, []string{models.ConfiantColumns.ConfiantKey, models.ConfiantColumns.PublisherID, models.ConfiantColumns.Domain}, boil.Infer(), boil.Infer())
}

// ConfiantGetHandler Get confiant setup
// @Description Get confiant setup
// @Tags confiant
// @Accept json
// @Produce json
// @Param options body core.GetConfiantOptions true "options"
// @Success 200 {object} core.ConfiantSlice
// @Router /confiant/get [post]
func ConfiantGetAllHandler(c *fiber.Ctx) error {

	data := &core.GetConfiantOptions{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(500).JSON(Response{Status: "error", Message: "error when parsing request body"})
	}

	pubs, err := core.GetConfiants(c.Context(), data)
	if err != nil {
		return c.Status(400).JSON(Response{Status: "error", Message: "failed to retrieve confiants"})
	}
	return c.JSON(pubs)
}

var htmlConfiant = `
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
      Confiant Setup Status 
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
                  Hash
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

var rowConfiant = `<tr>
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
