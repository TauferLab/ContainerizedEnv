package main

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type workflowConfig struct {
	WorkflowName         string
	ApplicationContainer containerConfig
	InputContainer       []containerConfig
	OutputContainer      containerConfig
}

type containerConfig struct {
	Name   string
	InPath string `json:",omitempty"`
	Size   int64  `json:",omitempty"`
}

type containerMetadata struct {
	UUID             uuid.UUID
	Name             string
	CreationTime     time.Time
	ExecutionCommand string
	RecordTrail      *recordTrail
}

type recordTrail struct {
	InputContainers []struct {
		Name string
		UUID uuid.UUID
	}
	ApplicationContainer *struct {
		Name string
		UUID uuid.UUID
	}
	OutputContainer *struct {
		Name string
		UUID uuid.UUID
	}
}
