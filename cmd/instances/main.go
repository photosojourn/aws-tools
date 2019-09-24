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

func main() {

	ec2session := ec2.New(session.New(), aws.NewConfig())

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
					aws.String("pending"),
				},
			},
		},
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "Instance ID\tAMI ID\tInstance Type\tState\tName\tLaunch Time\tAvailability Zone\tPrivate IP")

	result, err := ec2session.DescribeInstances(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	for idx := range result.Reservations {
		for _, inst := range result.Reservations[idx].Instances {
			name := awstools.GetName(*inst)
			fmt.Fprintf(w, "%v\t", *inst.InstanceId)
			fmt.Fprintf(w, "%v\t", *inst.ImageId)
			fmt.Fprintf(w, "%v\t", *inst.InstanceType)
			fmt.Fprintf(w, "%v\t", *inst.State.Name)
			fmt.Fprintf(w, "%v\t", name)
			fmt.Fprintf(w, "%v\t", *inst.LaunchTime)
			fmt.Fprintf(w, "%v\t", *inst.Placement.AvailabilityZone)
			fmt.Fprintf(w, "%v\t", *inst.PrivateIpAddress)
			fmt.Fprintf(w, "\n")
		}
	}

	w.Flush()
}
