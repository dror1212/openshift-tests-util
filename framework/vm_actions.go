package framework

import (
	"myproject/util"
	. "github.com/onsi/gomega"
)

// CreateTestVM creates a VM with default settings using the random name from the context
func (ctx *TestContext) CreateTestVM(vmName string, scriptPath string, templateName string) {
	resourceRequirements := util.ConvertCoreV1ToKubeVirtResourceRequirements(
		util.GenerateResourceRequirements("4000m", "4000m", "4Gi", "4Gi"))

	_, err := util.CreateVM(ctx.Config, ctx.Namespace, templateName, vmName, &resourceRequirements, nil, true, scriptPath, "")
	Expect(err).ToNot(HaveOccurred(), "Failed to create VM")
}
