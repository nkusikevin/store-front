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

type Order struct {
	PK              string      `dynamodbav:"PK"` // Attribute key names should match exactly
	SK              string      `dynamodbav:"SK"`
	OrderID         string      `dynamodbav:"OrderID"`
	OrderDate       string      `dynamodbav:"OrderDate"`
	TotalAmount     float64     `dynamodbav:"TotalAmount"`
	CustomerID      string      `dynamodbav:"CustomerID"`
	Items           []OrderItem `dynamodbav:"Items"` // List of ordered items
	ShippingAddress string      `dynamodbav:"ShippingAddress"`
	CreatedAt       time.Time   `dynamodbav:"CreatedAt"`
	UpdatedAt       time.Time   `dynamodbav:"UpdatedAt"`
}

type OrderItem struct {
	ProductID string  `dynamodbav:"ProductID"`
	Quantity  int     `dynamodbav:"Quantity"`
	Price     float64 `dynamodbav:"Price"`
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var requestBody RequestBody

	// Process the POST request
	if request.HTTPMethod == "POST" {

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
		var item Order
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

	// Return a default response for unsupported methods
	return events.APIGatewayProxyResponse{
		Body:       "Unsupported HTTP method",
		StatusCode: 405,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}