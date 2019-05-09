package sn

import (
	"bytes"
	"log"
	"text/template"
	"time"

	"github.com/gobuffalo/packr/v2"
)

const reportingDuration = (time.Hour * 6)
const stationaryDuration = (time.Minute * 10)

var (
	grTpl     *template.Template
	slTpl     *template.Template
	gearthTpl *template.Template
	rssTpl    *template.Template
)

func Initialize(box *packr.Box) {
	fmap := template.FuncMap{
		"isReporting":  isReporting,
		"isStationary": isStationary,
	}

	tplStr, err := box.FindString("gr_placefile.txt")
	if err != nil {
		log.Fatal(err)
	}
	grTpl = template.Must(template.New("gr_placefile").Funcs(fmap).Parse(tplStr))

	tplStr, err = box.FindString("sl_placefile.txt")
	if err != nil {
		log.Fatal(err)
	}
	slTpl = template.Must(template.New("sl_placefile").Funcs(fmap).Parse(tplStr))

	tplStr, err = box.FindString("gearth.xml")
	if err != nil {
		log.Fatal(err)
	}
	gearthTpl = template.Must(template.New("gearth").Funcs(fmap).Parse(tplStr))

	tplStr, err = box.FindString("rss.xml")
	if err != nil {
		log.Fatal(err)
	}
	rssTpl = template.Must(template.New("rss").Funcs(fmap).Parse(tplStr))
}

func GenerateGRPlacefile(title string, spotters []Spotter) (*bytes.Buffer, error) {
	data := struct {
		Title    string
		Spotters []Spotter
	}{}
	data.Title = title
	data.Spotters = spotters

	buf := bytes.Buffer{}
	err := grTpl.Execute(&buf, data)

	if err != nil {
		return nil, err
	}

	return &buf, nil
}

func GenerateSLPlacefile(title string, spotters []Spotter) (*bytes.Buffer, error) {
	data := struct {
		Title    string
		Spotters []Spotter
	}{}
	data.Title = title
	data.Spotters = spotters

	buf := bytes.Buffer{}
	err := slTpl.Execute(&buf, data)

	if err != nil {
		return nil, err
	}

	return &buf, nil
}

func GenerateGearthPlacefile(spotters []Spotter) (*bytes.Buffer, error) {
	data := struct {
		Spotters []Spotter
	}{}
	data.Spotters = spotters

	buf := bytes.Buffer{}
	err := gearthTpl.Execute(&buf, data)

	if err != nil {
		return nil, err
	}

	return &buf, nil
}

func GenerateRssPlacefile(title string, url string, spotters []Spotter) (*bytes.Buffer, error) {
	data := struct {
		Title    string
		Url      string
		Spotters []Spotter
	}{}
	data.Title = title
	data.Url = url
	data.Spotters = spotters

	buf := bytes.Buffer{}
	err := rssTpl.Execute(&buf, data)

	if err != nil {
		return nil, err
	}

	return &buf, nil
}

func isReporting(unix int64) (bool, error) {
	now := time.Now().UTC()
	lastReport := time.Unix(unix, 0).UTC()
	duration := now.Sub(lastReport)

	reporting := true
	if duration.Hours() >= reportingDuration.Hours() {
		reporting = false
	}

	return reporting, nil
}

func isStationary(unix int64) (bool, error) {
	now := time.Now().UTC()
	lastReport := time.Unix(unix, 0).UTC()
	duration := now.Sub(lastReport)

	stationary := false
	if duration.Minutes() >= stationaryDuration.Minutes() {
		stationary = true
	}

	return stationary, nil
}
