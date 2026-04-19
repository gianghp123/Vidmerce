package models

import "github.com/gianghp123/Vidmerce/backend/internal/core/enums"

type JobEntity struct {
	BaseItem
	TargetID  string          `dynamodbav:"targetId,omitempty"`
	Type      enums.JobType   `dynamodbav:"type"`
	Status    enums.JobStatus `dynamodbav:"status"`
	Payload   map[string]any  `dynamodbav:"payload,omitempty"`
	ErrorLog  string          `dynamodbav:"errorLog,omitempty"`
	CreatedAt string          `dynamodbav:"createdAt"`
}
