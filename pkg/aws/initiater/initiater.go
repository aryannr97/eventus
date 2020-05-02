package initiater

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

var AWSSNSClient *sns.SNS
var AWSSQSClient *sqs.SQS
var AWSSession *session.Session

func GetSNSClientInstance(region string) (*sns.SNS, error) {
	if AWSSNSClient != nil {
		return AWSSNSClient, nil
	} else {
		AWSSession, err := awsSessionInit(region)
		if err != nil {
			log.Println("Error initiating AWS session")
		} else {
			AWSSNSClient = sns.New(AWSSession)
		}
		return AWSSNSClient, err
	}
}

func GetSQSClientInstance(region string) (*sqs.SQS, error) {
	if AWSSQSClient != nil {
		return AWSSQSClient, nil
	}
	AWSSession, err := awsSessionInit(region)
	if err != nil {
		log.Println("Error initiating AWS session")
	} else {
		AWSSQSClient = sqs.New(AWSSession)
	}
	return AWSSQSClient, err
}

func awsSessionInit(region string) (client.ConfigProvider, error) {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	if AWSSession != nil {
		return AWSSession, nil
	}
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	log.Println(err)

	return sess, err
}
