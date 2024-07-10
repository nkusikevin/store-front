package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	// Create a new session, allowing SDK to use the default credential chain
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a new DynamoDB client
	svc := dynamodb.New(sess)

	// Define the table schema
	input := &dynamodb.CreateTableInput{
		// Defines the attributes used as keys (PK and SK) and their types (S for String).
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("PK"), // Partition Key
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("SK"), // Sort Key
				AttributeType: aws.String("S"),
			},
		},
		//  Specifies the primary key schema for the table using PK as the partition key and SK as the sort key.
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       aws.String("RANGE"),
			},
		},
		// Sets the read and write capacity units to 5 each.
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			// The table can handle up to 5 strongly consistent reads per second for items up to 4 KB each.
			ReadCapacityUnits: aws.Int64(5),
			// The table can handle up to 5 writes per second for items up to 1 KB each.
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String("OnlineStore"),
	}

	// Create the table
	result, err := svc.CreateTable(input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
