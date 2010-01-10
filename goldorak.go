package goldorak

import (
	"rand"
	"time"
	"web"
)

func Initialize(filename string) {
	ReadConfig(filename)
	rand.Seed(time.Nanoseconds())
	web.SetStaticDir(GetConfig("docroot"))
}

func Start() {
	addr := GetConfig("interface") + ":" + GetConfig("port")
	web.Run(addr)
}

