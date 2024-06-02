package an

type ResultItemByMinutes struct {
	Index int
	DT    uint64
	DTStr string
	Value float64
}

type ResultItemTable struct {
	Text   string
	Values []float64
}

type Result struct {
	Code            string
	Count           int
	CurrentDateTime string
	ItemsByMinutes  []*ResultItemByMinutes
	ItemsTable      []*ResultItemTable
}
