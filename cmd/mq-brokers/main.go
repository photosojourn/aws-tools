package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mq"
)

func main() {

	//Set up buffer and sort headings
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(w, "Name\tType\tVPC ID\tAZs\tDNS Name\t")

	//Create session
	mqSession := mq.New(session.New(), aws.NewConfig())
	input := &mq.ListBrokersInput{}

	result, err := mqSession.ListBrokers(input)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}
