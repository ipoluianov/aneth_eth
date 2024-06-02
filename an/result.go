package an

type ResultParameter struct {
	Text  string
	Value string
}

type Result struct {
	Code            string
	Type            string
	Count           int
	CurrentDateTime string

	Parameters []*ResultParameter
	TimeChart  ResultTimeChart
	Table      ResultTable
}
