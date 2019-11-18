package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/photosojourn/aws-tools/pkg/awstools"
)

func main() {

	ec2session := ec2.New(session.New(), aws.NewConfig())
	params := &ec2.DescribeVolumesInput{}

	result, err := ec2session.DescribeVolumes(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(w, "Volume ID\tName\tInstance ID\tDevice\tAvailabilityZone\tSize (GB)\tType\tIOPS\tEncrypted")

	for _, volume := range result.Volumes {
		name := awstools.GetVolumeName(*volume)
		fmt.Fprintf(w, "%v\t", *volume.VolumeId)
		fmt.Fprintf(w, "%v\t", name)
		if len(volume.Attachments) > 0 {
			for _, attachment := range volume.Attachments {
				fmt.Fprintf(w, "%v\t", *attachment.InstanceId)
				fmt.Fprintf(w, "%v\t", *attachment.Device)
			}
		} else {
			fmt.Fprintf(w, "\t")
			fmt.Fprintf(w, "\t")
		}
		fmt.Fprintf(w, "%v\t", *volume.AvailabilityZone)
		fmt.Fprintf(w, "%v\t", *volume.Size)
		fmt.Fprintf(w, "%v\t", *volume.VolumeType)
		fmt.Fprintf(w, "%v\t", *volume.Iops)
		fmt.Fprintf(w, "%v\t", *volume.Encrypted)
		fmt.Fprintf(w, "\n")
	}

	w.Flush()

}
