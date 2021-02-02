package cloudtrail

import "time"

type UserIdentity struct {
	Type        string `json:"type"`
	PrincipalID string `json:"principalId"`
	ARN         string `json:"arn"`
	AccountID   string `json:"accountId"`
	UserName    string `json:"userName"`
}

type CloudTrailEvent struct {
	EventVersion        string                 `json:"eventVersions"`
	UserIdentity        UserIdentity           `json:userIdentity"`
	EventTime           time.Time              `json:"eventTime"`
	EventSource         string                 `json:"eventSource"`
	EventName           string                 `json:"eventName"`
	AwsRegion           string                 `json:"awsRegion"`
	SourceIPAddress     string                 `json:"sourceIPAddress"`
	UserAgent           string                 `json:"userAgent"`
	RequestParameters   map[string]interface{} `json:"requestParameters"`
	ResponseElements    map[string]interface{} `json:"responseElements"`
	AdditionalEventData map[string]interface{} `json:"additionalEventData"`
	EventID             string                 `json:"eventID"`
	EventType           string                 `json:"eventType"`
	RecipientAccountID  string                 `json:"recipientAccountId"`
	ErrorCode           *string                `json:"ErrorCode"`
	ErrorMessage        *string                `json:"errorMessage"`
	Rest                map[string]interface{} `json:"-"`
}
