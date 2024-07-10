package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, request events.LambdaFunctionURLRequest) (Response, error) {
	return Response{Body: "It works!", StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
