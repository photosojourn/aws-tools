package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func main() {

	r53session := route53.New(session.New(), aws.NewConfig())
	params := &route53.ListHostedZonesInput{}

	result, err := r53session.ListHostedZones(params)
	if err != nil {
		fmt.Println(err)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "Hosted Zone ID\tName\tPrivate Zone\tComment")

	for _, zone := range result.HostedZones {
		id := strings.SplitN(*zone.Id, "/", 3)
		fmt.Fprintf(w, "%v\t", id[2])
		fmt.Fprintf(w, "%v\t", *zone.Name)
		fmt.Fprintf(w, "%v\t", *zone.Config.PrivateZone)
		fmt.Fprintf(w, "%v\t", *zone.Config.Comment)
		fmt.Fprintf(w, "\n")
	}

	w.Flush()
}
