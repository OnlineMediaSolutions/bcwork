package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type TableData struct {
	Publisher string
	Domain    string
}

type Request struct {
	Title   string      `json:"title"`
	Data    []TableData `json:"data"`
	Columns []string    `json:"columns"`
}

func GenerateHtml(c *fiber.Ctx) error {
	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).SendString("Failed to parse request")
	}

	firstTableHTML := fmt.Sprintf("<h3>%s</h3><table>", req.Title)
	firstTableHTML += `<tr class="table-head"><td class="border-right">Publisher</td><td class="border-right">Domain</td></tr>`

	for i, row := range req.Data {
		className := "even"
		if i%2 == 0 {
			className = "odd"
		}
		firstTableHTML += fmt.Sprintf(
			`<tr class="table-row %s"><td class="border-right">%s</td></tr>`,
			className, row.Publisher, row.Domain)
	}

	firstTableHTML += "</table>"

	err := SendEmail(EmailRequest{
		To:      "sonai@onlinemediasolutions.com",
		Bcc:     "sonai@onlinemediasolutions.com",
		Subject: "Hello",
		Body:    firstTableHTML,
		IsHTML:  true,
	})

	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"html": firstTableHTML,
	})
}
