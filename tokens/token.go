package tokens

type Token struct {
	Name     string
	Symbol   string
	Address  string
	Decimals int
}

func NewToken(symbol string, decimals int, address string, name string) *Token {
	var c Token
	c.Name = name
	c.Symbol = symbol
	c.Address = address
	c.Decimals = decimals
	return &c
}
