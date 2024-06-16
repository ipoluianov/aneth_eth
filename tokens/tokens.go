package tokens

import (
	"sync"
)

type Tokens struct {
	mtx   sync.Mutex
	items []*Token
}

var Instance *Tokens

func init() {
	Instance = NewTokens()
}

func NewTokens() *Tokens {
	var c Tokens
	c.items = append(c.items, NewToken("USDT", 6, "0xdAC17F958D2ee523a2206206994597C13D831ec7", "Tether USD", ""))
	c.items = append(c.items, NewToken("BNB", 18, "0xB8c77482e45F1F44dE1745F52C74426C631bDD52", "BNB", "BNBUSDT"))
	c.items = append(c.items, NewToken("TONCOIN", 9, "0x582d872A1B094FC48F5DE31D3B73F2D9bE47def1", "Wrapped TON Coin", "TONUSDT"))
	c.items = append(c.items, NewToken("SHIB", 18, "0x95aD61b0a150d79219dCF64E1E6Cc01f0B64C4cE", "SHIBA INU", "SHIBUSDT"))
	c.items = append(c.items, NewToken("LINK", 18, "0x514910771AF9Ca656af840dff83E8264EcF986CA", "ChainLink Token", "LINKUSDT"))

	c.items = append(c.items, NewToken("WETH", 18, "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2", "Wrapped Ether", "ETHUSDT"))
	c.items = append(c.items, NewToken("WBTC", 8, "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599", "Wrapped BTC", "BTCUSDT"))
	c.items = append(c.items, NewToken("UNI", 18, "0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984", "Uniswap", "UNIUSDT"))
	c.items = append(c.items, NewToken("NEAR", 24, "0x85F17Cf997934a597031b2E18a9aB6ebD4B9f6a4", "NEAR", "NEARUSDT"))
	c.items = append(c.items, NewToken("PEPE", 18, "0x6982508145454Ce325dDbE47a25d4ec3d2311933", "Pepe", "PEPEUSDT"))
	c.items = append(c.items, NewToken("DAI", 18, "0x6B175474E89094C44Da98b954EedeAC495271d0F", "Dai Stablecoin", "DAIUSDT"))
	c.items = append(c.items, NewToken("FET", 18, "0xaea46A60368A7bD060eec7DF8CBa43b7EF41Ad85", "Fetch", "FETUSDT"))

	return &c
}

func (c *Tokens) Start() {

}

func (c *Tokens) GetTokens() []*Token {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.items
}
