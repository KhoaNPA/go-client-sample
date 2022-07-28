# Part 1: Architectural Challenge
See answer in below:
```
$ cat./part-1-architecture-challenge.md
```

# Part 2: Technical Challenge: go-client-sample
Using k8s.io/go-client to interact with Kubernetes cluster.

```
% tree                         
.
├── README.md
├── batch-job
│   ├── go.mod
│   ├── go.sum
│   └── main.go
└── part-1-architecture-challenge.md
```

---
Schedule batch jobs
We are running a cluster with multiple running nodes and have a list of batch jobs containing multiple jobs that need to run. Each job is defined with specific resource requests and resource limits. 

Write a script using kubernetes client (prefer Golang but not required) to schedule the batchJob1 jobs on a minimum number of nodes, based on the 3 scenarios provided below.

Note that the goal is to always schedule the jobs on the minimum amount of nodes so that once the source code is only downloaded the least amount of times and mounted into as many jobs as possible to ensure the fastest possible execution time that also reduces the stress on the repoService.
Batch Job for Scenarios

{
 "batchJob1": [
     {
       "jobName": "job1",
       "image": "docker_img_1", #can change to nginx to test
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
}





Scenario 1

Name:               node1
Labels:                                                             
 kubernetes.io/hostname=ip-172-100-1-1.ec2.internal
Allocatable:
 cpu:                         1000m
 memory:                      4000Mi
Allocated resources:
 cpu                          1200m
 memory                       3000Mi

---

Name:               node2
Labels:    
 kubernetes.io/hostname=ip-172-100-1-2.ec2.internal
Allocatable:
 cpu:                         1000m
 memory:                      4000Mi
Allocated resources:
 cpu                          0m
 memory                       0Mi




In this scenario, all jobs should be scheduled on node2.

Scenario 2 

Name:               node1
Labels:                                                             
 kubernetes.io/hostname=ip-172-100-1-1.ec2.internal
Allocatable:
 cpu:                         1000m
 memory:                      4000Mi
Allocated resources:
 cpu                          0m
 memory                       00Mi

---

Name:               node2
Labels:    
 kubernetes.io/hostname=ip-172-100-1-2.ec2.internal
Allocatable:
 cpu:                         1000m
 memory:                      4000Mi
Allocated resources:
 cpu                          0m
 memory                       0Mi

In this scenario the jobs should be scheduled on either node1 or node 2, but it has to be scheduled on one node.



Scenario 3

Name:               node1
Labels:                                                             
 kubernetes.io/hostname=ip-172-100-1-1.ec2.internal
Allocatable:
 cpu:                         1000m
 memory:                      4000Mi
Allocated resources:
 cpu                          100m
 memory                       2500Mi

---

Name:               node2
Labels:    
 kubernetes.io/hostname=ip-172-100-1-2.ec2.internal
Allocatable:
 cpu:                         1000m
 memory:                      3000Mi
Allocated resources:
 cpu                          0m
 memory                       2000Mi

---

Name:               node2
Labels:    
 kubernetes.io/hostname=ip-172-100-1-2.ec2.internal
Allocatable:
 cpu:                         1000m
 memory:                      3000Mi
Allocated resources:
 cpu                          0m
 memory                       0Mi

In this scenario schedule job1 to node1 and the rest of the jobs to node3. There are different possible combinations, but the jobs should only run on 2 nodes.

---
How to run:

1. Start a Minikube cluster with 2 nodes (MacOS)
```
$ minikube start --driver=docker
$ minikube node add

$ kubectl get node
NAME           STATUS   ROLES    AGE   VERSION
minikube       Ready    master   17h   v1.18.3
minikube-m02   Ready    <none>   15h   v1.18.3
```

2. Create batch job in scenario 1 & 2:
- Uncomment the Function `launchJobScenario1n2(clientset, job.JobName, job.Image, job.RequestMem, job.RequestCpu)` in main() and run below command

- Comment the Function `launchJobScenario3(clientset, job.JobName, job.Image, job.RequestMem, job.RequestCpu)` in main() and run below command

```
$ cd ./batch-job
$ pwd 
~/go-client-sample/batch-job
$ go run main.go
$ kubectl get job
NAME   COMPLETIONS   DURATION   AGE
job1   0/1           5s         5s
job2   0/1           5s         5s
job3   0/1           5s         5s

$ kubectl get pod 
NAME         READY   STATUS    RESTARTS   AGE
job1-4s6wp   1/1     Running   0          15s
job2-98pns   1/1     Running   0          15s
job3-zm52w   1/1     Running   0          15s

$ kubectl get pods --all-namespaces -o wide --field-selector spec.nodeName=minikube
...

# Verify scenario
# minikube node (node 1) is control plan and doesnt have enough space for new jobs => all jobs jump to minikube-m02 and go together by PodAffinity policy
$ kubectl get pods --all-namespaces -o wide --field-selector spec.nodeName=minikube-m02
NAMESPACE     NAME               READY   STATUS    RESTARTS   AGE   IP           NODE           NOMINATED NODE   READINESS GATES
default       job1-7m6ql         1/1     Running   0          6s    172.18.0.2   minikube-m02   <none>           <none>
default       job2-jrf9l         1/1     Running   0          6s    172.18.0.3   minikube-m02   <none>           <none>
default       job3-j89c4         1/1     Running   0          6s    172.18.0.4   minikube-m02   <none>           <none>
kube-system   kindnet-hn9dx      1/1     Running   0          15h   172.17.0.4   minikube-m02   <none>           <none>
kube-system   kube-proxy-rfgvs   1/1     Running   0          15h   172.17.0.4   minikube-m02   <none>           <none>

# Clean up
$ kubectl delete job --all
job.batch "job1" deleted
job.batch "job2" deleted
job.batch "job3" deleted
khoanpa@khoanpas-MacBook-P
```

- Comment out the Function `launchJobScenario1n2(clientset, job.JobName, job.Image, job.RequestMem, job.RequestCpu)` again

3. Create the batch job in scenario 3
- Uncomment the Function `launchJobScenario3(clientset, job.JobName, job.Image, job.RequestMem, job.RequestCpu)` and run below command
```
Add 3rd node to cluster to simulate the scenario 3
$ pwd 
~/go-client-sample/batch-job

$ minikube node add
$ kubectl get node
NAME           STATUS   ROLES    AGE   VERSION
minikube       Ready    master   18h   v1.18.3
minikube-m02   Ready    <none>   15h   v1.18.3
minikube-m03   Ready    <none>   35s   v1.18.3

$ go run main.go
$ kubectl get job
NAME   COMPLETIONS   DURATION   AGE
job1   0/1           5s         5s
job2   0/1           5s         5s
job3   0/1           5s         5s

$ kubectl get pod 
NAME         READY   STATUS    RESTARTS   AGE
job1-ct74f   1/1     Running   0          12s
job2-6j5z5   1/1     Running   0          12s
job3-8gjb8   1/1     Running   0          12s

# Verify scenario
# job1 is running alone in a node
# job2 & job3 are running together in same node (different node with job1)
$ kubectl get pods --all-namespaces -o wide --field-selector spec.nodeName=minikube-m02
NAMESPACE     NAME               READY   STATUS    RESTARTS   AGE   IP           NODE           NOMINATED NODE   READINESS GATES
default       job2-6j5z5         1/1     Running   0          52s   172.18.0.2   minikube-m02   <none>           <none>
default       job3-8gjb8         1/1     Running   0          52s   172.18.0.3   minikube-m02   <none>           <none>
kube-system   kindnet-hn9dx      1/1     Running   0          15h   172.17.0.4   minikube-m02   <none>           <none>
kube-system   kube-proxy-rfgvs   1/1     Running   0          15h   172.17.0.4   minikube-m02   <none>           <none>

$ kubectl get pods --all-namespaces -o wide --field-selector spec.nodeName=minikube-m03
NAMESPACE     NAME               READY   STATUS    RESTARTS   AGE     IP           NODE           NOMINATED NODE   READINESS GATES
default       job1-ct74f         1/1     Running   0          76s     172.18.0.2   minikube-m03   <none>           <none>
kube-system   kindnet-qw52j      1/1     Running   0          3m36s   172.17.0.5   minikube-m03   <none>           <none>
kube-system   kube-proxy-bs45k   1/1     Running   0          3m36s   172.17.0.5   minikube-m03   <none>           <none>

$ kubectl delete job --all
```

4. Stop minikube cluster
```
$ minikube stop
```