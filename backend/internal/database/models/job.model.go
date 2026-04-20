package models

import "github.com/gianghp123/Vidmerce/backend/internal/core/enums"

type JobEntity struct {
	BaseItem
	// Creation timestamp
	CreatedAt string `dynamodbav:"createdAt"`
	// Error log if failed
	ErrorLog *string `dynamodbav:"errorLog,omitempty"`
	// JobEntity payload data
	Payload map[string]interface{} `dynamodbav:"payload,omitempty"`
	Status  enums.JobStatus        `dynamodbav:"status"`
	// Target entity ID
	TargetID *string       `dynamodbav:"targetId,omitempty"`
	Type     enums.JobType `dynamodbav:"type"`
}
