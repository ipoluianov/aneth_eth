package common

import "github.com/ipoluianov/aneth_eth/db"

type Task struct {
	Code  string
	Type  string
	Fn    func(task *Task, result *Result, txsByMin *db.TxsByMinutes, txs []*db.Tx)
	State TaskState

	Name        string
	Description string
	Text        string
	Ticker      string
	Symbol      string
}

type TaskState struct {
	Code                   string
	LastExecTime           string
	LastExecTimeDurationMs int
}
