package goob

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"strings"

	"github.com/gopherjs/gopherjs/js"
)

var tmpl = template.New("").Funcs(funcMap)

// var data interface{}
var templatestmp = map[string]string{}
var isrunning bool

func AddTemplate(name, src string) {
	if !isrunning {
		templatestmp[name] = src
	} else {
		var err error
		tmpl, err = tmpl.New(name).Parse(src)
		if err != nil {
			println(err)
		}
	}
}

// func BindVariables(v interface{}) {
// 	data = v
// }

func Run() {
	// data = v
	var err error
	for k, v := range templatestmp {
		tmpl, err = tmpl.New(k).Parse(v)
		if err != nil {
			println(err)
		}
	}
	isrunning = true
	templatestmp = nil

	defineVariables()
	runTemplate("main")
	Tick()
}

// Tick renders all not processed templates in html and put them in a source code
func Tick() {
	for {
		every := js.Global.Get("document").Call("querySelectorAll", "["+goobBind+"]:not(["+goobProcessed+"])")
		if every.Length() == 0 {
			break
		}
		for i := 0; i < every.Length(); i++ {
			templatename := every.Index(i).Call("getAttribute", goobBind).String()

			var buf bytes.Buffer
			var err error

			// if there is a goob-data defined: pass a data related to it to template defined in goob-bind: //TODO: maybe try to understand if someone want to pass a defined value in array/slice/map or smth?
			if every.Index(i).Call("hasAttribute", goobData).Bool() {
				dataname := every.Index(i).Call("getAttribute", goobData).String()
				if dat, ok := Data[dataname]; ok {
					//if goob-foreach is defined, then rather than try to pass a data to template, we need to pass each element to different templates and renders it inside
					if every.Index(i).Call("hasAttribute", goobForeach).Bool() {
						valueArray := reflect.ValueOf(Data[dataname])

						switch valueArray.Kind() {
						case reflect.Array, reflect.Slice, reflect.String:
							for j := 0; j < valueArray.Len(); j++ {
								if foreachExpressionCheck(every.Index(i), valueArray.Index(j)) {
									err = tmpl.ExecuteTemplate(&buf, templatename, valueArray.Index(j).Interface())
									if err != nil {
										println(err)
									}
								}
								buf.WriteString("<div style=\"display:none;\" goob-index></div>")
							}
						case reflect.Map:
							allkeys := valueArray.MapKeys()
							for _, v := range allkeys {
								if foreachExpressionCheck(every.Index(i), valueArray.MapIndex(v)) {
									err = tmpl.ExecuteTemplate(&buf, templatename, valueArray.MapIndex(v).Interface())
									if err != nil {
										println(err)
									}
								}
								buf.WriteString("<div style=\"display:none;\" " + goobIndex + "=\"" + fmt.Sprintf("%v", v) + "\"></div>")
							}
						default:
							err = fmt.Errorf("Data used for foreach is neither array/slice/map nor string")
						}
					} else {
						err = tmpl.ExecuteTemplate(&buf, templatename, dat)
					}
				} else {
					err = fmt.Errorf("Could not found data %s to pass to template %s!", every.Index(i).Call("getAttribute", goobData).String(), every.Index(i).Call("getAttribute", goobBind).String())
				}
			} else {
				err = tmpl.ExecuteTemplate(&buf, templatename, Data)
			}

			if err != nil {
				println(err)
				// continue
			}

			every.Index(i).Set("innerHTML", buf.String())
			every.Index(i).Call("setAttribute", goobProcessed, "")
		}
	}

	for {
		every := js.Global.Get("document").Call("querySelectorAll", "["+goobIfVar+"]:not(["+goobIfExpressionProcessed+"])")
		if every.Length() == 0 {
			break
		}
		for i := 0; i < every.Length(); i++ {
			variablename := every.Index(i).Call("getAttribute", goobIfVar).String()
			checkOneExpression(reflect.ValueOf(Data[variablename]), every.Index(i))
		}
	}
}

func UpdateTemplate(name string) {
	runTemplate(name)
	Tick()
}

func Debug() {
	println(strings.Replace(fmt.Sprintf("%#v", tmpl.Tree), ",", ",\n", -1))
}

func runTemplate(name string) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, name, Data)
	if err != nil {
		println(err)
		return
	}

	if name == "main" {
		js.Global.Get("document").Call("getElementById", "app").Set("innerHTML", buf.String())
	} else {
		every := js.Global.Get("document").Call("querySelectorAll", "[goob-bind=\""+name+"\"]")
		for i := 0; i < every.Length(); i++ {
			every.Index(i).Set("innerHTML", buf.String())
			every.Index(i).Call("setAttribute", goobProcessed, "")
		}
	}
}

func updateValueTemplates(name string) {

}
