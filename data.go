package goob

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jsbuiltin"
)

var Data = make(map[string]interface{})

func DebugData() {
	println(strings.Replace(fmt.Sprintf("%#v", Data), ",", ",\n", -1))
}

func expressionApply(obj *js.Object, trueorfalse bool) {
	todo := ""

	//TODO: optymalize this code, so it wont run every time a variable change, we can just do it when expresion changes, but anyway... it works I suppose?

	if trueorfalse {
		if jsbuiltin.TypeOf(obj.Get("hasAttribute")) == "function" && obj.Call("hasAttribute", goobIfTrue).Bool() {
			todo = obj.Call("getAttribute", goobIfTrue).String()
		}
	} else {
		if jsbuiltin.TypeOf(obj.Get("hasAttribute")) == "function" && obj.Call("hasAttribute", goobIfFalse).Bool() {
			todo = obj.Call("getAttribute", goobIfFalse).String()
		}
	}

	if todo != "" {
		arr := strings.Split(todo, ";")
		for _, v := range arr {
			arr2 := strings.Split(v, " ")
			switch arr2[0] {
			case "addclass":
				for k, v2 := range arr2 {
					if k == 0 {
						continue
					}
					obj.Get("classList").Call("add", v2)
				}
			case "removeclass":
				for k, v2 := range arr2 {
					if k == 0 {
						continue
					}
					obj.Get("classList").Call("remove", v2)
				}
			}
		}
	}
}

func convertToInt64(value reflect.Value) int64 {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int()
	}
	println("Warning: trying to convert type " + value.Kind().String() + " to type int64!")
	return 0
}

func checkOneExpression(value reflect.Value, ei *js.Object) {
	switch value.Kind() {
	case reflect.Bool:
		expressionApply(ei, value.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		number := convertToInt64(value)
		exp := ei.Call("getAttribute", goobIfExpression).String()
		exparr := strings.Split(exp, " ")
		if len(exparr) != 2 {
			println("Warning: wrong expression: " + exp)
		} else {
			number2, err := strconv.ParseInt(exparr[1], 10, 64)
			if err != nil {
				println("Warning: error parsing int: " + exparr[1])
			}
			// println(fmt.Sprintf("%d %d", number, number2))
			switch exparr[0] {
			case "eq":
				expressionApply(ei, number == number2)
			case "ne":
				expressionApply(ei, number != number2)
			case "lt":
				expressionApply(ei, number < number2)
			case "le":
				expressionApply(ei, number <= number2)
			case "gt":
				expressionApply(ei, number > number2)
			case "ge":
				expressionApply(ei, number >= number2)
			}
		}
	}

	ei.Call("setAttribute", goobIfExpressionProcessed, "")
}

func checkExpressions(name string) {
	//CHECK EXPRESSIONS ON DOCUMENT:

	value := reflect.ValueOf(Data[name])

	every := js.Global.Get("document").Call("querySelectorAll", "["+goobIfVar+"=\""+name+"\"]")
	for i := 0; i < every.Length(); i++ {
		checkOneExpression(value, every.Index(i))
	}
}

func foreachExpressionCheck(ei *js.Object, valueA reflect.Value) bool {
	if jsbuiltin.TypeOf(ei.Get("hasAttribute")) == "function" && ei.Call("hasAttribute", goobForeachIf).Bool() {
		if valueA.Kind() != reflect.Struct {
			println("goob-foreach-if for not-structs are not supported yet")
			return false
		}

		ifele := ei.Call("getAttribute", goobForeachIf).String()
		ifelearr := strings.Split(ifele, " ")

		if len(ifelearr) == 0 {
			return true
		} else if len(ifelearr) == 1 {
			ifelearr = append(ifelearr, "")
		}

		element := valueA.FieldByName(ifelearr[0])

		if fmt.Sprintf("%v", element.Interface()) == strings.Join(ifelearr[1:], " ") {
			return true
		}
		return false
	}
	return true
}

func Get(name string) interface{} {
	return Data[name]
}

func Set(name string, value interface{}) {
	Data[name] = value

	//UPDATE ON DOCUMENT:
	every := js.Global.Get("document").Call("querySelectorAll", "["+goobData+"=\""+name+"\"], ["+goobWatch+"=\""+name+"\"]")
	for i := 0; i < every.Length(); i++ {
		every.Index(i).Call("removeAttribute", goobProcessed)
	}

	every = js.Global.Get("document").Call("querySelectorAll", "["+goobIfVar+"='"+name+"']")
	for i := 0; i < every.Length(); i++ {
		every.Index(i).Call("removeAttribute", goobIfExpressionProcessed)
	}

	Tick()
	// checkExpressions(name)
}

func GetIndex(name string, index interface{}) interface{} {
	if _, ok := Data[name]; !ok {
		return nil
	}

	valueArray := reflect.ValueOf(Data[name])
	valueIndex := reflect.ValueOf(index)

	switch valueArray.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		if valueArray.Len() <= int(valueIndex.Int()) {
			return nil //, fmt.Errorf("Index out of range.")
		}
		return valueArray.Index(int(valueIndex.Int())).Interface() //, nil
	case reflect.Map:
		return valueArray.MapIndex(valueIndex).Interface() //, nil
	default:
		return nil //, fmt.Errorf("Type isnt array, slice, map or string.")
	}
}

func SetIndex(name string, index, value interface{}) (interface{}, error) {
	if _, ok := Data[name]; !ok {
		return nil, fmt.Errorf("You need to define string/array/slice/map first")
	}

	valueArray := reflect.ValueOf(Data[name])
	valueIndex := reflect.ValueOf(index)
	valueA := reflect.ValueOf(value)

	if valueArray.Kind() != reflect.Array && valueArray.Kind() != reflect.Slice && valueArray.Kind() != reflect.String && valueArray.Kind() != reflect.Map {
		return nil, fmt.Errorf("Type isnt array, slice, map or string.")
	}

	if valueArray.Kind() == reflect.Map {
		valueArray.SetMapIndex(valueIndex, valueA)
		Data[name] = valueArray.Interface()
	} else {
		if valueArray.Len() <= int(valueIndex.Int()) {
			return nil, fmt.Errorf("Index out of range.")
		}
		valueArray.Index(int(valueIndex.Int())).Set(valueA)
		Data[name] = valueArray.Interface()
	}

	//UPDATE ON DOCUMENT:
	every := js.Global.Get("document").Call("querySelectorAll", "["+goobData+"=\""+name+"\"]["+goobForeach+"]")
	for i := 0; i < every.Length(); i++ {
		var everyindex *js.Object
		var ei *js.Object

		if valueArray.Kind() == reflect.Map {
			everyindex = every.Index(i).Call("querySelectorAll", "["+goobIndex+"=\""+fmt.Sprintf("%v", valueIndex.Interface())+"\"]")
			if everyindex.Length() > 0 {
				ei = everyindex.Index(0)
			}
		} else {
			everyindex = every.Index(i).Call("querySelectorAll", "["+goobIndex+"]")
			ei = everyindex.Index(int(valueIndex.Int()))
		}

		var buf bytes.Buffer
		var err error
		if foreachExpressionCheck(every.Index(i), valueA) {
			err = tmpl.ExecuteTemplate(&buf, every.Index(i).Call("getAttribute", goobBind).String(), value)
			if err != nil {
				println(err)
			}
		}

		wheretoadd := "beforeend"

		if ei != nil && ei != js.Undefined && jsbuiltin.TypeOf(ei) != "undefined" {
			//get a previous element as long it isn't a goob-index, so we will not need to create more childs like div or something..
			ei = ei.Get("previousSibling")
			for ei != nil && ei != js.Undefined && jsbuiltin.TypeOf(ei) != "undefined" && !(jsbuiltin.TypeOf(ei.Get("hasAttribute")) == "function" && ei.Call("hasAttribute", goobIndex).Bool()) {
				einew := ei.Get("previousSibling")
				ei.Call("remove") //TODO: check browser compatybility
				ei = einew
			}

			wheretoadd = "beforebegin"
		} else {
			if valueArray.Kind() == reflect.Map {
				buf.WriteString("<div style=\"display:none;\" " + goobIndex + "=\"" + fmt.Sprintf("%v", valueIndex.Interface()) + "\"></div>")
			} else {
				buf.WriteString("<div style=\"display:none;\" " + goobIndex + "></div>")
			}
		}

		if valueArray.Kind() == reflect.Map {
			if everyindex.Length() == 0 {
				// println(buf.String(), wheretoadd)
				every.Index(i).Call("insertAdjacentHTML", wheretoadd, buf.String())
			} else {
				for j := 0; j < everyindex.Length(); j++ {
					everyindex.Index(j).Call("insertAdjacentHTML", wheretoadd, buf.String())
				}
			}
		} else {
			everyindex.Index(int(valueIndex.Int())).Call("insertAdjacentHTML", wheretoadd, buf.String())
		}

	}
	every = js.Global.Get("document").Call("querySelectorAll", "["+goobData+"=\""+name+"\"]:not(["+goobForeach+"]), ["+goobWatch+"=\""+name+"\"]:not(["+goobForeach+"])")
	for i := 0; i < every.Length(); i++ {
		every.Index(i).Call("removeAttribute", goobProcessed)
	}
	Tick()

	return valueArray.Interface(), nil
}

//Push is pushing value at the end of slice or string
func Push(name string, value interface{}) (interface{}, error) {
	if _, ok := Data[name]; !ok {
		return nil, fmt.Errorf("You need to define slice first")
	}

	valueArray := reflect.ValueOf(Data[name])
	valueA := reflect.ValueOf(value)

	switch valueArray.Kind() {
	case reflect.Slice:
		Data[name] = reflect.Append(valueArray, valueA).Interface()
	case reflect.String:
		switch valueA.Kind() {
		case reflect.String:
			Data[name] = valueArray.String() + valueA.String()
		default:
			return nil, fmt.Errorf("Value needs to be a type of string.")
		}
	default:
		return nil, fmt.Errorf("Type isnt slice or string.")
	}

	//UPDATE ON DOCUMENT:
	every := js.Global.Get("document").Call("querySelectorAll", "["+goobData+"=\""+name+"\"]["+goobForeach+"]")
	for i := 0; i < every.Length(); i++ {
		var buf bytes.Buffer
		if foreachExpressionCheck(every.Index(i), valueA) {
			err := tmpl.ExecuteTemplate(&buf, every.Index(i).Call("getAttribute", goobBind).String(), value)
			if err != nil {
				println(err)
			}
		}

		buf.WriteString("<div style=\"display:none;\" " + goobIndex + "></div>")
		every.Index(i).Call("insertAdjacentHTML", "beforeend", buf.String())
	}
	every = js.Global.Get("document").Call("querySelectorAll", "["+goobData+"=\""+name+"\"]:not(["+goobForeach+"]), ["+goobWatch+"=\""+name+"\"]:not(["+goobForeach+"])")
	for i := 0; i < every.Length(); i++ {
		every.Index(i).Call("removeAttribute", goobProcessed)
	}
	Tick()

	return Data[name], nil
}

// TODO
// //Join appends a slice to a slice value. The slices must have the same element type.
// func Join(name string, value interface{}) (interface{}, error) {
// 	if _, ok := data[name]; !ok {
// 		return nil, fmt.Errorf("You need to define slice first")
// 	}
//
// 	valueArray := reflect.ValueOf(data[name])
// 	valueA := reflect.ValueOf(value)
//
// 	if valueArray.Kind() == reflect.Slice && valueA.Kind() == reflect.Slice {
// 		data[name] = reflect.AppendSlice(valueArray, valueA).Interface()
// 	} else {
// 		return nil, fmt.Errorf("Both types needs to be a slice.")
// 	}
//
// 	return data[name], nil
// }

func InsertAt(name string, index int, value interface{}) (interface{}, error) {
	if _, ok := Data[name]; !ok {
		return nil, fmt.Errorf("You need to define slice first")
	}

	valueArray := reflect.ValueOf(Data[name])
	valueA := reflect.ValueOf(value)

	switch valueArray.Kind() {
	case reflect.Slice:
		valueArray = reflect.Append(valueArray, valueA)
		reflect.Copy(valueArray.Slice(index+1, valueArray.Len()), valueArray.Slice(index, valueArray.Len()-1))
		valueArray.Index(index).Set(valueA)
		Data[name] = valueArray.Interface()
	case reflect.String:
		switch valueA.Kind() {
		case reflect.String:
			Data[name] = valueArray.Slice(0, index).String() + valueA.String() + valueArray.Slice(index, valueArray.Len()).String()
		default:
			return nil, fmt.Errorf("Value needs to be a type of string.")
		}
	default:
		return nil, fmt.Errorf("Type isnt slice or string.")
	}

	//UPDATE ON DOCUMENT:
	every := js.Global.Get("document").Call("querySelectorAll", "["+goobData+"=\""+name+"\"]["+goobForeach+"]")
	for i := 0; i < every.Length(); i++ {
		var ei *js.Object
		wheretoadd := "afterend"

		if index == 0 {
			ei = every.Index(i)
			wheretoadd = "afterbegin"
		} else {
			ei = every.Index(i).Call("querySelectorAll", "["+goobIndex+"]").Index(index - 1)
		}

		var buf bytes.Buffer
		if foreachExpressionCheck(every.Index(i), valueA) {
			err := tmpl.ExecuteTemplate(&buf, every.Index(i).Call("getAttribute", goobBind).String(), value)
			if err != nil {
				println(err)
			}
		}

		buf.WriteString("<div style=\"display:none;\" " + goobIndex + "></div>")
		ei.Call("insertAdjacentHTML", wheretoadd, buf.String())
	}
	every = js.Global.Get("document").Call("querySelectorAll", "["+goobData+"=\""+name+"\"]:not(["+goobForeach+"]), ["+goobWatch+"=\""+name+"\"]:not(["+goobForeach+"])")
	for i := 0; i < every.Length(); i++ {
		every.Index(i).Call("removeAttribute", goobProcessed)
	}
	Tick()

	return Data[name], nil
}

//TODO: map
func RemoveAt(name string, index int) (interface{}, error) {
	if _, ok := Data[name]; !ok {
		return nil, fmt.Errorf("You need to define slice first")
	}

	valueArray := reflect.ValueOf(Data[name])

	switch valueArray.Kind() {
	case reflect.Slice:
		valueArray = reflect.AppendSlice(valueArray.Slice(0, index), valueArray.Slice(index+1, valueArray.Len()))
		Data[name] = valueArray.Interface()
	default:
		return nil, fmt.Errorf("Type isnt slice.")
	}

	//UPDATE ON DOCUMENT:
	every := js.Global.Get("document").Call("querySelectorAll", "["+goobData+"=\""+name+"\"]["+goobForeach+"]")
	for i := 0; i < every.Length(); i++ {
		ei := every.Index(i).Call("querySelectorAll", "["+goobIndex+"]").Index(index)

		if ei != nil && ei != js.Undefined && jsbuiltin.TypeOf(ei) != "undefined" {
			//get a previous element as long it isn't a goob-index, so we will not need to create more childs like div or something..
			// ei = ei.Get("previousSibling")
			einew := ei.Get("previousSibling")
			ei.Call("remove") //TODO: check browser compatybility
			ei = einew
			for ei != nil && ei != js.Undefined && jsbuiltin.TypeOf(ei) != "undefined" && !(jsbuiltin.TypeOf(ei.Get("hasAttribute")) == "function" && ei.Call("hasAttribute", goobIndex).Bool()) {
				einew = ei.Get("previousSibling")
				ei.Call("remove") //TODO: check browser compatybility
				ei = einew
			}
		}
	}
	every = js.Global.Get("document").Call("querySelectorAll", "["+goobData+"=\""+name+"\"]:not(["+goobForeach+"]), ["+goobWatch+"=\""+name+"\"]:not(["+goobForeach+"])")
	for i := 0; i < every.Length(); i++ {
		every.Index(i).Call("removeAttribute", goobProcessed)
	}
	Tick()

	return Data[name], nil
}

func Pop(name string) (interface{}, error) {
	if _, ok := Data[name]; !ok {
		return nil, fmt.Errorf("You need to define slice first")
	}
	//TODO: add outputing last element

	valueArray := reflect.ValueOf(Data[name])

	_, err := RemoveAt(name, valueArray.Len()-1)
	return nil, err
}

//TODO: map
func Empty(name string) (interface{}, error) {
	if _, ok := Data[name]; !ok {
		return nil, fmt.Errorf("You need to define slice first")
	}

	valueArray := reflect.ValueOf(Data[name])

	switch valueArray.Kind() {
	case reflect.Slice:
		Data[name] = valueArray.Slice(0, 0).Interface()
	default:
		return nil, fmt.Errorf("Type isnt slice.")
	}

	//UPDATE ON DOCUMENT:
	every := js.Global.Get("document").Call("querySelectorAll", "["+goobData+"=\""+name+"\"], ["+goobWatch+"=\""+name+"\"]")
	for i := 0; i < every.Length(); i++ {
		every.Index(i).Call("removeAttribute", goobProcessed)
	}
	Tick()

	return Data[name], nil
}
