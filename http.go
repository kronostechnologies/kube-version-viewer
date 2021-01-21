package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func getHtmlTemplate() *template.Template {
	templateString := `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Versions</title>
</head>
<body>
<table border="1">
  <thead>
    <tr><th>Application</th><th>Version</th></tr>
  </thead>
  <tbody>
  {{- range $key, $value := .Versions }}
     <tr><td>{{ $key }}</td><td>{{ $value }}</td></tr>
  {{- end }}
  </tbody>
</table>
{{- if .ShowDebug }}
<h2>Debug</h2>
{{- range .Debug }}
{{ . }}<br/>
{{- end }}
{{- end }}
</body>
</html>`

	t, err := template.New("htmlVersions").Parse(templateString)
	if err != nil {
		panic(err.Error())
	}

	return t
}

func httpHealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func httpVersionHandlerFactory(t *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		versions, debug := getDeploymentVersions()
		err := t.Execute(w, map[string]interface{}{
			"Versions": versions,
			"Debug": debug,
			"ShowDebug": *globalConfig.debug,
		})
		if err != nil {
			panic(err.Error())
		}
	}
}

func httpServe(listen string) {
	t := getHtmlTemplate()

	http.HandleFunc("/health", httpHealthHandler)
	http.HandleFunc("/versions", httpVersionHandlerFactory(t))
	fmt.Printf("Listening on %s\n", listen)
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		panic(err.Error())
	}
}