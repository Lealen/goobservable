package goob

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jsbuiltin"
)

var data = make(map[string]interface{})

func Get(name string) interface{} {
	return data[name]
}

func Set(name string, value interface{}) {
	data[name] = value

	//UPDATE ON DOCUMENT:
	every := js.Global.Get("document").Call("querySelectorAll", "["+goobData+"=\""+name+"\"], ["+goobWatch+"=\""+name+"\"]")
	for i := 0; i < every.Length(); i++ {
		every.Index(i).Call("removeAttribute", goobProcessed)
	}
	Tick()
}

func GetIndex(name string, index interface{}) interface{} {
	if _, ok := data[name]; !ok {
		return nil
	}

	valueArray := reflect.ValueOf(data[name])
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
	if _, ok := data[name]; !ok {
		return nil, fmt.Errorf("You need to define string/array/slice/map first")
	}

	valueArray := reflect.ValueOf(data[name])
	valueIndex := reflect.ValueOf(index)
	valueA := reflect.ValueOf(value)

	if valueArray.Kind() != reflect.Array && valueArray.Kind() != reflect.Slice && valueArray.Kind() != reflect.String && valueArray.Kind() != reflect.Map {
		return nil, fmt.Errorf("Type isnt array, slice, map or string.")
	}

	if valueArray.Kind() == reflect.Map {
		valueArray.SetMapIndex(valueIndex, valueA)
		data[name] = valueArray.Interface()
	} else {
		if valueArray.Len() <= int(valueIndex.Int()) {
			return nil, fmt.Errorf("Index out of range.")
		}
		valueArray.Index(int(valueIndex.Int())).Set(valueA)
		data[name] = valueArray.Interface()
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
		err := tmpl.ExecuteTemplate(&buf, every.Index(i).Call("getAttribute", goobBind).String(), value)
		if err != nil {
			println(err)
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
	if _, ok := data[name]; !ok {
		return nil, fmt.Errorf("You need to define slice first")
	}

	valueArray := reflect.ValueOf(data[name])
	valueA := reflect.ValueOf(value)

	switch valueArray.Kind() {
	case reflect.Slice:
		data[name] = reflect.Append(valueArray, valueA).Interface()
	case reflect.String:
		switch valueA.Kind() {
		case reflect.String:
			data[name] = valueArray.String() + valueA.String()
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
		err := tmpl.ExecuteTemplate(&buf, every.Index(i).Call("getAttribute", goobBind).String(), value)
		if err != nil {
			println(err)
		}

		buf.WriteString("<div style=\"display:none;\" " + goobIndex + "></div>")
		every.Index(i).Call("insertAdjacentHTML", "beforeend", buf.String())
	}
	every = js.Global.Get("document").Call("querySelectorAll", "["+goobData+"=\""+name+"\"]:not(["+goobForeach+"]), ["+goobWatch+"=\""+name+"\"]:not(["+goobForeach+"])")
	for i := 0; i < every.Length(); i++ {
		every.Index(i).Call("removeAttribute", goobProcessed)
	}
	Tick()

	return data[name], nil
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
	if _, ok := data[name]; !ok {
		return nil, fmt.Errorf("You need to define slice first")
	}

	valueArray := reflect.ValueOf(data[name])
	valueA := reflect.ValueOf(value)

	switch valueArray.Kind() {
	case reflect.Slice:
		valueArray = reflect.Append(valueArray, valueA)
		reflect.Copy(valueArray.Slice(index+1, valueArray.Len()), valueArray.Slice(index, valueArray.Len()-1))
		valueArray.Index(index).Set(valueA)
		data[name] = valueArray.Interface()
	case reflect.String:
		switch valueA.Kind() {
		case reflect.String:
			data[name] = valueArray.Slice(0, index).String() + valueA.String() + valueArray.Slice(index, valueArray.Len()).String()
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
		err := tmpl.ExecuteTemplate(&buf, every.Index(i).Call("getAttribute", goobBind).String(), value)
		if err != nil {
			println(err)
		}

		buf.WriteString("<div style=\"display:none;\" " + goobIndex + "></div>")
		ei.Call("insertAdjacentHTML", wheretoadd, buf.String())
	}
	every = js.Global.Get("document").Call("querySelectorAll", "["+goobData+"=\""+name+"\"]:not(["+goobForeach+"]), ["+goobWatch+"=\""+name+"\"]:not(["+goobForeach+"])")
	for i := 0; i < every.Length(); i++ {
		every.Index(i).Call("removeAttribute", goobProcessed)
	}
	Tick()

	return data[name], nil
}

//TODO: map
func RemoveAt(name string, index int) (interface{}, error) {
	if _, ok := data[name]; !ok {
		return nil, fmt.Errorf("You need to define slice first")
	}

	valueArray := reflect.ValueOf(data[name])

	switch valueArray.Kind() {
	case reflect.Slice:
		valueArray = reflect.AppendSlice(valueArray.Slice(0, index), valueArray.Slice(index+1, valueArray.Len()))
		data[name] = valueArray.Interface()
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

	return data[name], nil
}

//TODO: map
func Empty(name string) (interface{}, error) {
	if _, ok := data[name]; !ok {
		return nil, fmt.Errorf("You need to define slice first")
	}

	valueArray := reflect.ValueOf(data[name])

	switch valueArray.Kind() {
	case reflect.Slice:
		data[name] = valueArray.Slice(0, 0).Interface()
	default:
		return nil, fmt.Errorf("Type isnt slice.")
	}

	//UPDATE ON DOCUMENT:
	every := js.Global.Get("document").Call("querySelectorAll", "["+goobData+"=\""+name+"\"], ["+goobWatch+"=\""+name+"\"]")
	for i := 0; i < every.Length(); i++ {
		every.Index(i).Call("removeAttribute", goobProcessed)
	}
	Tick()

	return data[name], nil
}
