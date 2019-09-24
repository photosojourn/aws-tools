package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func main() {

	rdssession := rds.New(session.New(), aws.NewConfig())

	params := &rds.DescribeDBClustersInput{}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(w, "Instance Name\tCluster Name\tMulti AZ\tEngine\tEngine Version\tInstance Type\tCluster Writer\tPromotion Tier")

	//First job is to deal with the instances that are covered under clusters
	clusterResult, err := rdssession.DescribeDBClusters(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, cluster := range clusterResult.DBClusters {
		clusterName := *cluster.DBClusterIdentifier

		for _, clusterInst := range cluster.DBClusterMembers {

			instParams := &rds.DescribeDBInstancesInput{
				DBInstanceIdentifier: aws.String(*clusterInst.DBInstanceIdentifier),
			}

			dbinst, err := rdssession.DescribeDBInstances(instParams)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Fprintf(w, "%v\t", *clusterInst.DBInstanceIdentifier)
			fmt.Fprintf(w, "%v\t", clusterName)
			fmt.Fprintf(w, "%v\t", *cluster.MultiAZ)
			fmt.Fprintf(w, "%v\t", *dbinst.DBInstances[0].Engine)
			fmt.Fprintf(w, "%v\t", *dbinst.DBInstances[0].EngineVersion)
			fmt.Fprintf(w, "%v\t", *dbinst.DBInstances[0].DBInstanceClass)
			fmt.Fprintf(w, "%v\t", *clusterInst.IsClusterWriter)
			fmt.Fprintf(w, "%v\t", *clusterInst.PromotionTier)
			fmt.Fprintf(w, "\n")
		}
	}

	instanceSearch := &rds.DescribeDBInstancesInput{}

	//Now lets grab all instances, strip Aurora so we can grab the rest
	instanceResults, err := rdssession.DescribeDBInstances(instanceSearch)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, insts := range instanceResults.DBInstances {
		if *insts.Engine != "aurora" {
			fmt.Fprintf(w, "%v\t", *insts.DBInstanceIdentifier)
			fmt.Fprintf(w, "%v\t", "N/A")
			fmt.Fprintf(w, "%v\t", *insts.MultiAZ)
			fmt.Fprintf(w, "%v\t", *insts.Engine)
			fmt.Fprintf(w, "%v\t", *insts.EngineVersion)
			fmt.Fprintf(w, "%v\t", *insts.DBInstanceClass)
			fmt.Fprintf(w, "%v\t", "N/A")
			fmt.Fprintf(w, "%v\t", "N/A")
			fmt.Fprintf(w, "\n")
		}
	}

	w.Flush()
}
