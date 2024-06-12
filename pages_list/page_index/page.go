package page_index

import (
	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/views"
)

func New() *common.Page {
	var c common.Page
	c.Code = "index"
	c.Name = common.GlobalSiteName
	c.Description = common.GlobalSiteDescription
	c.Fn = Run
	return &c
}

func Run(page *common.Page, result *common.PageRunResult) {
	content := ""
	p1, _, _ := views.GetView("eth-price", "", "", "instance1", 200, false, false, false, true)
	content += p1
	p2, _, _ := views.GetView("summary", "", "", "instance2", 200, false, false, false, true)
	content += p2
	/*p2, _, _ := views.GetView("btc-price", "", "", "instance2", 200, false, false, false, true)
	content += p2*/

	content += `<h1>Content</h1>`
	content += `<div><a href="/p/tokens">Tokens</a></div>`
	result.Content = content
}
