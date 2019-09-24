package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/photosojourn/aws-tools/pkg/awstools"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

type options struct {
	name             string
	elbtype          string
	netInfo          bool
	listenerInfo     bool
	loggingInfo      bool
	tagsInfo         bool
	instancesInfo    bool
	targetGroupsInfo bool
	rulesInfo        bool
	allInfo          bool
}

func elbinfov1(opts options) {

	elbsession := elb.New(session.New(), aws.NewConfig())

	//Get the base ELB Info
	input := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{
			aws.String(opts.name),
		},
	}

	result, err := elbsession.DescribeLoadBalancers(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Get ELB's additional Attributes
	attrInput := &elb.DescribeLoadBalancerAttributesInput{
		LoadBalancerName: aws.String(opts.name),
	}

	attributes, err := elbsession.DescribeLoadBalancerAttributes(attrInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	inst := result.LoadBalancerDescriptions[0]
	fmt.Println("ELB Information")
	fmt.Println("===============")
	fmt.Println("ELB Name: " + *inst.LoadBalancerName)
	fmt.Println("ELB Type: classic")
	fmt.Println("DNS Name: " + *inst.DNSName)
	fmt.Println("VPC ID: " + *inst.VPCId)
	fmt.Println("Scheme: " + *inst.Scheme)

	if opts.listenerInfo || opts.allInfo {
		fmt.Println("\nListeners")
		fmt.Println("===========")
		for _, listener := range inst.ListenerDescriptions {
			fmt.Println("Listener Port: " + strconv.FormatInt(*listener.Listener.LoadBalancerPort, 10))
			fmt.Println("  Protocol: " + *listener.Listener.InstanceProtocol)
			fmt.Println("  Instance Port: " + strconv.FormatInt(*listener.Listener.InstancePort, 10))
			if *listener.Listener.Protocol == "HTTPS" || *listener.Listener.Protocol == "SSL" {
				fmt.Println("  ACM Cert: " + *listener.Listener.SSLCertificateId)
			}
		}
	}

	if opts.instancesInfo || opts.allInfo {
		fmt.Println("\nInstances")
		fmt.Println("=========")
		for _, ec2 := range inst.Instances {
			fmt.Println(*ec2.InstanceId)
		}
	}

	if opts.loggingInfo || opts.allInfo {
		if *attributes.LoadBalancerAttributes.AccessLog.Enabled == true {
			fmt.Println("\nLogging Info")
			fmt.Println("========")
			fmt.Println("Logging Status: true")
			fmt.Println("Log Bucket: " + *attributes.LoadBalancerAttributes.AccessLog.S3BucketName)
			fmt.Println("Log Bucket Prefix: " + *attributes.LoadBalancerAttributes.AccessLog.S3BucketPrefix)
		} else {
			fmt.Println("\nLogging Info")
			fmt.Println("========")
			fmt.Println("Logging Status: false")
		}
	}

	if opts.tagsInfo || opts.allInfo {

		tagInput := &elb.DescribeTagsInput{
			LoadBalancerNames: []*string{
				aws.String(opts.name),
			},
		}

		tags, err := elbsession.DescribeTags(tagInput)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("\nTags")
		fmt.Println("====")
		for _, tagList := range tags.TagDescriptions {
			for _, tag := range tagList.Tags {
				fmt.Println(*tag.Key + " : " + *tag.Value)
			}
		}

	}

}

func elbinfov2(opts options) {

	elbsession := elbv2.New(session.New(), aws.NewConfig())

	input := &elbv2.DescribeLoadBalancersInput{
		Names: []*string{
			aws.String(opts.name),
		},
	}

	result, err := elbsession.DescribeLoadBalancers(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	inputAttr := &elbv2.DescribeLoadBalancerAttributesInput{
		LoadBalancerArn: aws.String(*result.LoadBalancers[0].LoadBalancerArn),
	}

	attributes, err := elbsession.DescribeLoadBalancerAttributes(inputAttr)
	if err != nil {
		fmt.Println(err)
		return
	}

	attrMap := map[string]string{}
	for _, attr := range attributes.Attributes {
		attrMap[*attr.Key] = *attr.Value
	}

	inst := result.LoadBalancers[0]
	fmt.Println("ELB Information")
	fmt.Println("===============")
	fmt.Println("ELB Name: " + *inst.LoadBalancerName)
	fmt.Println("ELB Type: " + *inst.Type)
	fmt.Println("DNS Name: " + *inst.DNSName)
	fmt.Println("Deletion Protection: " + attrMap["deletion_protection.enabled"])
	if *inst.Type == "application" {
		fmt.Println("\nALB Specific Information")
		fmt.Println("=====================")
		fmt.Println("Idle Timeout: " + attrMap["idle_timeout.timeout_seconds"])
		fmt.Println("HTTP2 Enabled: " + attrMap["routing.http2.enabled"])
	}

	if opts.netInfo || opts.allInfo {
		fmt.Println("\nNetwork Information")
		fmt.Println("===================")
		fmt.Println("Scheme: " + *inst.Scheme)
		fmt.Println("VPC ID: " + *inst.VpcId)
		fmt.Println("AZ's: " + strings.Join(awstools.AzsToStringv2(*inst), ","))
	}

	if opts.listenerInfo || opts.allInfo {

		inputListener := &elbv2.DescribeListenersInput{
			LoadBalancerArn: aws.String(*result.LoadBalancers[0].LoadBalancerArn),
		}

		listeners, err := elbsession.DescribeListeners(inputListener)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("\nListener Information")
		fmt.Println("======================")
		for _, listener := range listeners.Listeners {
			fmt.Println("Listener ARN: " + *listener.ListenerArn)
			fmt.Println("  Port: " + strconv.FormatInt(*listener.Port, 10))
			fmt.Println("  Protocol: " + *listener.Protocol)
			if *listener.Protocol == "HTTPS" {
				for _, cert := range listener.Certificates {
					fmt.Println("  ACM Cert ARN: " + *cert.CertificateArn)
				}
			}
		}
	}

	if opts.targetGroupsInfo || opts.allInfo {

		inputTargetgroups := &elbv2.DescribeTargetGroupsInput{
			LoadBalancerArn: aws.String(*result.LoadBalancers[0].LoadBalancerArn),
		}

		targetgroups, err := elbsession.DescribeTargetGroups(inputTargetgroups)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("\nTarget Group Information")
		fmt.Println("======================")
		for _, target := range targetgroups.TargetGroups {
			fmt.Println("Target Group Name: " + *target.TargetGroupName)
			fmt.Println("  Target Group ARN: " + *target.TargetGroupArn)
		}
	}

	if opts.loggingInfo || opts.allInfo {
		if attrMap["access_logs.s3.enabled"] == "true" {
			fmt.Println("\nLogging Information")
			fmt.Println("=====================")
			fmt.Println("Logging Status: true")
			fmt.Println("Log Bucket: " + attrMap["access_logs.s3.bucket"])
			fmt.Println("Log Bucket Prefix: " + attrMap["access_logs.s3.prefix"])
		} else {
			fmt.Println("\nLogging Information")
			fmt.Println("=====================")
			fmt.Println("Logging Status: false")
		}
	}

}

func main() {

	opts := new(options)
	flag.StringVar(&opts.name, "n", "", "ELB Name")
	flag.StringVar(&opts.elbtype, "t", "", "ELB Type")
	flag.BoolVar(&opts.instancesInfo, "i", false, "List Instances (type classic only)")
	flag.BoolVar(&opts.listenerInfo, "l", false, "List Listeners (type application/network only)")
	flag.BoolVar(&opts.loggingInfo, "L", false, "Show Log Information")
	flag.BoolVar(&opts.netInfo, "N", false, "List Network Information (type application/network only")
	flag.BoolVar(&opts.targetGroupsInfo, "tg", false, "List Target Groups (type application/network only")
	flag.BoolVar(&opts.tagsInfo, "T", false, "List Tags")
	flag.BoolVar(&opts.allInfo, "a", false, "Show all information")

	flag.Parse()

	switch opts.elbtype {
	case "classic":
		elbinfov1(*opts)
	case "application":
		elbinfov2(*opts)
	default:
		fmt.Println("Error: ELB must be of type classic or application")
	}

}
