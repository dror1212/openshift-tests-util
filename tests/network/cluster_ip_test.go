package network_test

import (
	"time"
	"myproject/framework"
	"myproject/util"
	"myproject/consts"
	. "github.com/onsi/ginkgo/v2"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("Service type ClusterIP access from same namespace", func() {
	var (
		ctx           *framework.TestContext
		serverPodName string
		clientPodName string
		serviceName   string
		serviceIP     string
		imageClient = consts.ClientImage
		image       = consts.HttpdImage
	)

	BeforeEach(func() {
		// Initialize the TestContext and setup environment in the current namespace
		ctx = framework.Setup("core")

		// Generate names for the pod, test pod, and service using the random name from context
		serverPodName = consts.TestPrefix + "-server-" + ctx.RandomName
		clientPodName = consts.TestPrefix + "-client-" + ctx.RandomName
		serviceName = consts.TestPrefix + "-clusterip-" + ctx.RandomName

		// Define the pod to be exposed by the ClusterIP service
		containers := []util.ContainerConfig{
			util.CreateContainerConfig("test-container", image, nil, util.GenerateResourceRequirements("250m", "1000m", "1Gi", "1Gi")),
		}

		// Create the main test pod
		ctx.CreateTestPodHelper(serverPodName, containers, 3)

		// Create a ClusterIP service for the pod
		servicePorts := []corev1.ServicePort{
			util.GeneratePort("http", 80, 80, "TCP"),
		}
		ctx.CreateServiceHelper(serviceName, "ClusterIP", servicePorts, map[string]string{"app": serverPodName})
	})

	It("should allow access to the ClusterIP service from the same namespace", func() {
		serviceIP = ctx.WaitForServiceIP(serviceName, 2*time.Minute, 10*time.Second)
	
		// Log the retrieved ClusterIP
		util.LogInfo("ClusterIP for service %s: %s", serviceName, serviceIP)
	
		// Define the test pod that will access the service in the same namespace
		testContainers := []util.ContainerConfig{
			util.CreateContainerConfig("curl-container", imageClient, []string{"curl", "--fail", "--retry", "5", "-w", "HTTP Response Code: %{http_code}\n", "http://" + serviceIP}, util.GenerateResourceRequirements("100m", "400m", "200Mi", "200Mi")),
		}
	
		// Create the test pod in the same namespace
		ctx.CreateTestPodHelper(clientPodName, testContainers, 3)
	
		// Verify access to the service
		ctx.VerifyPodResponse(clientPodName, "HTTP Response Code: 200")
	})

	AfterEach(func() {
		// Clean up resources in both namespaces
		ctx.CleanupResource(serverPodName, "pod")
		ctx.CleanupResource(clientPodName, "pod")
		ctx.CleanupResource(serviceName, "service")
	})
})
