package db

type DbStateBlockRange struct {
	DtStr1  string
	Number1 uint64
	DtStr2  string
	Number2 uint64
	Count   int
}

type DbState struct {
	Status                string
	SubStatus             string
	CountOfBlocks         int
	ReceiptsReceivedCount int
	ReceiptsReceivedError int
	ReceiptsMismatchError int
	LoadedBlocks          []DbStateBlockRange
	LoadedBlocksTimeRange string
}
