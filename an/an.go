package an

import (
	"sync"
	"time"

	"github.com/ipoluianov/aneth_eth/db"
	"github.com/ipoluianov/aneth_eth/tokens"
	"github.com/ipoluianov/gomisc/logger"
)

type An struct {
	analytics  map[string]*Result
	mtx        sync.Mutex
	tasks      []*Task
	cache      *Cache
	taskGroups []*TaskGroup
}

type AnState struct {
	Tasks []*TaskState
	Cache *CacheState
}

type TaskGroup struct {
	Code  string
	Name  string
	Tasks []string
}

var Instance *An

func init() {
	Instance = NewAn()
}

func NewAn() *An {
	var c An
	c.analytics = make(map[string]*Result)
	c.cache = NewCache()

	return &c
}

func (c *An) Start() {
	groupBase := &TaskGroup{}
	groupBase.Code = "task-group-base"
	groupBase.Name = "Base reports"
	groupBase.Tasks = append(groupBase.Tasks, "number-of-transactions-per-minute")
	groupBase.Tasks = append(groupBase.Tasks, "eth-transfer-volume-per-minute")
	groupBase.Tasks = append(groupBase.Tasks, "number-of-rejected-transactions-per-minute")
	groupBase.Tasks = append(groupBase.Tasks, "number-of-new-contracts-per-minute")
	groupBase.Tasks = append(groupBase.Tasks, "number-of-erc20-transfers-per-minute")
	c.taskGroups = append(c.taskGroups, groupBase)

	// TimeCharts
	c.tasks = append(c.tasks, NewTask("number-of-transactions-per-minute", "timechart", "Number of transactions per minute", c.taskMinutesCount, `
	Number of transactions by minute on the chart. Only successful transactions are displayed. This data indicates the overall activity of the network.`, `
	`))
	c.tasks = append(c.tasks, NewTask("eth-transfer-volume-per-minute", "timechart", "ETH transfer volume per minute", c.taskMinutesValues, `
	The graph shows the total volume of ETH transfers. These can be either regular transfers between accounts or transfers to smart merchant addresses.`, `
	`))
	c.tasks = append(c.tasks, NewTask("number-of-rejected-transactions-per-minute", "timechart", "Number of rejected transactions per minute", c.taskMinutesRejected, `
	Displays the number of unsuccessful transactions recently. An increase in the number of such transactions indicates possible unsuccessful attacks on the network.`, `
	The graph shows the number of transactions that, after being included in a block, were rejected as a result of executing a smart contract, per minute. Possible reasons for rejection include incorrect call parameters, insufficient funds to complete the operation, errors in the smart contract logic, and failure to meet contract conditions.`))
	c.tasks = append(c.tasks, NewTask("number-of-new-contracts-per-minute", "timechart", "Number of new contracts per minute", c.taskMinutesNewContracts, `
	The graph displays the number of transactions that create new smart contracts on the network.`, `
	`))
	c.tasks = append(c.tasks, NewTask("number-of-erc20-transfers-per-minute", "timechart", "Number of ERC20 transfers by minute", c.taskMinutesERC20Transfers, `
	The graph shows the number of ERC-20 transfers`, `
	`))

	// Tables
	c.tasks = append(c.tasks, NewTask("accounts-by-send-count", "table", "Top ETH FROM", c.taskAccountsBySendCount, `
	Top 10 addresses participating in transactions as a sender`, `
	`))
	c.tasks = append(c.tasks, NewTask("accounts-by-recv-count", "table", "Top ETH TO", c.taskAccountsByRcvCount, `
	Top 10 addresses participating in transactions as a receiver`, `
	`))
	c.tasks = append(c.tasks, NewTask("new-eth-contracts-list", "table", "New ETH Contracts - Last 24 hours", c.taskNewContracts, `
	List of new smart contracts`, `
	`))

	// Tokens
	for _, token := range tokens.Instance.GetTokens() {
		groupToken := &TaskGroup{}
		groupToken.Code = "task-group-token-" + token.Symbol
		groupToken.Name = "Token " + token.Symbol

		// Volume
		codeVolume := token.Symbol + "-token-transfers-volume-per-minute"
		c.tasks = append(c.tasks, NewTask(codeVolume, "timechart", "Volume of "+token.Symbol+" token transfers by minute", c.taskMinutesTokenTransfers, `
		Displaying volume of `+token.Symbol+` token transfers per minute on the network`, `
		`))
		groupToken.Tasks = append(groupToken.Tasks, codeVolume)

		// Number
		codeNumber := token.Symbol + "-token-transfers-number-per-minute"
		c.tasks = append(c.tasks, NewTask(codeNumber, "timechart", "Number of "+token.Symbol+" token transfers by minute", c.taskMinutesTokenTransfersNumber, `
		Displaying number of `+token.Symbol+` token transfers per minute on the network`, `
		`))
		groupToken.Tasks = append(groupToken.Tasks, codeNumber)

		c.taskGroups = append(c.taskGroups, groupToken)
	}

	c.cache.Start()
	go c.ThAn()
}

func (c *An) GetState() *AnState {
	var st AnState
	c.mtx.Lock()
	for _, t := range c.tasks {
		state := t.State
		st.Tasks = append(st.Tasks, &state)
	}
	st.Cache = c.cache.GetState()
	c.mtx.Unlock()
	return &st
}

func (c *An) GetTask(code string) *Task {
	var task *Task
	c.mtx.Lock()
	for _, t := range c.tasks {
		if t.Code == code {
			task = t
		}
	}
	c.mtx.Unlock()
	return task
}

func (c *An) GetTasks() []Task {
	result := make([]Task, 0)
	c.mtx.Lock()
	for _, task := range c.tasks {
		var t Task
		t.Code = task.Code
		t.Name = task.Name
		t.Description = task.Description
		t.Type = task.Type
		result = append(result, t)
	}
	c.mtx.Unlock()
	return result
}

func (c *An) GetTaskGroups() []*TaskGroup {
	return c.taskGroups
}

func (c *An) GetResult(code string) *Result {
	var res *Result
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if r, ok := c.analytics[code]; ok {
		res = r
	}
	return res
}

func (c *An) ThAn() {
	for {
		dt1 := time.Now()
		logger.Println("")
		logger.Println("---------- an ------------")
		logger.Println("reading transactions")
		txsByMinutes, txs := db.Instance.GetLatestTransactions()
		count := 0
		for _, item := range txsByMinutes.Items {
			count += len(item.TXS)
		}
		logger.Println("An::an txs:", len(txs))

		for _, task := range c.tasks {
			var res Result

			res.Code = task.Code
			res.Type = task.Type
			res.Parameters = make([]*ResultParameter, 0)
			res.Table.Items = make([]*ResultTableItem, 0)
			res.Table.Columns = make([]*ResultTableColumn, 0)
			res.TimeChart.Items = make([]*ResultTimeChartItem, 0)

			dtBegin := time.Now()
			task.Fn(&res, txsByMinutes, txs)
			dtEnd := time.Now()
			duration := dtEnd.Sub(dtBegin).Milliseconds()
			task.State.Code = task.Code
			task.State.LastExecTime = time.Now().Format("2006-01-02 15:04:05")
			task.State.LastExecTimeDurationMs = int(duration)
			c.mtx.Lock()
			c.analytics[task.Code] = &res
			c.mtx.Unlock()
		}

		dt2 := time.Now()
		logger.Println("execution time:", dt2.Sub(dt1).Milliseconds(), "ms")
		logger.Println("--------------------------")
		logger.Println("")

		time.Sleep(3 * time.Second)
	}
}
