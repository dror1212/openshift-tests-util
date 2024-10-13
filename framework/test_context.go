package framework

import (
	"myproject/util"
	. "github.com/onsi/gomega"
	"github.com/openshift/client-go/route/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	kubecli "kubevirt.io/client-go/kubecli"
)

// TestContext holds reusable values across tests, including a random name for resources
type TestContext struct {
	Config      *rest.Config
	KubeClient   *kubernetes.Clientset
	VirtClient  kubecli.KubevirtClient
	RouteClient *versioned.Clientset
	Namespace   string
	RandomName  string
}

// Setup initializes the environment (e.g., auth, logging) and sets the random name for each test
func Setup(namespace string) *TestContext {
	var err error

	err = util.SetLogLevel("debug")
	Expect(err).ToNot(HaveOccurred(), "Failed to initiate the logger")

	kubeclient, config, err := util.Authenticate()
	Expect(err).ToNot(HaveOccurred(), "Failed to authenticate with Kubernetes")

	virtClient, err := kubecli.GetKubevirtClientFromRESTConfig(config)
	Expect(err).ToNot(HaveOccurred(), "Failed to authenticate with KubeVirt")

	routeClient, err := versioned.NewForConfig(config)
	Expect(err).ToNot(HaveOccurred(), "Failed to create Route client")

	randomName := util.GenerateRandomName()

	return &TestContext{
		KubeClient:   kubeclient,
		Config:      config,
		VirtClient:  virtClient,
		RouteClient: routeClient,
		Namespace:   namespace,
		RandomName:  randomName,
	}
}