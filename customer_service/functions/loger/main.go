package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)


var logger *zap.Logger

func init () {
	l, _ := zap.NewProduction()
	logger = l
	defer logger.Sync()
}

type Event struct {
	Name string `json:"name"`
}

func Handler(ctx context.Context, event Event) (string, error) {

	logger.Info("logger construction succeeded", zap.Any("event", event))

	return "", nil
}

func main() {
	lambda.Start(Handler)
}