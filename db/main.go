package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type Response events.APIGatewayProxyResponse

var awsSession = session.Must(session.NewSession(&aws.Config{}))
var dynamoSvc = dynamodb.New(awsSession)

func main() {
	lambda.Start(Handler)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	details, err := (&Api{
		DynamoDBAPI: dynamoSvc,
	}).Get("example")
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	var buf bytes.Buffer
	body, err := json.Marshal(&details)
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

type Api struct {
	dynamodbiface.DynamoDBAPI
}

type Details struct {
	TableName string
	Created   *time.Time
}

func (d *Api) Get(tableName string) (details *Details, err error) {
	var out *dynamodb.DescribeTableOutput

	req, out := d.DescribeTableRequest(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})

	if err = req.Send(); err != nil {
		return nil, err
	}

	dateCreated := out.Table.CreationDateTime

	return &Details{
		TableName: tableName,
		Created:   dateCreated,
	}, nil
}