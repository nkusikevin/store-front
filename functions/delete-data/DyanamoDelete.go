package main

import (
	"fmt"
	"log"

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

func main() {

	primaryKey := map[string]string{
		"PK": "PRODUCT#1",
		"SK": "PRODUCT#2",
	}

	pk, err := dynamodbattribute.MarshalMap(primaryKey)
	if err != nil {
		log.Fatal(err)
	}

	// Define the table schema
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("OnlineStore"),
		Key:       pk,
	}
	// Delete the item
	_, err = svc.DeleteItem(input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Table deleted")
}
