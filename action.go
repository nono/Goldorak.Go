package goldorak

import (
	"gostache"
	"log"
	"strings"
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
func Get(route string, handler func(*Action, []string)) {
	web.Get(route, func (ctx *web.Context) {
		action := NewAction()
		action.layout = defaultLayout
		params := strings.Split(ctx.Request.URL.Path, "/", 0)
		handler(&action, params[1:])
		ctx.StartResponse(200)
		ctx.WriteString(action.Render())
	});
}

func DefaultLayout(handler func(*Action)) {
	action := NewAction()
	handler(&action)
	defaultLayout = &action
}

func NewAction() Action {
	return Action{nil, "", make(map[string]string)}
}

func (this *Action) Template(template string) {
	this.template = template
}

func (this *Action) Assign(key string, value string) {
	this.locals[key] = value
}

func (this *Action) Layout(handler func(*Action)) {
	action := NewAction()
	handler(&action)
	this.layout = &action
}

func (this *Action) NoLayout() {
	this.layout = nil
}

func (this *Action) Render() string {
	log.Stdoutf("Rendering %s", this.template)
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

