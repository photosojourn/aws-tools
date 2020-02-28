package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/photosojourn/aws-tools/pkg/awstools"
)

func main() {

	zoneID := flag.String("i", "", "Zone ID")
	getRecords := flag.Bool("r", false, "List RecordSets")
	flag.Parse()

	r53session := route53.New(session.New(), aws.NewConfig())
	params := &route53.GetHostedZoneInput{
		Id: aws.String(*zoneID),
	}

	result, err := r53session.GetHostedZone(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	zone := result.HostedZone
	delegationSet := result.DelegationSet
	fmt.Println("Hosted Zone Information")
	fmt.Println("=======================")
	fmt.Println("Hosted Zone ID: " + *zoneID)
	fmt.Println("Name: " + *zone.Name)
	fmt.Println("Resource Record Set Count: " + strconv.FormatInt(*zone.ResourceRecordSetCount, 10))
	fmt.Println("")

	fmt.Println("Delegation Set")
	fmt.Println("==============")
	for _, record := range delegationSet.NameServers {
		fmt.Println(*record)
	}
	fmt.Println("")

	if *getRecords {
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
		fmt.Fprintln(w, "Name\tType\tTTL (Secs)\tValues")
		recordInput := &route53.ListResourceRecordSetsInput{
			HostedZoneId: aws.String(*zoneID),
		}
		recordSets, err := r53session.ListResourceRecordSets(recordInput)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Record Sets")
		fmt.Println("===========")
		fmt.Println("")
		for _, record := range recordSets.ResourceRecordSets {

			fmt.Println(*record)
			values := awstools.RecordSetValuestoString(*record)
			fmt.Fprintf(w, "%v\t", *record.Name)
			fmt.Fprintf(w, "%v\t", *record.Type)
			fmt.Fprintf(w, "%v\t", *record.TTL)
			fmt.Fprintf(w, "%v\t", strings.Join(values, " "))
			fmt.Fprintf(w, "\n")
		}
		w.Flush()
	}

}
