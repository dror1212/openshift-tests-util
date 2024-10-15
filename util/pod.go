package util

import (
	"context"
	"time"
	"bytes"
	"io"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"myproject/consts"
)

// ContainerConfig defines the configuration for a container in a Pod
type ContainerConfig struct {
	Name      string
	Image     string
	Command   []string
	Args      []string
	Resources corev1.ResourceRequirements
}

// CreatePod creates a Pod with multiple containers specified by the ContainerConfig list.
func CreatePod(config *rest.Config, namespace, podName string, containerConfigs []ContainerConfig, labels map[string]string, waitForCreation bool) (*corev1.Pod, error) {
	if podName == "" {
		podName = GenerateRandomName()
		LogInfo("Generated random Pod name: %s", podName)
	}

	if labels == nil {
		labels = consts.DefaultLabels
		labels["app"] = podName
		LogInfo("Using default labels for VM: %s", podName)
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		LogError("Failed to create Kubernetes client: %v", err)
		return nil, err
	}

	// Generate containers from the container configurations
	containers := []corev1.Container{}
	for _, config := range containerConfigs {
		containers = append(containers, GenerateContainerFromConfig(config))
	}

	// Define the Pod object
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers:    containers,
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	LogInfo("Starting Pod %s creation", podName)

	// Create the Pod in Kubernetes
	createdPod, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		LogError("Failed to create Pod: %v", err)
		return nil, err
	}

	LogInfo("Pod %s created", podName)

	// Wait for the Pod to be ready, if requested
	if waitForCreation {
		err = WaitForPodState(clientset, namespace, podName, 10*time.Second, 120*time.Second, 0, true)
		if err != nil {
			LogError("Error waiting for Pod to be ready: %v", err)
			return nil, err
		}
		LogInfo("Pod %s is ready", podName)
	}

	return createdPod, nil
}

// GetPod fetches the pod object from the Kubernetes API
func GetPod(clientset *kubernetes.Clientset, namespace, podName string) (*corev1.Pod, error) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		LogError("Failed to fetch Pod %s in namespace %s: %v", podName, namespace, err)
		return nil, fmt.Errorf("failed to get pod %s: %v", podName, err)
	}

	LogInfo("Fetched Pod %s in namespace %s", podName, namespace)
	return pod, nil
}

// RetryPodCreationWithWait creates a pod and waits for it to either run or complete successfully.
func RetryPodCreationWithWait(clientset *kubernetes.Clientset, config *rest.Config, namespace, podName string, containers []ContainerConfig, labels map[string]string, interval, timeout time.Duration, retries int) (*corev1.Pod, error) {
	var createdPod *corev1.Pod
	var err error

	// Define the pod creation and check function
	createAndCheckPod := func() (bool, error) {
		// Try creating the pod
		LogInfo("Start Pod %s Creation", podName)
		createdPod, err = CreatePod(config, namespace, podName, containers, labels, false)
		if err != nil {
			LogWarn("Failed to create pod %s. Error: %v", podName, err)
			return false, err
		}

		LogInfo("Pod %s created successfully.", podName)

		// Wait for the pod to reach the Running state or succeed
		err = WaitForPodState(clientset, namespace, podName, interval, timeout, 0, true)
		if err == nil {
			LogInfo("Pod %s is running or completed successfully.", podName)
			return true, nil
		}

		LogWarn("Pod %s failed or did not reach running/completed state. Error: %v", podName, err)

		// Delete the pod if it fails
		delErr := clientset.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
		if delErr != nil {
			LogError("Failed to delete pod %s after failure: %v", podName, delErr)
		}

		return false, err
	}

	// Use WaitFor to try creating the pod and waiting for it to reach running state
	err = WaitFor(createAndCheckPod, interval, timeout, retries)
	if err != nil {
		LogError("Failed to create pod %s", podName, err)
		return nil, err
	}

	return createdPod, nil
}

// GetPodLogs fetches the logs from a specific pod
func GetPodLogs(clientset *kubernetes.Clientset, namespace, podName string) (string, error) {
    podLogOpts := corev1.PodLogOptions{}
    podLogRequest := clientset.CoreV1().Pods(namespace).GetLogs(podName, &podLogOpts)
    podLogs, err := podLogRequest.Stream(context.TODO())
    if err != nil {
        return "", err
    }
    defer podLogs.Close()

    buf := new(bytes.Buffer)
    _, err = io.Copy(buf, podLogs)
    if err != nil {
        return "", err
    }
    return buf.String(), nil
}

// CreateContainerConfig creates a container configuration for a pod
func CreateContainerConfig(name, image string, command []string, resources corev1.ResourceRequirements) ContainerConfig {
	return ContainerConfig{
		Name:      name,
		Image:     image,
		Command:   command,
		Resources: resources,
	}
}

// GenerateContainerFromConfig creates a container spec from the given ContainerConfig
func GenerateContainerFromConfig(config ContainerConfig) corev1.Container {
	return corev1.Container{
		Name:      config.Name,
		Image:     config.Image,
		Command:   config.Command,
		Args:      config.Args,
		Resources: config.Resources,
	}
}
