package main

import (
	"fmt"
	"net/http"
)

func corsHandler(fn http.HandlerFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "accept, content-type")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method != "OPTIONS" {
			fn(w, r)
		}
	}
}

func createServer(listen string) {
	http.HandleFunc("/", corsHandler(httpHealthHandler))
	http.HandleFunc("/versions", httpVersionHandlerFactory())
	http.HandleFunc("/query", corsHandler(grafanaQueryHandler))
	http.HandleFunc("/search", corsHandler(grafanaSearchHandler))
	http.HandleFunc("/annotations", corsHandler(grafanaAnnotationHandler))
	fmt.Printf("Listening on %s\n", listen)
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		panic(err.Error())
	}
}