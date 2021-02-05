package ctl

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/stretchr/testify/assert"
)

type mockCFClient struct {
	cloudformationiface.CloudFormationAPI

	ListStackAsListResult []*cloudformation.StackSummary
	ListStackAsListError  error

	ListStackWithNameAsListResult *cloudformation.StackSummary
	ListStackWithNameAsListError  error

	DescribeStackEventsAsListResult []*cloudformation.StackEvent
	DescribeStackEventsAsListError  error

	DescribeStackResourcesAsListResult []*cloudformation.StackResource
	DescribeStackResourcesAsListError  error
}

func (c *mockCFClient) ListStackAsList(ctx context.Context, stackStatus []*string) ([]*cloudformation.StackSummary, error) {
	return c.ListStackAsListResult, c.ListStackAsListError
}

func (c *mockCFClient) ListStackWithNameAsList(ctx context.Context, stackStatus []*string, stackName string) (*cloudformation.StackSummary, error) {
	return c.ListStackWithNameAsListResult, c.ListStackWithNameAsListError
}

func (c *mockCFClient) DescribeStackEventsAsList(ctx context.Context, stackID string) ([]*cloudformation.StackEvent, error) {
	return c.DescribeStackEventsAsListResult, c.DescribeStackEventsAsListError
}

func (c *mockCFClient) DescribeStackResourcesAsList(ctx context.Context, stackID string) ([]*cloudformation.StackResource, error) {
	return c.DescribeStackResourcesAsListResult, c.DescribeStackResourcesAsListError
}

func TestIsStartAndEndTimeInBound_StartWithin15Mins_False(t *testing.T) {
	startTime := time.Now().Add(time.Minute * -10)
	endTime := time.Now()
	actual := isStartAndEndTimeInBound(startTime, endTime)
	assert.False(t, actual, "isStartAndEndTimeInBound should return false when start is < 15 min from now")
}

func TestIsStartAndEndTimeInBound_StartGreater15MinsAgo_True(t *testing.T) {
	startTime := time.Now().Add(time.Minute * -30)
	endTime := time.Now().Add(time.Minute * -20)
	actual := isStartAndEndTimeInBound(startTime, endTime)
	assert.True(t, actual, "isStartAndEndTimeInBound should return true when start is > 15 min from now")
}

func TestIsStartAndEndTimeInBound_StartLaterEnd_False(t *testing.T) {
	startTime := time.Now().Add(time.Minute * -10)
	endTime := time.Now().Add(time.Minute * -20)
	actual := isStartAndEndTimeInBound(startTime, endTime)
	assert.False(t, actual, "isStartAndEndTimeInBound should return false when start is > end")
}

func TestIsStartAndEndTimeInBound_EndLater90Days_False(t *testing.T) {
	startTime := time.Now().Add(time.Minute * -10)
	endTime := time.Now().Add(time.Hour * -24 * 91)
	actual := isStartAndEndTimeInBound(startTime, endTime)
	assert.False(t, actual, "isStartAndEndTimeInBound should return false when end is > 90 days from now")
}

func TestFindCFStackEvents_NoStacks_nil(t *testing.T) {
	mockClient := mockCFClient{}
	ctx := context.Background()

	actual, err := findCFStackEvents(ctx, &mockClient, "test")
	assert.Nil(t, err, "Error should be nil")
	assert.Nil(t, actual, "findCFStackEvents should return nil if CF client returned nil")
}

func TestFindCFStackEvents_ListStackWithNameAsListError_nil(t *testing.T) {
	mockClient := mockCFClient{
		ListStackWithNameAsListError: errors.New("fail"),
	}
	ctx := context.Background()

	actual, err := findCFStackEvents(ctx, &mockClient, "test")
	assert.NotNil(t, err, "Should return an error if findCFStackEvents failed")
	assert.Nil(t, actual, "findCFStackEvents should return nil if CF client returned nil")
}

func TestFindCFStackEvents_FoundStackAndEvents_Events(t *testing.T) {
	stackID := "1234"
	stack := cloudformation.StackSummary{
		StackId: &stackID,
	}
	events := []*cloudformation.StackEvent{}
	mockClient := mockCFClient{
		ListStackWithNameAsListResult:   &stack,
		DescribeStackEventsAsListResult: events,
	}

	ctx := context.Background()

	actual, err := findCFStackEvents(ctx, &mockClient, "test")
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, actual, "findCFStackEvents should return events")
}
