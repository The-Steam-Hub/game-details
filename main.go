package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	mongoURI = ""
	games    = 0
	request  = 0
)

const (
	mongoDB   = "steam_hub"
	mongoColl = "games"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
}

func init() {
	// Load the environment and default to "development" if not set
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	logrus.Info("Loading configurations...")

	// Load environment variables from the .env file based on the APP_ENV
	if err := godotenv.Load(".env." + env); err != nil {
		logrus.Warnf("Error loading .env file: %s. Checking for environment variables...", err)
	}

	// Retrieve MongoDB connection string
	mongoURI := os.Getenv("MONGO_CONNECTION_STRING")
	if mongoURI == "" {
		logrus.Fatal("Environment variable 'MONGO_CONNECTION_STRING' has not been initialized")
	}

	// Continue with the rest of your application...
	logrus.Info("MongoDB connection string loaded")
}

func main() {
	apps, err := AppsList()
	if err != nil {
		fmt.Println(err)
	}

	var wg sync.WaitGroup

	jobs := make(chan AppList, len(*apps))
	results := make(chan AppData, len(*apps))

	for _, app := range *apps {
		jobs <- app
	}
	close(jobs)

	for w := 0; w <= 5; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	wg.Wait()
	close(results)

}

func worker(id int, jobs <-chan AppList, results chan<- AppData, wg *sync.WaitGroup) {
	defer wg.Done()
	for app := range jobs {
		if app.AppID != 0 && app.Name != "" {
			request = request + 1

			proxyURL, _ := url.Parse("http://164.90.255.167:3128")
			client := &http.Client{
				Timeout: time.Second * 10,
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyURL),
				},
			}

			appData, err := AppDetailedData(app.AppID, client)
			if err != nil {
				logrus.Error(err)
			}

			if appData.Type == "game" {
				games = games + 1
			}

			logrus.WithFields(logrus.Fields{
				"game-count":    games,
				"request-count": request,
				"appID":         appData.SteamAppid,
				"type":          appData.Type,
				"name":          appData.Name,
				"jobID":         id,
			}).Info()

			results <- *appData
			time.Sleep(time.Millisecond * 1350)
		}
	}
}
