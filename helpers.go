package goldorak

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
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

func Parameterize(s string) string {
	runes := strings.Runes(s)
	t := make([]int, len(runes))
	for i := 0; i < len(runes); i++ {
		rune := runes[i]
		if (unicode.IsDigit(rune) || unicode.IsLetter(rune)) {
			t[i] = rune
		} else {
			t[i] = '-'
		}
	}
	return string(t)
}

