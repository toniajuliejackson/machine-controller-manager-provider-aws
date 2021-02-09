package helpers

/**
	Orphaned Resources
	- VMs:
		Describe instances with specified tag name:<cluster-name>
		Report/Print out instances found
		Describe volumes attached to the instance (using instance id)
		Report/Print out volumes found
		Delete attached volumes found
		Terminate instances found
	- Disks:
		Describe volumes with tag status:available
		Report/Print out volumes found
		Delete identified volumes
**/

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var _ aws.Config

// DescribeInstancesWithTag describes the instance with the specified tag
func DescribeInstancesWithTag(tagName string, tagValue string) error {
	svc := ec2.New(session.New())
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String(tagName),
				Values: []*string{
					aws.String(tagValue),
				},
			},
		},
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return err.(awserr.Error)
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return err
		}
		return err
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			fmt.Println(*instance.InstanceId)
			// describe volumes attached to instance & delete them
			DescribeVolumesAttached(*instance.InstanceId)

			// terminate the instance
			TerminateInstance(*instance.InstanceId)
		}
	}
	return nil
}

// TerminateInstance terminates the specified EC2 instance.
func TerminateInstance(instanceID string) error {
	svc := ec2.New(session.New())
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	result, err := svc.TerminateInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return err.(awserr.Error)
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return err
		}
		return err
	}

	fmt.Println(result)
	return nil
}

// DescribeAvailableVolumes describes volumes with the specified tag
func DescribeAvailableVolumes() error {
	svc := ec2.New(session.New())
	input := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("status"),
				Values: []*string{
					aws.String("available"),
				},
			},
		},
	}

	result, err := svc.DescribeVolumes(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return err.(awserr.Error)
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return err
		}
		return err
	}

	for _, volume := range result.Volumes {
		fmt.Printf("available volume: %s\n", *volume.VolumeId)

		// delete the volume
		//DeleteVolume(*volume.VolumeId)
	}

	return nil
}

// DescribeVolumesAttached describes volumes that are attached to a specific instance
func DescribeVolumesAttached(InstanceID string) error {
	svc := ec2.New(session.New())
	input := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("attachment.instance-id"),
				Values: []*string{
					aws.String(InstanceID),
				},
			},
			{
				Name: aws.String("attachment.delete-on-termination"),
				Values: []*string{
					aws.String("true"),
				},
			},
		},
	}

	result, err := svc.DescribeVolumes(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return err.(awserr.Error)
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return err
		}
		return err
	}

	for _, volume := range result.Volumes {
		fmt.Println(*volume.VolumeId)

		// delete the volume
		DeleteVolume(*volume.VolumeId)
	}

	return nil
}

// DeleteVolume deletes the specified volume
func DeleteVolume(VolumeID string) error {
	// TO-DO: deletes an available volume with the specified volume ID
	// If the command succeeds, no output is returned.
	svc := ec2.New(session.New())
	input := &ec2.DeleteVolumeInput{
		VolumeId: aws.String(VolumeID),
	}

	result, err := svc.DeleteVolume(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return err.(awserr.Error)
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return err
		}
		return err
	}

	fmt.Println(result)
	return nil 
}
