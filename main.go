//go:generate packr2
package main

import (
	"bytes"
	"fmt"
	"html/template"
	htmltemplate "html/template"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/seanlatimer/placefileserver/sn"

	"github.com/fsnotify/fsnotify"
	"github.com/go-chi/chi"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/viper"
)

var (
	infoTpl      *htmltemplate.Template
	infoPage     = new(bytes.Buffer)
	infoPageLock = new(sync.RWMutex)

	kmlTpl  *template.Template
	kmlFile = new(bytes.Buffer)
	kmlLock = new(sync.RWMutex)

	grPlacefile     = new(bytes.Buffer)
	grPlacefileLock = new(sync.RWMutex)

	slPlacefile     = new(bytes.Buffer)
	slPlacefileLock = new(sync.RWMutex)

	gearthPlacefile     = new(bytes.Buffer)
	gearthPlacefileLock = new(sync.RWMutex)

	rssFile     = new(bytes.Buffer)
	rssFileLock = new(sync.RWMutex)

	reloadCounter = new(ReloadCounter)
)

func initialize(box *packr.Box) {
	tplStr, err := box.FindString("index.html")
	if err != nil {
		log.Fatal(err)
	}
	infoTpl = htmltemplate.Must(htmltemplate.New("index").Parse(tplStr))

	tplStr, err = box.FindString("gearth.kml")
	if err != nil {
		log.Fatal(err)
	}
	kmlTpl = template.Must(template.New("gearth.kml").Parse(tplStr))
}

func main() {
	log.Println("SVL Placefile Generator")
	setupConfig()

	if !viper.IsSet("Port") {
		log.Fatal("Port is not set in the config")
	}

	box := packr.New("templates", "./templates")
	initialize(box)
	sn.Initialize(box)

	go updateInfoPage()
	go updateKml()
	go updatePlacefile()

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		infoPageLock.RLock()
		w.Write(infoPage.Bytes())
		infoPageLock.RUnlock()
	})
	r.Get("/gr", func(w http.ResponseWriter, r *http.Request) {
		grPlacefileLock.RLock()
		w.Write(grPlacefile.Bytes())
		grPlacefileLock.RUnlock()
	})
	r.Get("/sl", func(w http.ResponseWriter, r *http.Request) {
		slPlacefileLock.RLock()
		w.Write(slPlacefile.Bytes())
		slPlacefileLock.RUnlock()
	})
	r.Get("/gearth.kml", func(w http.ResponseWriter, r *http.Request) {
		kmlLock.RLock()
		w.Write(kmlFile.Bytes())
		kmlLock.RUnlock()
	})
	r.Get("/gearth.xml", func(w http.ResponseWriter, r *http.Request) {
		gearthPlacefileLock.RLock()
		w.Write(gearthPlacefile.Bytes())
		gearthPlacefileLock.RUnlock()
	})
	r.Get("/rss", func(w http.ResponseWriter, r *http.Request) {
		rssFileLock.RLock()
		w.Write(rssFile.Bytes())
		rssFileLock.RUnlock()
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
	viper.SetDefault("BaseUrl", "http://localhost:5000")
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
				go updateKml()
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

		log.Println("Updating placefiles")
		spotters, err := getSpotters()
		if err != nil {
			log.Println("Failed getting spotters ", err)
			return
		}

		wg := new(sync.WaitGroup)
		wg.Add(4)
		go func() {
			buf, err := sn.GenerateGRPlacefile(viper.GetString("Title"), spotters)
			if err != nil {
				log.Fatal("Failed to generate GR placefile ", err)
			}
			grPlacefileLock.Lock()
			grPlacefile = buf
			grPlacefileLock.Unlock()
			wg.Done()
		}()

		go func() {
			buf, err := sn.GenerateSLPlacefile(viper.GetString("Title"), spotters)
			if err != nil {
				log.Fatal("Failed to generate StormLab placefile ", err)
			}
			slPlacefileLock.Lock()
			slPlacefile = buf
			slPlacefileLock.Unlock()
			wg.Done()
		}()

		go func() {
			buf, err := sn.GenerateGearthPlacefile(spotters)
			if err != nil {
				log.Fatal("Failed to generate Google Earth placefile ", err)
			}
			gearthPlacefileLock.Lock()
			gearthPlacefile = buf
			gearthPlacefileLock.Unlock()
			wg.Done()
		}()

		go func() {
			buf, err := sn.GenerateRssPlacefile(viper.GetString("Title"), viper.GetString("BaseUrl"), spotters)
			if err != nil {
				log.Fatal("Failed to generate RSS placefile ", err)
			}
			rssFileLock.Lock()
			rssFile = buf
			rssFileLock.Unlock()
			wg.Done()
		}()

		wg.Wait()
		log.Println("Done updating placefiles")

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
		Title    string
		Spotters []sn.Spotter
	}{}
	data.Title = viper.GetString("Title")
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

func updateKml() {
	log.Println("Updating KML")
	data := struct {
		Title string
		Url   string
	}{}
	data.Title = viper.GetString("Title")
	data.Url = viper.GetString("BaseUrl")

	buf := bytes.Buffer{}
	err := kmlTpl.Execute(&buf, data)
	if err != nil {
		log.Fatal("Failed to generate KML ", err)
	}

	kmlLock.Lock()
	kmlFile = &buf
	kmlLock.Unlock()

	log.Println("Done updating KML")
}
