package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Product struct {
	Title        string
	Price        string
	Availability string
	URL          string
}

var products []Product
var mu sync.Mutex

func parseProduct(url string) Product {
	time.Sleep(300 * time.Millisecond) 

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Ошибка загрузки:", err)
		return Product{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("Статус страницы:", resp.StatusCode)
		return Product{}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println("Ошибка парсинга:", err)
		return Product{}
	}

	title := strings.TrimSpace(doc.Find("h1.title__font").Text())
	price := strings.TrimSpace(doc.Find("p.product-price__big-color-red").First().Text())
	availability := strings.TrimSpace(doc.Find("p.status-label--green").Text())

	return Product{
		Title:        title,
		Price:        price,
		Availability: availability,
		URL:          url,
	}
}

func fetchProductLinks(pageCount int) []string {
	var links []string

	for page := 1; page <= pageCount; page++ {
		url := fmt.Sprintf("https://hard.rozetka.com.ua/videocards/c80087/page=%d/", page)
		resp, err := http.Get(url)
		if err != nil {
			log.Println("Ошибка загрузки страницы:", err)
			continue
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Println("Ошибка парсинга страницы:", err)
			continue
		}

		doc.Find("a.tile-title").Each(func(i int, s *goquery.Selection) {
			link, exists := s.Attr("href")
			if exists {
				if !strings.HasPrefix(link, "http") {
					link = "https://hard.rozetka.com.ua" + link
				}
				links = append(links, link)
			}
		})
	}

	return links
}

func saveToCSV(data []Product, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Название", "Цена", "Наличие", "Ссылка"})
	for _, p := range data {
		writer.Write([]string{p.Title, p.Price, p.Availability, p.URL})
	}

	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index.html").ParseFiles("index.html"))
	tmpl.Execute(w, products)
}

func main() {
	log.Println("Сбор ссылок на товары...")
	productLinks := fetchProductLinks(2) 

	log.Printf("Найдено ссылок: %d\n", len(productLinks))

	var wg sync.WaitGroup
	productChan := make(chan Product)

	for _, link := range productLinks {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			p := parseProduct(url)
			productChan <- p
		}(link)
	}

	go func() {
		wg.Wait()
		close(productChan)
	}()

	for p := range productChan {
		if p.Title != "" {
			mu.Lock()
			products = append(products, p)
			mu.Unlock()
		}
	}

	log.Println("Сохранение CSV...")
	if err := saveToCSV(products, "products.csv"); err != nil {
		log.Println("Ошибка CSV:", err)
	} else {
		log.Println("CSV успешно сохранён: products.csv")
	}

	log.Println("Запуск веб-сервера на http://localhost:8080")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
