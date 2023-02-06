package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func cloudNativeSort(N []int, clientset *kubernetes.Clientset) []int {
	if len(N) < 2 {
		return N
	} else {
		mid := len(N) / 2
		/*
			Here, instead of calling cloudNativeSort() recursively,
			we start a new Job, wait for it to finish, and then read its logs
			to see the result
		*/

		/* just a little test */
		jobsClient := clientset.BatchV1().Jobs(apiv1.NamespaceDefault)
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		left := cloudNativeSort(N[:mid], clientset)
		right := cloudNativeSort(N[mid:], clientset)
		return cloudNativeMerge(left, right)
	}

}

func cloudNativeMerge(M []int, L []int) []int {
	N := make([]int, len(M)+len(L))
	m := 0
	l := 0
	i := 0
	for m < len(M) && l < len(L) {
		if M[m] <= L[l] {
			N[i] = M[m]
			m += 1
		} else {
			N[i] = L[l]
			l += 1
		}
		i += 1
	}

	for m < len(M) {
		N[i] = M[m]
		m += 1
		i += 1
	}

	for l < len(L) {
		N[i] = L[l]
		l += 1
		i += 1
	}

	return N
}

func main() {
	var N = make([]int, len(os.Args)-1)

	for i, v := range os.Args {
		if i == 0 {
			continue
		}
		integer, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		N[i-1] = integer
	}

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	N = cloudNativeSort(N, clientset)
	for _, v := range N {
		fmt.Printf("%d ", v)
	}
	fmt.Println()
}
