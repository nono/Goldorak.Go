package goldorak

import (
	"gostache"
	"log"
	"web"
)

type Action struct {
	layout    *Action
	template  string
	locals    map[string]string
}

var defaultLayout *Action = nil

// TODO not found page?
// TODO func Restful() ?
// TODO what about POST/PUT/DELETE?
func Get(route string, handler func(Action)) {
	web.Get(route, func (ctx *web.Context) {
		action := Action{defaultLayout, "", map[string]string {}}
		// TODO args.Insert(action)
		handler(action)
		ctx.WriteString(action.Render())
	});
}

func DefaultLayout(handler func(Action)) {
	action := Action{nil, "", map[string]string {}}
	handler(action)
	defaultLayout = &action
}

func (this *Action) Template(template string) {
	this.template = template
}

func (this *Action) Assign(key string, value string) {
	this.locals[key] = value
}

func (this *Action) Layout(handler func(Action)) {
	action := Action{nil, "", map[string]string {}}
	handler(action)
	this.layout = &action
}

func (this *Action) NoLayout() {
	this.layout = nil
}

func (this *Action) Render() string {
	filename := GetConfig("templates") + "/" + this.template + ".mustache"
	output, err := gostache.RenderFile(filename, this.locals)
	if err != nil {
		log.Stderrf("Error on rendering %s", filename, err)
		return "" // TODO error page
	}
	if this.layout != nil {
		this.layout.locals["yield"] = output
		output = this.layout.Render()
		this.layout.locals["yield"] = "", false
	}
	return output
}

