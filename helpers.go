package cleaner

import (
	"fmt"
	"reflect"
	"runtime"
)

func whisper(msg string) {
	fmt.Println("\033[32m[closer]\033[0m: " + msg)
}

func functionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
