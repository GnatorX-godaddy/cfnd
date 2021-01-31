package ctl

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"cloudformation-error/pkg/aws/services"
)

func Find(ctx context.Context, stackName string, region string) string {
	awsSess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	cfClient := services.NewCloudFormation(awsSess)
	stackStatus := []*string{}
	cfStack, err := cfClient.ListStackWithNameAsList(ctx, stackStatus, stackName)
	cfStackEvents, err := cfClient.DescribeStackEventsAsList(ctx, *cfStack.StackName)
	fmt.Println(cfStackEvents)
	if err != nil {
		log.Fatal(err)
	}
	return "sdf"
}
