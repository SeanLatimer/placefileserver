package sn

import (
	"bytes"
	"text/template"
	"time"
)

const reportingDuration = (time.Hour * 6)
const stationaryDuration = (time.Minute * 10)

var tpl *template.Template

func init() {
	fmap := template.FuncMap{
		"isReporting":  isReporting,
		"isStationary": isStationary,
	}
	tpl = template.Must(template.New("placefile").Funcs(fmap).Parse(PlaceFileTemplate))
}

func GeneratePlacefile(title string, spotters []Spotter) (*bytes.Buffer, error) {
	data := struct {
		Title    string
		Spotters []Spotter
	}{}
	data.Title = title
	data.Spotters = spotters

	buf := bytes.Buffer{}
	err := tpl.Execute(&buf, data)

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
