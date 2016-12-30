package goob

import "github.com/gopherjs/gopherjs/js"

func defineVariables() {
	Set("WindowLocation", WindowLocation{
		Host:     js.Global.Get("window").Get("location").Get("host").String(),
		HostName: js.Global.Get("window").Get("location").Get("hostname").String(),
		Href:     js.Global.Get("window").Get("location").Get("href").String(),
		Origin:   js.Global.Get("window").Get("location").Get("origin").String(),
		PathName: js.Global.Get("window").Get("location").Get("pathname").String(),
		Port:     js.Global.Get("window").Get("location").Get("port").String(),
		Protocol: js.Global.Get("window").Get("location").Get("protocol").String(),
		Search:   js.Global.Get("window").Get("location").Get("search").String(),
	})
}

type WindowLocation struct {
	Host,
	HostName,
	Href,
	Origin,
	PathName,
	Port,
	Protocol,
	Search string
}
