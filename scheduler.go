package main

import (
	"log"
	"net/http"
	"os"

	services "github.com/flynn/go-discoverd"
)

func error400(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return true
	}
	return false
}

func error500(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return true
	}
	return false
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}

	ip := os.Getenv("EXTERNAL_IP")
	if ip == "" {
		ip = "127.0.0.1"
	}
	err := services.Connect(ip + ":1111") // TODO: fix this
	if err != nil {
		panic(err)
	}
	services.Register("simple-scheduler", port)

	http.HandleFunc("/batch/", batchHandler)
	http.HandleFunc("/services/", servicesHandler)
	http.HandleFunc("/pty/", ptyHandler)

	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
