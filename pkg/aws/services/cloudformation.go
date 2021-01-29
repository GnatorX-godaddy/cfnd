package services

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
)

// Cloudformation client wrapper
type Cloudformation interface {
	cloudformationiface.CloudFormationAPI

	// wrapper for ListStackInstancesWithContext and aggregate result into list
	ListStackAsList(ctx context.Context, stackStatus []*string) ([]*cloudformation.StackInstanceSummary, error)

	// wrapper for DescribeStackEvents and aggregate result into list
	DescribeStackEventsAsList(ctx context.Context, stackName string) ([]*cloudformation.StackEvent, error)
}

func NewCloudFormation(session *session.Session) Cloudformation {
	return &defaultCloudformation{
		CloudFormationAPI: cloudformation.New(session),
	}
}

type defaultCloudformation struct {
	cloudformationiface.CloudFormationAPI
}

func (c *defaultCloudformation) ListStackAsList(ctx context.Context, stackStatus []*string) ([]*cloudformation.StackInstanceSummary, error) {
	var result []*cloudformation.StackInstanceSummary

	input := cloudformation.ListStacksInput{}

	input.SetStackStatusFilter(stackStatus)

	if err := c.ListStackInstancesPagesWithContext(ctx, &input, func(output *cloudformation.ListStackInstancesOutput, _ bool) bool {
		result = append(result, output.Summaries...)
		return true
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *defaultCloudformation) DescribeStackEventsAsList(ctx context.Context, stackName string) ([]*cloudformation.StackEvent, error) {
	var result []*cloudformation.StackEvent
	input := cloudformation.DescribeStackEventsInput{
		StackName: &stackName,
	}
	if err := c.DescribeStackEventsPagesWithContext(ctx, &input, func(output *cloudformation.DescribeStackEventsOutput, _ bool) bool {
		result = append(result, output.StackEvents...)
		return true
	}); err != nil {
		return nil, err
	}
	return result, nil
}
