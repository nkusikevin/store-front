package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	l, _ := zap.NewDevelopment()
	logger = l
	defer logger.Sync()
}

type DefaultResponse struct {
	StatusCode string `json:"statusCode"`
	Message    string `json:"message"`
}

func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse
	logger.Info("recieved request", zap.Any("method", event.HTTPMethod), zap.Any("path", event.Path), zap.Any("body", event.Body))

	if event.Path == "/hello" {
		body, _ := json.Marshal(&DefaultResponse{
			StatusCode: string(http.StatusOK),
			Message:    "hello world",
		})
		res = &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}
	} else {
		body, _ := json.Marshal(&DefaultResponse{
			StatusCode: string(http.StatusOK),
			Message:    "default response",
		})
		res = &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}
	}

	return res, nil
}

func main() {
	lambda.Start(Handler)
}
