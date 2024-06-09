package pages

import (
	"strings"

	"github.com/ipoluianov/aneth_eth/an"
	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/pages_list/page_index"
	"github.com/ipoluianov/aneth_eth/pages_list/page_token_summary"
	"github.com/ipoluianov/aneth_eth/pages_list/page_tokens"
	"github.com/ipoluianov/aneth_eth/tokens"
)

type Pages struct {
	itemsMap map[string]*common.Page
	items    []*common.Page
}

var Instance *Pages

func init() {
	Instance = NewPages()
}

func NewPages() *Pages {
	var c Pages
	c.itemsMap = make(map[string]*common.Page)
	c.items = make([]*common.Page, 0)

	c.addPage(page_index.New())
	c.addPage(page_tokens.New())

	for _, token := range tokens.Instance.GetTokens() {
		c.addPage(page_token_summary.New(token))
	}

	return &c
}

func (c *Pages) addPage(page *common.Page) {
	c.itemsMap[page.Code] = page
	c.items = append(c.items, page)
}

func (c *Pages) GetPage(pageCode string) *common.PageRunResult {
	var result *common.PageRunResult
	if page, ok := c.itemsMap[pageCode]; ok {
		var res common.PageRunResult
		res.Page = page
		res.Code = pageCode
		res.Name = page.Name
		res.Description = page.Description
		res.Content = ""
		page.Fn(page, &res)
		res.Content = `<h1>` + res.Name + `</h1>` + "\r\n" + res.Content
		result = &res
	}
	return result
}

func (c *Pages) GetPages() []*common.Page {
	return c.items
}

func BuildMap() string {
	result := ""

	fAddItem := func(name string, url string) {
		tmp := `    <li><a href="%URL%">%NAME%</a></li>` + "\r\n"
		tmp = strings.ReplaceAll(tmp, "%URL%", url)
		tmp = strings.ReplaceAll(tmp, "%NAME%", name)
		result += tmp
	}

	fAddHeader := func(name string) {
		tmp := `    <h2>%NAME%</h2>` + "\r\n"
		tmp = strings.ReplaceAll(tmp, "%NAME%", name)
		result += tmp
	}

	fAddHeader("Main")
	fAddItem("INDEX", "/")
	fAddItem("SITE MAP", "/v/map")

	tasks := an.Instance.GetTasks()
	fAddHeader("Views")
	for _, task := range tasks {
		fAddItem(task.Name, "/v/"+task.Code)
	}

	fAddHeader("Pages")
	pgs := Instance.GetPages()
	for _, p := range pgs {
		fAddItem(p.Name, "/p/"+p.Code)
	}

	fAddHeader("JSON-REST")
	fAddItem("STATE", "/d/state")

	for _, task := range tasks {
		fAddItem(task.Code, "/d/"+task.Code)
	}

	return result
}
