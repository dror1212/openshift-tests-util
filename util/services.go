package util

import (
	"context"
	"time"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// GeneratePort generates a ServicePort object based on the provided values
func GeneratePort(name string, port, targetPort int, protocol string) corev1.ServicePort {
	var protocolType corev1.Protocol
	switch protocol {
	case "TCP":
		protocolType = corev1.ProtocolTCP
	case "UDP":
		protocolType = corev1.ProtocolUDP
	default:
		protocolType = corev1.ProtocolTCP // Default to TCP if not recognized
	}

	return corev1.ServicePort{
		Name:       name,
		Port:       int32(port),
		TargetPort: intstr.FromInt(targetPort),
		Protocol:   protocolType,
	}
}

// CreateService creates a Kubernetes service of a specified type (ClusterIP, NodePort, LoadBalancer, or Headless).
func CreateService(clientset *kubernetes.Clientset, namespace, serviceName string, serviceType string, ports []corev1.ServicePort, labels map[string]string) (*corev1.Service, error) {
	// Set default labels if not provided
	if labels == nil {
		labels = map[string]string{
			"app": serviceName,
		}
	}

	// Define the service object
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports:    ports,
			Selector: labels,
		},
	}

	// Handle the service type
	switch serviceType {
	case "Headless":
		service.Spec.ClusterIP = "None"
		LogInfo("Creating headless service %s", serviceName)
	case "ClusterIP":
		service.Spec.Type = corev1.ServiceTypeClusterIP
		LogInfo("Creating ClusterIP service %s", serviceName)
	case "NodePort":
		service.Spec.Type = corev1.ServiceTypeNodePort
		LogInfo("Creating NodePort service %s", serviceName)
	case "LoadBalancer":
		service.Spec.Type = corev1.ServiceTypeLoadBalancer
		LogInfo("Creating LoadBalancer service %s", serviceName)
	default:
		LogError("Unsupported service type: %s", serviceType)
		return nil, fmt.Errorf("unsupported service type: %s", serviceType)
	}

	// Create the service in Kubernetes
	service, err := clientset.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		LogError("Failed to create service %s: %v", serviceName, err)
		return nil, err
	}

	// If the service type is LoadBalancer, wait for the external IP to be assigned
	if serviceType == "LoadBalancer" {
		LogInfo("Waiting for the LoadBalancer service %s to get an external IP...", serviceName)
		err := WaitForServiceReady(clientset, namespace, serviceName, 5*time.Second, 120*time.Second)
		if err != nil {
			LogError("Error waiting for LoadBalancer service %s to be ready: %v", serviceName, err)
			return nil, err
		}

		// Fetch the service again to get the updated external IP
		service, err = clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
		if err != nil {
			LogError("Failed to get service %s after waiting for external IP: %v", serviceName, err)
			return nil, err
		}

		LogInfo("Service %s is ready with external IP: %s", serviceName, service.Status.LoadBalancer.Ingress[0].IP)
	}

	LogInfo("Service %s of type %s created successfully", serviceName, serviceType)
	return service, nil
}

// GetServiceIP fetches the IP of a service, handling both LoadBalancer and ClusterIP types
func GetServiceIP(clientset *kubernetes.Clientset, namespace, serviceName string) (string, error) {
	// Fetch the service object
	service, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		LogError("Failed to retrieve service %s: %v", serviceName, err)
		return "", err
	}

	// Handle LoadBalancer service: return external IP
	if service.Spec.Type == corev1.ServiceTypeLoadBalancer {
		if len(service.Status.LoadBalancer.Ingress) > 0 {
			return service.Status.LoadBalancer.Ingress[0].IP, nil
		}
		LogError("LoadBalancer service %s has no external IP assigned", serviceName)
		return "", nil
	}

	// Handle ClusterIP service: return ClusterIP
	if service.Spec.Type == corev1.ServiceTypeClusterIP {
		LogInfo("Service %s has a ClusterIP: %s", serviceName, service.Spec.ClusterIP)
		return service.Spec.ClusterIP, nil
	}

	// Handle unsupported service types
	LogError("Unsupported service type %s for service %s", service.Spec.Type, serviceName)
	return "", nil
}

// GetServiceDNSName retrieves the DNS name of a service in the cluster dynamically
func GetServiceDNSName(clientset *kubernetes.Clientset, namespace, serviceName string) (string, error) {
	// Fetch the service from the Kubernetes cluster
	service, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		LogError("Failed to get service %s: %v", serviceName, err)
		return "", fmt.Errorf("failed to get service %s: %v", serviceName, err)
	}

	// Construct the DNS name dynamically
	dnsName := fmt.Sprintf("%s.%s.svc.cluster.local", service.Name, service.Namespace)
	LogInfo("Service DNS Name: %s", dnsName)

	return dnsName, nil
}