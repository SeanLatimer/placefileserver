package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"text/template"
	"time"

	"github.com/seanlatimer/placefileserver/sn"

	"github.com/fsnotify/fsnotify"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"
)

var (
	infoTpl      *template.Template
	infoPage     = new(bytes.Buffer)
	infoPageLock = new(sync.RWMutex)

	placefile     = new(bytes.Buffer)
	placefileLock = new(sync.RWMutex)

	reloadCounter = new(ReloadCounter)
)

func init() {
	infoTpl = template.Must(template.New("infoPage").Parse(infoPageTemplate))
}

func main() {
	log.Println("SVL Placefile Generator")
	setupConfig()

	if !viper.IsSet("Port") {
		log.Fatal("Port is not set in the config")
	}

	go updateInfoPage()
	go updatePlacefile()

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		infoPageLock.RLock()
		w.Write(infoPage.Bytes())
		infoPageLock.RUnlock()
	})
	r.Get("/gr", func(w http.ResponseWriter, r *http.Request) {
		placefileLock.RLock()
		w.Write(placefile.Bytes())
		placefileLock.RUnlock()
	})

	port := viper.GetString("Port")
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func setupConfig() {
	viper.AddConfigPath("./cfg")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetDefault("Port", 5000)
	viper.SetDefault("PlacefileTitle", "SVL Spotters")
	viper.SetDefault("SpotterNetworkAppID", "")
	viper.SetDefault("SpotterIds", []int{})

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		current := reloadCounter.Increment()
		go func() {
			time.Sleep(time.Millisecond * 500)
			if current == reloadCounter.Current() {
				log.Println("Reloading Configuration")
				go updateInfoPage()
			}
		}()
	})
}

func getSpotters() ([]sn.Spotter, error) {
	appID := viper.GetString("SpotterNetworkAppID")
	spotterIds := GetIntSlice("SpotterIds")

	spotters, err := sn.GetSpotters(appID, spotterIds)
	if err, ok := err.(*url.Error); ok {
		if err.Timeout() {
			log.Println("Failed getting spotters ", err)
			return nil, err
		}
		log.Fatal("Failed getting spotters ", err)
	}

	return spotters, nil
}

func updatePlacefile() {
	for {
		if !viper.IsSet("SpotterNetworkAppID") && !viper.IsSet("SpotterIds") {
			log.Fatal("Spotter Network App ID/Spotter IDs are missing from the config")
		}

		log.Println("Updating placefile")
		spotters, err := getSpotters()
		if err != nil {
			log.Println("Failed getting spotters ", err)
			return
		}

		buf, err := sn.GeneratePlacefile(viper.GetString("PlacefileTitle"), spotters)
		if err != nil {
			log.Fatal("Failed to generate placefile ", err)
		}

		placefileLock.Lock()
		placefile = buf
		placefileLock.Unlock()
		log.Println("Done updating placefile")

		time.Sleep(time.Minute * 2)
	}
}

func updateInfoPage() {
	log.Println("Updating info page")
	spotters, err := getSpotters()
	if err != nil {
		log.Println("Failed getting spotters ", err)
		return
	}

	data := struct {
		Spotters []sn.Spotter
	}{}
	data.Spotters = spotters

	buf := bytes.Buffer{}
	err = infoTpl.Execute(&buf, data)
	if err != nil {
		log.Fatal("Failed to generate info page ", err)
	}

	infoPageLock.Lock()
	infoPage = &buf
	infoPageLock.Unlock()
	log.Println("Done updating info page")
}
