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
	result, err := cfClient.ListStackWithNameAsList(ctx, stackStatus, stackName)
	fmt.Println(result)
	if err != nil {
		log.Fatal(err)
	}
	return "sdf"
}
