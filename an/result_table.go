package an

type ResultTableItem struct {
	Values []string
}

type ResultTableColumn struct {
	Name string
}

type ResultTable struct {
	Columns []*ResultTableColumn
	Items   []*ResultTableItem
}
