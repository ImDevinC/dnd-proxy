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
	mux.HandleFunc("/api/character", charHandler)
	mux.HandleFunc("/api/items", itemHandler)
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
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Println("[ERROR]", err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func charHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	character := r.URL.Query().Get("character")
	if len(character) == 0 {
		log.Println("[ERROR] invalid request")
		log.Println(r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	destinationUrl := "https://character-service.dndbeyond.com/character/v5/character/" + character
	forwardRequest(&w, destinationUrl)
}

func itemHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	destinationUrl := "https://character-service.dndbeyond.com/character/v5.1/game-data/items?" + r.URL.Query().Encode()
	forwardRequest(&w, destinationUrl)
}

func forwardRequest(w *http.ResponseWriter, destinationUrl string) error {
	resp, err := http.DefaultClient.Get(destinationUrl)
	if err != nil {
		log.Println("[ERROR]", err.Error())
		(*w).WriteHeader(http.StatusInternalServerError)
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR]", err.Error())
		(*w).WriteHeader(http.StatusInternalServerError)
		return err
	}
	_, err = (*w).Write(body)
	if err != nil {
		log.Println("[ERROR]", err.Error())
		(*w).WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}
