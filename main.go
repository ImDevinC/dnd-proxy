package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const defaultPort = 8080

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api", handler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("[ERROR] invalid request")
		log.Println(r.URL.String())
		w.WriteHeader(http.StatusNotFound)
	})
	rawPort := os.Getenv("DND_PORT")
	port, err := strconv.ParseInt(rawPort, 10, 64)
	if err != nil {
		port = defaultPort
	}
	log.Println("[INFO] listening on port ", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func handler(w http.ResponseWriter, r *http.Request) {
	character := r.URL.Query().Get("character")
	if len(character) == 0 {
		log.Println("[ERROR] invalid request")
		log.Println(r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	destinationUrl := "https://character-service.dndbeyond.com/character/v5/character/" + character
	resp, err := http.DefaultClient.Get(destinationUrl)
	if err != nil {
		log.Println("[ERROR]", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR]", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		log.Println("[ERROR]", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
