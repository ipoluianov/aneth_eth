package page_tokens_price_minus_btc_price

import (
	"strings"

	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/tokens"
	"github.com/ipoluianov/aneth_eth/views"
)

func New() *common.Page {
	var c common.Page
	c.Code = "tokens-prices-minus-btc-price"
	c.Name = "Token prices minus BTC price"
	c.Description = "Token prices minus BTC price"
	c.Fn = Run
	c.Symbol = ""
	c.Ticket = ""
	return &c
}

func Run(page *common.Page, result *common.PageRunResult) {
	content := ""
	for index, token := range tokens.Instance.GetTokens() {
		showTimeScale := index == len(tokens.Instance.GetTokens())-1
		p1, _, _ := views.GetView("token-"+strings.ToLower(token.Symbol)+"-price-minus-btc-price", "", "", token.Symbol, 200, false, false, false, showTimeScale)
		content += p1
	}
	result.Content = content
}
