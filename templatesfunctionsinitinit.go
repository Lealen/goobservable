package goob

//copyright Mateusz Tomasz Doroszko for project s3n

import (
	"html/template"
	"strings"
)

var funcMap = template.FuncMap{
	//package specified:
	"global": getGlobalFunc,

	//Generic:
	"gettype":        getTypeFunc, //in: interface{} ; out: string
	"log":            logFunc,     //in: string
	"howmanylikefor": howManyLikeForFunc,

	//Memory banks:
	"newmb": newMemoryBankFunc,
	"setmb": setToMemoryBankFunc,
	"getmb": getFromMemoryBankFunc,

	//Conversions:
	"deref": derefFunc,

	"donotparse":    doNotParseFunc,    //in: string ; out: template.HTML
	"donotparsejs":  doNotParseJSFunc,  //in: string ; out: template.JS
	"donotparsecss": doNotParseCSSFunc, //in: string ; out: template.CSS
	"donotparseurl": doNotParseURLFunc, //in: string ; out: template.URL

	"bool":    boolFunc,
	"int":     intFunc,
	"uint":    uintFunc,
	"int64":   int64Func,
	"uint64":  uint64Func,
	"float32": float32Func,
	"float64": float64Func,
	"string":  stringFunc,

	"tojson": toJSONFunc,

	//Arrays:
	"setindex": setIndexFunc, // (array index something) interface{}
	"getindex": getIndexFunc, // (array index) interface{}

	//Math:
	"add":      addFunc,      //[int|uint|float|string]
	"subtract": subtractFunc, //[int|uint|float]
	"multiply": multiplyFunc, //[int]uint|float|string+int|string+uint]
	"divide":   divideFunc,   //[int|uint|float]

	"ceil":  ceilFunc,  //[int|uint|float]
	"floor": floorFunc, //[int|uint|float]
	"mod":   modFunc,   //[int|uint|float]
	"min":   minFunc,   //[int|uint|float]
	"max":   maxFunc,   //[int|uint|float]

	//Strings:
	"contains":     strings.Contains,
	"containsany":  strings.ContainsAny,
	"count":        strings.Count,
	"equalfold":    strings.EqualFold,
	"fields":       strings.Fields,
	"hasprefix":    strings.HasPrefix,
	"hassuffix":    strings.HasSuffix,
	"index":        strings.Index,
	"indexany":     strings.IndexAny,
	"indexbyte":    strings.IndexByte,
	"join":         strings.Join,
	"lastindex":    strings.LastIndex,
	"lastindexany": strings.LastIndexAny,
	"repeat":       strings.Repeat,
	"replace":      strings.Replace,
	"split":        strings.Split,
	"splitafter":   strings.SplitAfter,
	"splitaftern":  strings.SplitAfterN,
	"splitn":       strings.SplitN,
	"title":        strings.Title,
	"tolower":      strings.ToLower,
	"totitle":      strings.ToTitle,
	"toupper":      strings.ToUpper,
	"trim":         strings.Trim,
	"trimleft":     strings.TrimLeft,
	"trimprefix":   strings.TrimPrefix,
	"trimright":    strings.TrimRight,
	"trimspace":    strings.TrimSpace,
	"trimsuffix":   strings.TrimSuffix,

	//Date:
	"timenow":           timenowFunc,
	"timeutc":           timeutcFunc,
	"timeparseduration": timeparsedurationFunc,
	"timeadd":           timeaddFunc,
	"timeformat":        timeformatFunc,
	"timeyear":          timeyearFunc,
	"timemonth":         timemonthFunc,
	"timeday":           timedayFunc,
	"timehour":          timehourFunc,
	"timeminute":        timeminuteFunc,
	"timesecond":        timesecondFunc,
	"timeparseunix":     timeparseunixFunc,
}

func RegisterTemplateFunction(name string, function interface{}) {
	funcMap[name] = function
	tmpl = tmpl.Funcs(funcMap)
	// funcMap[name] = function
}
