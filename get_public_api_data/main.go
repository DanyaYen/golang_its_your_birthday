package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

type NumberFact struct {
	Text   string  `json:"text"`
	Found  bool    `json:"found"`
	Number float64 `json:"number"`
	Type   string  `json:"type"`
	Date   string  `json:"date,omitempty"`
	Year   string  `json:"year,omitempty"`
}

var templates *template.Template

func main() {
	templates = template.Must(template.ParseGlob("*.html"))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/fact", factHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fetchFact(number string) (NumberFact, error) {
	url := "http://numbersapi.com/random?json"
	if number != "" {
		url = fmt.Sprintf("http://numbersapi.com/%s?json", number)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return NumberFact{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return NumberFact{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return NumberFact{}, fmt.Errorf("numbersapi returned status %d", resp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return NumberFact{}, err
	}
	var fact NumberFact
	if err := json.Unmarshal(bodyBytes, &fact); err != nil {
		return NumberFact{}, err
	}
	return fact, nil
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return templates.ExecuteTemplate(w, name, data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := renderTemplate(w, "index.html", nil); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}

func factHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	number := r.URL.Query().Get("number")
	if number != "" {
		var tmp float64
		if _, err := fmt.Sscanf(number, "%f", &tmp); err != nil {
			http.Error(w, "Invalid number format", http.StatusBadRequest)
			return
		}
	}
	fact, err := fetchFact(number)
	if err != nil {
		http.Error(w, "Could not fetch fact: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := renderTemplate(w, "fact.html", fact); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}
