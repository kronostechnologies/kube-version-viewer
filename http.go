package main

import (
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

	w.Header().Set("Access-Control-Allow-Headers", "accept, content-type")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(http.StatusOK)
}

func httpVersionHandlerFactory() func(w http.ResponseWriter, r *http.Request) {
	t := getHtmlTemplate()

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