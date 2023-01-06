package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Team struct {
	ID       string `dynamodbav:"id" json:"id"`
	TeamName string `dynamodbav:"teamName" json:"teamName"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess, err := session.NewSession()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	db := dynamodb.New(sess)
	params := &dynamodb.QueryInput{
		TableName: aws.String("sample"),
		ExpressionAttributeNames: map[string]*string{
			"#dataType": aws.String("dataType"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":dataType": {
				S: aws.String("team"),
			},
		},
		KeyConditionExpression: aws.String("#dataType = :dataType"),
		IndexName:              aws.String("dataType-index"),
	}

	result, err := db.Query(params)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	var teams []Team
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &teams)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	bytes, err := json.Marshal(teams)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(bytes),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
