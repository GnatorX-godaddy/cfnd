package ctl

import (
	"context"
	pjson "encoding/json"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"k8s.io/apimachinery/pkg/util/json"

	ctmodel "cloudformation-error/pkg/aws/model/cloudtrail"
	"cloudformation-error/pkg/aws/services"
)

// Find helps find cloudtrail event of failed cloudformation stacks
func Find(ctx context.Context, stackName string, region string) string {
	awsSess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	cfClient := services.NewCloudFormation(awsSess)
	ctClient := services.NewCloudTrail(awsSess)

	stackStatus := []*string{}
	cfStack, err := cfClient.ListStackWithNameAsList(ctx, stackStatus, stackName)
	cfStackEvents, err := cfClient.DescribeStackEventsAsList(ctx, *cfStack.StackId)
	if err != nil {
		log.Fatal(err)
	}

	for i, stackEvent := range cfStackEvents {
		// Look for events where resource failed, search for when the resource started the update
		if strings.Contains(*stackEvent.ResourceStatus, "FAILED") && !strings.EqualFold(*stackEvent.LogicalResourceId, stackName) {
			// fmt.Println(stackEvent.String())
			status := strings.ReplaceAll(*stackEvent.ResourceStatus, "FAILED", "IN_PROGRESS")
			for j := i + 1; j < len(cfStackEvents); j++ {
				if *cfStackEvents[j].PhysicalResourceId == *stackEvent.PhysicalResourceId &&
					status == *cfStackEvents[j].ResourceStatus {
					startTime := cfStackEvents[j].Timestamp
					endTime := stackEvent.Timestamp
					attributeKey := "ReadOnly"
					attributeValue := "false"
					lookup := []*cloudtrail.LookupAttribute{
						&cloudtrail.LookupAttribute{
							AttributeKey:   &attributeKey,
							AttributeValue: &attributeValue,
						},
					}
					events, err := ctClient.LookupEventsAsList(ctx, startTime, endTime, lookup)
					if err != nil {
						log.Fatal(err)
					}
					f, err := os.Create("test.json")
					defer f.Close()
					for _, event := range events {
						cte := ctmodel.CloudTrailEvent{}
						err := json.Unmarshal([]byte(*event.CloudTrailEvent), &cte)
						if err != nil {
							log.Fatal(err)
						}
						if cte.ErrorCode != nil {
							prettyJSON, err := pjson.MarshalIndent(cte, "", "    ")
							_, err = f.WriteString(string(prettyJSON))
							if err != nil {
								log.Fatal(err)
							}
						}
					}
					f.WriteString(startTime.Local().String() + "\n")
					f.WriteString(endTime.Local().String())

					return ""
				}
			}
		}
	}

	return "sdf"
}
