package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

type RequestBody struct {
	OrderID         string  `json:"id"`
	TotalAmount     float64 `json:"TotalAmount"`
	CustomerID      string  `json:"CustomerID"`
	Items           []Item  `json:"Items"`
	ShippingAddress string  `json:"ShippingAddress"`
}

type Item struct {
	ProductID string  `json:"ProductID"`
	Quantity  int     `json:"Quantity"`
	Price     float64 `json:"Price"`
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody RequestBody

	// Parse the request body
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid request body",
		}, nil
	}

	order := Order{
		PK:              "ORDER#" + requestBody.OrderID,
		SK:              "ORDER#" + requestBody.OrderID,
		OrderID:         requestBody.OrderID,
		OrderDate:       time.Now().Format("2006-01-02"),
		TotalAmount:     requestBody.TotalAmount,
		CustomerID:      requestBody.CustomerID,
		Items:           convertItems(requestBody.Items),
		ShippingAddress: requestBody.ShippingAddress,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	av, err := dynamodbattribute.MarshalMap(order)
	if err != nil {
		log.Fatalf("Got error marshalling map: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("OnlineStore"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully added order to table")

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Order was processed successfully",
	}, nil
}

func convertItems(items []Item) []OrderItem {
	var orderItems []OrderItem
	for _, item := range items {
		orderItems = append(orderItems, OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}
	return orderItems
}

func main() {
	lambda.Start(HandleRequest)
}
