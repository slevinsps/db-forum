package api

import (
	"db_forum/internal/utils"
	"encoding/json"
	"net/http"
)

func printResult(catched error, number int, place string) {
	if catched != nil {
		utils.PrintDebug("api/"+place+" failed(code:", number, "). Error message:"+catched.Error())
	} else {
		utils.PrintDebug("api/"+place+" success(code:", number, ")")
	}
}

func sendJSON(rw http.ResponseWriter, result interface{}, place string) {
	json.NewEncoder(rw).Encode(result)
}
