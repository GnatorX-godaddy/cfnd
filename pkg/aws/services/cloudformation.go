package services

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
)

// Cloudformation client wrapper
type Cloudformation interface {
	cloudformationiface.CloudFormationAPI

	// wrapper for ListStackInstancesWithContext and aggregate result into list
	ListStackAsList(ctx context.Context, stackStatus []*string) ([]*cloudformation.StackSummary, error)

	// Searchs for stacks with stackname
	ListStackWithNameAsList(ctx context.Context, stackStatus []*string, stackName string) (*cloudformation.StackSummary, error)

	// wrapper for DescribeStackEvents and aggregate result into list
	DescribeStackEventsAsList(ctx context.Context, stackID string) ([]*cloudformation.StackEvent, error)

	// DescribeStackResources(ctx context.Context, )
}

//NewCloudFormation creates a new CloudFormatinon wrapper
func NewCloudFormation(session *session.Session) Cloudformation {
	return &defaultCloudformation{
		CloudFormationAPI: cloudformation.New(session),
	}
}

type defaultCloudformation struct {
	cloudformationiface.CloudFormationAPI
}

func (c *defaultCloudformation) ListStackAsList(ctx context.Context, stackStatus []*string) ([]*cloudformation.StackSummary, error) {
	var result []*cloudformation.StackSummary

	input := cloudformation.ListStacksInput{}

	input.SetStackStatusFilter(stackStatus)

	if err := c.ListStacksPagesWithContext(ctx, &input, func(output *cloudformation.ListStacksOutput, _ bool) bool {
		result = append(result, output.StackSummaries...)
		return true
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *defaultCloudformation) ListStackWithNameAsList(ctx context.Context, stackStatus []*string, stackName string) (*cloudformation.StackSummary, error) {
	var result *cloudformation.StackSummary

	input := cloudformation.ListStacksInput{}

	input.SetStackStatusFilter(stackStatus)

	if err := c.ListStacksPagesWithContext(ctx, &input, func(output *cloudformation.ListStacksOutput, _ bool) bool {
		for _, summary := range output.StackSummaries {
			if strings.Contains(*summary.StackName, stackName) {
				result = summary
				return false
			}
		}
		return true
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *defaultCloudformation) DescribeStackEventsAsList(ctx context.Context, stackID string) ([]*cloudformation.StackEvent, error) {
	var result []*cloudformation.StackEvent
	input := cloudformation.DescribeStackEventsInput{
		StackName: &stackID,
	}
	if err := c.DescribeStackEventsPagesWithContext(ctx, &input, func(output *cloudformation.DescribeStackEventsOutput, _ bool) bool {
		result = append(result, output.StackEvents...)
		return true
	}); err != nil {
		return nil, err
	}
	return result, nil
}
