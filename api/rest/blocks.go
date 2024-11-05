package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasttemplate"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"golang.org/x/text/message"
)

// DemandReportGetRequest contains filter parameters for retrieving events
type BlockUpdateRequest struct {
	Publisher string   `json:"publisher"`
	Domain    string   `json:"domain"`
	BCAT      []string `json:"bcat"`
	BADV      []string `json:"badv"`
}

type BlockGetRequest struct {
	Types     []string `json:"types"`
	Publisher string   `json:"publisher"`
	Domain    string   `json:"domain"`
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
// @Tags MetaData
// @Accept json
// @Produce json
// @Param options body BlockUpdateRequest true "Block update Options"
// @Success 200 {object} BlockUpdateRespose
// @Security ApiKeyAuth
// @Router /block [post]
func (o *OMSNewPlatform) BlockPostHandler(c *fiber.Ctx) error {
	data := &BlockUpdateRequest{}

	if err := c.BodyParser(&data); err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to parse metadata update payload")
		return c.Status(http.StatusBadRequest).JSON(Response{Status: "error", Message: "Failed to parse metadata update payload"})
	}

	if err := validateData(data); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Response{Status: "error", Message: err.Error()})
	}

	if data.BCAT != nil {
		if err := updateDB(c.Context(), "bcat", data.Publisher, data.Domain, data.BCAT); err != nil {
			return handleError(err, c)
		}
	}

	if data.BADV != nil {
		if err := updateDB(c.Context(), "badv", data.Publisher, data.Domain, data.BADV); err != nil {
			return handleError(err, c)
		}
	}

	return c.JSON(BlockUpdateRespose{
		Status: "ok",
	})
}

// BlockGetAllHandler Get publisher block list (bcat and badv) setup
// @Description Get publisher block list (bcat and badv) setup
// @Tags MetaData
// @Accept json
// @Produce json
// @Param options body BlockGetRequest true "Block update Options"
// @Success 200 {object} BlockUpdateRespose
// @Security ApiKeyAuth
// @Router /block/get [post]
func (o *OMSNewPlatform) BlockGetAllHandler(c *fiber.Ctx) error {

	request := &BlockGetRequest{}

	errMessage := validateRequest(c, request)
	if len(errMessage) != 0 {
		return c.Status(http.StatusInternalServerError).JSON(Response{Status: "error", Message: errMessage})
	}

	key := createKeyForQuery(request)
	records := models.MetadataQueueSlice{}
	err := queries.Raw(query+key+sortQuery).Bind(c.Context(), bcdb.DB(), &records)

	if err != nil {
		log.Error().Err(err).Msg("failed to fetch all price factors")
		return c.Status(http.StatusInternalServerError).JSON(Response{Status: "error", Message: "Failed to fetch all price factors"})
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

func updateDB(ctx context.Context, businessType, publisher, domain string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	mod := models.MetadataQueue{
		Key:           fmt.Sprintf("%s:%s", businessType, publisher),
		TransactionID: bcguid.NewFromf(publisher, domain, businessType, time.Now()),
		Value:         b,
	}

	if domain != "" {
		mod.Key += ":" + domain
	}

	if err := mod.Insert(ctx, bcdb.DB(), boil.Infer()); err != nil {
		return err
	}
	return nil
}

func handleError(err error, c *fiber.Ctx) error {
	log.Error().Err(err).Str("body", string(c.Body())).Msg("Failed to insert metadata update to queue")
	return c.Status(http.StatusInternalServerError).JSON(Response{Status: "error", Message: "Failed to insert metadata update to queue"})
}

func validateRequest(c *fiber.Ctx, request *BlockGetRequest) string {

	if err := c.BodyParser(&request); err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to parse metadata update payload")
		c.SendStatus(http.StatusInternalServerError)
		return "Failed to parse metadata update payload"
	}

	return ""
}

func validateData(data *BlockUpdateRequest) error {
	if data.Publisher == "" {
		return errors.New("Publisher is mandatory")
	}
	return nil
}

func createKeyForQuery(request *BlockGetRequest) string {
	types := request.Types
	publisher := request.Publisher
	domain := request.Domain

	var query bytes.Buffer

	//If no publisher or no business types or empty body than return all
	if len(publisher) == 0 || len(types) == 0 {
		query.WriteString(` and 1=1 `)
		return query.String()
	}

	for index, btype := range types {
		if index == 0 {
			query.WriteString("AND (")
		}
		if len(publisher) != 0 && len(domain) != 0 {
			query.WriteString(" (metadata_queue.key = '" + btype + ":" + publisher + ":" + domain + "')")

		} else if len(publisher) != 0 {
			query.WriteString(" (metadata_queue.key = '" + btype + ":" + publisher + "')")
		}
		if index < len(types)-1 {
			query.WriteString(" OR")
		}
	}
	query.WriteString(")")
	return query.String()
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
