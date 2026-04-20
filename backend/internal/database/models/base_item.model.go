package models

type BaseItem struct { // Global secondary index 1 partition key
	Gsi1Pk string `dynamodbav:"Gsi1Pk,omitempty"`
	// Global secondary index 1 sort key
	Gsi1Sk string `dynamodbav:"Gsi1Sk,omitempty"`
	// Partition key
	Pk string `dynamodbav:"Pk"`
	// Sort key
	Sk string `dynamodbav:"Sk"`
}
