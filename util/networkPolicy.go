package util

import (
	"context"
	"fmt"

	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// CreateNetworkPolicyWithNamespaceAllow creates a NetworkPolicy that allows ingress traffic from other namespaces on a specified port.
func CreateNetworkPolicyWithNamespaceAllow(clientset *kubernetes.Clientset, namespace, policyName string, allowPorts []netv1.NetworkPolicyPort) (*netv1.NetworkPolicy, error) {
    policy := &netv1.NetworkPolicy{
        ObjectMeta: metav1.ObjectMeta{
            Name:      policyName,
            Namespace: namespace,
        },
        Spec: netv1.NetworkPolicySpec{
            PodSelector: metav1.LabelSelector{
                MatchLabels: map[string]string{}, // Apply to all pods in this namespace
            },
            PolicyTypes: []netv1.PolicyType{
                netv1.PolicyTypeIngress,
            },
            Ingress: []netv1.NetworkPolicyIngressRule{
                {
                    Ports: allowPorts, // Allow traffic only on the specified ports
                    From: []netv1.NetworkPolicyPeer{
                        {
                            // Allow traffic from all pods in other namespaces
                            NamespaceSelector: &metav1.LabelSelector{
                                MatchExpressions: []metav1.LabelSelectorRequirement{
                                    {
                                        Key:      "kubernetes.io/metadata.name",
                                        Operator: metav1.LabelSelectorOpNotIn,
                                        Values:   []string{namespace}, // Exclude the current namespace
                                    },
                                },
                            },
                        },
                    },
                },
            },
        },
    }

    // Create the NetworkPolicy in Kubernetes
    networkPolicy, err := clientset.NetworkingV1().NetworkPolicies(namespace).Create(context.TODO(), policy, metav1.CreateOptions{})
    if err != nil {
        LogError("Failed to create network policy: %v", err)
        return nil, err
    }

    LogInfo("Successfully created NetworkPolicy %s in namespace %s", policyName, namespace)
    return networkPolicy, nil
}

// DeleteNetworkPolicy deletes a NetworkPolicy
func DeleteNetworkPolicy(clientset *kubernetes.Clientset, namespace, policyName string) error {
	err := clientset.NetworkingV1().NetworkPolicies(namespace).Delete(context.TODO(), policyName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete network policy: %v", err)
	}
	return nil
}

// CreateNetworkPolicyPort defines a NetworkPolicyPort to restrict access to a specific port and protocol as a string
func CreateNetworkPolicyPort(port int32, protocolStr string) []netv1.NetworkPolicyPort {
	var protocol corev1.Protocol

	// Convert string protocol to corev1.Protocol type
	switch protocolStr {
	case "TCP":
		protocol = corev1.ProtocolTCP
	case "UDP":
		protocol = corev1.ProtocolUDP
	default:
		protocol = corev1.ProtocolTCP // Default to TCP if no match
	}

	return []netv1.NetworkPolicyPort{
		{
			Port:     &intstr.IntOrString{IntVal: port},
			Protocol: &protocol,
		},
	}
}