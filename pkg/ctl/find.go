package ctl

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"cloudformation-error/pkg/aws/services"
)

func Find(ctx context.Context, stackName string, region string) string {
	awsSess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	cfClient := services.NewCloudFormation(awsSess)
	result, err := cfClient.ListStackAsList(ctx, "input")
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range result {
		print(r.StackId)
	}
	return "sdf"
}
