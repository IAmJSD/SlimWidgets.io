package main

import "io/ioutil"

var (
	WidgetTemplate string
	InviteTemplate string
)

func FSInit() {
	w, err := ioutil.ReadFile("./templates/widget.template.html")
	if err != nil {
		panic(err)
	}
	WidgetTemplate = string(w)

	w, err = ioutil.ReadFile("./templates/invite.template.html")
	if err != nil {
		panic(err)
	}
	InviteTemplate = string(w)
}
