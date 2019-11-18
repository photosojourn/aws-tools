package awstools

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

//GetName pulls the Name tag from the instance
func GetName(inst ec2.Instance) string {
	for _, tag := range inst.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return ""
}

//GetVpcName pull the name tag from a VPC
func GetVpcName(vpc ec2.Vpc) string {
	for _, tag := range vpc.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return ""
}

//GetSubnetName pull the name tag from a Subnet
func GetSubnetName(subnet ec2.Subnet) string {
	for _, tag := range subnet.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return ""
}

//GetVolumeName pull the name tag from a Volume
func GetVolumeName(volume ec2.Volume) string {
	for _, tag := range volume.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return ""
}

//AzsToStringv1 creates a string on AZ's from the list for elb api
func AzsToStringv1(azs elb.LoadBalancerDescription) []string {
	var azsString []string
	for _, az := range azs.AvailabilityZones {
		azsString = append(azsString, *az)
	}

	return azsString
}

//AzsToStringv2 creates a string on AZ's from the list for elbv2 api
func AzsToStringv2(azs elbv2.LoadBalancer) []string {
	var azsString []string
	for _, az := range azs.AvailabilityZones {
		azsString = append(azsString, *az.ZoneName)
	}

	return azsString
}
