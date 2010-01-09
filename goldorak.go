package goldorak

import (
	"web"
)

func Initialize(filename string) {
	ReadConfig(filename)
	web.SetStaticDir(GetConfig("docroot"))
}

func Start() {
	addr := GetConfig("interface") + ":" + GetConfig("port")
	web.Run(addr)
}

