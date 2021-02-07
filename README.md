# Cloudformation Detective

CloudFormation detective is a tool that helps you figure out why your CloudFomration stack updates, create or delete might have failed when Cloudformation events are not clear.

![Unit Test](https://github.com/GnatorX/cfnd/workflows/Go/badge.svg?branch=main)


## Setup 

Requires go 1.14
Run `go get -u github.com/GnatorX/cfnd`
Then you can try `cfnd help` to see what options are available

## Usage

```bash
cfnd help
NAME:
   cfnd - Find Cloudtrail events for your CloudFormation errors

USAGE:
   cfnd [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --stackname value, -s value  Name of the stack
   --output value, -o value     Output file name (default: cf_error.json)
   --readonly, --ro             Return readonly events from CloudTrail. Add the flag if you want readonly to be true (default: false)
   --all, -a                    Return all events from CloudTrail. By default, only Events with error is returned. Add the flag if you want all events (default: false)
   --region value, -r value     AWS region for the search (default: us-west-2)
   --help, -h                   show help (default: false)
```

Let say you have a stack that failed with an unknown error. Stack name is `example-stack`. You would run `cfnd -s example-stack -o cf-error.json` (assuming you are in us-west-2 and this failure happened > 15 mins ago). This will dump the cloudtrail logs into cf-error.json similar to the follow example json.

```javascript
//example-stack: Failure reason: You are not authorized to launch instances with this launch template. Not authorized for images: [ami-x] (Service: AmazonEKS; Status Code: 400; Error Code: InvalidRequestException; Request ID: 12345-1cb2-4b15-8f10-7f9f8cfb7ada; Proxy: null)
//2021-02-02 16:03:44.452 -0800 PST
//2021-02-02 16:03:46.611 -0800 PST
//Found 1 CloudTrail Events
{
    "eventVersions": "",
    "UserIdentity": {
        "type": "AssumedRole",
        "principalId": "",
        "arn": "",
        "accountId": "",
        "userName": ""
    },
    "eventTime": "2021-02-02T16:03:46-08:00",
    "eventSource": "ec2.amazonaws.com",
    "eventName": "RunInstances",
    "awsRegion": "us-west-2",
    "sourceIPAddress": "apigateway.amazonaws.com",
    "userAgent": "apigateway.amazonaws.com",
    "requestParameters": {
        "blockDeviceMapping": {},
        "clientToken": "",
        "disableApiTermination": false,
        "instanceType": "t3.medium",
        "instancesSet": {
            "items": [
                {
                    "imageId": "x",
                    "maxCount": 1,
                    "minCount": 1
                }
            ]
        },
        "launchTemplate": {
            "launchTemplateId": "lt-x",
            "version": "4"
        },
        "monitoring": {
            "enabled": false
        },
        "subnetId": "subnet-x"
    },
    "responseElements": null,
    "additionalEventData": null,
    "eventID": "12345-f145-4a45-bca8-123456",
    "eventType": "AwsApiCall",
    "recipientAccountId": "",
    "ErrorCode": "Client.AuthFailure",
    "errorMessage": "Not authorized for images: [ami-x]"
}

```
