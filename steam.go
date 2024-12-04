package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type AppList struct {
	AppID int    `json:"appid"`
	Name  string `json:"name"`
}

type AppData struct {
	Type                string   `json:"type"`
	Name                string   `json:"name"`
	SteamAppid          int      `json:"steam_appid"`
	RequiredAge         int      `json:"required_age"`
	IsFree              bool     `json:"is_free"`
	ControllerSupport   string   `json:"controller_support"`
	DLC                 []int    `json:"dlc"`
	DetailedDescription string   `json:"detailed_description"`
	AboutTheGame        string   `json:"about_the_game"`
	ShortDescription    string   `json:"short_description"`
	SupportedLanguages  string   `json:"supported_languages"`
	Developers          []string `json:"developers"`
	Publishers          []string `json:"publishers"`
	PriceOverview       struct {
		Currency         string `json:"currency"`
		Initial          int    `json:"initial"`
		Final            int    `json:"final"`
		DiscountPercent  int    `json:"discount_percent"`
		InitialFormatted string `json:"initial_formatted"`
		FinalFormatted   string `json:"final_formatted"`
	} `json:"price_overview"`
	Platforms struct {
		Windows bool `json:"windows"`
		Mac     bool `json:"mac"`
		Linux   bool `json:"linux"`
	} `json:"platforms"`
	Categories []struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	} `json:"categories"`
	Genres []struct {
		ID          string `json:"id"`
		Description string `json:"description"`
	} `json:"genres"`
	Recommendations struct {
		Total int `json:"total"`
	} `json:"recommendations"`
	ReleaseDate struct {
		ComingSoon bool   `json:"coming_soon"`
		Date       string `json:"date"`
	} `json:"release_date"`
}

const (
	SteamWebAPI           = "http://api.steampowered.com/"
	SteamPoweredAPI       = "https://store.steampowered.com/"
	SteamWebAPIISteamApps = SteamWebAPI + "ISteamApps/"
)

func AppsList() (*[]AppList, error) {
	baseURL, _ := url.Parse(SteamWebAPIISteamApps)
	baseURL.Path += "GetAppList/v2/"

	resp, err := http.Get(baseURL.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		AppList struct {
			Apps []AppList `json:"apps"`
		} `json:"applist"`
	}

	json.Unmarshal(b, &response)
	return &response.AppList.Apps, nil
}

func AppDetailedData(appID int, client *http.Client) (*AppData, error) {
	baseURL, _ := url.Parse(SteamPoweredAPI)
	baseURL.Path += "api/appdetails"

	params := url.Values{}
	params.Add("appids", strconv.Itoa(appID))
	params.Add("l", "english")
	params.Add("cc", "US")
	baseURL.RawQuery = params.Encode()

	// Use the provided client or default to http.DefaultClient
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Get(baseURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response map[string]struct {
		AppData AppData `json:"data"`
	}

	json.Unmarshal(b, &response)
	appData := response[strconv.Itoa(appID)].AppData
	return &appData, nil
}
