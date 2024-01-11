package log

import (
	"fmt"
	"time"
)

func Info(str string) {
	iso8601 := time.Now().UTC().Format(time.RFC3339)
	fmt.Printf("%s INFO - %s\n", iso8601, str)
}

func Infof(format string, a ...any) {
	iso8601 := time.Now().UTC().Format(time.RFC3339)
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s INFO - %s", iso8601, msg)
}
