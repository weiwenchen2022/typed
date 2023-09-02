## Typed

*Making `map[string]any and []any` slightly less painful*

It isn't always desirable or even possible to decode JSON into a structure.
That, unfortunately, leaves us with writing ugly code around `map[string]any` and `[]any`.

This library hopes to make that slightly less painful.

## Install

```sh
go get github.com/weiwenchen2022/typed
```

## Note

Although the library can be used with any `map[string]any` and `[]any`, it *is* tailored to be used with `encoding/json`.

Specifically, it wrapped `map[string]any` and `[]any` to `M` and `A` recurse down. If the given key is multiple keys conjunction with dot "." character, corresponding methods also recurse down.

## Usage:

A typed wrapper around a `map[string]any` and `[]any` can be created in one of two ways:

```go
// Wrapped from a map[string]any or []any
m := typed.Wrap(originalM).(M)

a := typed.Wrap(originalA).(A)

// Decoded from a json data
var m typed.M
err := json.Unmarshal(data, &m)

var a typed.A
err := json.Unmarshal(data, &a)
```

Once we have a wrapper `M`, we can use various methods to navigate the structure:

- `(M) Bool(key string) bool`
- `(M) BoolOK(key string) (bool, bool)`

- `(M) AsInt(key string) int`
- `(M) AsIntOK(key string) (int, bool)`

- `(M) AsInt64(key string) int64`
- `(M) AsInt64OK(key string) (int64, bool)`

- `(M) Float(key string) float64`
- `(M) FloatOK(key string) (float64, bool)`

- `(M) StringValue(key string) string`
- `(M) StringValueOK(key string) (string, bool)`

- `(M) Any(key string) any`
- `(M) AnyOK(key string) (any, bool)`

We can also extract arrays via wrapper `A`:

- `(M) Array(key string) A`

`A` has various methods to convert to concrete array:

- `(A) Bools(key string) []bool`
- `(A) BoolsOK(key string) ([]bool, bool)`

- `(A) AsInts(key string) []int`
- `(A) AsIntsOK(key string) ([]int, bool)`

- `(A) AsInt64s(key string) []int64`
- `(A) AsInt64sOK(key string) ([]int64, bool)`

- `(A) Floats(key string) []float64`
- `(A) FloatsOK(key string) ([]float64, bool)`

- `(A) Strings(key string) []string`
- `(A) StringsOK(key string) ([]string, bool)`

We can extract nested document, as a nested wrapper, or as a original `map[string]any`:

- `(M) Document(key string) M`
- `(M) DocumentOK(key string) (M, bool)`

- `(M) Map(key string) map[string]any`
- `(M) MapOK(key string) (map[string]any, bool)`

- `(A) Documents(key string) []M`
- `(A) DocumentsOK(key string) ([]M, bool)`

- `(A) Maps(key string) []map[string]any`
- `(A) MapsOK(key string) ([]map[string]any, bool)`

## Example

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/weiwenchen2022/typed"
)

func main() {
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
}
```

# Misc

## AsTime and AsDuration
`(M) AsTime(key string) time.Time` can be used to get the string value, as a time.Time.

Alternatively, `(M) AsDuration(key string) time.Duration` can be used to get the string value, as a time.Duration.

## Root Array

JSON array at root is supported. Use the `A` to decoded, or wrapped from []any.

```go
var js = []byte(`[6, {"Name": "Wednesday"}, "Gomez"]`)
var a typed.A
err := json.Unmarshal(js, &a)
if err != nil {
	panic(err)
}

fmt.Println(a[0].(float64))
fmt.Println(a[1].(typed.M).StringValue("Name"))
fmt.Println(a[2].(string))
```

## Doc
GoDoc: [https://godoc.org/github.com/weiwenchen2022/typed](https://godoc.org/github.com/weiwenchen2022/typed)
