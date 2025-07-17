package cmd

import (
	"encoding/json"
	"fmt"
)

func Echo_Error(err error) {
	fmt.Println(fmt.Printf(`{"items": [
    {
        "uid": "error",
        "title": "Error",
        "subtitle": "%s",
        "arg": "error"
    }
]}`, err.Error()))
}

func Echo_Success(data any) {
	marshal, err := json.Marshal(data)
	if err != nil {
		Echo_Error(err)
	}
	fmt.Println(string(marshal))
}
