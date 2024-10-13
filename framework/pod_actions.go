package framework

import (
	"time"
	"strings"
	"myproject/util"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

// CreateTestPodHelper creates a test pod with a retry mechanism
func (ctx *TestContext) CreateTestPodHelper(podName string, containers []util.ContainerConfig, retries int) {
	util.LogInfo("Creating test pod %s with retry mechanism", podName)
	_, err := util.RetryPodCreationWithWait(ctx.KubeClient, ctx.Config, ctx.Namespace, podName, containers, nil, retries, 15*time.Second, 5*time.Minute)
	if err != nil {
		util.LogError("Failed to create test pod %s: %v", podName, err)
		Expect(errors.Wrap(err, "failed to create test client pod after retries")).ToNot(HaveOccurred())
	}
	util.LogInfo("Successfully created test pod %s", podName)
}

// WaitForPodAndCheckLogs waits for the pod to complete and checks logs for a substring
func (ctx *TestContext) WaitForPodAndCheckLogs(podName, logSubstring string, checkInterval, timeout time.Duration) error {
	util.LogInfo("Waiting for pod %s to complete and checking logs for substring: %s", podName, logSubstring)
	Eventually(func() (bool, error) {
		err := util.WaitForPodCompletionOrFailure(ctx.KubeClient, ctx.Namespace, podName, checkInterval, timeout)
		if err != nil {
			util.LogError("Failed to wait for pod %s completion: %v", podName, err)
			return false, err
		}

		podLogs, err := util.GetPodLogs(ctx.KubeClient, ctx.Namespace, podName)
		if err != nil {
			util.LogError("Failed to fetch logs for pod %s: %v", podName, err)
			return false, err
		}

		util.LogInfo("Fetched logs for pod %s. Checking for substring: %s", podName, logSubstring)
		return strings.Contains(podLogs, logSubstring), nil
	}, timeout, checkInterval).Should(BeTrue(), "Expected to find %s in pod logs for pod %s", logSubstring, podName)

	util.LogInfo("Successfully found log substring %s in pod %s", logSubstring, podName)
	return nil
}

// VerifyPodResponse verifies that the pod logs contain the expected response
func (ctx *TestContext) VerifyPodResponse(podName, expectedResponse string) {
	err := ctx.WaitForPodAndCheckLogs(podName, expectedResponse, 5*time.Second, 5*time.Minute)
	Expect(err).ToNot(HaveOccurred(), "Pod %s returned different response", podName)
}
