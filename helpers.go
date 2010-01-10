package goldorak

import (
	"fmt"
	"log"
	"os"
)

// TODO http://golang.org/src/pkg/template/format.go#L38

func StaticUrl(filename string) string {
	file, err := os.Open(GetConfig("docroot") + "/" + filename, os.O_RDONLY, 0444)
	if err != nil {
		log.Exitf("Impossible to open %s", filename, err)
	}
	dir, err := file.Stat()
	if err != nil {
		log.Exitf("Impossible to stat %s", filename, err)
	}
	secs := dir.Mtime_ns / 1e9
	return fmt.Sprintf("/%s?%d", filename, secs)
}
