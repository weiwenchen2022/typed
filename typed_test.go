package typed

import (
	"encoding/json"
	"math"
	"sort"
	"strconv"
	"testing"
	"time"

	"slices"
)

func TestM_Exists(t *testing.T) {
	t.Parallel()

	var j = []byte(`{
		"Name": "Wednesday",
		"Age": 6,
		"Parents": ["Gomez", "Morticia"]
	}`)

	var m M
	err := json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		key  string
		want bool
	}{
		{"NotExistsKey", false},
		{"Name", true},
	}
	for _, tc := range tests {
		equal(t, tc.want, m.Exists(tc.key))
	}
}

func TestM_IsNumber(t *testing.T) {
	t.Parallel()

	var j = []byte(`{
		"Name": "Wednesday",
		"Age": 6,
		"Parents": ["Gomez", "Morticia"]
	}`)
	var m M
	err := json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		key  string
		want bool
	}{
		{"Name", false},
		{"Age", true},
	}
	for _, tc := range tests {
		equal(t, tc.want, m.IsNumber(tc.key))
	}
}

func TestM(t *testing.T) {
	t.Parallel()

	var j = []byte(`{
		"Name": "Wednesday",
		"Age": 6,
		"Parents": ["Gomez", "Morticia"]
	}`)
	var m M
	err := json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}

	equal(t, "Wednesday", m.StringValue("Name"))
	equal(t, 6, m.AsInt("Age"))
	equalSlice(t, []string{"Gomez", "Morticia"}, m.Array("Parents").Strings())
}

func TestM_Keys(t *testing.T) {
	t.Parallel()

	var j = []byte(`{
		"Name": "Wednesday",
		"Age": 6,
		"Parents": ["Gomez","Morticia"]
	}`)
	var m M
	err := json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}

	keys := m.Keys()
	equal(t, true, sort.StringsAreSorted(keys))

	expected := []string{"Name", "Age", "Parents"}
	sort.Strings(expected)
	equalSlice(t, expected, keys)
}

func TestM_Document(t *testing.T) {
	t.Parallel()

	var j = []byte(`{
		"Profile": {
			"Name": "Wednesday",
			"Age": 6,
			"Parents": ["Gomez", "Morticia"]
		}
	}`)

	var m M
	err := json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}

	doc := m.Document("Profile")
	equal(t, "Wednesday", doc.StringValue("Name"))
	equal(t, 6, doc.AsInt("Age"))
	equalSlice(t, []string{"Gomez", "Morticia"}, doc.Array("Parents").Strings())
}

func TestM_Map(t *testing.T) {
	t.Parallel()

	var j = []byte(`{
		"Inner": {
			"Name": "Wednesday",
			"Age": 6,
			"Parents": ["Gomez", "Morticia"]
		}
	}`)

	var m M
	err := json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}

	inner := m.Map("Inner")
	equal(t, "Wednesday", inner["Name"])
	equal(t, 6, inner["Age"].(float64))
	equalSlice(t, []any{"Gomez", "Morticia"}, inner["Parents"].([]any))
}

func TestM_Any(t *testing.T) {
	t.Parallel()

	var j = []byte(`{
		"Name": "Wednesday",
		"Age": 6,
		"Parents": ["Gomez", "Morticia"]
	}`)

	var m M
	err := json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}

	equal(t, "Wednesday", m.Any("Name"))
	equal(t, 6, m.Any("Age").(float64))
	equalSlice(t, []any{"Gomez", "Morticia"}, m.Any("Parents").([]any))
}

func TestM_AsTime(t *testing.T) {
	t.Parallel()

	now := time.Now()
	jnow, _ := json.Marshal(now)
	var j = []byte(`{
		"datetime": ` + string(jnow) + ` }`,
	)

	var m M
	err := json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	equal(t, true, now.Equal(m.AsTime("datetime")))
}

func TestM_AsDuration(t *testing.T) {
	t.Parallel()

	var j = []byte(`{"timeout": "300ms"}`)
	var m M
	err := json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}
	equal(t, 300*time.Millisecond, m.AsDuration("timeout"))
}

func TestM_RawMessage(t *testing.T) {
	t.Parallel()

	var j = []byte(`{
		"Inner": {
			"Name": "Wednesday",
			"Age": 6,
			"Parents": ["Gomez","Morticia"]}
		}
	`)

	var m M
	err := json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}

	equal(t, `{"Age":6,"Name":"Wednesday","Parents":["Gomez","Morticia"]}`, string([]byte(m.RawMessage("Inner"))))
}

func TestM_RawMessageNull(t *testing.T) {
	t.Parallel()

	var j = []byte(`{"body": null}`)
	var m M
	err := json.Unmarshal(j, &m)
	if err != nil {
		t.Fatal(err)
	}

	equal(t, "null", string([]byte(m.RawMessage("body"))))
}

func TestA(t *testing.T) {
	t.Parallel()

	var j = []byte(`["Gomez", "Morticia"]`)
	var a A
	err := json.Unmarshal(j, &a)
	if err != nil {
		t.Fatal(err)
	}
	equalSlice(t, []string{"Gomez", "Morticia"}, a.Strings())

	err = json.Unmarshal([]byte(`[]`), &a)
	if err != nil {
		t.Fatal(err)
	}
	equalSlice(t, A(nil), a)
}

func TestA_AsInts(t *testing.T) {
	t.Parallel()

	var j = []byte(`[1, 2, 3]`)
	var a A
	err := json.Unmarshal(j, &a)
	if err != nil {
		t.Fatal(err)
	}
	equalSlice(t, []int{1, 2, 3}, a.AsInts())
}

func TestA_AsInt64s(t *testing.T) {
	t.Parallel()

	var j = []byte(`[1, 2, 3]`)
	var a A
	err := json.Unmarshal(j, &a)
	if err != nil {
		t.Fatal(err)
	}
	equalSlice(t, []int64{1, 2, 3}, a.AsInt64s())
}

func TestA_AsInts_WithFloats(t *testing.T) {
	t.Parallel()

	var j = []byte(`[3, 3.14159]`)
	var a A
	err := json.Unmarshal(j, &a)
	if err != nil {
		t.Fatal(err)
	}
	equalSlice(t, []int{3, 3}, a.AsInts())
}

func TestA_Floats(t *testing.T) {
	t.Parallel()

	var j = []byte(`[` + strconv.FormatFloat(math.Pi, 'f', -1, 64) + `]`)
	var a A
	err := json.Unmarshal(j, &a)
	if err != nil {
		t.Fatal(err)
	}
	equalSlice(t, []float64{math.Pi}, a.Floats())
}

func TestA_Strings(t *testing.T) {
	t.Parallel()

	var j = []byte(`["Gomez", "Morticia"]`)
	var a A
	err := json.Unmarshal(j, &a)
	if err != nil {
		t.Fatal(err)
	}
	equalSlice(t, []string{"Gomez", "Morticia"}, a.Strings())
}

func TestA_Documents(t *testing.T) {
	t.Parallel()

	var j = []byte(`[{"id": 1}, {"id": 2}]`)
	var a A
	err := json.Unmarshal(j, &a)
	if err != nil {
		t.Fatal(err)
	}
	equal(t, 1, a.Documents()[0].AsInt("id"))
	equal(t, 2, a.Documents()[1].AsInt("id"))
}

func TestA_Maps(t *testing.T) {
	t.Parallel()

	var j = []byte(`[{"id": 1}, {"id": 2}]`)
	var a A
	err := json.Unmarshal(j, &a)
	if err != nil {
		t.Fatal(err)
	}
	equal(t, 1, a.Maps()[0]["id"].(float64))
	equal(t, 2, a.Maps()[1]["id"].(float64))
}

func TestA_Primitives(t *testing.T) {
	t.Parallel()

	var j = []byte(`[-41, {"id": 1}, 2]`)
	var a A
	err := json.Unmarshal(j, &a)
	if err != nil {
		t.Fatal(err)
	}
	equal(t, -41, a[0].(float64))
	equal(t, 1, a[1].(M).AsInt("id"))
	equal(t, 2, a[2].(float64))
}

func equal[T comparable](tb testing.TB, expected, actual T) {
	tb.Helper()

	if expected != actual {
		tb.Errorf("want: %v (%[1]T); got: %v (%[2]T)", expected, actual)
	}
}

func equalSlice[S ~[]E, E comparable](t *testing.T, expected, actual S) {
	t.Helper()

	if !slices.Equal(expected, actual) {
		t.Errorf("want: %v; got: %v", expected, actual)
	}
}

func panics(f func()) (didPanic bool) {
	defer func() { _ = recover() }()
	didPanic = true
	f()
	return false
}
