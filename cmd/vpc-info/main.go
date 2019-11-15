package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/photosojourn/aws-tools/pkg/awstools"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getSubnets(vpc ec2.Vpc) {
	ec2session := ec2.New(session.New(), aws.NewConfig())

	subnetFilter := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []*string{
					aws.String(*vpc.VpcId),
				},
			},
		},
	}

	resultSub, err := ec2session.DescribeSubnets(subnetFilter)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Subnet Information")
	fmt.Println("==================")

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(w, "Subnet ID\tName\tCIDR Block\tAvailable Ip's\tAvailability Zone\tPublic IP")

	for _, subnet := range resultSub.Subnets {
		fmt.Fprintf(w, "%v\t", *subnet.SubnetId)
		fmt.Fprintf(w, "%v\t", awstools.GetSubnetName(*subnet))
		fmt.Fprintf(w, "%v\t", *subnet.CidrBlock)
		fmt.Fprintf(w, "%v\t", *subnet.AvailableIpAddressCount)
		fmt.Fprintf(w, "%v\t", *subnet.AvailabilityZone)
		fmt.Fprintf(w, "%v\t", *subnet.MapPublicIpOnLaunch)
		fmt.Fprintf(w, "\n")
	}

	w.Flush()
	fmt.Println("")

}

func getNacls(vpc ec2.Vpc) {
	ec2session := ec2.New(session.New(), aws.NewConfig())

	naclFilter := &ec2.DescribeNetworkAclsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []*string{
					aws.String(*vpc.VpcId),
				},
			},
		},
	}

	resultNacl, err := ec2session.DescribeNetworkAcls(naclFilter)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("NACL Information")
	fmt.Println("==================")

	for _, nacl := range resultNacl.NetworkAcls {
		fmt.Println(*nacl.NetworkAclId)
	}
	fmt.Println("")
}

func main() {

	vpcID := flag.String("i", "", "ID of VPC")
	vpcSubnet := flag.Bool("s", false, "Subnet info")
	vpcNacl := flag.Bool("n", false, "NACL Information")
	vpcAll := flag.Bool("a", false, "List of info")
	flag.Parse()

	ec2session := ec2.New(session.New(), aws.NewConfig())
	params := &ec2.DescribeVpcsInput{
		VpcIds: []*string{
			aws.String(*vpcID),
		},
	}

	result, err := ec2session.DescribeVpcs(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	vpc := result.Vpcs[0]

	//subnets := getSubnets(*vpc)
	name := awstools.GetVpcName(*vpc)

	fmt.Println("VPC Info")
	fmt.Println("========")
	fmt.Println("VPC Id: " + *vpc.VpcId)
	fmt.Println("Name: " + name)
	fmt.Println("")

	if *vpcSubnet || *vpcAll {
		getSubnets(*vpc)
	}

	if *vpcNacl || *vpcAll {
		getNacls(*vpc)
	}

}
