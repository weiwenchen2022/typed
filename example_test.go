package typed_test

import (
	"encoding/json"
	"fmt"

	"github.com/weiwenchen2022/typed"
)

func Example() {
	var j = []byte(`{
		"Name": "Wednesday",
		"Age": 6,
		"Parents": ["Gomez","Morticia"]
	}`)

	var m typed.M
	err := json.Unmarshal(j, &m)
	if err != nil {
		panic(err)
	}

	fmt.Println(m.StringValue("Name"))
	fmt.Println(m.AsInt("Age"))
	fmt.Println(m.Array("Parents").Strings())

	var js = []byte(`[6, {"Name": "Wednesday"}, "Gomez"]`)
	var a typed.A
	err = json.Unmarshal(js, &a)
	if err != nil {
		panic(err)
	}

	fmt.Println(a[0].(float64))
	fmt.Println(a[1].(typed.M).StringValue("Name"))
	fmt.Println(a[2].(string))

	// Output:
	// Wednesday
	// 6
	// [Gomez Morticia]
	// 6
	// Wednesday
	// Gomez
}
