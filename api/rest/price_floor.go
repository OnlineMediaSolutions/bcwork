package rest

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasttemplate"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/text/message"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func PriceFloorSetHandler(c *fiber.Ctx) error {

	publisher := c.Query("publisher")
	if publisher == "" {
		c.SendString("'publisher' is mandatory")
		return c.SendStatus(http.StatusBadRequest)
	}

	floor := c.Query("floor")
	if floor == "" {
		c.SendString("'floor' is mandatory")
		return c.SendStatus(http.StatusBadRequest)
	}

	domain := c.Query("domain")

	_, err := strconv.ParseFloat(floor, 64)
	if err != nil {
		c.SendString("failed to parse floor")
		return c.SendStatus(http.StatusBadRequest)
	}
	mod := models.MetadataQueue{
		Key:           "price:floor:" + publisher,
		TransactionID: bcguid.NewFromf(publisher, domain, time.Now()),
		Value:         []byte(floor),
	}

	if domain != "" {
		mod.Key = mod.Key + ":" + domain
	}

	if strings.ToLower(c.Query("mobile")) == "true" {
		mod.Key = "mobile:" + mod.Key
	}

	//log.Info().Interface("update", data).Msg("metadata update parsed")

	err = mod.Insert(c.Context(), bcdb.DB(), boil.Infer())
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to insert metadata update to queue")
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendStatus(http.StatusOK)
}

func PriceFloorGetHandler(c *fiber.Ctx) error {

	publisher := c.Query("publisher")
	if publisher == "" {
		c.SendString("'publisher' is mandatory")
		return c.SendStatus(http.StatusBadRequest)
	}

	domain := c.Query("domain")

	key := "price:floor:" + publisher
	if domain != "" {
		key = key + ":" + domain
	}
	if strings.ToLower(c.Query("mobile")) == "true" {
		key = "mobile:" + key
	}

	//log.Info().Interface("update", data).Msg("metadata update parsed")

	meta, err := models.MetadataQueues(models.MetadataQueueWhere.Key.EQ(key), qm.OrderBy("created_by desc")).One(c.Context(), bcdb.DB())
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch " + key)
		c.SendString("failed to fetch " + key)
		return c.SendStatus(http.StatusInternalServerError)
	}

	c.Set("Content-Type", "application/json")

	return c.JSON(meta)
}

func PriceFloorGetAllHandler(c *fiber.Ctx) error {

	query := `select metadata_queue.*
from metadata_queue,(select key,max(created_at) created_at from metadata_queue where key like '%price:floor%' group by key) last
where last.created_at=metadata_queue.created_at and last.key=metadata_queue.key order by metadata_queue.key`

	records := models.MetadataQueueSlice{}
	err := queries.Raw(query).Bind(c.Context(), bcdb.DB(), &records)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch all price floors")
		c.SendString("failed to fetch")
		return c.SendStatus(http.StatusInternalServerError)
	}

	c.Set("Content-Type", "text/html")
	b := bytes.Buffer{}
	p := message.NewPrinter(message.MatchLanguage("en"))

	for _, rec := range records {
		b.WriteString(p.Sprintf(rowFloorPrice, rec.Key, rec.Value, rec.CreatedAt.Format("2006-01-02 15:04"), rec.CommitedInstances))
	}
	t := fasttemplate.New(htmlFloorPrice, "{{", "}}")
	s := t.ExecuteString(map[string]interface{}{
		"data": b.String(),
	})

	return c.SendString(s)
}

var htmlFloorPrice = `
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
      Current Publisher Floors 
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
                  Floor
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

var rowFloorPrice = `<tr>
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
