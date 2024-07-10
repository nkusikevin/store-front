package main

import (
	//time package is used to get the current time
	"time"
)

type Order struct {
	PK              string      `dynamodbav:"pk"`
	SK              string      `dynamodbav:"sk"`
	OrderID         string      `dynamodbav:"orderId"`
	OrderDate       string      `dynamodbav:"orderDate"`
	TotalAmount     float64     `dynamodbav:"totalAmount"`
	CustomerID      string      `dynamodbav:"customerId"`
	Items           []OrderItem `dynamodbav:"items"` // List of ordered items
	ShippingAddress string      `dynamodbav:"shippingAddress"`
	CreatedAt       time.Time   `dynamodbav:"createdAt"`
	UpdatedAt       time.Time   `dynamodbav:"updatedAt"`
}

type OrderItem struct {
	ProductID string  `dynamodbav:"productId"`
	Quantity  int     `dynamodbav:"quantity"`
	Price     float64 `dynamodbav:"price"`
}

// type Customer struct {
// 	PK         string    `dynamodbav:"pk"`
// 	SK         string    `dynamodbav:"sk"`
// 	CustomerID string    `dynamodbav:"customerId"`
// 	Name       string    `dynamodbav:"name"`
// 	Email      string    `dynamodbav:"email"`
// 	Phone      string    `dynamodbav:"phone"`
// 	Address    string    `dynamodbav:"address"`
// 	CreatedAt  time.Time `dynamodbav:"createdAt"`
// 	UpdatedAt  time.Time `dynamodbav:"updatedAt"`
// }
