package helpers

import ()

// CheckForOrphanedResources will search the cloud provider for orphaned resources that are left behind after the test cases
func CheckForOrphanedResources() error{
	// Check for VM instances with matching tags/labels
	// Describe volumes attached to VM instance & delete the volumes
	// Finally delete the VM instance
	err := DescribeInstancesWithTag("tag:mcmtest", "integration-test")
	if err != nil {
		return err
	}

	// Check for available volumes in cloud provider with tag/label [Status:available]
	err = DescribeAvailableVolumes()
	if err != nil {
		return err
	}

	return nil
}