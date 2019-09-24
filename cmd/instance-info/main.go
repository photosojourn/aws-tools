package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/photosojourn/aws-tools/pkg/awstools"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

/*func getName(inst ec2.Instance) string {
	for _, tag := range inst.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return ""
}*/

func main() {

	instName := flag.String("i", "", "Instance ID")
	instNet := flag.Bool("n", false, "Network information")
	instVol := flag.Bool("v", false, "Volumes")
	instTags := flag.Bool("t", false, "List Tags")
	instAll := flag.Bool("a", false, "analogous to -nvt")

	flag.Parse()

	ec2session := ec2.New(session.New(), aws.NewConfig())

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(*instName),
		},
	}

	result, err := ec2session.DescribeInstances(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	inst := result.Reservations[0].Instances[0]
	name := awstools.GetName(*inst)

	fmt.Println("Instance Information")
	fmt.Println("====================")
	fmt.Println("Instance Id: " + *inst.InstanceId)
	fmt.Println("Name: " + name)
	fmt.Println("Instance Type: " + *inst.InstanceType)
	fmt.Println("Launch Time: " + (*inst.LaunchTime).Format(time.RFC3339))
	fmt.Println("AMI ID: " + *inst.ImageId)

	if *instNet || *instAll {

		fmt.Printf("\n")
		fmt.Println("Network Information")
		fmt.Println("===================")
		fmt.Println("VPC ID: " + *inst.VpcId)
		fmt.Println("Subnet ID: " + *inst.SubnetId)
		if inst.PublicIpAddress != nil {
			fmt.Println("Public IP: " + *inst.PublicIpAddress)
		} else {
			fmt.Println("Public IP: N/A")
		}

		for _, ifaces := range inst.NetworkInterfaces {
			fmt.Println("\nENI: " + *ifaces.NetworkInterfaceId)
			for _, ip := range ifaces.PrivateIpAddresses {
				if *ip.Primary == false {
					fmt.Println(*ip.PrivateDnsName + " : " + *ip.PrivateIpAddress)
				} else {
					fmt.Println("* " + *ip.PrivateDnsName + " : " + *ip.PrivateIpAddress)
				}
			}
		}
	}

	if *instVol || *instAll {

		fmt.Printf("\n")
		fmt.Println("Volume Information")
		fmt.Println("==================")
		fmt.Println("Root Device Name: " + *inst.RootDeviceName)
		fmt.Println("Root Device Type: " + *inst.RootDeviceType)

		fmt.Println("\nEBS Volumes:")
		for _, vols := range inst.BlockDeviceMappings {
			fmt.Println(*vols.DeviceName + " : " + *vols.Ebs.VolumeId)
		}
	}

	if *instTags || *instAll {
		fmt.Printf("\n")
		fmt.Println("Tags")
		fmt.Println("==================")

		for _, tag := range inst.Tags {
			fmt.Println(*tag.Key + " : " + *tag.Value)
		}
	}
}
