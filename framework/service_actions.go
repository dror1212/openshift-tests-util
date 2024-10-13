package framework

import (
	"time"
	corev1 "k8s.io/api/core/v1"
	. "github.com/onsi/gomega"
	"myproject/util"
)

// CreateServiceHelper creates a Kubernetes service of a specified type (ClusterIP, LoadBalancer, etc.)
func (ctx *TestContext) CreateServiceHelper(serviceName string, serviceType corev1.ServiceType, servicePorts []corev1.ServicePort, labels map[string]string) {
	_, err := util.CreateService(ctx.KubeClient, ctx.Namespace, serviceName, serviceType, servicePorts, labels)
	Expect(err).ToNot(HaveOccurred(), "Failed to create service %s of type %s", serviceName, serviceType)
}

// WaitForServiceIP waits for a service to get a valid IP within a specified timeout
func (ctx *TestContext) WaitForServiceIP(serviceName string, timeout, interval time.Duration) string {
	var serviceIP string

	Eventually(func() (string, error) {
		var err error
		serviceIP, err = util.GetServiceIP(ctx.KubeClient, ctx.Namespace, serviceName)
		return serviceIP, err
	}, timeout, interval).ShouldNot(BeEmpty(), "Expected service to get a service IP")

	return serviceIP
}