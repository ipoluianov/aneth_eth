package common

type Page struct {
	Code        string
	Name        string
	Description string
	Symbol      string
	Ticket      string
	Fn          func(page *Page, result *PageRunResult)
}

type PageRunResult struct {
	Page        *Page
	Code        string
	Name        string
	Description string
	Content     string
}
