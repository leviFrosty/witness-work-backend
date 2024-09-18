package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	queryParams := request.QueryStringParameters

	// TODO: setup aws credential store to handle env files setup

	// Retrieve environment variable
	hereApiKey := os.Getenv("HERE_API_KEY")

	if hereApiKey == "" {
		fmt.Println("MY_VARIABLE is not set")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, nil
	} else {
		fmt.Println("MY_VARIABLE:", hereApiKey)
	}

	query := queryParams["q"]
	if query == "" {
		fmt.Println("received empty query")
		return events.APIGatewayProxyResponse{
			Body:       "Expected query string 'q', but was not not provided.",
			StatusCode: 500,
		}, nil
	}

	fmt.Println("query param", query)
	baseUrl := "https://geocode.search.hereapi.com/v1/geocode"
	params := url.Values{}
	params.Add("apiKey", hereApiKey)
	params.Add("q", query)
	fullUrl := baseUrl + "?" + params.Encode()
	res, err := http.Get(fullUrl)

	if err != nil {
		fmt.Println("Error fetching here api geocode")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error fetching HERE API geocode",
		}, nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body.")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error reading response body.",
		}, nil
	}

	fmt.Println("Body: ", string(body))
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(handler)
}
