package main

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var store = make(map[string]string)
var mu sync.RWMutex
var templates *template.Template

func main() {
	templates = template.Must(template.ParseGlob("*.html"))
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/", redirectHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func generateKey(n int) string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		result[i] = alphabet[idx.Int64()]
	}
	return string(result)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := templates.ExecuteTemplate(w, "form.html", nil)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}
		err := r.ParseForm()
	if err != nil {http.Error(w, "Wrong request", 400); return}
	orig := r.FormValue("url")
	
	u, err := url.Parse(orig)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		http.Error(w, "URL invalid", 400)
		return 
	}
	var key string
	for {
		key = generateKey(6)
		mu.RLock()
		_, exist := store[key]
		mu.RUnlock()
		if !exist {
			break
		}
	}
	mu.Lock()
	store[key] = orig
	mu.Unlock()
	shortURL := fmt.Sprintf("http://localhost:8080/%s", key)
	err = templates.ExecuteTemplate(w, "form.html", shortURL)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/")
	mu.RLock()
	orig, ok := store[key]
	mu.RUnlock()
	
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, orig, http.StatusFound)
}