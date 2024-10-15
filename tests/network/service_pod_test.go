package network_test

import (
	"time"
	"myproject/framework"
	"myproject/util"
	"myproject/consts"
	. "github.com/onsi/ginkgo/v2"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("Service type LoadBalancer on Pod", func() {
	var (
		ctx           *framework.TestContext
		serverPodName string
		clientPodName string
		serviceName   string
		imageClient = consts.ClientImage
		image       = consts.HttpdImage
		serviceIP     string
	)

	BeforeEach(func() {
		// Initialize the TestContext and setup environment
		ctx = framework.Setup("core")

		// Generate names for the pod, test pod, and service using the random name from context
		serverPodName = consts.TestPrefix + "-server-" + ctx.RandomName
		clientPodName = consts.TestPrefix + "-client-" + ctx.RandomName
		serviceName = consts.TestPrefix + "-lb-" + ctx.RandomName

		// Define the pod to be exposed by the LoadBalancer
		containers := []util.ContainerConfig{
			util.CreateContainerConfig("test-container", image, nil, util.GenerateResourceRequirements("250m", "1000m", "1Gi", "1Gi")),
		}

		// Create the main test pod
		ctx.CreateTestPodHelper(serverPodName, containers, 3)

		// Create a LoadBalancer service for the pod
		servicePorts := []corev1.ServicePort{
			util.GeneratePort("http", 80, 80, "TCP"),
		}
		ctx.CreateServiceHelper(serviceName, "LoadBalancer", servicePorts, map[string]string{"app": serverPodName})
	})

	It("should expose the service with a LoadBalancer and allow access to a pod from another pod", func() {
		
		// Wait for the service to get an external IP
		serviceIP = ctx.WaitForServiceIP(serviceName, 2*time.Minute, 10*time.Second)

		// Define the test pod that will access the service
		testContainers := []util.ContainerConfig{
			util.CreateContainerConfig("curl-container", imageClient, []string{"curl", "--fail", "--retry", "5", "-w", "HTTP Response Code: %{http_code}\n", "http://" + serviceIP}, util.GenerateResourceRequirements("100m", "400m", "200Mi", "200Mi")),
		}

		// Create the test pod using the retry mechanism
		ctx.CreateTestPodHelper(clientPodName, testContainers, 3)

		// Wait for the test pod to complete and verify its status
		ctx.VerifyPodResponse(clientPodName, "HTTP Response Code: 200")
	})

	AfterEach(func() {
		// Clean up resources: Delete the test pods and the service
		ctx.CleanupResource(serverPodName, "pod")
		ctx.CleanupResource(clientPodName, "pod")
		ctx.CleanupResource(serviceName, "service")
	})
})
