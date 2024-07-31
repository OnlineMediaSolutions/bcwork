package rest

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
)

// ConfiantGetHandler Get confiant setup
// @Description Get confiant setup
// @Tags Confiant
// @Accept json
// @Produce json
// @Param options body core.GetConfiantOptions true "options"
// @Success 200 {object} core.ConfiantSlice
// @Router /confiant/get [post]
func ConfiantGetAllHandler(c *fiber.Ctx) error {

	data := &core.GetConfiantOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body parsing error")
	}

	pubs, err := core.GetConfiants(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve confiants")
	}

	return c.JSON(pubs)
}

// ConfiantPostHandler Update and enable Confiant setup
// @Description Update and enable Confiant setup (publisher is mandatory, domain is optional)
// @Tags Confiant
// @Accept json
// @Produce json
// @Param options body core.ConfiantUpdateRequest true "Confiant update Options"
// @Success 200 {object} core.ConfiantUpdateRespose
// @Security ApiKeyAuth
// @Router /confiant [post]
func ConfiantPostHandler(c *fiber.Ctx) error {

	data := &core.ConfiantUpdateRequest{}
	err := c.BodyParser(&data)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Confiant payload parsing error")
	}

	err = core.UpdateConfiant(c, data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Confiant table")
	}

	err = core.UpdateMetaDataQueue(c, data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("Failed to update metadata table for confiant, %s", err))
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Confiant and Metadata tables successfully updated")
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
