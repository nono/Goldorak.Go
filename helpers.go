package goldorak

import (
	"log"
	"os"
)

// TODO http://golang.org/src/pkg/template/format.go#L38

func StaticUrl(filename string) string {
	file, err := os.Open(filename, os.O_RDONLY, 0444)
	if err != nil {
		log.Stderrf("Impossible to open %s", filename, err)
	}
	dir, err := file.Stat()
	if err != nil {
		log.Stderrf("Impossible to stat %s", filename, err)
	}
	secs := dir.Mtime_ns / 1e9
	return filename + "?" + string(secs)
}
