/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	_ "bufio"
	_ "flag"
	"fmt"
	_ "os"
	"time"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	_ "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	_ "k8s.io/client-go/tools/clientcmd"
)

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 获取namespace
	namespaces, err := clientset.Core().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d namespaces in the cluster\n", len(namespaces.Items))
	for _, namespace := range namespaces.Items {
		fmt.Println("namespaces name: %s ", namespace.Name)
	}

	services, err := clientset.CoreV1().Services("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d services in the cluster\n", len(services.Items))
	for _, service := range services.Items {
		fmt.Println("services name: %s ", service.Name)
	}

	deps, err := clientset.AppsV1beta1().Deployments("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d deps in the cluster\n", len(deps.Items))
	for _, dep := range deps.Items {
		fmt.Println("deployment name : %s ", dep.Name)
	}

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	for _, pod := range pods.Items {
		fmt.Println("pod name: %s ", pod.Name)
	}

	deploymentsClient := clientset.AppsV1beta1().Deployments("yifan")

	deployment := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"io.daocloud.dce.template": "e2e-fabric",
						"io.daocloud.dce.app":      "e2e-fabric",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.13",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
							ImagePullPolicy: apiv1.PullAlways,
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	servicesClient := clientset.CoreV1().Services("yifan")
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
			Labels: map[string]string{
				"io.daocloud.dce.template": "e2e-fabric",
				"io.daocloud.dce.app":      "e2e-fabric",
			},
		},

		Spec: apiv1.ServiceSpec{
			Type: apiv1.ServiceTypeClusterIP,
			Ports: []apiv1.ServicePort{
				{
					Name: "7051",
					Port: 7051,
				},
				{
					Name: "7053",
					Port: 7053,
				},
			},
			Selector: map[string]string{
				"service": "demo-deployment",
			},
		},
	}
	fmt.Println("Creating service...")
	result2, err := servicesClient.Create(service)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created service %q.\n", result2.GetObjectMeta().GetName())

	time.Sleep(10 * time.Second)

	deps2, err := clientset.AppsV1beta1().Deployments("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d deps in the cluster\n", len(deps2.Items))
	for _, dep := range deps2.Items {
		fmt.Println("deployment name : %s ", dep.Name)
	}
	/*
		// 创建pod
		pod := new(v1.Pod)
		pod.TypeMeta = unversioned.TypeMeta{Kind: "Pod", APIVersion: "v1"}
		pod.ObjectMeta = v1.ObjectMeta{Name: "testapi", Namespace: "yifan", Labels: map[string]string{"name": "testapi"}}
		pod.Spec = v1.PodSpec{
			RestartPolicy: v1.RestartPolicyAlways,
			Containers: []v1.Container{
				v1.Container{
					Name:  "testapi",
					Image: "nginx",
					Ports: []v1.ContainerPort{
						v1.ContainerPort{
							ContainerPort: 80,
							Protocol:      v1.ProtocolTCP,
						},
					},
				},
			},
		}
		_, err = clientset.Core().Pods("yifan").Create(pod)
		if err != nil {
			panic(err.Error())
		}
		// 获取现有的pod数量
		pods, err := clientset.Core().Pods("").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items)) */
}

func int32Ptr(i int32) *int32 { return &i }
