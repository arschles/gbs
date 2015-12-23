package handlers

import (
	"bufio"
	"io"
	"net/http"
)

func buildStatusURL(containerID string) string {
	return "/status/" + containerID
}

// httpFlushWriter calls Flush on w if it is an http.Flusher
func httpFlushWriter(w io.Writer) {
	fl, ok := w.(http.Flusher)
	if ok {
		fl.Flush()
	}
}

// streamReader uses a bufio.Scanner to read from r and write to w. After each successful call
// to scanner.Scan(), calls afterEach and passes in the same writer as passed to this func
func streamReader(r io.Reader, w io.Writer, afterEach func(w io.Writer)) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		txt := scanner.Text()
		w.Write([]byte(txt))
		if err := scanner.Err(); err != nil {
			continue
		}
		afterEach(w)
	}
}
