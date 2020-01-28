package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func getItem(isbn string) (*Book, error) {

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))}),
	)
	var db = dynamodb.New(sess)

	fmt.Println(isbn, os.Getenv("TABLE_NAME"))

	// Prepare the input for the query
	input := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"ISBN": {
				S: aws.String(isbn),
			},
		},
	}

	// Retrieve input from DynamoDB
	result, err := db.GetItem(input)
	if err != nil {
		fmt.Println("Error Getting itmes ")
		fmt.Println(err.Error())
		return nil, err
	}

	if result.Item == nil {
		// fmt.Println(err.Error())
		fmt.Println("No Books Found.")
		return nil, nil
	}

	// The result.Item Object is map[string]*AttributeValue
	// Unmarshal Map helper to parse into structure

	bk := new(Book)
	err = dynamodbattribute.UnmarshalMap(result.Item, bk)
	if err != nil {
		fmt.Println("Error Unmarshalling")
		fmt.Println(err.Error())
		return nil, err
	}

	return bk, nil
}

// Add a book record to DynamoDB.
func putItem(bk *Book) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))}),
	)
	var db = dynamodb.New(sess)

	input := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Item: map[string]*dynamodb.AttributeValue{
			"ISBN": {
				S: aws.String(bk.ISBN),
			},
			"Title": {
				S: aws.String(bk.Title),
			},
			"Author": {
				S: aws.String(bk.Author),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}

// Get all book records.
func scanAll() ([]Book, error) {

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))}),
	)
	var db = dynamodb.New(sess)

	params := &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	results, err := db.Scan(params)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	books := []Book{}

	// Unmarshall all Items
	err = dynamodbattribute.UnmarshalListOfMaps(results.Items, &books)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return books, nil

}
