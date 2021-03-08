package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	mu = sync.Mutex{}
	lastUpdate time.Time
	cachedMetrics map[string]string
	cachedDebug []string
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

func updateMetrics() bool {
	if time.Since(lastUpdate).Seconds() > 30 {
		lastUpdate = time.Now()
		return true
	}
	return false
}

func getDeploymentVersions() (map[string]string, []string) {
	mu.Lock()
	defer mu.Unlock()
	if updateMetrics() {
		deployments, err := globalConfig.clientSet.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		cachedMetrics = make(map[string]string)
		cachedDebug = []string{}

		for _, item := range deployments.Items {
			version := item.Labels["app.kubernetes.io/version"]
			if strings.TrimSpace(version) == "" {
				cachedDebug = append(cachedDebug, fmt.Sprintf("'%s' has no 'app.kubernetes.io/version' label", item.Name))
				continue
			}

			instance := item.Labels["app.kubernetes.io/instance"]
			name := item.Labels["app.kubernetes.io/name"]
			component := item.Labels["app.kubernetes.io/component"]

			name = getMainComponentName(instance, name)

			// Is main component
			if item.Name != name && item.Name != fmt.Sprintf("%s-%s", name, component){
				cachedDebug = append(cachedDebug, fmt.Sprintf("'%s' mismatches app name '%s'", item.Name, name))
				continue
			}

			cachedMetrics[name] = sanitizeVersion(version)
		}
	}

	return cachedMetrics, cachedDebug
}