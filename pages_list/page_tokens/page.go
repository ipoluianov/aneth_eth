package page_tokens

import (
	"strings"

	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/tokens"
)

func New() *common.Page {
	var c common.Page
	c.Code = "tokens"
	c.Name = "Tokens"
	c.Description = "Tokens - ETH"
	c.Fn = Run
	return &c
}

func Run(page *common.Page, result *common.PageRunResult) {
	content := ""

	fAddItem := func(name string, url string, desc string, viewCode string) {
		tmp := `
		<a href="%URL%" class="menu-block">
			<h2>%NAME%</h2>
			<img src="/images/%VIEW_CODE%"/>
		</a>
		`

		if url != "" {
			tmp = strings.ReplaceAll(tmp, "%URL%", url)
		} else {
			tmp = strings.ReplaceAll(tmp, "%URL%", "#")
		}
		tmp = strings.ReplaceAll(tmp, "%NAME%", name)

		if viewCode == "" {
			viewCode = "none"
		}

		tmp = strings.ReplaceAll(tmp, "%VIEW_CODE%", viewCode)
		content += tmp
	}

	fAddHeader2 := func(name string) {
		tmp := `    <h2>%NAME%</h2>` + "\r\n"
		tmp = strings.ReplaceAll(tmp, "%NAME%", name)
		content += tmp
	}

	for _, t := range tokens.Instance.GetTokens() {
		fAddHeader2(t.Name)

		content += `<div class="menu-container">`
		if t.Symbol != "USDT" {
			fAddItem("Price", "/v/token-"+strings.ToLower(t.Symbol)+"-price", "", "token-"+strings.ToLower(t.Symbol)+"-price")
		} else {
			fAddItem("Price", "", "", "")
		}
		fAddItem("Transfer Amount", "/v/token-"+strings.ToLower(t.Symbol)+"-transfer-amount", "", "token-"+strings.ToLower(t.Symbol)+"-transfer-amount")
		fAddItem("Number of transfers", "/v/token-"+strings.ToLower(t.Symbol)+"-number-of-transactions", "", "token-"+strings.ToLower(t.Symbol)+"-number-of-transactions")
		fAddItem("Summary", "/p/token-"+strings.ToLower(t.Symbol)+"-summary", "", "")
		content += `</div>`
	}

	result.Content = content
}
