package an

import "github.com/ipoluianov/aneth_eth/db"

type Task struct {
	Code  string
	Type  string
	Fn    func(result *Result, txsByMin *db.TxsByMinutes, txs []*db.Tx)
	State TaskState

	Name        string
	Description string
}

type TaskState struct {
	Code                   string
	LastExecTime           string
	LastExecTimeDurationMs int
}

func NewTask(code string, tp string, name string, fn func(result *Result, txsByMin *db.TxsByMinutes, txs []*db.Tx), desc string) *Task {
	var c Task
	c.Code = code
	c.Type = tp
	c.Name = name
	c.Description = desc
	c.Fn = fn
	c.State.Code = code
	return &c
}
