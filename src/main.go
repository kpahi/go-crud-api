package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"encoding/json"
)

type Book struct {
	ISBN   string `json:isbn`
	Title  string `json:title`
	Author string `json:string`
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return show(req)
	case "POST":
		return create(req)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       http.StatusText(http.StatusMethodNotAllowed),
		}, nil
	}
}

func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	isbn := req.QueryStringParameters["isbn"]

	fmt.Println("ISBN: ", isbn)

	bk, err := getItem(isbn)

	js, err := json.Marshal(bk)

	if err != nil {
		fmt.Println("Error Marshalling")
		fmt.Println(err.Error())
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	bk := new(Book)
	err := json.Unmarshal([]byte(req.Body), bk)
	if err != nil {
		fmt.Println("Failed to unmarsh")
		fmt.Println(err.Error())
	}

	err = putItem(bk)
	if err != nil {
		fmt.Println("Failed to put item")
		fmt.Println(err.Error())
	}

	resp, _ := json.Marshal(map[string]string{
		"id": bk.ISBN,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers:    map[string]string{"Location": fmt.Sprintf("/books?isbn=%s", bk.ISBN)},
		Body:       string(resp),
	}, nil

}

func main() {
	lambda.Start(router)
}
