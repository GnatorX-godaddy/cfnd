package ctl

import (
	"context"
	pjson "encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"k8s.io/apimachinery/pkg/util/json"

	ctmodel "cfnd/pkg/aws/model/cloudtrail"
	"cfnd/pkg/aws/services"
)

// Find helps find cloudtrail event of failed cloudformation stacks
func Find(ctx context.Context, stackName string, region string, outputFile string) string {
	awsSess := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	cfClient := services.NewCloudFormation(awsSess)
	ctClient := services.NewCloudTrail(awsSess)

	stackStatus := []*string{}
	cfStack, err := cfClient.ListStackWithNameAsList(ctx, stackStatus, stackName)
	cfStackEvents, err := cfClient.DescribeStackEventsAsList(ctx, *cfStack.StackId)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(outputFile)
	defer f.Close()

	for i, stackEvent := range cfStackEvents {
		// Look for events where resource failed, search for when the resource started the update
		if strings.Contains(*stackEvent.ResourceStatus, "FAILED") && !strings.EqualFold(*stackEvent.LogicalResourceId, stackName) &&
			!strings.Contains(*stackEvent.ResourceStatusReason, "Resource creation cancelled") {
			f.WriteString(stackName + ": Failure reason: " + *stackEvent.ResourceStatusReason + "\n")
			status := strings.ReplaceAll(*stackEvent.ResourceStatus, "FAILED", "IN_PROGRESS")
			for j := i + 1; j < len(cfStackEvents); j++ {
				if *cfStackEvents[j].PhysicalResourceId == *stackEvent.PhysicalResourceId &&
					status == *cfStackEvents[j].ResourceStatus {
					startTime := cfStackEvents[j].Timestamp.Add(time.Second)
					endTime := stackEvent.Timestamp

					f.WriteString(startTime.Local().String() + "\n")
					f.WriteString(endTime.Local().String() + "\n")

					// https://docs.aws.amazon.com/awscloudtrail/latest/userguide/how-cloudtrail-works.html
					// Cloudtrail only tracks for last 90 days + within 15 min of current time
					if time.Now().Sub(*endTime).Hours()/24 > 90 {
						fmt.Println("Your stack failure happened > 90 days ago and we don't have information on it from CloudTrail")
						return ""
					}
					if time.Now().Sub(startTime).Minutes() < 15 {
						fmt.Println("Your stack failed too recent, Cloudtrail only supports within the last 15 mins of events")
						return ""
					}

					attributeKey := "ReadOnly"
					attributeValue := "false"
					lookup := []*cloudtrail.LookupAttribute{
						&cloudtrail.LookupAttribute{
							AttributeKey:   &attributeKey,
							AttributeValue: &attributeValue,
						},
					}

					events, err := ctClient.LookupEventsAsList(ctx, &startTime, endTime, lookup)

					if err != nil {
						log.Fatal(err)
					}

					f.WriteString(fmt.Sprintf("Found %d CloudTrail Events \n", len(events)))
					fmt.Println(len(events))

					// Sort it so we have earliest events first
					sort.Slice(events, func(i, j int) bool {
						return events[i].EventTime.Before(*events[j].EventTime)
					})

					for _, event := range events {
						cte := ctmodel.CloudTrailEvent{}
						err := json.Unmarshal([]byte(*event.CloudTrailEvent), &cte)
						cte.EventTime = cte.EventTime.Local()
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

					return ""
				}
			}
		}
	}

	return ""
}
