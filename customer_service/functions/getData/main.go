package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var svc *dynamodb.DynamoDB
var sess *session.Session

func init() {
	// Create a new session, allowing SDK to use the default credential chain
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a new DynamoDB client
	svc = dynamodb.New(sess)
}

type RequestBody struct {
	PK string `json:"pk"`
	SK string `json:"sk"`
}

type Product struct {
	PK          string    `dynamodbav:"PK"` // Attribute key names should match exactly
	SK          string    `dynamodbav:"SK"`
	ProductID   string    `dynamodbav:"ProductID"`
	Name        string    `dynamodbav:"Name"`
	Price       float64   `dynamodbav:"Price"`
	Category    string    `dynamodbav:"Category"`
	Stock       int       `dynamodbav:"Stock"`
	Description string    `dynamodbav:"Description"`
	CreatedAt   time.Time `dynamodbav:"CreatedAt"`
	UpdatedAt   time.Time `dynamodbav:"UpdatedAt"`
}

func Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody RequestBody

	// Log the incoming request body for debugging
	log.Printf("Received request body ===>>>> : %s", request)

	// Parse the request body
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		log.Printf("Error parsing request body: %s", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Invalid request body: %s", err.Error()),
		}, nil
	}

	// Ensure PK and SK are not empty
	if requestBody.PK == "" || requestBody.SK == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "PK and SK must be provided",
		}, nil
	}

	// Define the primary key of the item to retrieve
	primaryKey := map[string]*dynamodb.AttributeValue{
		"PK": {
			S: aws.String(requestBody.PK),
		},
		"SK": {
			S: aws.String(requestBody.SK),
		},
	}

	// Retrieve the item by primary key
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("OnlineStore"),
		Key:       primaryKey,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error retrieving item: %s", err.Error()),
		}, nil
	}

	// Check if the item was found
	if result.Item == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Item not found",
		}, nil
	}

	// Unmarshal the item's attributes into a struct
	var item Product
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error unmarshalling item: %s", err.Error()),
		}, nil
	}

	// Return the item
	itemJson, err := json.Marshal(item)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error marshalling item to JSON: %s", err.Error()),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(itemJson),
	}, nil
}

func main() {
	lambda.Start(Handle)
}
