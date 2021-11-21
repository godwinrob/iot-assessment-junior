package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Declare a new DynamoDB instance. Note that this is safe for concurrent
// use.
var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))

func getItem(email string) (*user, error) {
	// Prepare the input for the query.
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Users"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	}

	// Retrieve the item from DynamoDB. If no matching item is found
	// return nil.
	result, err := db.GetItem(input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}

	// The result.Item object returned has the underlying type
	// map[string]*AttributeValue. We can use the UnmarshalMap helper
	// to parse this straight into the fields of a struct. Note:
	// UnmarshalListOfMaps also exists if you are working with multiple
	// items.
	usr := new(user)
	err = dynamodbattribute.UnmarshalMap(result.Item, usr)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

// Add item to DynamoDB
func putItem(usr *user) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("Users"),
		Item: map[string]*dynamodb.AttributeValue{

			"email": {
				S: aws.String(usr.Email),
			},
			"updatedAt": {
				S: aws.String(usr.UpdatedAt),
			},
			"hogwartsHouse": {
				S: aws.String(usr.HogwartsHouse),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}

// Add item to DynamoDB
func updateItem(usr *user) (*user, error) {

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":hh": {
				S: aws.String(usr.HogwartsHouse),
			},
			":ua": {
				S: aws.String(usr.UpdatedAt),
			},
		},
		TableName: aws.String("Users"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(usr.Email),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("SET hogwartsHouse = :hh, updatedAt = :ua"),
	}

	_, err := db.UpdateItem(input)

	if err != nil {
		return nil, err
	}

	returnUser, err := getItem(usr.Email)

	return returnUser, nil
}
