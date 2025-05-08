package main

import (
	"fmt"
	"net/http"
	// "time"
	// "io"
)

type Site struct {
	URL    string
	Online bool
}

var sites = []Site{}


func checkURL(url string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == http.StatusOK, nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Site monitoring</title>
		</head>
		<body>
			<h1>Add site</h1>
			<form method="POST" action="/add-url">
				<label for="url">Site URL:</label>
				<input type="text" id="url" name="url_to_check" size="50">
				<input type="submit" value="Add">
			</form>
			<h2>Site to check:</h2>
			<ul>
	`
	for _, site := range sites {
		status := "Offline"
		if site.Online {
			status = "Online"
		}
		html += fmt.Sprintf("<li>%s — %s</li>", site.URL, status)
	}
	html += `
			</ul>
		</body>
		</html>
	`
	w.Header().Set("Content-Type", "text/html; charset=utf-8") 
	fmt.Fprint(w, html)  
}

func addUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error while parsing", http.StatusBadRequest)
			return
		}
		
		newUrl := r.FormValue("url_to_check")
		if newUrl != "" {
			online, _ := checkURL(newUrl)
			sites = append(sites, Site{URL: newUrl, Online: online})
			fmt.Printf("Added URL: %s (Online: %v)\n", newUrl, online)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/add-url", addUrlHandler)
	
	fmt.Println("Server is running on http://localhost:8080 ...")
	err := http.ListenAndServe(":8080", nil) 
	if err != nil {
		fmt.Println("Error while starting:", err)
	}
}