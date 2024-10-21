
# OpenShift Tests Util

A Go-based utility for managing tests related to Virtual Machines (VMs), services, network policies, and routes on OpenShift clusters. This project provides a framework for creating, deploying, and managing VMs, pods, and other resources using Kubernetes and KubeVirt, while also offering network and route testing capabilities.

## Table of Contents
- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Testing](#testing)

## Overview

This project provides a utility to simplify:
- Creating and managing VMs on KubeVirt using predefined templates.
- Testing network policies, services (ClusterIP, LoadBalancer), and routes.
- Running tests using multiple namespaces.

## Prerequisites

Ensure you have the following installed:
- Go (v1.16+)
- Kubernetes cluster with KubeVirt and OpenShift installed.
- Access to a KubeVirt-enabled cluster with the required permissions.
- Git for version control.

### Environment Requirements:
- `kubeconfig`: Ensure that your kubeconfig file is accessible for authentication.
- Access to OpenShift and KubeVirt APIs for VM and network testing.

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/dror1212/openshift-tests-util.git
   cd openshift-tests-util
   ```

2. **Install dependencies**:
   Ensure that the Go modules are installed:
   ```bash
   go mod tidy
   ```

3. **Set up KubeConfig**:
   - Ensure that your `kubeconfig` is set up and that you have the necessary access to interact with the Kubernetes API.
   - You can specify your `KUBECONFIG` environment variable if needed:
     ```bash
     export KUBECONFIG=/path/to/your/kubeconfig
     ```

## Usage

To run specific tests or use utilities from the framework, follow the instructions below.

### Example Commands:

- **Run tests**: 
  To run tests, use the Ginkgo test suite (default testing framework):
  ```bash
  ginkgo -r --focus "<Test Name>"
  ```

### Key Files and Directories:

- **`util/`**: Contains utility files like:
  - VM creation logic (`vm.go`)
  - Pod management (`pod.go`)
  - Network policy handling (`networkPolicy.go`)
  - Route creation and handling (`route.go`)
- **`framework/`**: Contains higher-level test helpers, such as:
  - `pod_actions.go` for pod-related test utilities.
  - `service_actions.go` for service creation and management.
  - `route_actions.go` for managing OpenShift routes in tests.
  - `test_context.go` for managing reusable test context (namespace, clients, etc.).
- **`tests/`**: Contains Ginkgo-based test cases, including:
  - `network/cluster_ip_test.go`: Tests for ClusterIP service access.
  - `network/network_policy_test.go`: Tests for NetworkPolicy restrictions and access.
  - `network/route_test.go`: Tests for routes in OpenShift.
- **`consts/`**: Holds constant variables such as default memory, namespace settings, and more.

## Configuration

Default settings like CPU, memory, and namespaces can be configured in the `consts/constants.go` file.

To adjust cloud-init scripts or custom startup scripts, modify or add new shell scripts to the `scripts/` directory.

## Testing

To run tests, you can utilize the Ginkgo test suite. For example running all the network relating tests:

```bash
ginkgo -v tests/network/
```

Each test is designed to validate the functionality of OpenShift/Kubernetes components in an isolated environment using test helpers from the `framework` package.