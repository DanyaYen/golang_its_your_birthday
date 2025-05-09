package main

import(
	"fmt"
	"net/http"
	"sync"
	"log"
)

var storage map[string][]byte
var mu sync.RWMutex


func set(key string, value []byte) {
	mu.Lock()
	defer mu.Unlock()
	storage[key] = value
	log.Printf("SET: key=%s, value=%s\n", key, string(value))
}

func get(key string) ([]byte, bool) {
	mu.RLock()
	defer mu.RUnlock()
	value, found := storage[key]
	if found {
		log.Printf("GET: key=%s value=%s\n", key, string(value))
	} else {
		log.Printf("GET: key=%s (not found)\n", key)
	}
	return value, found
}

func deleteKey(key string) bool {
	mu.Lock()
	defer mu.Unlock()
	
	_, found := storage[key]
	if found {
		delete(storage, key)
		log.Printf("DELETE: key=%s\n", key)
		return true
	}
	log.Printf("DELETE: key=%s (not found for deletion)\n", key)
	return false
}

func kvHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key") 
	
	if key == "" {
		http.Error(w, "Missing key parameter", http.StatusBadRequest)
		return
	}
	
	switch r.Method {
	case http.MethodGet:
		value, found := get(key)
		if !found {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(value) 
	case http.MethodPost:
		value := r.URL.Query().Get("value")
		if value == "" {
			http.Error(w, "Missing value parameter for POST", http.StatusBadRequest)
			return
		}
		
		set(key, []byte(value))
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "OK")
	case http.MethodDelete:
		wasDeleted := deleteKey(key)
		if wasDeleted { 
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Key '", key, "' deleted successfully")
		} else {
			http.Error(w, "Key not found for deletion", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	storage = make(map[string][]byte)
	
	log.Println("Run server on http://localhost:8080")
	log.Println("Try GET: curl http://localhost:8080/kv?key=somekey")
	log.Println("Try SET: curl -X POST \"http://localhost:8080/kv?key=somekey&value=somevalue\"")
	
	http.HandleFunc("/kv", kvHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
