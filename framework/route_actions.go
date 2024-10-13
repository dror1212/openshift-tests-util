package framework

import (
	"time"
	"myproject/util"
	. "github.com/onsi/gomega"
)

// CreateRouteHelper creates a route for the given service with the given port and hostname
func (ctx *TestContext) CreateRouteHelper(routeName, serviceName string, targetPort interface{}, hostname string) {
    err := util.CreateRoute(ctx.RouteClient, ctx.Namespace, routeName, serviceName, targetPort, hostname)
    Expect(err).ToNot(HaveOccurred(), "Failed to create route %s", routeName)
}

// GetRouteURLHelper fetches the URL of a route and returns it
func (ctx *TestContext) GetRouteURLHelper(routeName string) string {
	routeURL, err := util.GetRouteURL(ctx.RouteClient, ctx.Namespace, routeName)
	Expect(err).ToNot(HaveOccurred(), "Failed to get URL for Route %s", routeName)
	return routeURL
}

// WaitForRouteURL waits for a route to get a valid URL within a specified timeout
func (ctx *TestContext) WaitForRouteURL(routeName string, timeout, interval time.Duration) string {
	var routeURL string

	Eventually(func() (string, error) {
		routeURL = ctx.GetRouteURLHelper(routeName)
		return routeURL, nil
	}, timeout, interval).ShouldNot(BeEmpty(), "Expected route to get a valid URL")

	return routeURL
}