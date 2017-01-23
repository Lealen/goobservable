package goob

//copyright Mateusz Tomasz Doroszko for project s3n

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func getTypeFunc(i interface{}) string {
	return reflect.TypeOf(i).String()
}

//package specified:
func getGlobalFunc() map[string]interface{} {
	return Data
}

//Generic:

func logFunc(str interface{}) string {
	println(fmt.Sprintf("%v", str))
	return ""
}

func howManyLikeForFunc(number int) []int {
	var ret []int
	ret = make([]int, number)
	for i := 0; i < number; i++ {
		ret[i] = i
	}
	return ret
}

//Memory banks:
type TemplateMemoryBank struct {
	V *map[interface{}]interface{}
}

func newMemoryBankFunc() TemplateMemoryBank {
	V := make(map[interface{}]interface{})
	return TemplateMemoryBank{&V}
}

func setToMemoryBankFunc(tmb TemplateMemoryBank, key, value interface{}) string {
	v := *tmb.V
	v[key] = value
	return ""
}

func getFromMemoryBankFunc(tmb TemplateMemoryBank, key interface{}) interface{} {
	v := *tmb.V
	return v[key]
}

//Conversions:

func derefFunc(a interface{}) interface{} {
	valueA := reflect.ValueOf(a)
	return reflect.Indirect(valueA)
}

func doNotParseFunc(in string) template.HTML {
	return template.HTML(in)
}

func doNotParseJSFunc(in string) template.JS {
	return template.JS(in)
}

func doNotParseCSSFunc(in string) template.CSS {
	return template.CSS(in)
}

func doNotParseURLFunc(in string) template.URL {
	return template.URL(in)
}

func boolFunc(a interface{}) (bool, error) {
	valueA := reflect.ValueOf(a)

	switch valueA.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return valueA.Int() != 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return valueA.Uint() != 0, nil
	case reflect.Float32, reflect.Float64:
		return valueA.Float() > 0, nil
	case reflect.String:
		if valueA.String() == "false" || valueA.String() == "0" || valueA.String() == "" || valueA.String() == " " {
			return false, nil
		}
		return true, nil
	default:
		return false, fmt.Errorf("Type cant be converted to bool.")
	}
}

func intFunc(a interface{}) (int, error) {
	valueA := reflect.ValueOf(a)

	switch valueA.Kind() {
	case reflect.Bool:
		if valueA.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(valueA.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int(valueA.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int(valueA.Float()), nil
	case reflect.String:
		amount, err := strconv.ParseInt(valueA.String(), 10, 0)
		return int(amount), err
	default:
		return 0, fmt.Errorf("Type cant be converted to int.")
	}
}

func uintFunc(a interface{}) (uint, error) {
	valueA := reflect.ValueOf(a)

	switch valueA.Kind() {
	case reflect.Bool:
		if valueA.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint(valueA.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint(valueA.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return uint(valueA.Float()), nil
	case reflect.String:
		amount, err := strconv.ParseInt(valueA.String(), 10, 0)
		return uint(amount), err
	default:
		return 0, fmt.Errorf("Type cant be converted to uint.")
	}
}

func int64Func(a interface{}) (int64, error) {
	valueA := reflect.ValueOf(a)

	switch valueA.Kind() {
	case reflect.Bool:
		if valueA.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int64(valueA.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(valueA.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int64(valueA.Float()), nil
	case reflect.String:
		return strconv.ParseInt(valueA.String(), 10, 64)
	default:
		return 0, fmt.Errorf("Type cant be converted to int64.")
	}
}

func uint64Func(a interface{}) (uint64, error) {
	valueA := reflect.ValueOf(a)

	switch valueA.Kind() {
	case reflect.Bool:
		if valueA.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(valueA.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint64(valueA.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return uint64(valueA.Float()), nil
	case reflect.String:
		return strconv.ParseUint(valueA.String(), 10, 64)
	default:
		return 0, fmt.Errorf("Type cant be converted to uint64.")
	}
}

func float32Func(a interface{}) (float32, error) {
	valueA := reflect.ValueOf(a)

	switch valueA.Kind() {
	case reflect.Bool:
		if valueA.Bool() {
			return 1.0, nil
		}
		return 0.0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float32(valueA.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float32(valueA.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return float32(valueA.Float()), nil
	case reflect.String:
		amount, err := strconv.ParseFloat(valueA.String(), 32)
		return float32(amount), err
	default:
		return 0.0, fmt.Errorf("Type cant be converted to float32.")
	}
}

func float64Func(a interface{}) (float64, error) {
	valueA := reflect.ValueOf(a)

	switch valueA.Kind() {
	case reflect.Bool:
		if valueA.Bool() {
			return 1.0, nil
		}
		return 0.0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(valueA.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(valueA.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return float64(valueA.Float()), nil
	case reflect.String:
		return strconv.ParseFloat(valueA.String(), 64)
	default:
		return 0.0, fmt.Errorf("Type cant be converted to float64.")
	}
}

func stringFunc(a interface{}) (string, error) {
	valueA := reflect.ValueOf(a)

	switch valueA.Kind() {
	case reflect.Bool:
		if valueA.Bool() {
			return "true", nil
		}
		return "false", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(valueA.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(valueA.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(valueA.Float(), 'f', 64, 64), nil
	case reflect.String:
		return valueA.String(), nil
	default:
		return "", fmt.Errorf("Type cant be converted to string.")
	}
}

func toJSONFunc(a interface{}) (string, error) {
	tmp, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(tmp), nil
}

//Arrays:

func setIndexFunc(array, index, a interface{}) (interface{}, error) {
	valueArray := reflect.ValueOf(array)
	valueIndex := reflect.ValueOf(index)
	valueA := reflect.ValueOf(a)

	valueA.Pointer()

	switch valueArray.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		if valueArray.Len() <= int(valueIndex.Int()) {
			return nil, fmt.Errorf("Index out of range.")
		}
		valueArray.Elem().Field(int(valueIndex.Int())).Set(valueA)
		return valueArray.Interface(), nil
	case reflect.Map:
		valueArray.SetMapIndex(valueIndex, valueA)
		return valueArray.Interface(), nil
	default:
		return nil, fmt.Errorf("Type isnt array, slice, map or string.")
	}
}

func getIndexFunc(array, index interface{}) (interface{}, error) {
	valueArray := reflect.ValueOf(array)
	valueIndex := reflect.ValueOf(index)

	switch valueArray.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		if valueArray.Len() <= int(valueIndex.Int()) {
			return nil, fmt.Errorf("Index out of range.")
		}
		return valueArray.Index(int(valueIndex.Int())).Interface(), nil
	case reflect.Map:
		return valueArray.MapIndex(valueIndex).Interface(), nil
	default:
		return nil, fmt.Errorf("Type isnt array, slice, map or string.")
	}
}

//Math:

func addFunc(a, b interface{}) (interface{}, error) {
	valueA := reflect.ValueOf(a)
	valueB := reflect.ValueOf(b)

	if valueA.Kind() != valueB.Kind() {
		return nil, fmt.Errorf("Different kinds, can't add them.")
	}

	switch valueA.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return valueA.Int() + valueB.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return valueA.Uint() + valueB.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return valueA.Float() + valueB.Float(), nil
	case reflect.String:
		return valueA.String() + valueB.String(), nil
	default:
		return nil, fmt.Errorf("Type does not support addition.")
	}
}

func subtractFunc(a, b interface{}) (interface{}, error) {
	valueA := reflect.ValueOf(a)
	valueB := reflect.ValueOf(b)

	if valueA.Kind() != valueB.Kind() {
		return nil, fmt.Errorf("Different kinds, can't subtract them.")
	}

	switch valueA.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return valueA.Int() - valueB.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return valueA.Uint() - valueB.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return valueA.Float() - valueB.Float(), nil
	default:
		return nil, fmt.Errorf("Type does not support subtraction.")
	}
}

func multiplyFunc(a, b interface{}) (interface{}, error) {
	valueA := reflect.ValueOf(a)
	valueB := reflect.ValueOf(b)

	if valueA.Kind() != valueB.Kind() && valueA.Kind() != reflect.String {
		return nil, fmt.Errorf("Different kinds, can't multiply them.")
	}

	switch valueA.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return valueA.Int() * valueB.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return valueA.Uint() * valueB.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return valueA.Float() * valueB.Float(), nil
	case reflect.String:
		switch valueB.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return strings.Repeat(valueA.String(), int(valueB.Int())), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return strings.Repeat(valueA.String(), int(valueB.Uint())), nil
		default:
			return nil, fmt.Errorf("Second Variable needs to by kind of int.")
		}
	default:
		return nil, fmt.Errorf("Type does not support multiplication.")
	}
}

func divideFunc(a, b interface{}) (interface{}, error) {
	valueA := reflect.ValueOf(a)
	valueB := reflect.ValueOf(b)

	if valueA.Kind() != valueB.Kind() {
		return nil, fmt.Errorf("Different kinds, can't divide them.")
	}

	switch valueA.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return valueA.Int() / valueB.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return valueA.Uint() / valueB.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return valueA.Float() / valueB.Float(), nil
	default:
		return nil, fmt.Errorf("Type does not support dividation.")
	}
}

func modFunc(a, b interface{}) (interface{}, error) {
	valueA := reflect.ValueOf(a)
	valueB := reflect.ValueOf(b)

	if valueA.Kind() != valueB.Kind() {
		return nil, fmt.Errorf("Different kinds, can't divide them.")
	}

	switch valueA.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return valueA.Int() % valueB.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return valueA.Uint() % valueB.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return math.Mod(valueA.Float(), valueB.Float()), nil
	default:
		return nil, fmt.Errorf("Type does not support modulation.")
	}
}

func ceilFunc(a interface{}) (int, error) {
	valueA := reflect.ValueOf(a)

	switch valueA.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(math.Ceil(float64(valueA.Int()))), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int(math.Ceil(float64(valueA.Uint()))), nil
	case reflect.Float32, reflect.Float64:
		return int(math.Ceil(valueA.Float())), nil
	default:
		return 0, fmt.Errorf("Type does not support ceil.")
	}
}

func floorFunc(a interface{}) (int, error) {
	valueA := reflect.ValueOf(a)

	switch valueA.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(math.Floor(float64(valueA.Int()))), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int(math.Floor(float64(valueA.Uint()))), nil
	case reflect.Float32, reflect.Float64:
		return int(math.Floor(valueA.Float())), nil
	default:
		return 0, fmt.Errorf("Type does not support floor.")
	}
}

func minFunc(a, b interface{}) (interface{}, error) {
	valueA := reflect.ValueOf(a)
	valueB := reflect.ValueOf(b)

	if valueA.Kind() != valueB.Kind() {
		return nil, fmt.Errorf("Different kinds, can't min them.")
	}

	switch valueA.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if valueA.Int() < valueB.Int() {
			return valueA.Int(), nil
		}
		return valueB.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if valueA.Uint() < valueB.Uint() {
			return valueA.Uint(), nil
		}
		return valueB.Uint(), nil
	case reflect.Float32, reflect.Float64:
		if valueA.Float() < valueB.Float() {
			return valueA.Float(), nil
		}
		return valueB.Float(), nil
	default:
		return nil, fmt.Errorf("Type does not support min.")
	}
}

func maxFunc(a, b interface{}) (interface{}, error) {
	valueA := reflect.ValueOf(a)
	valueB := reflect.ValueOf(b)

	if valueA.Kind() != valueB.Kind() {
		return nil, fmt.Errorf("Different kinds, can't max them.")
	}

	switch valueA.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if valueA.Int() > valueB.Int() {
			return valueA.Int(), nil
		}
		return valueB.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if valueA.Uint() > valueB.Uint() {
			return valueA.Uint(), nil
		}
		return valueB.Uint(), nil
	case reflect.Float32, reflect.Float64:
		if valueA.Float() > valueB.Float() {
			return valueA.Float(), nil
		}
		return valueB.Float(), nil
	default:
		return nil, fmt.Errorf("Type does not support max.")
	}
}

//Date:

func timenowFunc() time.Time {
	return time.Now()
}

func timeutcFunc(t time.Time) time.Time {
	return t.UTC()
}

func timeparsedurationFunc(duration string) (time.Duration, error) {
	return time.ParseDuration(duration)
}

func timeaddFunc(t time.Time, d time.Duration) time.Time {
	return t.Add(d)
}

func timeformatFunc(t time.Time, layout string) string {
	return t.Format(layout)
}

func timeyearFunc(t time.Time) int {
	return t.Year()
}

func timemonthFunc(t time.Time) int {
	return int(t.Month())
}

func timedayFunc(t time.Time) int {
	return t.Day()
}

func timehourFunc(t time.Time) int {
	return t.Hour()
}

func timeminuteFunc(t time.Time) int {
	return t.Minute()
}

func timesecondFunc(t time.Time) int {
	return t.Second()
}

func timeparseunixFunc(u int64) time.Time {
	return time.Unix(u, 0)
}
