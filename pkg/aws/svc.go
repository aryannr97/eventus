package svc

import (
	"encoding/json"
	"log"

	"github.com/aryannr97/eventus/pkg/aws/initiater"
	"github.com/aryannr97/eventus/pkg/domain"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type AWSServiceProvider struct {
	AWSId        string
	AWSRegion    string
	AWSSNSClient *sns.SNS
	AWSSQSClient *sqs.SQS
}

//Create instance of AWS Provider
func NewAWSServiceProvider(id string, region string) *AWSServiceProvider {
	snsClient, _ := initiater.GetSNSClientInstance(region)
	sqsClient, _ := initiater.GetSQSClientInstance(region)
	return &AWSServiceProvider{
		AWSId:        id,
		AWSRegion:    region,
		AWSSNSClient: snsClient,
		AWSSQSClient: sqsClient,
	}
}

//Publish Event on SNS
func (p *AWSServiceProvider) PublishEvent(subject string, region string, eventData interface{}) {
	if _, valid := eventData.(string); !valid {
		log.Println("Event Data error")
		return
	}

	input := &sns.PublishInput{
		//Type assertion for interface value to concrete value
		Message:  aws.String(eventData.(string)),
		TopicArn: aws.String("arn:aws:sns:" + region + ":" + p.AWSId + ":" + subject),
	}
	result, err := p.AWSSNSClient.Publish(input)
	if err != nil {
		log.Println("Publish error:", err)
		return
	}
	log.Println("Message:", result)
}

//Subscribe Event
func (p *AWSServiceProvider) SubscribeEvent(subject, region string, onEventHandler func(interface{})) {
	queueName, queueUrl, queueErr := p.CreateQueue(subject)

	if queueErr != nil {
		log.Println("Create Queue Call Error:", queueErr)
		return
	}

	topicArn := "arn:aws:sns:" + region + ":" + p.AWSId + ":" + subject
	queueArn := "arn:aws:sqs:" + region + ":" + p.AWSId + ":" + queueName

	result, err := p.AWSSNSClient.Subscribe(&sns.SubscribeInput{
		Endpoint:              &queueArn,
		Protocol:              aws.String("sqs"),
		ReturnSubscriptionArn: aws.Bool(true), // Return the ARN, even if user has yet to confirm
		TopicArn:              &topicArn,
	})

	if err != nil {
		log.Println("Error:", err)
		return
	}

	log.Println("Result:", result)
	p.PollQueueRealTime(queueName, queueUrl, onEventHandler)
}

//Create Queue
func (p *AWSServiceProvider) CreateQueue(subject string) (string, *string, error) {
	queueName := subject + "_Q"
	response, queueErr := p.AWSSQSClient.CreateQueue(&sqs.CreateQueueInput{
		QueueName: &queueName,
	})

	if queueErr != nil {
		log.Println("Queue creation error:", queueErr)
		return "", nil, queueErr
	}

	return queueName, response.QueueUrl, queueErr
}

// Real Time SQS Poller
func (p *AWSServiceProvider) PollQueueRealTime(queueName string, queueUrl *string, onEventHandler func(interface{})) {
	var waitTime int64
	waitTime = 20
	for {
		log.Println("Polling SQS", queueName)
		res, messageErr := p.AWSSQSClient.ReceiveMessage(&sqs.ReceiveMessageInput{QueueUrl: queueUrl, WaitTimeSeconds: &waitTime})

		if messageErr != nil {
			log.Println("Receive error:", messageErr)
		} else {
			log.Println(res.Messages)
			p.processQueueMessages(res.Messages, queueUrl, onEventHandler)
		}

		log.Println("Polled queue messages processed")
	}
}

//Process Queue Messages
func (p *AWSServiceProvider) processQueueMessages(queueMessages []*sqs.Message, queueUrl *string, onEventHandler func(interface{})) {
	messageBody := domain.SQSMessageBody{}
	for _, message := range queueMessages {
		log.Println("Message ID:", *message.MessageId)
		json.Unmarshal([]byte(*message.Body), &messageBody)
		onEventHandler(messageBody.Message)
		p.deleteMessagefromQueue(message, queueUrl)
	}
}

//Delete Queue Messages through ReceiptHandle
func (p *AWSServiceProvider) deleteMessagefromQueue(message *sqs.Message, queueUrl *string) {
	_, deleteErr := p.AWSSQSClient.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      queueUrl,
		ReceiptHandle: message.ReceiptHandle,
	})

	if deleteErr != nil {
		log.Println("DeleteMessage on SQS failed", deleteErr)
	} else {
		log.Println("DeleteMessage on SQS successful")
	}
}
