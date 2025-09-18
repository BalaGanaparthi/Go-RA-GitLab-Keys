package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// JWK represents a JSON Web Key structure
type JWK struct {
	Kty string `json:"kty"` // Key Type
	Kid string `json:"kid"` // Key ID
	E   string `json:"e"`   // Exponent (for RSA keys)
	N   string `json:"n"`   // Modulus (for RSA keys)
	Use string `json:"use"` // Public Key Use
	Alg string `json:"alg"` // Algorithm
}

// JWKS represents a JSON Web Key Set
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// getJWKS returns the JSON Web Key Set with the two predefined RSA keys
func getJWKS() JWKS {
	// First RSA key - stored as a constant within the code
	key1 := JWK{
		Kty: "RSA",
		Kid: "Q93ngDAyTNgdjaPblfooyeT00CFLJV5pWulhwAEg8Sw",
		E:   "AQAB",
		N:   "xTcVMXDqrYepPqPmA7cp2QnEe1tGMvb_Y8PAhVt8iJBEydtzP4_lDWTCvwYD65gzcUr0qpU3-RE1qMz2mwCgIE-ngUeIsOJDoYjf_Snx7P5mqv6ADKI3BTJ88IU_kTFIUR97eun3mxQF4ddZUA6VE39ReJSCSNpsMTyqG9NF5WiLrJNkkhPQ1cdd35Jhku1i4Cg--x7RAFIjmw0E4yGEiKJ0scD7cvQFHTRL6fa2M5DVv_cjcGCOjiCzM0hvTA6jqRB2KZOw2YtpGEo_Vi1-68filWvphbplbFsk4nvpx_NpRiQHjeSgtKVUEBg9OrBa6sqTk0kZUvKDwCIAhl2PTVfbN1WuK8FkSSy4z9DMLNNwaQ3Y73objP6mrx7YGZSH6I5z3jlvLENv9jSrdy-NhiEAHagYQVZ_MtAUecutnO_l4cpNGYs6TOGc4OkI782Mgzwxm4acY2szaVgMeozuqFDvlIp7O00F267Cce693ztQdkvT2R7GXnHON8SqMQaCnuCZis431DoCsSKDkzUdiZdLPmkVjQuYVgfERkIM9X7tYyN5Zr2PizVLsT8NSjJy9hAJPcp3Ufh69kXS385NggFhkdkmU1BD72T4e_0E7OImGO7zH5PHrVmbJCNULCO1fBeVYHDKmywHeA6pO9q5WyrCLnWdC9T9xBMWurbmBeU",
		Use: "sig",
		Alg: "RS256",
	}

	// Second RSA key - stored as a constant within the code
	key2 := JWK{
		Kty: "RSA",
		Kid: "CjZ2bP54fm1lEkeHGo_E4UVyc4MbN4fye1i6DrxFaqQ",
		E:   "AQAB",
		N:   "nN8D8DxzqoQJsctdkAGmwEU9aL6zsyLjHP_OkPD2yB3swCEiSbKgvxVDIDh-ClGERR1VHL9tFWbyYfWzOZkLCQT4YmvuilSiQ1wu0j_CEH0dqIvXWIyueoz7ZVyjghpkuuW9RVm84C3zvFm6S7d_kLHKYk0SCpP4RbjcHz-1EGA_p4Vl_nTdz16KOuwkNBxH1LJBqq2e3WXczadYvh-EZGrLip1wIUotjebWZaRJWTaw1jn6ssdaRZhZD3JlRUyGO2clK-ULrl-VCSh_FIi_hf4SGTVMhFK1iAnEGxJ5zJ8FaGHxzu25RAaMtwmN1guW02NH1TPOv1eJtXqQoaykfQ",
		Use: "sig",
		Alg: "RS256",
	}

	// Return the JWKS containing both keys
	return JWKS{
		Keys: []JWK{key1, key2},
	}
}

// handler is the main Lambda function handler for the JWKS endpoint
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received request: %s %s", request.HTTPMethod, request.Path)

	// Handle CORS preflight requests
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type",
			},
		}, nil
	}

	// Only allow GET requests for the JWKS endpoint
	if request.HTTPMethod != "GET" {
		return events.APIGatewayProxyResponse{
			StatusCode: 405,
			Headers: map[string]string{
				"Content-Type":                 "application/json",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, OPTIONS",
			},
			Body: `{"error": "Method not allowed"}`,
		}, nil
	}

	// Get the JWKS data
	jwks := getJWKS()

	// Marshal the JWKS to JSON
	jsonData, err := json.Marshal(jwks)
	if err != nil {
		log.Printf("Error marshaling JWKS to JSON: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type":                "application/json",
				"Access-Control-Allow-Origin": "*",
			},
			Body: `{"error": "Internal server error"}`,
		}, fmt.Errorf("failed to marshal JWKS: %w", err)
	}

	// Return successful response with JWKS data
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
			"Cache-Control":               "public, max-age=3600", // Cache for 1 hour
		},
		Body: string(jsonData),
	}, nil
}

// main function starts the Lambda handler
func main() {
	log.Println("Starting JWKS serverless function...")
	lambda.Start(handler)
}
