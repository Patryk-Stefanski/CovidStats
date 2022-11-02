package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	http *http.Client
	key  string
}

type Results struct {
	Countries []string `json:"response"`
}

func (c *Client) GetCountries() (*Results, error) {
	url := "https://covid-193.p.rapidapi.com/countries"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", c.key)
	req.Header.Add("X-RapidAPI-Host", "covid-193.p.rapidapi.com")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

	res := &Results{}
	return res, json.Unmarshal(body, res)
}

func (c *Client) GetHistoricalStats(country, date string) (*StatResults, error) {
	url := fmt.Sprintf("https://covid-193.p.rapidapi.com/history?country=%s&day=%s", country, date)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", c.key)
	req.Header.Add("X-RapidAPI-Host", "covid-193.p.rapidapi.com")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

	res := &StatResults{}
	return res, json.Unmarshal(body, res)
}

func (c *Client) GetLiveStats(country string) (*StatResults, error) {
	url := fmt.Sprintf("https://covid-193.p.rapidapi.com/history?country=%s", country)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", c.key)
	req.Header.Add("X-RapidAPI-Host", "covid-193.p.rapidapi.com")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

	res := &StatResults{}
	return res, json.Unmarshal(body, res)
}

func NewClient(httpClient *http.Client, key string) *Client {
	return &Client{httpClient, key}
}

type StatResults struct {
	Entries []Stats `json:"response"`
}

type Stats struct {
	Country    string `json:"country"`
	Population int    `json:"population"`
	Cases      Cases  `json:"cases"`
	Deaths     Deaths `json:"deaths"`
	Tests      Tests  `json:"tests"`
	Day        string `json:"day"`
	Time       string `json:"time"`
}

type Cases struct {
	New       string `json:"new"`
	Active    int    `json:"active"`
	Critical  int    `json:"critical"`
	Recovered int    `json:"recovered"`
	Total     int    `json:"total"`
}

type Deaths struct {
	Total int `json:"total"`
}

type Tests struct {
	Total int `json:"total"`
}
