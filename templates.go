package goob

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"strings"

	"github.com/gopherjs/gopherjs/js"
)

var tmpl = template.Must(template.New("").Funcs(funcMap), nil)

// var data interface{}

func AddTemplate(name, src string) {
	var err error
	tmpl, err = tmpl.New(name).Parse(src)
	if err != nil {
		println(err)
		return
	}
}

// func BindVariables(v interface{}) {
// 	data = v
// }

func Run() {
	// data = v

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
				if dat, ok := data[dataname]; ok {
					//if goob-foreach is defined, then rather than try to pass a data to template, we need to pass each element to different templates and renders it inside
					if every.Index(i).Call("hasAttribute", goobForeach).Bool() {
						valueArray := reflect.ValueOf(data[dataname])

						switch valueArray.Kind() {
						case reflect.Array, reflect.Slice, reflect.String:
							for j := 0; j < valueArray.Len(); j++ {
								err = tmpl.ExecuteTemplate(&buf, templatename, valueArray.Index(j).Interface())
								if err != nil {
									println(err)
								}
								buf.WriteString("<div style=\"display:none;\" goob-index></div>")
							}
						case reflect.Map:
							allkeys := valueArray.MapKeys()
							for _, v := range allkeys {
								err = tmpl.ExecuteTemplate(&buf, templatename, valueArray.MapIndex(v).Interface())
								if err != nil {
									println(err)
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
				err = tmpl.ExecuteTemplate(&buf, templatename, data)
			}

			if err != nil {
				println(err)
				continue
			}

			every.Index(i).Set("innerHTML", buf.String())
			every.Index(i).Call("setAttribute", goobProcessed, "")
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
	err := tmpl.ExecuteTemplate(&buf, name, data)
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
