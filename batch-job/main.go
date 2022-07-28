package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	v1 "k8s.io/api/core/v1"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type MyJob struct {
	JobName    string
	Image      string
	RequestMem string
	RequestCpu string
}

func get_kube_config_path() string {
	var kube_config_path string
	home_dir := homedir.HomeDir()

	if _, err := os.Stat(home_dir + "/.kube/config"); err == nil {
		kube_config_path = home_dir + "/.kube/config"
	} else {
		fmt.Println("Enter kubernetes config directory: ")
		fmt.Scanf("%s", kube_config_path)
	}

	return kube_config_path
}

func main() {

	batchJob1Json := `[
		{
			"jobName": "job1",
			"image": "docker_img_1",
			"requestMem": "500Mi",
			"requestCpu": "200m"
		},
		{
			"jobName": "job2",
			"image": "docker_img_2",
			"requestMem": "1Gi",
			"requestCpu": "100m"
		},
		{
			"jobName": "job3",
			"image": "docker_img_3",
			"requestMem": "2Gi",
			"requestCpu": "200m"
		}
	]
	`

	// Test Batch Job in Scenario
	// batchJob1Json := `[
	// 	{
	// 		"jobName": "job1",
	// 		"image": "nginx",
	// 		"requestMem": "250Mi",
	// 		"requestCpu": "200m"
	// 	},
	// 	{
	// 		"jobName": "job2",
	// 		"image": "nginx",
	// 		"requestMem": "500Mi",
	// 		"requestCpu": "100m"
	// 	},
	// 	{
	// 		"jobName": "job3",
	// 		"image": "nginx",
	// 		"requestMem": "1Gi",
	// 		"requestCpu": "200m"
	// 	}
	// ]
	// `

	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		rules,
		&clientcmd.ConfigOverrides{},
	)

	config, err := kubeconfig.ClientConfig()

	if err != nil {
		panic(err)
	}

	clientset := kubernetes.NewForConfigOrDie(config)

	var batchJob1 []MyJob
	json.Unmarshal([]byte(batchJob1Json), &batchJob1)

	// Launch Batch Job in Scenario 1
	for _, job := range batchJob1 {
		fmt.Printf("%+v\n", job)

		// Launch Batch Job in Scenario 1 & 2
		// launchJobScenario1n2(clientset, job.JobName, job.Image, job.RequestMem, job.RequestCpu)

		// Launch Batch Job in Scenario 3
		launchJobScenario3(clientset, job.JobName, job.Image, job.RequestMem, job.RequestCpu)
	}

}

// CREATE Batch Job in Scenario 1 & 2
func launchJobScenario1n2(clientset *kubernetes.Clientset, jobName string, image string, requestMem string, requestCpu string) {
	jobs := clientset.BatchV1().Jobs("default")
	var backOffLimit int32 = 0

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: "default",
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Affinity: &v1.Affinity{
						PodAffinity: &v1.PodAffinity{
							PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
								{
									Weight: 100,
									PodAffinityTerm: v1.PodAffinityTerm{
										TopologyKey: "kubernetes.io/hostname",
										LabelSelector: &metav1.LabelSelector{
											MatchExpressions: []metav1.LabelSelectorRequirement{
												{
													Key:      "job-name",
													Operator: "In",
													Values:   []string{"job"},
												},
											},
										},
									},
								},
							},
						},
					},
					Containers: []v1.Container{
						{
							Name:            jobName,
							Image:           image,
							ImagePullPolicy: "IfNotPresent",
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									"cpu":    resource.MustParse(requestCpu),
									"memory": resource.MustParse(requestMem),
								},
								Limits: v1.ResourceList{
									"cpu":    resource.MustParse(requestCpu),
									"memory": resource.MustParse(requestMem),
								},
							},
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

	_, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln("FAILED to Launch Batch Job in Scenario 1 & 2")
	}

	//print job details
	log.Println("SUCCESSFULLY Launch Batch Job in Scenario 1 & 2")
}

// CREATE Batch Job in Scenario 3
func launchJobScenario3(clientset *kubernetes.Clientset, jobName string, image string, requestMem string, requestCpu string) {
	jobs := clientset.BatchV1().Jobs("default")
	var backOffLimit int32 = 0

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: "default",
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Affinity: &v1.Affinity{
						PodAntiAffinity: &v1.PodAntiAffinity{
							// IMPORTANT: USE PreferredDuringSchedulingIgnoredDuringExecution rather than RequiredDuringSchedulingIgnoredDuringExecution
							// PreferScheduling allows to deploy brand jobs, in otherhand RequireScheduling forces that there's available jobs to have condition
							PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
								{
									Weight: 100,
									PodAffinityTerm: v1.PodAffinityTerm{
										TopologyKey: "kubernetes.io/hostname",
										LabelSelector: &metav1.LabelSelector{
											MatchExpressions: []metav1.LabelSelectorRequirement{
												{
													Key:      "job-name",
													Operator: "In",
													Values:   []string{"job1"},
												},
											},
										},
									},
								},
							},
						},
						PodAffinity: &v1.PodAffinity{
							PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
								{
									Weight: 100,
									PodAffinityTerm: v1.PodAffinityTerm{
										TopologyKey: "kubernetes.io/hostname",
										LabelSelector: &metav1.LabelSelector{
											MatchExpressions: []metav1.LabelSelectorRequirement{
												{
													Key:      "job-name",
													Operator: "In",
													Values:   []string{"job2", "job3"},
												},
											},
										},
									},
								},
							},
						},
					},
					Containers: []v1.Container{
						{
							Name:            jobName,
							Image:           image,
							ImagePullPolicy: "IfNotPresent",
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									"cpu":    resource.MustParse(requestCpu),
									"memory": resource.MustParse(requestMem),
								},
								Limits: v1.ResourceList{
									"cpu":    resource.MustParse(requestCpu),
									"memory": resource.MustParse(requestMem),
								},
							},
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

	_, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln("FAILED to Launch Batch Job in Scenario 3.")
	}

	//print job details
	log.Println("SUCCESSFULLY Launch Batch Job in Scenario 3")
}
