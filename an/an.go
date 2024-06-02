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
	c.tasks = append(c.tasks, NewTask("minutes_count", "timechart", "Number of transactions by minute", c.taskMinutesCount, "desc"))
	c.tasks = append(c.tasks, NewTask("minutes_values", "timechart", "minutes_values", c.taskMinutesValues, "desc"))
	c.tasks = append(c.tasks, NewTask("minutes_count_of_usdt", "timechart", "minutes_count_of_usdt", c.taskMinutesCountOfUsdt, "desc"))
	c.tasks = append(c.tasks, NewTask("minutes_rejected", "timechart", "minutes_rejected", c.taskMinutesRejected, "desc"))
	c.tasks = append(c.tasks, NewTask("minutes_new_contracts", "timechart", "minutes_new_contracts", c.taskMinutesNewContracts, "desc"))
	c.tasks = append(c.tasks, NewTask("minutes_erc20_transfers", "timechart", "minutes_erc20_transfers", c.taskMinutesERC20Transfers, "desc"))
	c.tasks = append(c.tasks, NewTask("minutes_pepe_transfers", "timechart", "minutes_pepe_transfers", c.taskMinutesPepeTransfers, "desc"))

	// Tables
	c.tasks = append(c.tasks, NewTask("accounts_by_send_count", "table", "Top FROM", c.taskAccountsBySendCount, "desc"))
	c.tasks = append(c.tasks, NewTask("accounts_by_recv_count", "table", "Top TO", c.taskAccountsByRcvCount, "desc"))
	c.tasks = append(c.tasks, NewTask("new_contracts", "table", "New Contracts", c.taskNewContracts, "desc"))

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

func (c *An) GetResultsCodes() []string {
	result := make([]string, 0)
	c.mtx.Lock()
	for _, task := range c.tasks {
		result = append(result, task.Code)
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
