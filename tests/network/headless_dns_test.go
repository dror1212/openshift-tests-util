package network_test

import (
	"myproject/framework"
	"myproject/util"
	"myproject/consts"
	. "github.com/onsi/ginkgo/v2"
	corev1 "k8s.io/api/core/v1"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service type Headless access from same namespace", func() {
	var (
		ctx           *framework.TestContext
		serverPodName string
		clientPodName string
		serviceName   string
		headlessDNS   string
		imageClient   = consts.ClientImage
		image         = consts.HttpdImage
	)

	BeforeEach(func() {
		// Initialize the TestContext and setup environment in the current namespace
		ctx = framework.Setup("core")

		// Generate names for the server pod, client pod, and service using the random name from context
		serverPodName = consts.TestPrefix + "-server-" + ctx.RandomName
		clientPodName = consts.TestPrefix + "-client-" + ctx.RandomName
		serviceName = consts.TestPrefix + "-headless-" + ctx.RandomName

		// Define the pod to be exposed by the Headless service
		containers := []util.ContainerConfig{
			util.CreateContainerConfig("test-container", image, nil, util.GenerateResourceRequirements("250m", "1000m", "1Gi", "1Gi")),
		}

		// Create the server pod
		ctx.CreateTestPodHelper(serverPodName, containers, 3)

		// Create a Headless service for the server pod
		servicePorts := []corev1.ServicePort{
			util.GeneratePort("http", 80, 80, "TCP"),
		}
		ctx.CreateServiceHelper(serviceName, "Headless", servicePorts, map[string]string{"app": serverPodName})

		// Fetch the DNS name dynamically from the service object
		var err error
		headlessDNS, err = util.GetServiceDNSName(ctx.KubeClient, ctx.Namespace, serviceName)
		Expect(err).ToNot(HaveOccurred(), "Failed to get DNS name for service %s", serviceName)
	})

	It("should allow access to the Headless service from the same namespace using DNS", func() {
		// Define the test pod that will access the headless service via DNS in the same namespace
		testContainers := []util.ContainerConfig{
			util.CreateContainerConfig("curl-container", imageClient, []string{"curl", "--fail", "--retry", "5", "-w", "HTTP Response Code: %{http_code}\n", "http://" + headlessDNS}, util.GenerateResourceRequirements("100m", "400m", "200Mi", "200Mi")),
		}

		// Create the test pod in the same namespace
		ctx.CreateTestPodHelper(clientPodName, testContainers, 3)

		// Verify DNS access to the headless service
		ctx.VerifyPodResponse(clientPodName, "HTTP Response Code: 200")
	})

	AfterEach(func() {
		// Clean up resources in both namespaces
		ctx.CleanupResource(serverPodName, "pod")
		ctx.CleanupResource(clientPodName, "pod")
		ctx.CleanupResource(serviceName, "service")
	})
})
