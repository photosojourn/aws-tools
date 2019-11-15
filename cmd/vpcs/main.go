package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/photosojourn/aws-tools/pkg/awstools"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func getSubnets(vpc ec2.Vpc) []string {
	ec2session := ec2.New(session.New(), aws.NewConfig())

	var vpcSubnets []string
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
		return nil
	}

	for _, subnet := range resultSub.Subnets {
		vpcSubnets = append(vpcSubnets, *subnet.SubnetId)
	}

	return vpcSubnets
}

func main() {

	ec2session := ec2.New(session.New(), aws.NewConfig())
	params := &ec2.DescribeVpcsInput{}

	result, err := ec2session.DescribeVpcs(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "VPC ID\tName\tCIDR Block\tSubnets\tDefault")

	for _, vpc := range result.Vpcs {

		subnets := getSubnets(*vpc)
		name := awstools.GetVpcName(*vpc)
		fmt.Fprintf(w, "%v\t", *vpc.VpcId)
		fmt.Fprintf(w, "%v\t", name)
		fmt.Fprintf(w, "%v\t", *vpc.CidrBlock)
		fmt.Fprintf(w, "%v\t", subnets)
		fmt.Fprintf(w, "%v\t", *vpc.IsDefault)
		fmt.Fprintf(w, "\n")
	}

	w.Flush()

}
