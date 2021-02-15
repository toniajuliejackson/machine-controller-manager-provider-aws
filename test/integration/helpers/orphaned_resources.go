package helpers

import (
	v1alpha1 "github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

// CheckForOrphanedResources will search the cloud provider for orphaned resources that are left behind after the test cases
func CheckForOrphanedResources(machineClass *v1alpha1.MachineClass, secret *v1.Secret) error {
	// Check for VM instances with matching tags/labels
	// Describe volumes attached to VM instance & delete the volumes
	// Finally delete the VM instance

	clusterTag := "tag:kubernetes.io/cluster/shoot--mcm-test--tonia-aws"
	clusterTagValue := "1"

	err := DescribeInstancesWithTag("tag:mcmtest", "integration-test", machineClass, secret)
	if err != nil {
		return err
	}

	// Check for available volumes in cloud provider with tag/label [Status:available]
	err = DescribeAvailableVolumes(clusterTag, clusterTagValue, machineClass, secret)
	if err != nil {
		return err
	}

	// Check for available vpc and network interfaces in cloud provider with tag
	err = AdditionalResourcesCheck(clusterTag, clusterTagValue)
	if err != nil {
		return err
	}

	return nil
}
