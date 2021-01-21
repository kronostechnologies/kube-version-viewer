package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
	"strings"
)

func getMainComponentName(instance string, name string) string {
	r, _ := regexp.Compile(fmt.Sprintf("^%s(-.*)?$", instance))
	if !r.MatchString(name) {
		name = fmt.Sprintf("%s-%s", instance, name)
	}
	return name
}

func sanitizeVersion(version string) string {
	vr, e := regexp.Compile("^(?:version-|v)([0-9])")
	if e != nil {
		fmt.Println(e)
	}
	return vr.ReplaceAllString(version, "${1}")
}

func getDeploymentVersions() (map[string]string, []string) {
	deployments, err := globalConfig.clientSet.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	versions := make(map[string]string)
	debug := []string{}

	for _, item := range deployments.Items {
		version := item.Labels["app.kubernetes.io/version"]
		if strings.TrimSpace(version) == "" {
			debug = append(debug, fmt.Sprintf("'%s' has no 'app.kubernetes.io/version' label", item.Name))
			continue
		}

		instance := item.Labels["app.kubernetes.io/instance"]
		name := item.Labels["app.kubernetes.io/name"]
		component := item.Labels["app.kubernetes.io/component"]

		name = getMainComponentName(instance, name)

		// Is main component
		if item.Name != name && item.Name != fmt.Sprintf("%s-%s", name, component){
			debug = append(debug, fmt.Sprintf("'%s' mismatches app name '%s'", item.Name, name))
			continue
		}

		versions[name] = sanitizeVersion(version)
	}

	return versions, debug
}