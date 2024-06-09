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
	return &c
}

func (c *Tokens) Start() {

}

func (c *Tokens) GetTokens() []*Token {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.items
}
