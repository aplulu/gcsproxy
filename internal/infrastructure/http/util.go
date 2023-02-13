package http

import (
	"net/http"
	"strconv"
)

func writeStringHeader(w http.ResponseWriter, name string, value string) {
	if len(value) > 0 {
		w.Header().Set(name, value)
	}
}

func writeInt64Header(w http.ResponseWriter, name string, value int64) {
	if value > 0 {
		w.Header().Set(name, strconv.FormatInt(value, 10))
	}
}
