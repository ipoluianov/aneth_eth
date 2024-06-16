package an

import (
	"image/color"
	"strings"
	"sync"
	"time"

	"github.com/ipoluianov/aneth_eth/cache"
	"github.com/ipoluianov/aneth_eth/common"
	"github.com/ipoluianov/aneth_eth/db"
	"github.com/ipoluianov/aneth_eth/images"
	"github.com/ipoluianov/aneth_eth/tasks/task_table_new_contracts"
	"github.com/ipoluianov/aneth_eth/tasks/task_table_summary"
	"github.com/ipoluianov/aneth_eth/tasks/task_timechart_count"
	"github.com/ipoluianov/aneth_eth/tasks/task_timechart_erc20_transfers"
	"github.com/ipoluianov/aneth_eth/tasks/task_timechart_new_contracts"
	"github.com/ipoluianov/aneth_eth/tasks/task_timechart_price"
	"github.com/ipoluianov/aneth_eth/tasks/task_timechart_price_minus_btc"
	"github.com/ipoluianov/aneth_eth/tasks/task_timechart_rejected"
	"github.com/ipoluianov/aneth_eth/tasks/task_timechart_token_transfers_number"
	"github.com/ipoluianov/aneth_eth/tasks/task_timechart_token_transfers_values"
	"github.com/ipoluianov/aneth_eth/tasks/task_timechart_values"
	"github.com/ipoluianov/aneth_eth/tasks/task_timechart_volatility"
	"github.com/ipoluianov/aneth_eth/utils"

	"github.com/ipoluianov/aneth_eth/tokens"
	"github.com/ipoluianov/gomisc/logger"
)

type An struct {
	analytics  map[string]*common.Result
	mtx        sync.Mutex
	tasks      []*common.Task
	taskGroups []*TaskGroup
}

type AnState struct {
	Tasks []*common.TaskState
	Cache *cache.CacheState
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
	c.analytics = make(map[string]*common.Result)

	return &c
}

func (c *An) Start() {
	c.AddTask(task_timechart_count.New())
	c.AddTask(task_timechart_erc20_transfers.New())
	c.AddTask(task_timechart_new_contracts.New())
	c.AddTask(task_timechart_rejected.New())
	c.AddTask(task_timechart_values.New())

	c.AddTask(task_table_new_contracts.New())
	c.AddTask(task_table_summary.New())
	c.AddTask(task_timechart_price.New("BTC", "Bitcoin", "BTCUSDT"))
	c.AddTask(task_timechart_price.New("ETH", "ETH", "ETHUSDT"))
	c.AddTask(task_timechart_volatility.New("BTC", "Bitcoin", "BTCUSDT"))
	c.AddTask(task_timechart_volatility.New("ETH", "ETH", "ETHUSDT"))

	for _, token := range tokens.Instance.GetTokens() {
		c.AddTask(task_timechart_token_transfers_values.New(token.Symbol, token.Name))
		c.AddTask(task_timechart_token_transfers_number.New(token.Symbol, token.Name))
		if token.Symbol != "USDT" && token.Symbol != "LEO" {
			c.AddTask(task_timechart_price.New(token.Symbol, token.Name, token.Ticket))
			c.AddTask(task_timechart_volatility.New(token.Symbol, token.Name, token.Ticket))
			c.AddTask(task_timechart_price_minus_btc.New(token.Symbol, token.Name, token.Ticket))
		}
	}

	go c.ThAn()
}

func (c *An) AddTask(task *common.Task) {
	task.ViewProvider = c
	c.tasks = append(c.tasks, task)
}

func (c *An) GetState() *AnState {
	var st AnState
	c.mtx.Lock()
	for _, t := range c.tasks {
		state := t.State
		st.Tasks = append(st.Tasks, &state)
	}
	st.Cache = cache.Instance.GetState()
	c.mtx.Unlock()
	return &st
}

func (c *An) GetTask(code string) *common.Task {
	var task *common.Task
	c.mtx.Lock()
	for _, t := range c.tasks {
		if t.Code == code {
			task = t
		}
	}
	c.mtx.Unlock()
	return task
}

func (c *An) GetTasks() []common.Task {
	result := make([]common.Task, 0)
	c.mtx.Lock()
	for _, task := range c.tasks {
		var t common.Task
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

func (c *An) GetResult(code string) *common.Result {
	var res *common.Result
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

		{
			emptyData := make([]float64, 1)
			emptyData[0] = 0
			pngBS, err := utils.CreateSparkline(emptyData, 100, 50, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0})
			if err != nil {
				logger.Println("An::ThAn CreateSparkline error:", err)
			} else {
				images.Instance.Set("none", pngBS, "none")
			}
		}

		for _, task := range c.tasks {
			var res common.Result

			res.Code = task.Code
			res.Type = task.Type
			res.Parameters = make([]*common.ResultParameter, 0)
			res.Table.Items = make([]*common.ResultTableItem, 0)
			res.Table.Columns = make([]*common.ResultTableColumn, 0)
			res.TimeChart.Items = make([]*common.ResultTimeChartItem, 0)

			dtBegin := time.Now()
			task.Fn(task, &res, txsByMinutes, txs)
			res.Count = len(res.TimeChart.Items)
			res.CurrentDateTime = time.Now().UTC().Format("2006-01-02 15:04:05")

			{
				values := make([]float64, 0)
				for _, item := range res.TimeChart.Items {
					values = append(values, item.Value)
				}
				pngBS, err := utils.CreateSparkline(values, 300, 100, color.RGBA{0, 0, 0, 0}, color.RGBA{200, 200, 200, 255})
				if err == nil {
					images.Instance.Set(task.Code, pngBS, task.Name)
				}
			}

			dtEnd := time.Now()
			duration := dtEnd.Sub(dtBegin).Milliseconds()
			task.State.Code = task.Code
			task.State.LastExecTime = time.Now().Format("2006-01-02 15:04:05")
			task.State.LastExecTimeDurationMs = int(duration)

			if strings.Contains(task.State.Code, "price-minus-btc-price") {
				logger.Println("!!!!!!!!!!!", task.State.Code, res.Count)
			}

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
