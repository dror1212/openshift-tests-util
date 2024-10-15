package framework

import (
	"time"
	"strings"
	"myproject/util"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"fmt"
)

// CreateTestPodHelper creates a test pod with a retry mechanism
func (ctx *TestContext) CreateTestPodHelper(podName string, containers []util.ContainerConfig, retries int) {
	util.LogInfo("Creating test pod %s with retry mechanism", podName)
	_, err := util.RetryPodCreationWithWait(ctx.KubeClient, ctx.Config, ctx.Namespace, podName, containers, nil, 15*time.Second, 5*time.Minute, retries)
	if err != nil {
		util.LogError("Failed to create test pod %s: %v", podName, err)
		Expect(errors.Wrap(err, "failed to create test client pod after retries")).ToNot(HaveOccurred())
	}
	util.LogInfo("Successfully created test pod %s", podName)
}

// CreateTestPodExpectingFailureHelper creates a test pod and expects it to fail (e.g., due to NetworkPolicy restrictions).
func (ctx *TestContext) CreateTestPodExpectingFailureHelper(podName string, containers []util.ContainerConfig, retries int) {
	util.LogInfo("Creating test pod %s, expecting failure", podName)
	_, err := util.RetryPodCreationWithWait(ctx.KubeClient, ctx.Config, ctx.Namespace, podName, containers, nil, 10*time.Second, 1*time.Minute, retries)
	if err != nil {
		util.LogInfo("Pod %s failed as expected", podName)
	} else {
		util.LogError("Pod %s did not fail as expected", podName)
		Expect(errors.Wrap(err, "did not fail fail to create test client pod after retries")).ToNot(HaveOccurred())
	}
}

// WaitForPodAndCheckLogs waits for the pod to complete and checks logs for a substring.
func (ctx *TestContext) WaitForPodAndCheckLogs(podName, logSubstring string, checkInterval, timeout time.Duration, retries int) error {
	util.LogInfo("Waiting for pod %s to complete and checking logs for substring: %s", podName, logSubstring)

	// Wait for the pod to complete (successfully or with failure)
	err := util.WaitForPodState(ctx.KubeClient, ctx.Namespace, podName, checkInterval, timeout, retries, false)
	if err != nil {
		util.LogError("Failed to wait for pod %s completion/failure: %v", podName, err)
		return err
	}

	// Fetch the logs once the pod has completed
	podLogs, err := util.GetPodLogs(ctx.KubeClient, ctx.Namespace, podName)
	if err != nil {
		util.LogError("Failed to fetch logs for pod %s: %v", podName, err)
		return err
	}

	// Check if the logs contain the expected substring
	util.LogInfo("Fetched logs for pod %s. Checking for substring: %s", podName, logSubstring)
	if !strings.Contains(podLogs, logSubstring) {
		util.LogError("Expected to find %s in pod logs for pod %s, but did not", logSubstring, podName)
		return fmt.Errorf("substring %s not found in pod logs for pod %s", logSubstring, podName)
	}

	util.LogInfo("Successfully found log substring %s in pod %s", logSubstring, podName)
	return nil
}

// VerifyPodResponse verifies that the pod logs contain the expected response, with retries.
func (ctx *TestContext) VerifyPodResponse(podName, expectedResponse string, retries int) {
	err := ctx.WaitForPodAndCheckLogs(podName, expectedResponse, 10*time.Second, 3*time.Minute, retries)
	Expect(err).ToNot(HaveOccurred(), "Pod %s returned a different response", podName)
}