package log

import "fmt"

func Printf(fmtStr string, args ...interface{}) {
	s := fmt.Sprintf(fmtStr, args...)
	fmt.Printf("%s\n", s)
}

func Errf(fmtStr string, args ...interface{}) {
	s := fmt.Sprintf(fmtStr, args...)
	fmt.Printf("[ERROR] %s\n", s)
}
