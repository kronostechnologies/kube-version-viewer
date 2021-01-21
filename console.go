package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func printVersions() {
	templateString := `NAME                             VERSION
{{ range $key, $value := . }}
{{- FormatName $key }} {{ $value }}
{{ end }}`

	t, err := template.New("printVersions").Funcs(template.FuncMap{
		"FormatName": func(s string) string { return fmt.Sprintf("%-32s", s) },
	}).Parse(templateString)
	if err != nil {
		panic(err.Error())
	}

	versions, debug := getDeploymentVersions()
	err = t.Execute(os.Stdout, versions)
	if err != nil {
		panic(err.Error())
	}
	if *globalConfig.debug {
		fmt.Println(strings.Join(debug, "\n"))
	}
}
