package main

import "net/http"

// handler to show an ok status when /healthz is accessed
func handlerOkStatus(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	message := "OK"

	w.Write([]byte(message))
}
