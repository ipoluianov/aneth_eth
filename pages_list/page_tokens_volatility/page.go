package page_tokens_volatility

import (
	"strings"

	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/tokens"
	"github.com/ipoluianov/aneth_eth/views"
)

func New() *common.Page {
	var c common.Page
	c.Code = "tokens-volatility"
	c.Name = "Tokens volatility"
	c.Description = "Tokens volatility"
	c.Fn = Run
	c.Symbol = ""
	c.Ticket = ""
	return &c
}

func Run(page *common.Page, result *common.PageRunResult) {
	content := ""
	for _, token := range tokens.Instance.GetTokens() {
		//showTimeScale := index == len(tokens.Instance.GetTokens())-1
		p1, _, _ := views.GetView("token-"+strings.ToLower(token.Symbol)+"-volatility", "", "", token.Symbol, 200, false, false, false, false)
		content += p1
	}
	{
		p1, _, _ := views.GetView("btc-volatility", "", "", "BTC", 200, false, false, false, false)
		content += p1
	}
	{
		p1, _, _ := views.GetView("eth-volatility", "", "", "ETH", 200, false, false, false, false)
		content += p1
	}
	result.Content = content
}
