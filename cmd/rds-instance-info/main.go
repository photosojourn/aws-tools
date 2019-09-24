package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

type options struct {
	instanceName   string
	instanceEvents bool
}

func main() {

	opts := new(options)
	flag.StringVar(&opts.instanceName, "n", "", "Instance Name")
	flag.BoolVar(&opts.instanceEvents, "e", false, "Show Instance Events")
	flag.Parse()

	rdssession := rds.New(session.New(), aws.NewConfig())

	instanceInput := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(opts.instanceName),
	}

	result, err := rdssession.DescribeDBInstances(instanceInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	dbinst := result.DBInstances[0]

	fmt.Println("DB Instance Information")
	fmt.Println("=======================")
	fmt.Println("Instance Name: " + *dbinst.DBInstanceIdentifier)
	fmt.Println("Instance ARN: " + *dbinst.DBInstanceArn)
	fmt.Println("Instance Status: " + *dbinst.DBInstanceStatus)
	fmt.Println("Instance Type: " + *dbinst.DBInstanceClass)
	fmt.Println("Engine Type: " + *dbinst.Engine)
	fmt.Println("Engine Version: " + *dbinst.EngineVersion)

	if *dbinst.Engine == "aurora" {

		fmt.Println("\nAurora Information")
		fmt.Println("==================")

		params := &rds.DescribeDBClustersInput{
			DBClusterIdentifier: aws.String(*dbinst.DBClusterIdentifier),
		}

		dbcluster, err := rdssession.DescribeDBClusters(params)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Cluster Name: " + *dbcluster.DBClusters[0].DBClusterIdentifier)
		fmt.Println("Cluster Mode: " + *dbcluster.DBClusters[0].EngineMode)
		for _, member := range dbcluster.DBClusters[0].DBClusterMembers {
			if *member.DBInstanceIdentifier == *dbinst.DBInstanceIdentifier {
				if *member.IsClusterWriter == true {
					fmt.Println("Cluster Role: Writer")
					fmt.Println("Promotion Tier: " + strconv.FormatInt(*member.PromotionTier, 10))
				} else {
					fmt.Println("Cluster Role: Reader")
					fmt.Println("Promotion Tier: " + strconv.FormatInt(*member.PromotionTier, 10))
				}
			}
		}
		fmt.Println("Earliest Restorable Time: " + (*dbcluster.DBClusters[0].LatestRestorableTime).Format(time.RFC3339))
	}

	if opts.instanceEvents == true {
		params := &rds.DescribeEventsInput{
			Duration:         aws.Int64(10080),
			SourceIdentifier: aws.String(opts.instanceName),
			SourceType:       aws.String("db-instance"),
		}

		events, err := rdssession.DescribeEvents(params)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("\nEvents")
		fmt.Println("======")
		for _, event := range events.Events {
			fmt.Println("*" + (*event.Date).Format(time.RFC3339))
			fmt.Println("  " + *event.Message)
		}

		//fmt.Println(events)
	}
}
