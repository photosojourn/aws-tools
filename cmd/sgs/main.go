package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	ec2session := ec2.New(session.New(), aws.NewConfig())

	input := &ec2.DescribeSecurityGroupsInput{}

	result, err := ec2session.DescribeSecurityGroups(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(w, "SecurityGroupID\tGroup Name\tVPC ID")

	for _, sg := range result.SecurityGroups {
		fmt.Fprintf(w, "%v\t", *sg.GroupId)
		fmt.Fprintf(w, "%v\t", *sg.GroupName)
		fmt.Fprintf(w, "%v\t", *sg.VpcId)
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
}
