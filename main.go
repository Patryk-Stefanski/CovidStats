package main

import (
	"CovidStats2/client"
	"bytes"
	"fmt"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var templates = template.Must(template.ParseFiles("index.html"))

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	apiKey := os.Getenv("COVID_STATS_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}
	myClient := &http.Client{Timeout: 10 * time.Second}
	statsApi := client.NewClient(myClient, apiKey)

	http.HandleFunc("/", indexHandler(statsApi))
	http.HandleFunc("/search", CalculateDailyStats(statsApi))

	//http.HandleFunc("/search", searchHandler(statsApi))
	err = http.ListenAndServe(":3000", nil)
}

type listCountries struct {
	Results *client.Results
}

func indexHandler(statsApi *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := &bytes.Buffer{}

		results, err := statsApi.GetCountries()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		listCountries := &listCountries{
			Results: results,
		}

		err = templates.Execute(buf, listCountries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = buf.WriteTo(w)
		if err != nil {
			return
		}
	}
}

type Search struct {
	Query   string
	Date    string
	Results *client.StatResults
}

func searchHandler(statsApi *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := &bytes.Buffer{}
		err := templates.Execute(buf, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		results, err := statsApi.GetCountries()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		listCountries := &listCountries{
			Results: results,
		}

		err = templates.Execute(buf, listCountries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)
	}
}

func CalculateDailyStats(statsApi *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()
		searchQuery := params.Get("country")
		searchDate := params.Get("date")

		fmt.Println(searchDate)
		fmt.Println(searchQuery)

		if searchQuery != "" && searchDate != "" {
			results, err := statsApi.GetHistoricalStats(searchQuery, searchDate)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Println(results)
		}

	}
}
