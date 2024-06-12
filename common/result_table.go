package common

type ResultTableItem struct {
	Values []string
}

type ResultTableColumn struct {
	Name  string
	Align string
}

type ResultTable struct {
	Name    string
	Columns []*ResultTableColumn
	Items   []*ResultTableItem
}
