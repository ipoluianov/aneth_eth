package an

import (
	"sync"
	"time"

	"github.com/ipoluianov/aneth_eth/db"
	"github.com/ipoluianov/gomisc/logger"
)

type An struct {
	analytics map[string]*Result
	mtx       sync.Mutex
	tasks     []*Task
	cache     *Cache
}

type AnState struct {
	Tasks []*TaskState
	Cache *CacheState
}

var Instance *An

func init() {
	Instance = NewAn()
}

func NewAn() *An {
	var c An
	c.analytics = make(map[string]*Result)
	c.cache = NewCache()

	// TimeCharts
	c.tasks = append(c.tasks, NewTask("minutes_count", "timechart", "Number of transactions by minute", c.taskMinutesCount, `
	Number of transactions by minute on the chart. Only successful transactions are displayed. This data indicates the overall activity of the network.`, `
	`))
	c.tasks = append(c.tasks, NewTask("minutes_values", "timechart", "Values by minute", c.taskMinutesValues, `
	The graph shows the total volume of ETH transfers. These can be either regular transfers between accounts or transfers to smart merchant addresses.`, `
	`))
	c.tasks = append(c.tasks, NewTask("minutes_count_of_usdt", "timechart", "Number of USDT transfers per minute", c.taskMinutesCountOfUsdt, `
	A graph of the number of USDT transfers on the Ethereum network is displayed. Only successful transactions are displayed.`, `
	`))
	c.tasks = append(c.tasks, NewTask("minutes_rejected", "timechart", "Number of rejected transactions by minute", c.taskMinutesRejected, `
	Displays the number of unsuccessful transactions recently. An increase in the number of such transactions indicates possible unsuccessful attacks on the network.`, `
	The graph shows the number of transactions that, after being included in a block, were rejected as a result of executing a smart contract, per minute. Possible reasons for rejection include incorrect call parameters, insufficient funds to complete the operation, errors in the smart contract logic, and failure to meet contract conditions.`))
	c.tasks = append(c.tasks, NewTask("minutes_new_contracts", "timechart", "Number of new contracts by minute", c.taskMinutesNewContracts, `
	The graph displays the number of transactions that create new smart contracts on the network.`, `
	`))
	c.tasks = append(c.tasks, NewTask("minutes_erc20_transfers", "timechart", "ERC20 transfers by minute", c.taskMinutesERC20Transfers, `
	The graph shows the number of transfers using the USDT token smart contract. USDT is a stablecoin whose price is maintained by the issuing company.`, `
	`))
	c.tasks = append(c.tasks, NewTask("minutes_pepe_transfers", "timechart", "PEPE transfers by minute", c.taskMinutesPepeTransfers, `
	Displaying the volume of PEPE token transfers on the network`, `
	`))

	// Tables
	c.tasks = append(c.tasks, NewTask("accounts_by_send_count", "table", "Top ETH FROM", c.taskAccountsBySendCount, `
	Top 10 addresses participating in transactions as a sender`, `
	`))
	c.tasks = append(c.tasks, NewTask("accounts_by_recv_count", "table", "Top ETH TO", c.taskAccountsByRcvCount, `
	Top 10 addresses participating in transactions as a receiver`, `
	`))
	c.tasks = append(c.tasks, NewTask("new_contracts", "table", "New ETH Contracts - Last 24 hours", c.taskNewContracts, `
	List of new smart contracts`, `
	`))

	return &c
}

func (c *An) Start() {
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
