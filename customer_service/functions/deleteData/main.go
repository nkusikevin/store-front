package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var svc *dynamodb.DynamoDB
var sess *session.Session

type DeleteRequest struct {
	PK string `json:"PK"`
	SK string `json:"SK"`
}

type DefaultResponse struct {
	StatusCode string `json:"statusCode"`
	Message    string `json:"message"`
}

func init() {
	// Create a new session, allowing SDK to use the default credential chain
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a new DynamoDB client
	svc = dynamodb.New(sess)
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var deleteRequest DeleteRequest

	// Parse the request body
	err := json.Unmarshal([]byte(request.Body), &deleteRequest)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid request body",
		}, nil
	}

	// Construct the primary key map
	primaryKey := map[string]string{
		"PK": deleteRequest.PK,
		"SK": deleteRequest.SK,
	}

	pk, err := dynamodbattribute.MarshalMap(primaryKey)
	if err != nil {
		log.Printf("Error marshalling primary key: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error marshalling primary key",
		}, nil
	}

	// Define the table schema
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("OnlineStore"),
		Key:       pk,
	}

	// Delete the item
	_, err = svc.DeleteItem(input)
	if err != nil {
		log.Printf("Error deleting item: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error deleting item",
		}, nil
	}

	fmt.Println("Item deleted")
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Item deleted",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
