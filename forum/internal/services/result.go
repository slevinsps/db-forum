package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func printResult(catched error, number int, place string) {
	if catched != nil {
		fmt.Println("api/"+place+" failed(code:", number, "). Error message:"+catched.Error())
	} else {
		fmt.Println("api/"+place+" success(code:", number, ")")
	}
}

func sendJSON(rw http.ResponseWriter, result interface{}, place string) {
	json.NewEncoder(rw).Encode(result)
}
