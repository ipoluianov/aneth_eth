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

	fAddItem := func(name string, url string) {
		tmp := `    <li><a href="%URL%">%NAME%</a></li>` + "\r\n"
		tmp = strings.ReplaceAll(tmp, "%URL%", url)
		tmp = strings.ReplaceAll(tmp, "%NAME%", name)
		content += tmp
	}

	/*fAddText := func(text string) {
		tmp := `<div>%TEXT%</div>` + "\r\n"
		tmp = strings.ReplaceAll(tmp, "%TEXT%", text)
		content += tmp
	}*/

	fAddHeader2 := func(name string) {
		tmp := `    <h2>%NAME%</h2>` + "\r\n"
		tmp = strings.ReplaceAll(tmp, "%NAME%", name)
		content += tmp
	}

	for _, t := range tokens.Instance.GetTokens() {
		fAddHeader2(t.Name)
		if t.Symbol != "USDT" {
			fAddItem("Price", "/v/token-"+strings.ToLower(t.Symbol)+"-price")
		}

		fAddItem("Transfer Amount", "/v/token-"+strings.ToLower(t.Symbol)+"-transfer-amount")
		fAddItem("Number of transfers", "/v/token-"+strings.ToLower(t.Symbol)+"-number-of-transactions")
		fAddItem("Summary", "/p/token-"+strings.ToLower(t.Symbol)+"-summary")
	}

	result.Content = content
}
