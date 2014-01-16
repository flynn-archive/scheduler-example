package main

import (
	"net/http"
)

func ptyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pty"))
}
