package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"gitlab.com/russell.whelan/aws-tools/pkg/awstools"
)

func main() {

	//Set up buffer and sort headings
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(w, "Name\tType\tVPC ID\tAZs\tDNS Name\t")

	//Create session for elbv1 api aka Classic ELB
	elbsessionv1 := elb.New(session.New(), aws.NewConfig())
	inputv1 := &elb.DescribeLoadBalancersInput{}

	//Create session for elbv2 api aka ALB, NLB etc
	elbsessionv2 := elbv2.New(session.New(), aws.NewConfig())
	inputv2 := &elbv2.DescribeLoadBalancersInput{}

	//Now Lets collect and Display the v2 elbs
	resultv1, err := elbsessionv1.DescribeLoadBalancers(inputv1)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, lb := range resultv1.LoadBalancerDescriptions {
		fmt.Fprintf(w, "%v\t", *lb.LoadBalancerName+" ")
		fmt.Fprintf(w, "%v\t", "classic")
		fmt.Fprintf(w, "%v\t", *lb.VPCId)
		fmt.Fprintf(w, "%v\t", strings.Join(awstools.AzsToStringv1(*lb), ","))
		fmt.Fprintf(w, "%v\t", *lb.DNSName)
		fmt.Fprintf(w, "\n")

	}

	//Now Lets collect and Display the v2 elbs
	resultv2, err := elbsessionv2.DescribeLoadBalancers(inputv2)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, lb := range resultv2.LoadBalancers {
		fmt.Fprintf(w, "%v\t", *lb.LoadBalancerName+" ")
		fmt.Fprintf(w, "%v\t", *lb.Type)
		fmt.Fprintf(w, "%v\t", *lb.VpcId)
		fmt.Fprintf(w, "%v\t", strings.Join(awstools.AzsToStringv2(*lb), ","))
		fmt.Fprintf(w, "%v\t", *lb.DNSName)
		fmt.Fprintf(w, "\n")

	}

	w.Flush()
}
