package httpserver

import (
	"encoding/json"

	"github.com/ipoluianov/aneth_eth/an"
	"github.com/ipoluianov/aneth_eth/db"
)

func GetData(code string) []byte {
	if code == "state" {
		type MainState struct {
			DbState *db.DbState
			AnState *an.AnState
		}

		var mainState MainState
		mainState.DbState = db.Instance.GetState()
		mainState.AnState = an.Instance.GetState()
		res, _ := json.MarshalIndent(mainState, "", " ")
		return res
	}
	return nil
}
