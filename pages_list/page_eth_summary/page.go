package page_eth_summary

import (
	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/views"
)

func New() *common.Page {
	var c common.Page
	c.Code = "eth-summary"
	c.Name = "ETH - Summary"
	c.Description = "ETH charts: price, number of transaction, amount"
	c.Fn = Run
	c.Symbol = "ETH"
	c.Ticket = "ETHUSDT"
	return &c
}

func Run(page *common.Page, result *common.PageRunResult) {
	content := ""
	p1, _, _ := views.GetView("eth-transfer-volume-per-minute", "", "", "instance1", 200, false, false, false, false)
	content += p1
	p2, _, _ := views.GetView("number-of-transactions-per-minute", "", "", "instance2", 200, false, false, false, false)
	content += p2
	p3, _, _ := views.GetView("eth-price", "", "", "instance3", 200, false, false, false, false)
	content += p3
	p4, _, _ := views.GetView("eth-volatility", "", "", "instance4", 200, false, false, false, false)
	content += p4
	p5, _, _ := views.GetView("btc-price", "", "", "instance5", 200, false, false, false, true)
	content += p5
	result.Content = content
}
