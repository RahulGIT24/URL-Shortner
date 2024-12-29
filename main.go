package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL))
	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)
	return hash[:8]
}

func createURL(OriginalURL string) string {
	shortURL := generateShortURL(OriginalURL)
	id := shortURL
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  OriginalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}

	return shortURL
}

func getOriginalURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL Not Found")
	}
	return url, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	shortURL_ := createURL(data.URL)

	// fmt.Fprintf(w,shortURL)
	response := struct {
		Short_URL string `json: "short_url"`
	}{Short_URL: shortURL_}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectURLHandler(w http.ResponseWriter, r *http.Request){
	id := r.URL.Path[len("/redirect/"):]
	url,err := getOriginalURL(id)

	if err!=nil{
		http.Error(w,"Invalid Request",http.StatusNotFound)
	}

	http.Redirect(w,r,url.OriginalURL,http.StatusFound)
}

func main() {
	fmt.Println("Starting URL Shortner")
	fmt.Printf("Server Starting on port 8090\n")

	// Register the handler function to handle all requests
	http.HandleFunc("/", handler)
	http.HandleFunc("/convertURL", ShortURLHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println("Error while starting server: ", err)
	}
}
