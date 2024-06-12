package market

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

var Instance *Market

func init() {
	Instance = NewMarket()
}

type Market struct {
	mtx     sync.Mutex
	items   map[string]*TickerData
	tickers []*TickerInfo
}

func NewMarket() *Market {
	var c Market
	c.items = make(map[string]*TickerData)
	c.tickers = make([]*TickerInfo, 0)
	return &c
}

type TickerInfo struct {
	Symbol     string
	TickerCode string
}

type TickerData struct {
	Ticker string
	Data   []*Candle
}

func (c *TickerData) Summary() string {
	if len(c.Data) == 0 {
		return "empty"
	}
	str1 := c.Data[0].StartTime.UTC().Format("2006-01-02 15:04:05")
	str2 := c.Data[len(c.Data)-1].StartTime.UTC().Format("2006-01-02 15:04:05")
	return fmt.Sprint(str1 + " - " + str2 + " count: " + fmt.Sprint(len(c.Data)))
}

func (c *Market) addTicker(symbol string, tickerCode string) {
	var tickerInfo TickerInfo
	tickerInfo.Symbol = symbol
	tickerInfo.TickerCode = tickerCode
	c.tickers = append(c.tickers, &tickerInfo)
	var data TickerData
	data.Ticker = tickerCode
	data.Data = make([]*Candle, 0)
	c.items[tickerCode] = &data
}

func (c *Market) Start() {
	c.addTicker("BTC", "BTCUSDT")
	c.addTicker("ETH", "ETHUSDT")
	c.addTicker("BNB", "BNBUSDT")
	c.addTicker("TON", "TONUSDT")
	c.addTicker("SHIB", "SHIBUSDT")
	c.addTicker("LINK", "LINKUSDT")
	c.addTicker("UNI", "UNIUSDT")
	c.addTicker("PEPE", "PEPEUSDT")
	c.addTicker("NEAR", "NEARUSDT")
	c.addTicker("FET", "FETUSDT")
	go c.thUpdate()
}

func (c *Market) Stop() {
}

func (c *Market) GetData(symbol string) []*Candle {
	result := make([]*Candle, 0)
	endTime := time.Now().Unix()
	beginTime := endTime - 86400

	c.mtx.Lock()
	d := c.items[symbol]
	if d != nil {
		for _, item := range d.Data {
			if item.StartTime.Unix() >= beginTime && item.StartTime.Unix() < endTime {
				result = append(result, item)
			}
		}
	}
	c.mtx.Unlock()
	return result
}

func (c *Market) GetTickers() []*TickerInfo {
	return c.tickers
}

func (c *Market) thUpdate() {
	for {
		c.updateData()
		time.Sleep(1 * time.Second)
	}
}

func (c *Market) updateData() {
	fmt.Println("------------LOAD DATA-----------")
	for i := 0; i < len(c.tickers); i++ {
		c.loadTickerData(c.tickers[i].TickerCode)
	}
	time.Sleep(120 * time.Second)
}

func (c *Market) loadTickerData(tickerCode string) {
	now := time.Now()
	dt1 := now.Add(-24 * time.Hour)
	dt2 := now.Add(0 * time.Hour)
	data1 := c.LoadData(tickerCode, dt1)
	data2 := c.LoadData(tickerCode, dt2)
	data := make([]*Candle, 0)
	data = append(data, data1...)
	data = append(data, data2...)

	sort.Slice(data, func(i, j int) bool {
		return data[i].StartTime.Unix() < data[j].StartTime.Unix()
	})

	c.mtx.Lock()
	c.items[tickerCode].Data = data
	summary := c.items[tickerCode].Summary()
	c.mtx.Unlock()

	fmt.Println("Loaded", tickerCode, summary)
	time.Sleep(1000 * time.Millisecond)
}

type Candle struct {
	StartTime  time.Time
	OpenPrice  string
	HighPrice  string
	LowPrice   string
	ClosePrice string
	Volume     string
	Turnover   string
}

type HeaderResponse struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
}

type StringList []string

type GetCandlesResponseInt struct {
	Symbol   string       `json:"symbol"`
	Category string       `json:"category"`
	List     []StringList `json:"list"`
}

type GetCandlesResponse struct {
	HeaderResponse
	Result GetCandlesResponseInt `json:"result"`
}

func (c *Market) LoadData(symbol string, date time.Time) []*Candle {
	date1 := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	date1_end := date1.Add(12 * time.Hour).Add(-1 * time.Millisecond)
	date2 := date1.Add(12 * time.Hour)
	date2_end := date2.Add(12 * time.Hour).Add(-1 * time.Millisecond)

	res1 := c.GetCandles(symbol, date1, date1_end, "1")
	res2 := c.GetCandles(symbol, date2, date2_end, "1")
	res := make([]*Candle, 0)
	res = append(res, res1...)
	res = append(res, res2...)

	sort.Slice(res, func(i, j int) bool {
		return res[i].StartTime.UnixMilli() < res[j].StartTime.UnixMilli()
	})

	return res
}

func (c *Market) GetCandles(symbol string, startDT time.Time, endDT time.Time, interval string) []*Candle {
	time.Sleep(100 * time.Millisecond)
	start := fmt.Sprint(startDT.UnixMilli())
	end := fmt.Sprint(endDT.UnixMilli() - 1)
	requestLine := "https://api.bybit.com/v5/market/kline?category=spot&symbol=" + symbol + "&interval=" + interval + "&start=" + start + "&end=" + end + "&limit=1000"
	resp, err := http.Get(requestLine)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//fmt.Println("Status:", resp.StatusCode)
	buf := make([]byte, 10*1024*1024)
	data := make([]byte, 0)
	for {
		n, err := resp.Body.Read(buf)
		if n == 0 {
			//fmt.Println("0 received")
			break
		}
		if err != nil {
			fmt.Println("err:", err)
		}
		data = append(data, buf[:n]...)
		//buf = buf[:n]

	}
	//fmt.Println(string(data), err)

	var v GetCandlesResponse
	err = json.Unmarshal(data, &v)
	if err != nil {
		fmt.Println("Unmarshal error:", err)
	}
	result := make([]*Candle, 0)
	for _, item := range v.Result.List {
		var c Candle
		timeAsIntMs, _ := strconv.ParseInt(item[0], 10, 64)
		c.StartTime = time.Unix(timeAsIntMs/1000, 0)
		c.OpenPrice = item[1]
		c.HighPrice = item[2]
		c.LowPrice = item[3]
		c.ClosePrice = item[4]
		c.Volume = item[5]
		c.Turnover = item[6]
		result = append(result, &c)
	}
	return result
}
