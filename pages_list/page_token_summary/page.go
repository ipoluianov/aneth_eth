package page_token_summary

import (
	"strings"

	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/tokens"
	"github.com/ipoluianov/aneth_eth/views"
)

func New(token *tokens.Token) *common.Page {
	var c common.Page
	c.Code = "token-" + strings.ToLower(token.Symbol) + "-summary"
	c.Name = "Token " + token.Name + " - Summary"
	c.Description = "Token " + token.Name + " charts: price, number of transaction, amount"
	c.Fn = Run
	c.Symbol = token.Symbol
	c.Ticket = token.Ticket
	return &c
}

func Run(page *common.Page, result *common.PageRunResult) {
	content := ""
	p1, _, _ := views.GetView("token-"+strings.ToLower(page.Symbol)+"-transfer-amount", "", "", "instance1", 200, false, false, false, false)
	content += p1
	p2, _, _ := views.GetView("token-"+strings.ToLower(page.Symbol)+"-number-of-transactions", "", "", "instance2", 200, false, false, false, false)
	content += p2
	p3, _, _ := views.GetView("token-"+strings.ToLower(page.Symbol)+"-price", "", "", "instance3", 200, false, false, false, false)
	content += p3
	p4, _, _ := views.GetView("token-"+strings.ToLower(page.Symbol)+"-volatility", "", "", "instance4", 200, false, false, false, false)
	content += p4
	p5, _, _ := views.GetView("btc-price", "", "", "instance5", 200, false, false, false, true)
	content += p5
	result.Content = content
}
