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

type URL struct { //database
	ID           string    `json:"id"`
	OriginalUrl  string    `json:"original_url"`
	ShortUrl     string    `json:"short_url"`
	CreationDate time.Time `json:"create_date"`
}

var urlDb = make(map[string]URL) // generating shortURl by hashing

func GenerateShortUrl(original string) string {
	hasher := md5.New()
	hasher.Write([]byte(original))
	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)
	return hash[:8]

}
func createUrl(originalurl string) string { // storing url in db
	shorturl := GenerateShortUrl(originalurl)
	id := shorturl
	urlDb[id] = URL{
		ID:           id,
		OriginalUrl:  originalurl,
		ShortUrl:     shorturl,
		CreationDate: time.Now(),
	}
	return shorturl
}
func getUrl(id string) (URL, error) {
	url, ok := urlDb[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	} else {
		return url, nil
	}
}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GET READY TO SHORTEN UR URL")
}
func ShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Url string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	shortUrl := createUrl(data.Url)
	response := struct {
		Shorturl string `json:"short_url"`
	}{Shorturl: shortUrl}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectUrlhandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getUrl(id)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalUrl, http.StatusFound)
}
func main() {

	http.HandleFunc("/", handler)
	http.HandleFunc("/shortner", ShortUrlHandler)
	http.HandleFunc("/redirect/", redirectUrlhandler)
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		fmt.Println("Error while starting a server", err)
	}

}
