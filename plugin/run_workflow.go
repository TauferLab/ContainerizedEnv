package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/apptainer/sif/v2/pkg/sif"
	uuid "github.com/satori/go.uuid"
)

func execWorkflow(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var cfg workflowConfig
	if err := json.Unmarshal(file, &cfg); err != nil {
		return err
	}

	cmd := cfg.createRunCommand()

	if err := exec.Command(
		strings.Fields(cmd)[0],
		strings.Fields(cmd)[1:]...,
	).Run(); err != nil {
		return err
	}

	if err := cfg.annotateOutputContainer(); err != nil {
		return err
	}

	return nil
}

func (cfg workflowConfig) annotateOutputContainer() error {
	path := cfg.OutputContainer.Name + ".sif"

	rt, err := cfg.getRecordTrail()
	if err != nil {
		return err
	}

	outputContainerImg, err := sif.LoadContainerFromPath(path, sif.OptLoadWithFlag(os.O_RDWR))
	if err != nil {
		return err
	}

	cmd, err := cfg.getRunscript()
	if err != nil {
		return err
	}

	containerUuid, err := uuid.FromString(outputContainerImg.ID())
	if err != nil {
		return err
	}

	metadata := containerMetadata{
		UUID:             containerUuid,
		Name:             cfg.OutputContainer.Name,
		CreationTime:     outputContainerImg.CreatedAt(),
		ExecutionCommand: cmd,
		RecordTrail:      &rt,
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	metadataDescription, err := sif.NewDescriptorInput(sif.DataGenericJSON, bytes.NewReader(metadataJSON), sif.OptObjectName("metadata.json"))
	if err != nil {
		return err
	}

	if err := outputContainerImg.AddObject(metadataDescription); err != nil {
		return err
	}

	if err := outputContainerImg.UnloadContainer(); err != nil {
		return err
	}

	return nil
}

func (cfg workflowConfig) createRunCommand() string {
	var workflowCommand string

	workflowCommand += "apptainer run"
	workflowCommand += " "

	for _, inputContainer := range cfg.InputContainer {
		workflowCommand += "-B"
		workflowCommand += " "
		workflowCommand += inputContainer.Name + ".sif"
		workflowCommand += ":"
		workflowCommand += "/" + inputContainer.Name
		workflowCommand += ":image-src="
		workflowCommand += "/" + inputContainer.Name
		workflowCommand += " "
	}

	workflowCommand += "-B"
	workflowCommand += " "
	workflowCommand += cfg.OutputContainer.Name + ".sif"
	workflowCommand += ":"
	workflowCommand += "/" + cfg.OutputContainer.Name
	workflowCommand += ":image-src="
	workflowCommand += "/" + cfg.OutputContainer.Name
	workflowCommand += " "

	workflowCommand += cfg.ApplicationContainer.Name + ".sif"

	return workflowCommand
}

func (cfg workflowConfig) getRecordTrail() (recordTrail, error) {
	rt := recordTrail{}

	rt.InputContainers = make([]struct {
		Name string
		UUID uuid.UUID
	}, 0)
	for _, inputContainer := range cfg.InputContainer {
		name := inputContainer.Name
		img, err := sif.LoadContainerFromPath(name+".sif", sif.OptLoadWithFlag(os.O_RDONLY))
		if err != nil {
			return rt, err
		}

		containerUuid, err := uuid.FromString(img.ID())
		if err != nil {
			return rt, err
		}

		rt.InputContainers = append(rt.InputContainers, struct {
			Name string
			UUID uuid.UUID
		}{
			Name: name,
			UUID: containerUuid,
		})
		if err := img.UnloadContainer(); err != nil {
			return rt, err
		}
	}

	applicationContainerImg, err := sif.LoadContainerFromPath(cfg.ApplicationContainer.Name+".sif", sif.OptLoadWithFlag(os.O_RDONLY))
	if err != nil {
		return rt, err
	}

	containerUuid, err := uuid.FromString(applicationContainerImg.ID())
	if err != nil {
		return rt, err
	}

	rt.ApplicationContainer = &struct {
		Name string
		UUID uuid.UUID
	}{
		Name: cfg.ApplicationContainer.Name,
		UUID: containerUuid,
	}
	if err := applicationContainerImg.UnloadContainer(); err != nil {
		return rt, err
	}

	outputContainerImg, err := sif.LoadContainerFromPath(cfg.OutputContainer.Name+".sif", sif.OptLoadWithFlag(os.O_RDONLY))
	if err != nil {
		return rt, err
	}

	containerUuid, err = uuid.FromString(applicationContainerImg.ID())
	if err != nil {
		return rt, err
	}

	rt.OutputContainer = &struct {
		Name string
		UUID uuid.UUID
	}{
		Name: cfg.OutputContainer.Name,
		UUID: containerUuid,
	}
	if err := outputContainerImg.UnloadContainer(); err != nil {
		return rt, err
	}

	return rt, nil
}

func (cfg workflowConfig) getRunscript() (string, error) {
	command := ""

	fimg, err := sif.LoadContainerFromPath(cfg.ApplicationContainer.Name+".sif", sif.OptLoadWithFlag(os.O_RDONLY))
	if err != nil {
		return command, err
	}

	descriptions, err := fimg.GetDescriptors(sif.WithDataType(sif.DataDeffile))
	if err != nil {
		return "", fmt.Errorf("Could not retrive container descriptions: %v", err)
	}

	for _, descriptor := range descriptions {
		defFile, err := descriptor.GetData()
		if err != nil {
			return command, err
		}
		runscript := extractRunscript(string(defFile))
		command += runscript
		command += "\n"
	}

	if err := fimg.UnloadContainer(); err != nil {
		return command, err
	}

	return command[:len(command)-1], nil
}

func extractRunscript(str string) string {
	start := strings.Index(str, "%runscript")
	if start == -1 {
		return ""
	}
	start += len("%runscript")

	str = str[start:]

	end := strings.Index(str, "%")
	if end != -1 {
		str = str[:end]
	}

	str = strings.ReplaceAll(str, "\n", " ")
	str = strings.Trim(str, " ")

	return str
}
