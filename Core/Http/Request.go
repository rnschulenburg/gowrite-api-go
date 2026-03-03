package Http

import (
	"encoding/json"
	"net/http"
)

type ReturnValue struct {
	Status       string
	CustomStruct interface{}
}

func ParsePayload(class interface{}, r *http.Request) error {
	t := ReturnValue{Status: "success", CustomStruct: class}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t.CustomStruct)
	if err != nil {
		return err
	}

	class = t.CustomStruct
	return nil
}
