package api

import (
	"db_forum/internal/utils"
	"net/http"
)

func printResult(catched error, number int, place string) {
	if catched != nil {
		utils.PrintDebug("api/"+place+" failed(code:", number, "). Error message:"+catched.Error())
	} else {
		utils.PrintDebug("api/"+place+" success(code:", number, ")")
	}
}

func sendJSON(rw http.ResponseWriter, result []byte, place string) {
	// bytes,_ := 	result.MarshalJSON()
	rw.Write(result)
	//json.NewEncoder(rw).Encode(result)
}


