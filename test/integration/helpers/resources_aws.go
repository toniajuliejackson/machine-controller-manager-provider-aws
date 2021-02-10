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
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	v1 "k8s.io/api/core/v1"

	api "github.com/gardener/machine-controller-manager-provider-aws/pkg/aws/apis"
	"github.com/gardener/machine-controller-manager-provider-aws/pkg/spi"
	v1alpha1 "github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
)

var _ aws.Config

func newSession(machineClass *v1alpha1.MachineClass, secret *v1.Secret) *session.Session {

	var (
		providerSpec *api.AWSProviderSpec
		sPI          spi.PluginSPIImpl
	)

	err := json.Unmarshal([]byte(machineClass.ProviderSpec.Raw), &providerSpec)
	if err != nil {
		providerSpec = nil
		log.Printf("Error occured while performing unmarshal %s", err.Error())
	}
	sess, err := sPI.NewSession(secret, providerSpec.Region)
	if err != nil {
		log.Printf("Error occured while creating new session %s", err)
	}
	return sess
}

// DescribeInstancesWithTag describes the instance with the specified tag
func DescribeInstancesWithTag(tagName string, tagValue string, machineClass *v1alpha1.MachineClass, secret *v1.Secret) error {
	sess := newSession(machineClass, secret)
	svc := ec2.New(sess)
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
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			fmt.Println(*instance.InstanceId)
			// describe volumes attached to instance & delete them
			//DescribeVolumesAttached(*instance.InstanceId)

			// terminate the instance
			//TerminateInstance(*instance.InstanceId)
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
	}

	fmt.Println(result)
	return nil
}

// DescribeAvailableVolumes describes volumes with the specified tag
func DescribeAvailableVolumes(machineClass *v1alpha1.MachineClass, secret *v1.Secret) error {
	sess := newSession(machineClass, secret)
	svc := ec2.New(sess)
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
	}

	for _, volume := range result.Volumes {
		fmt.Println(*volume.VolumeId)

		// delete the volume
		//DeleteVolume(*volume.VolumeId)
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
	}

	fmt.Println(result)
	return nil
}
