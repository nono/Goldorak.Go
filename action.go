package goldorak

import (
	"gostache"
	"log"
	"os"
	"web"
)

type Action struct {
	ctx    *web.Context
	layout string
	locals map[string]string
}

// TODO func Restful() ?
// TODO what about POST/PUT/DELETE?
func Get(route string, handler func(Action)) {
	web.Get(route, func (ctx *web.Context) {
		action := Action{ctx, "layout", map[string]string {}}
		// TODO args.Insert(action)
		handler(action)
	});
}

func (this *Action) Assign(key string, value string) {
	this.locals[key] = value
}

func (this *Action) Layout(template string) {
	this.layout = template
}

func (this *Action) Render(template string) {
	filename:= template
	output, err := this.RenderFile(filename, this.locals)
	if err != nil {
		log.Stderrf("Error on rendering %s", filename, err)
	}
	if this.layout != "" {
		locals := map[string]string {"yield": output}
		output, err = this.RenderFile(this.layout, locals)
	}
	this.ctx.WriteString(output)
}

func (this *Action) RenderFile(filename string, context interface{}) (string, os.Error) {
	file := GetConfig("templates") + "/" + filename + ".mustache"
	return gostache.RenderFile(file, context)
}

