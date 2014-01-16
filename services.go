package main

import (
	"net/http"
)

func servicesHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("services"))
}
