// package typed defines types M and A, help to navigate arbitrary json data
// without knowing this data’s structure beforehand.
//
// Specifically, M and A wrapped map[string]any and []any recurse down.
// If the given key is multiple keys conjunction with dot "." character,
// corresponding methods also recurse down.
package typed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

// M is a representation of a JSON document.
type M map[string]any

var null = []byte("null")

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *M) UnmarshalJSON(data []byte) error {
	if bytes.Equal(null, data) {
		return nil
	}

	var mm map[string]any
	if err := json.Unmarshal(data, &mm); err != nil {
		return err
	}

	*m = wrapper(mm).(M)
	return nil
}

// An A is a representation of a JSON array.
type A []any

// UnmarshalJSON implements the json.Unmarshaler interface.
func (a *A) UnmarshalJSON(data []byte) error {
	if bytes.Equal(null, data) {
		return nil
	}

	var aa []any
	if err := json.Unmarshal(data, &aa); err != nil {
		return err
	}

	*a = wrapper(aa).(A)
	return nil
}

// Wrap wraps map[string]any and []any to M and A recurse down.
func Wrap(a any) any {
	return wrapper(a)
}

func wrapper(a any) any {
	switch x := a.(type) {
	default:
		return a
	case map[string]any:
		for k, v := range x {
			x[k] = wrapper(v)
		}
		return M(x)
	case []any:
		for i, v := range x {
			x[i] = wrapper(v)
		}
		return A(x)
	}
}

// Unwrap unwraps M and A to map[string]any and []any recurse down.
func Unwrap(a any) any {
	return unwrapper(a)
}

func unwrapper(a any) any {
	switch x := a.(type) {
	default:
		return a
	case M:
		for k, v := range x {
			x[k] = unwrapper(v)
		}
		return map[string]any(x)
	case A:
		for i, v := range x {
			x[i] = unwrapper(v)
		}
		return []any(x)
	}
}

// Exists reports whether key exists, potentially recursively for the given key. If
// there are multiple keys concatenated with ".", this method will recurse down, as long as the
// top and intermediate nodes are either documents or arrays. If an error
// occurs or if the value doesn't exist, false is returned.
func (m M) Exists(key string) bool {
	_, ok := lookupOK[any](m, key)
	return ok
}

// IsNumber reports whether the value represents for given key is a JSON number.
func (m M) IsNumber(key string) bool {
	_, ok := lookupOK[float64](m, key)
	return ok
}

// Bool returns the boolean value the value represents for given key. It panics if the
// value is a JSON type other than boolean.
func (m M) Bool(key string) bool {
	return lookup[bool](m, key)
}

// BoolOK is the same as Bool, except it returns a boolean instead of
// panicking.
func (m M) BoolOK(key string) (bool, bool) {
	return lookupOK[bool](m, key)
}

// AsInt returns the int value the value represents for given key. It panics if the
// value is JSON type other than number.
func (m M) AsInt(key string) int {
	return int(lookup[float64](m, key))
}

// AsIntOK is the same as AsInt, except that it returns a boolean instead of
// panicking.
func (m M) AsIntOK(key string) (int, bool) {
	f, ok := lookupOK[float64](m, key)
	return int(f), ok
}

// AsInt64 returns a JSON number as an int64 for given key. It panics if the
// value type is JSON type other than number.
func (m M) AsInt64(key string) int64 {
	return int64(lookup[float64](m, key))
}

// AsInt64OK is the same as AsInt64, except that it returns a boolean instead of
// panicking.
func (m M) AsInt64OK(key string) (int64, bool) {
	f, ok := lookupOK[float64](m, key)
	return int64(f), ok
}

// Float returns the float64 value the value represents for given key. It panics if the
// value is JSON type other than number.
func (m M) Float(key string) float64 {
	return lookup[float64](m, key)
}

// FloatOK is the same as Float, but returns a boolean instead of panicking.
func (m M) FloatOK(key string) (float64, bool) {
	return lookupOK[float64](m, key)
}

// StringValue returns the string value the value represents for given key. It panics if the
// value is JSON type other than string.
//
// NOTE: This method is called StringValue to avoid a collision with the String method which
// implements the fmt.Stringer interface.
func (m M) StringValue(key string) string {
	return lookup[string](m, key)
}

// StringValueOK is the same as StringValue, but returns a boolean instead of
// panicking.
func (m M) StringValueOK(key string) (string, bool) {
	return lookupOK[string](m, key)
}

// AsTime returns the time.Time value the value represents for given key. It panics if the
// value not represents time.AsTime.
func (m M) AsTime(key string) time.Time {
	s := lookup[string](m, key)

	var t time.Time
	if err := t.UnmarshalJSON([]byte(strconv.Quote(s))); err != nil {
		panic(err)
	}
	return t
}

// AsTimeOK is the same as AsTime, except it returns a boolean instead of
// panicking.
func (m M) AsTimeOK(key string) (time.Time, bool) {
	s, ok := lookupOK[string](m, key)
	if !ok {
		return time.Time{}, false
	}

	var t time.Time
	err := t.UnmarshalJSON([]byte(strconv.Quote(s)))
	return t, err == nil
}

// AsDuration returns the time.Duration value the value represents for given key. It panics if the
// value can not parsed by time.ParseDuration.
func (m M) AsDuration(key string) time.Duration {
	s := lookup[string](m, key)

	var d time.Duration
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}

// AsDurationOK is the same as AsDuration, except it returns a boolean instead of
// panicking.
func (m M) AsDurationOK(key string) (time.Duration, bool) {
	s, ok := lookupOK[string](m, key)
	if !ok {
		return 0, false
	}

	d, err := time.ParseDuration(s)
	return d, err == nil
}

// Array returns the JSON array the value represents for given key. It panics if the
// value is a JSON type other than array.
func (m M) Array(key string) A {
	return lookup[A](m, key)
}

// ArrayOK is the same as Array, except it returns a boolean instead of
// panicking.
func (m M) ArrayOK(key string) (A, bool) {
	return lookupOK[A](m, key)
}

// Document returns the JSON document the value represents for given key. It panics if the
// value is a JSON type other than document.
func (m M) Document(key string) M {
	return lookup[M](m, key)
}

// DocumentOK is the same as Document, except it returns a boolean instead of
// panicking.
func (m M) DocumentOK(key string) (M, bool) {
	return lookupOK[M](m, key)
}

// Map is the same as Document, except it returns a map[string]any
// instead of M.
func (m M) Map(key string) map[string]any {
	return unwrapper(lookup[M](m, key)).(map[string]any)
}

// MapOK is the same as Map, except it returns a boolean instead of
// panicking.
func (m M) MapOK(key string) (map[string]any, bool) {
	m2, ok := lookupOK[M](m, key)
	return unwrapper(m2).(map[string]any), ok
}

var nullRawMessage = json.RawMessage([]byte("null"))

// RawMessage returns the raw encoded JSON value the value represents for given key. It returns 'null' if the
// value doesn't exist.
func (m M) RawMessage(key string) json.RawMessage {
	v, ok := lookupOK[any](m, key)
	if !ok {
		return nullRawMessage
	}

	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return json.RawMessage(b)
}

// Any search the document, potentially recursively, for the given key. If
// there are multiple keys concatenated with ".", this method will recurse down, as long as the
// top and intermediate nodes are either documents or arrays. If an error
// occurs or if the value doesn't exist, this method panics.
func (m M) Any(key string) any {
	return unwrapper(lookup[any](m, key))
}

// AnyOK is the same as Any, except it returns a boolean instead of
// panicking.
func (m M) AnyOK(key string) (any, bool) {
	a, ok := lookupOK[any](m, key)
	return unwrapper(a), ok
}

// Keys returns all sorted keys within document.
func (m M) Keys() []string {
	if m == nil {
		return nil
	}

	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

// Bools returns the slice of boolean the array represents. It panics if the
// one of the elements is JSON type other than boolean.
func (a A) Bools() []bool {
	return array[bool](a)
}

// BoolsOK is the same as Bools, except it returns a boolean instead of
// panicking.
func (a A) BoolsOK() ([]bool, bool) {
	return arrayOK[bool](a)
}

// AsInts returns the slice of int the array represents. It panics if one
// of elements is a JSON type other than number.
func (a A) AsInts() []int {
	return asNumericArray[int](a)
}

// AsIntsOK is the same as AsInts, except is returns a boolean instead of
// panicking.
func (a A) AsIntsOK() ([]int, bool) {
	return asNumericArrayOK[int](a)
}

// AsInt64s returns the slice of int64 the array represents. It panics if one
// of elements is a JSON type other than number.
func (a A) AsInt64s() []int64 {
	return asNumericArray[int64](a)
}

// AsInt64sOK is the same as AsInt64s, except that it returns a boolean instead of
// panicking.
func (a A) AsInt64sOK() ([]int64, bool) {
	return asNumericArrayOK[int64](a)
}

// Floats returns the slice of float64 value the array represents. It panics if one
// of elements is a JSON type other than number.
func (a A) Floats() []float64 {
	return asNumericArray[float64](a)
}

// FloatsOK is the same as Floats, except that it returns a boolean instead of
// panicking.
func (a A) FloatsOK() ([]float64, bool) {
	return asNumericArrayOK[float64](a)
}

// Strings returns the slice of string the array represents. It panics if one
// of elements is a JSON type other than string.
func (a A) Strings() []string {
	return array[string](a)
}

// StringsOK is the same as Strings, except it returns a boolean instead of
// panicking.
func (a A) StringsOK() ([]string, bool) {
	return arrayOK[string](a)
}

// Documents returns the slice of JSON document the array represents. It panics if one
// of elements is a JSON type other than document.
func (a A) Documents() []M {
	return array[M](a)
}

// DocumentsOK is the same as Documents, except that it returns a boolean instead of
// panicking.
func (a A) DocumentsOK() ([]M, bool) {
	return arrayOK[M](a)
}

// Maps returns the slice of JSON document the array represents. It panics if one
// of elements is a JSON type other than document.
func (a A) Maps() []map[string]any {
	documents := array[M](a)
	s := make([]map[string]any, len(documents))
	for i, document := range documents {
		s[i] = unwrapper(document).(map[string]any)
	}
	return s
}

// MapsOK is the same as Maps, but returns a boolean instead of
// panicking.
func (a A) MapsOK() ([]map[string]any, bool) {
	documents, ok := arrayOK[M](a)
	if !ok {
		return nil, false
	}

	s := make([]map[string]any, len(documents))
	for i, document := range documents {
		s[i] = unwrapper(document).(map[string]any)
	}
	return s, true
}

func asNumericArray[E constraints.Integer | constraints.Float](a []any) []E {
	if a == nil {
		return nil
	}

	s := make([]E, len(a))
	for i, v := range a {
		s[i] = E(v.(float64))
	}
	return s
}

func asNumericArrayOK[E constraints.Integer | constraints.Float](a []any) (s []E, ok bool) {
	if a == nil {
		return nil, true
	}

	s = make([]E, len(a))
	var f float64
	for i, v := range a {
		f, ok = v.(float64)
		if !ok {
			return nil, false
		}
		s[i] = E(f)
	}
	return s, true
}

func array[E any](a []any) []E {
	if a == nil {
		return nil
	}

	s := make([]E, len(a))
	for i, v := range a {
		s[i] = v.(E)
	}
	return s
}

func arrayOK[E any](a []any) (s []E, ok bool) {
	if a == nil {
		return nil, true
	}

	s = make([]E, len(a))
	for i, v := range a {
		s[i], ok = v.(E)
		if !ok {
			return nil, false
		}
	}
	return s, true
}

func lookup[E any](a any, key string) E {
	e, err := lookupErr[E](a, key)
	if err != nil {
		panic(err)
	}
	return e
}

func lookupOK[E any](a any, key string) (e E, ok bool) {
	var err error
	e, err = lookupErr[E](a, key)
	return e, err == nil
}

func lookupErr[E any](a any, key string) (e E, err error) {
	keys := strings.Split(key, ".")
	for i, k := range keys[:len(keys)-1] {
		switch x := a.(type) {
		default:
			err = fmt.Errorf("unknown type %T", a)
			return
		case M:
			var ok bool
			a, ok = x[k]
			if !ok {
				err = fmt.Errorf("not found key %q", strings.Join(keys[:i+1], "."))
				return
			}
		case A:
			var j int
			j, err = strconv.Atoi(k)
			if err != nil {
				return
			}
			if j < 0 || j >= len(x) {
				err = fmt.Errorf("not found key %q", strings.Join(keys[:i+1], "."))
				return
			}
			a = x[j]
		}
	}

	lastKey := keys[len(keys)-1]
	switch x := a.(type) {
	default:
		err = fmt.Errorf("unknown type %T", a)
		return
	case M:
		v, ok := x[lastKey]
		if !ok {
			err = fmt.Errorf("not found key %q", key)
			return
		}

		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v", r)
			}
		}()
		return v.(E), nil
	case A:
		var i int
		i, err = strconv.Atoi(lastKey)
		if err != nil {
			return
		}
		if i < 0 || i >= len(x) {
			err = fmt.Errorf("not found key %q", key)
			return
		}

		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v", r)
			}
		}()
		return x[i].(E), nil
	}
}
