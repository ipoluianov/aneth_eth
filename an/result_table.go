package an

type ResultTableItem struct {
	Values []string
}

type ResultTableColumn struct {
	Name string
}

type ResultTable struct {
	Name    string
	Columns []*ResultTableColumn
	Items   []*ResultTableItem
}
