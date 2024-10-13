package framework

import (
	"context"
	"myproject/util"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Cleanup cleans up resources such as VMs and Pods
func (ctx *TestContext) CleanupResource(resourceName string, resourceType string) {
	switch resourceType {
	case "pod":
		err := ctx.KubeClient.CoreV1().Pods(ctx.Namespace).Delete(context.TODO(), resourceName, metav1.DeleteOptions{})
		Expect(err).ToNot(HaveOccurred(), "Failed to delete pod %s", resourceName)
	case "vm":
		err := ctx.VirtClient.VirtualMachine(ctx.Namespace).Delete(context.TODO(), resourceName, metav1.DeleteOptions{})
		Expect(err).ToNot(HaveOccurred(), "Failed to delete VM %s", resourceName)
	case "service":
		err := ctx.KubeClient.CoreV1().Services(ctx.Namespace).Delete(context.TODO(), resourceName, metav1.DeleteOptions{})
		Expect(err).ToNot(HaveOccurred(), "Failed to delete service %s", resourceName)
	case "route":
		err := ctx.RouteClient.RouteV1().Routes(ctx.Namespace).Delete(context.TODO(), resourceName, metav1.DeleteOptions{})
		Expect(err).ToNot(HaveOccurred(), "Failed to delete route %s", resourceName)
	case "networkPolicy":
		err := util.DeleteNetworkPolicy(ctx.KubeClient, ctx.Namespace, resourceName)
		Expect(err).ToNot(HaveOccurred(), "Failed to delete NetworkPolicy %s", resourceName)
	}
}