package datatypes

import (
	"encoding/json"
	"fmt"
	"testing"
)

type TestStruct struct {
	ID    int      `json:"id"`
	OTime OnlyTime `json:"otime"`
}

func TestUnmarshalOnlyTime(t *testing.T) {
	fmt.Println("\n\nTestJsonBind ***")

	jsonstr := `{"id":1,"otime":"16:04:21"}`
	var dtest TestStruct
	if err := json.Unmarshal([]byte(jsonstr), &dtest); err != nil {
		t.Error(err)
	}
	fmt.Printf("struct: %#v\n", dtest)
}

//
func TestUnmarshalOnlyTime2(t *testing.T) {
	fmt.Println("\n\nTestJsonBind ***")

	jsonstr := `{"id":1,"otime":"2022-08-20T15:00:00.000Z"}`
	var dtest TestStruct
	if err := json.Unmarshal([]byte(jsonstr), &dtest); err != nil {
		t.Error(err)
	}
	fmt.Printf("struct: %#v\n", dtest)
}
