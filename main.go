package main

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

type Config struct {
	debug     *bool
	clientSet *kubernetes.Clientset
}

var globalConfig Config

func setupKubernetesClient() {
	var restConfig *rest.Config
	var err error

	if _, err = os.Stat("/var/run/secrets/kubernetes.io"); err != nil {
		var kubeConfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeConfig = flag.String("kube-config", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kube config file")
		} else {
			kubeConfig = flag.String("kube-config", "", "absolute path to the kube config file")
		}

		restConfig, err = clientcmd.BuildConfigFromFlags("", *kubeConfig)
	} else {
		restConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		panic(err.Error())
	}

	globalConfig.clientSet, err = kubernetes.NewForConfig(restConfig)
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	globalConfig.debug = flag.Bool("debug", false, "debug output")
	serve := flag.Bool("serve", false, "http serve")
	flag.Parse()

	setupKubernetesClient()

	if *serve {
		httpServe(":8080")
	} else {
		printVersions()
	}
}

