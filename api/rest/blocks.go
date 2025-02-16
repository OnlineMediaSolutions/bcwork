package rest

import (
	"bytes"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
	"github.com/valyala/fasttemplate"
	"golang.org/x/text/message"
)

// BlockPostHandler Update bidder addomain and categories blocks
// @Description Update bidder addomain and categories blocks.
// @Tags MetaData
// @Accept json
// @Produce json
// @Param options body dto.BlockUpdateRequest true "Block update Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /block [post]
func (o *OMSNewPlatform) BlockPostHandler(c *fiber.Ctx) error {
	data := &dto.BlockUpdateRequest{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse blocks metadata update payload", err)
	}

	err := o.blocksService.UpdateBlocks(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to update blocks", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "blocks successfully updated")
}

// BlockGetAllHandler Get publisher block list (bcat and badv) setup
// @Description Get publisher block list (bcat and badv) setup
// @Tags MetaData
// @Accept json
// @Produce json
// @Param options body dto.BlockGetRequest true "Block update Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /block/get [post]
func (o *OMSNewPlatform) BlockGetAllHandler(c *fiber.Ctx) error {
	request := &dto.BlockGetRequest{}
	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for getting blocks", err)
	}

	blocks, err := o.blocksService.GetBlocks(c.Context(), request)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve blocks data", err)
	}

	if c.Query("format") == "html" {
		c.Set("Content-Type", "text/html")
		b := bytes.Buffer{}
		p := message.NewPrinter(message.MatchLanguage("en"))

		for _, rec := range blocks {
			b.WriteString(p.Sprintf(rowBlock, rec.Key, rec.Value, rec.CreatedAt.Format("2006-01-02 15:04"), rec.CommitedInstances))
		}
		t := fasttemplate.New(htmlBlock, "{{", "}}")
		s := t.ExecuteString(map[string]interface{}{
			"data": b.String(),
		})

		return c.SendString(s)
	} else {
		return c.JSON(blocks)
	}
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
