package rest

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasttemplate"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"golang.org/x/text/message"
	"net/http"
	"time"
)

// DemandReportGetRequest contains filter parameters for retrieving events
type BlockUpdateRequest struct {
	Publisher string   `json:"publisher"`
	Domain    string   `json:"domain"`
	BCAT      []string `json:"bcat"`
	BADV      []string `json:"badv"`
}

type BlockGetRequest struct {
	Publisher string `json:"publisher"`
	Domain    string `json:"domain"`
}

// BlockUpdateRespose
type BlockUpdateRespose struct {
	// in: body
	Status string `json:"status"`
}

var query = `SELECT metadata_queue.*
FROM metadata_queue,(select key,max(created_at) created_at FROM metadata_queue WHERE key LIKE 'bcat:%' OR key like 'badv:%' group by key) last
WHERE last.created_at=metadata_queue.created_at
    AND last.key=metadata_queue.key `

var sortQuery = ` ORDER by metadata_queue.key`

// BlockPostHandler Update bidder addomain and categories blocks
// @Description Update bidder addomain and categories blocks.
// @Tags metadata
// @Accept json
// @Produce json
// @Param options body BlockUpdateRequest true "Block update Options"
// @Success 200 {object} BlockUpdateRespose
// @Security ApiKeyAuth
// @Router /block [post]
func BlockPostHandler(c *fiber.Ctx) error {

	data := &BlockUpdateRequest{}
	if err := c.BodyParser(&data); err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to parse metadata update payload")

		return c.SendStatus(http.StatusBadRequest)
	}

	if data.Publisher == "" {
		c.SendString("'publisher' is mandatory")
		return c.SendStatus(http.StatusBadRequest)
	}

	if data.BCAT != nil {
		b, err := json.Marshal(data.BCAT)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		mod := models.MetadataQueue{
			Key:           "bcat:" + data.Publisher,
			TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, "bcat", time.Now()),
			Value:         b,
		}

		if data.Domain != "" {
			mod.Key += ":" + data.Domain
		}

		err = mod.Insert(c.Context(), bcdb.DB(), boil.Infer())
		if err != nil {
			log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to insert metadata update to queue")
			return c.SendStatus(http.StatusInternalServerError)
		}
	} else if data.BADV != nil {
		b, err := json.Marshal(data.BADV)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		mod := models.MetadataQueue{
			Key:           "badv:" + data.Publisher,
			TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, "badv", time.Now()),
			Value:         b,
		}
		if data.Domain != "" {
			mod.Key += ":" + data.Domain
		}

		err = mod.Insert(c.Context(), bcdb.DB(), boil.Infer())
		if err != nil {
			log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to insert metadata update to queue")
			return c.SendStatus(http.StatusInternalServerError)
		}
	}

	//log.Info().Interface("update", data).Msg("metadata update parsed")

	return c.JSON(BlockUpdateRespose{
		Status: "ok",
	})
}

// BlockGetAllHandler Get publisher block list (bcat and badv) setup
// @Description Get publisher block list (bcat and badv) setup
// @Tags metadata
// @Accept json
// @Produce json
// @Param options body BlockGetRequest true "Block update Options"
// @Success 200 {object} BlockGetResponse
// @Security ApiKeyAuth
// @Router /block/get [post]
func BlockGetAllHandler(c *fiber.Ctx) error {

	request := &BlockGetRequest{}

	if err := c.BodyParser(&request); err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to parse metadata update payload")

		return c.SendStatus(http.StatusBadRequest)
	}

	key := createKeyForQuery(request)
	records := models.MetadataQueueSlice{}
	err := queries.Raw(query+key+sortQuery).Bind(c.Context(), bcdb.DB(), &records)

	if err != nil {
		log.Error().Err(err).Msg("failed to fetch all price factors")
		c.SendString("failed to fetch")
		return c.SendStatus(http.StatusInternalServerError)
	}

	if c.Query("format") == "html" {
		c.Set("Content-Type", "text/html")
		b := bytes.Buffer{}
		p := message.NewPrinter(message.MatchLanguage("en"))

		for _, rec := range records {
			b.WriteString(p.Sprintf(rowBlock, rec.Key, rec.Value, rec.CreatedAt.Format("2006-01-02 15:04"), rec.CommitedInstances))
		}
		t := fasttemplate.New(htmlBlock, "{{", "}}")
		s := t.ExecuteString(map[string]interface{}{
			"data": b.String(),
		})
		return c.SendString(s)
	} else {
		c.Set("Content-Type", "application/json")
		return c.JSON(records)
	}

}

func createKeyForQuery(request *BlockGetRequest) string {
	publisher := request.Publisher
	domain := request.Domain

	if len(publisher) != 0 && len(domain) != 0 {
		return " and metadata_queue.key = '" + publisher + ":" + domain + "'"
	}

	if len(publisher) != 0 {
		return " and last.key = '" + publisher + "'"
	}

	return ` and 1=1 `
}

var htmlBlock = `
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
      Current Publisher Blocks 
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
                  List
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

var rowBlock = `<tr>
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
