package handlers

import (
	"fmt"
	"net/http"
)

func httpErrf(w http.ResponseWriter, code int, fmtMsg string, args ...interface{}) {
	msg := fmt.Sprintf(fmtMsg, args...)
	http.Error(w, msg, code)
}
