package goldorak

import (
	"mustache"
	"log"
	"strings"
	"web"
)

type Action struct {
	responded bool
	layout    *Action
	template  string
	locals    map[string]string
	context   *web.Context
}

var defaultLayout *Action = nil

func ctxHandler(handler func(*Action, []string)) (func(ctx *web.Context)) {
	f := func(ctx *web.Context) {
		action := NewAction()
		action.layout = defaultLayout
		action.context = ctx
		params := strings.Split(ctx.Request.URL.Path[1:], "/", 0)
		handler(&action, params)
		if !action.responded {
			action.responded = true
			ctx.StartResponse(200)
			ctx.WriteString(action.Render())
		}
	}
	return f
}

// TODO not found page?
// TODO func Restful() ?
// Note for myself: the last route is the most important one (it's the opposite of Rails)
func Get(route string, handler func(*Action, []string)) {
	web.Get(route, ctxHandler(handler))
}

func Post(route string, handler func(*Action, []string)) {
	web.Post(route, ctxHandler(handler))
}

func DefaultLayout(handler func(*Action)) {
	action := NewAction()
	handler(&action)
	defaultLayout = &action
}

func NewAction() Action {
	return Action{false, nil, "", make(map[string]string), nil}
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
	output, err := mustache.RenderFile(filename, this.locals)
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

func (this *Action) Redirect(path string) {
	this.responded = true
	url := "http://" + GetConfig("domain") + path
	this.context.Redirect(302, url)
}

