package main

import (
	"encoding/json"
	"net/http"
)

type QueryTableColumn struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type QueryTableRow [2]string

type QueryTableResponse struct {
	Columns []QueryTableColumn `json:"columns"`
	Rows    []QueryTableRow    `json:"rows"`
	Type    string             `json:"type"`
}

func grafanaQueryHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	versions, _ := getDeploymentVersions()

	response := QueryTableResponse{
		Columns: []QueryTableColumn{
			{
				Text: "Application",
				Type: "string",
			},
			{
				Text: "Version",
				Type: "string",
			},
		},
		Type: "table",
		Rows: []QueryTableRow{},
	}

	for app, version := range versions {
		response.Rows = append(response.Rows, QueryTableRow{app, version})
	}

	jsonResponse, je := json.Marshal([]QueryTableResponse{response})
	if je != nil {
		panic(je.Error())
	}

	if _, we := w.Write(jsonResponse); we != nil {
		panic(we.Error())
	}
}

func grafanaSearchHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write([]byte("[\"versions\"]")); err != nil {
		panic(err.Error())
	}
}

func grafanaAnnotationHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write([]byte("[]")); err != nil {
		panic(err.Error())
	}
}

func grafanaFullVersionHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	versions, _ := getDeploymentVersions()

	jsonResponse, je := json.Marshal(versions)
	if je != nil {
		panic(je.Error())
	}

	if _, we := w.Write(jsonResponse); we != nil {
		panic(we.Error())
	}
}
