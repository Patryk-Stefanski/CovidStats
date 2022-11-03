package main

import (
	"CovidStats2/client"
	"bytes"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var templates = template.Must(template.ParseFiles("index.html"))
var templates2 = template.Must(template.ParseFiles("stats.html"))

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

	http.HandleFunc("/searchHistorical", HandleHistorical(statsApi))
	http.HandleFunc("/searchLive", HandleLive(statsApi))

	err = http.ListenAndServe(":3000", nil)
}

type Search struct {
	ChosenCountry string
	ChosenDate    string
	Results       *client.Results
}

func indexHandler(statsApi *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := &bytes.Buffer{}

		results, err := statsApi.GetCountries()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		listCountries := &Search{
			Results:       results,
			ChosenDate:    "",
			ChosenCountry: "",
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

func HandleLive(statsApi *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := &bytes.Buffer{}

		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()
		searchQuery := params.Get("country")

		if searchQuery != "" {
			results, err := statsApi.GetLiveStats(searchQuery)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			newCases, _ := strconv.Atoi(results.Entries[0].Cases.New)
			newDeaths, _ := strconv.Atoi(results.Entries[0].Deaths.New)

			stats := &PresentStats{
				Country:        results.Entries[0].Country,
				Date:           results.Entries[0].Day,
				Population:     results.Entries[0].Population,
				NewCases:       newCases,
				TotalCases:     results.Entries[0].Cases.Total,
				ActiveCases:    results.Entries[0].Cases.Active,
				CriticalCases:  results.Entries[0].Cases.Critical,
				RecoveredCases: results.Entries[0].Cases.Recovered,
				NewDeaths:      newDeaths,
				TotalDeaths:    results.Entries[0].Deaths.Total,
				TotalTests:     results.Entries[0].Tests.Total,
			}

			err = templates2.Execute(buf, stats)
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
}

type StatResults struct {
	StatResults *client.StatResults
}

type PresentStats struct {
	Country        string
	Date           string
	Population     int
	NewCases       int
	TotalCases     int
	ActiveCases    int
	CriticalCases  int
	RecoveredCases int
	NewDeaths      int
	TotalDeaths    int
	TotalTests     int
}

func HandleHistorical(statsApi *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var country, date = "", ""
		var population, totalCases, newCasesTotal, activeCases, criticalCases, recoveredCases, totalDeaths, totalTests, newDeathsTotal = 0, 0, 0, 0, 0, 0, 0, 0, 0

		buf := &bytes.Buffer{}

		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()
		searchQuery := params.Get("country")
		searchDate := params.Get("date")

		if searchQuery != "" && searchDate != "" {
			results, err := statsApi.GetHistoricalStats(searchQuery, searchDate)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			historyStatsAll := StatResults{
				StatResults: results,
			}

			layout := "2006-01-02T15:04:05"
			mostRecentTimestamp, _ := time.Parse(layout, historyStatsAll.StatResults.Entries[0].Time[0:19])

			for i := 0; i < len(historyStatsAll.StatResults.Entries); i++ {

				country = historyStatsAll.StatResults.Entries[i].Country
				date = historyStatsAll.StatResults.Entries[i].Day
				population = historyStatsAll.StatResults.Entries[i].Population

				newCases, _ := strconv.Atoi(historyStatsAll.StatResults.Entries[i].Cases.New)
				newDeaths, _ := strconv.Atoi(historyStatsAll.StatResults.Entries[i].Deaths.New)

				newCasesTotal += newCases
				newDeathsTotal += newDeaths

				currentTimestamp, _ := time.Parse(layout, historyStatsAll.StatResults.Entries[i].Time[0:19])

				if mostRecentTimestamp.Before(currentTimestamp) || mostRecentTimestamp.Equal(currentTimestamp) {
					mostRecentTimestamp = currentTimestamp
					totalCases = historyStatsAll.StatResults.Entries[i].Cases.Total
					activeCases = historyStatsAll.StatResults.Entries[i].Cases.Active
					criticalCases = historyStatsAll.StatResults.Entries[i].Cases.Critical
					recoveredCases = historyStatsAll.StatResults.Entries[i].Cases.Recovered
					totalDeaths = historyStatsAll.StatResults.Entries[i].Deaths.Total
					totalTests = historyStatsAll.StatResults.Entries[i].Tests.Total
				}

			}

			stats := &PresentStats{
				Country:        country,
				Date:           date,
				Population:     population,
				NewCases:       newCasesTotal,
				TotalCases:     totalCases,
				ActiveCases:    activeCases,
				CriticalCases:  criticalCases,
				RecoveredCases: recoveredCases,
				NewDeaths:      newDeathsTotal,
				TotalDeaths:    totalDeaths,
				TotalTests:     totalTests,
			}

			err = templates2.Execute(buf, stats)
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
}
