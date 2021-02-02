package services

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/aws/aws-sdk-go/service/cloudtrail/cloudtrailiface"
)

// CloudTrail client wrapper
type CloudTrail interface {
	cloudtrailiface.CloudTrailAPI

	// wrapper for LookupEventsPagesWithContext and aggregate result into list
	LookupEventsAsList(ctx context.Context, startTime *time.Time, endTime *time.Time, lookupAttributes []*cloudtrail.LookupAttribute) ([]*cloudtrail.Event, error)
}

//NewCloudTrail creates a new CloudFormatinon wrapper
func NewCloudTrail(session *session.Session) CloudTrail {
	return &defaultCloudTrail{
		CloudTrailAPI: cloudtrail.New(session),
	}
}

type defaultCloudTrail struct {
	cloudtrailiface.CloudTrailAPI
}

func (c *defaultCloudTrail) LookupEventsAsList(ctx context.Context, startTime *time.Time, endTime *time.Time, lookupAttributes []*cloudtrail.LookupAttribute) ([]*cloudtrail.Event, error) {
	var result []*cloudtrail.Event
	input := cloudtrail.LookupEventsInput{
		StartTime:        startTime,
		EndTime:          endTime,
		LookupAttributes: lookupAttributes,
	}
	if err := c.LookupEventsPagesWithContext(ctx, &input, func(output *cloudtrail.LookupEventsOutput, _ bool) bool {
		result = append(result, output.Events...)
		return true
	}); err != nil {
		return nil, err
	}

	return result, nil
}
