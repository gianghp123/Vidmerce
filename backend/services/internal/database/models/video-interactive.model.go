package models

type VideoInteractiveEntity struct {
	BaseItem

	Hotspots []Hotspot `dynamodbav:"hotspots"`
}
