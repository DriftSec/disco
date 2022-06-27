package core

import (
	"fmt"
	"runtime"
	"strconv"
)

const (
	Reset        = "\033[0m"
	InfoColor    = "\033[1;34m"
	NoticeColor  = "\033[1;36m"
	WarningColor = "\033[1;33m"
	ErrorColor   = "\033[1;31m"
	DebugColor   = "\033[0;36m"
	SuccessColor = "\033[1;32m"
)

var DebugOn bool

// Dprint is a wrapper around fmt.Println() that prepends [DEBUG] and colorizes. Only prints if DebugOn is true
func Dprint(a ...interface{}) {
	if DebugOn {
		aa := make([]interface{}, 0, 2+len(a))
		aa = append(aa, NoticeColor)
		aa = append(append(aa, "[DEBUG]"), a...)
		aa = append(aa, Reset)
		// fmt.Println()
		fmt.Println(aa...)
	}
}

// Dprint is a wrapper around fmt.Println for printing debug messages
// func Dprint(a ...interface{}) {
// 	if DebugOn {
// 		fpcs := make([]uintptr, 1)
// 		var cname string
// 		var cf string
// 		var cline int
// 		_ = runtime.Callers(2, fpcs)

// 		caller := runtime.FuncForPC(fpcs[0] - 1)
// 		if caller == nil {
// 			cname = "MSG CALLER WAS NIL"
// 			cf = ""
// 			cline = -1
// 		} else {
// 			cname = caller.Name()
// 			cf, cline = caller.FileLine(fpcs[0] - 1)
// 		}

// 		aa := make([]interface{}, 0, 2+len(a))
// 		aa = append(aa, DebugColor)
// 		aa = append(append(aa, "("+cf+":"+strconv.Itoa(cline)+") "+cname+":"), a...)
// 		aa = append(aa, Reset)
// 		fmt.Println(aa...)
// 	}

// }

// Eprint is a wrapper around fmt.Println for printing errors
func Eprint(a ...interface{}) {
	if DebugOn {
		fpcs := make([]uintptr, 1)
		var cname string
		var cf string
		var cline int
		_ = runtime.Callers(2, fpcs)

		caller := runtime.FuncForPC(fpcs[0] - 1)
		if caller == nil {
			cname = "MSG CALLER WAS NIL"
			cf = ""
			cline = -1
		} else {
			cname = caller.Name()
			cf, cline = caller.FileLine(fpcs[0] - 1)
		}

		aa := make([]interface{}, 0, 2+len(a))
		aa = append(aa, ErrorColor)
		aa = append(append(aa, "("+cf+":"+strconv.Itoa(cline)+") "+cname+":"), a...)
		aa = append(aa, Reset)
		fmt.Println(aa...)
	}

}
