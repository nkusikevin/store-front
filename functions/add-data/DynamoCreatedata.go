package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// type Product struct {
// 	PK          string    `dynamodbav:"PK"` // Attribute key names should match exactly
// 	SK          string    `dynamodbav:"SK"`
// 	ProductID   string    `dynamodbav:"ProductID"`
// 	Name        string    `dynamodbav:"Name"`
// 	Price       float64   `dynamodbav:"Price"`
// 	Category    string    `dynamodbav:"Category"`
// 	Stock       int       `dynamodbav:"Stock"`
// 	Description string    `dynamodbav:"Description"`
// 	CreatedAt   time.Time `dynamodbav:"CreatedAt"`
// 	UpdatedAt   time.Time `dynamodbav:"UpdatedAt"`
// }

type Customer struct {
	PK         string    `dynamodbav:"PK"` // Attribute key names should match exactly
	SK         string    `dynamodbav:"SK"`
	CustomerID string    `dynamodbav:"CustomerID"`
	Name       string    `dynamodbav:"Name"`
	Email      string    `dynamodbav:"Email"`
	Phone      string    `dynamodbav:"Phone"`
	Address    string    `dynamodbav:"Address"`
	CreatedAt  time.Time `dynamodbav:"CreatedAt"`
	UpdatedAt  time.Time `dynamodbav:"UpdatedAt"`
}

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

// Create a new product
// product := Product{
// 	PK:          "PRODUCT#1",
// 	SK:          "PRODUCT#2",
// 	ProductID:   "1",
// 	Name:        "Product 2",
// 	Price:       100.00,
// 	Category:    "Category 2",
// 	Stock:       100,
// 	Description: "Description of product 2",
// 	CreatedAt:   time.Now(),
// 	UpdatedAt:   time.Now(),
// }

// customer := Customer{
// 	PK:         "CUSTOMER#2",
// 	SK:         "CUSTOMER#2",
// 	CustomerID: "2",
// 	Name:       "Kabera john Doe",
// 	Email:      "johndoe@example.com",
// 	Phone:      "123-555-1234",
// 	Address:    "123 Main St, Anytown, USA",
// 	CreatedAt:  time.Now(),
// 	UpdatedAt:  time.Now(),
// }

func main() {
	order := Order{
		PK:          "ORDER#2",
		SK:          "ORDER#2",
		OrderID:     "2",
		OrderDate:   "2024-07-11",
		TotalAmount: 150.00,
		CustomerID:  "CUSTOMER#2",
		Items: []OrderItem{
			{ProductID: "PRODUCT#3", Quantity: 1, Price: 50.00},
			{ProductID: "PRODUCT#4", Quantity: 2, Price: 50.00},
		},
		ShippingAddress: "456 Oak Ave, Springfield, USA",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	av, err := dynamodbattribute.MarshalMap(order)
	if err != nil {
		log.Fatal(err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("OnlineStore"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully added product to table")

}
