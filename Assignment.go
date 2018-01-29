package main

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
)
var (
	master         = flag.String("master", "", "Master URL to build a client config from. Either this or kubeconfig needs to be set ")
	kubeconfig     = flag.String("kubeconfig", "", "Absolute path to the kubeconfig file. Either this or master needs to be set ")
	namespace      = flag.String("namespace", "default", "Namespace from which you Want to list pods")
)
func main() {
	flag.Parse()

	// Create the client according to whether we are running in or out-of-cluster
	outOfCluster := *master != "" || *kubeconfig != ""
	var config *rest.Config
	/*fmt.Println("Namespace ",*namespace)
	fmt.Println("Master ",*master)
	fmt.Println("Kubeconfig ",*kubeconfig)*/
	var err error
	if outOfCluster {
		config, err = clientcmd.BuildConfigFromFlags(*master, *kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		fmt.Errorf("Failed to create config: %v", err)
	}
	config.Insecure =true
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Errorf("Failed to create client: %v", err)
	}
	pods, err := clientset.CoreV1().Pods(*namespace).List(metav1.ListOptions{})
	if err != nil {
		fmt.Println(err.Error())
	}
	if len(pods.Items) > 0 {
		for _, pod := range pods.Items {
			fmt.Printf("Pod %s\n", pod.GetName())
			fmt.Println("SR.NO \t container name \t cpulimit \t cpurequest \t memorylimit \t memoryequest")
			//container wise resource info
			for i,c:= range pod.Spec.Containers{
				limit:=	c.Resources.Limits
				req:=c.Resources.Requests
				fmt.Println(i,"\t",c.Name,"\t",  limit[v1.ResourceCPU],"\t", req[v1.ResourceCPU],"\t", limit[v1.ResourceMemory] ,"\t",req[v1.ResourceMemory])
			}
		}
	} else {
		fmt.Println("No pods in given namespace")
	}
}
