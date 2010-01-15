package goldorak

import (
	"fmt"
	"log"
	"os"
	"reflect"
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

// Pluralize(1, 'piano', 'pianos') => "1 piano"
// Pluralize(2, 'piano') => "2 pianos"
func Pluralize(count int, a ...) string {
	var singular, plural string
	v := reflect.NewValue(a).(*reflect.StructValue)
	singular = v.Field(0).(*reflect.StringValue).Get()
	if v.NumField() > 1 {
		plural = v.Field(1).(*reflect.StringValue).Get()
	} else {
		plural = singular + "s"
	}
	if count <= 1 {
		return string(count) + " " + singular
	}
	return string(count) + " " + plural
}

