package ctl

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"cloudformation-error/pkg/aws/model/cloudformation"
	"cloudformation-error/pkg/aws/services"
)

func Find(ctx context.Context, stackName string, region string) string {
	awsSess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	cfClient := services.NewCloudFormation(awsSess)
	create := cloudformation.StackStatusUPDATEROLLBACKCOMPLETE.String()
	stackStatus := []*string{&create}
	result, err := cfClient.ListStackAsList(ctx, stackStatus)
	for _, summary := range result {
		if strings.Contains(*summary.StackName, stackName) {
			return *summary.StackName
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	return "sdf"
}
