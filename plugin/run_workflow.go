package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/sylabs/sif/pkg/sif"
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

	outputContainerImg, err := sif.LoadContainer(path, false)
	if err != nil {
		return err
	}

	cmd, err := cfg.getRunscript()
	if err != nil {
		return err
	}

	metadata := containerMetadata{
		UUID:             outputContainerImg.Header.ID,
		Name:             cfg.OutputContainer.Name,
		CreationTime:     time.Unix(outputContainerImg.Header.Ctime, 0),
		ExecutionCommand: cmd,
		RecordTrail:      &rt,
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	metadataDescription := sif.DescriptorInput{
		Datatype: sif.DataGenericJSON,
		Groupid:  sif.DescrDefaultGroup,
		Link:     sif.DescrUnusedLink,
		Data:     metadataJSON,
		Fname:    "metadata.json",
		Size:     int64(len(metadataJSON)),
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

	workflowCommand += "singularity run"
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
		img, err := sif.LoadContainer(name+".sif", true)
		if err != nil {
			return rt, err
		}
		rt.InputContainers = append(rt.InputContainers, struct {
			Name string
			UUID uuid.UUID
		}{
			Name: name,
			UUID: img.Header.ID,
		})
		if err := img.UnloadContainer(); err != nil {
			return rt, err
		}
	}

	applicationContainerImg, err := sif.LoadContainer(cfg.ApplicationContainer.Name+".sif", true)
	if err != nil {
		return rt, err
	}
	rt.ApplicationContainer = &struct {
		Name string
		UUID uuid.UUID
	}{
		Name: cfg.ApplicationContainer.Name,
		UUID: applicationContainerImg.Header.ID,
	}
	if err := applicationContainerImg.UnloadContainer(); err != nil {
		return rt, err
	}

	outputContainerImg, err := sif.LoadContainer(cfg.OutputContainer.Name+".sif", true)
	if err != nil {
		return rt, err
	}
	rt.OutputContainer = &struct {
		Name string
		UUID uuid.UUID
	}{
		Name: cfg.OutputContainer.Name,
		UUID: outputContainerImg.Header.ID,
	}
	if err := outputContainerImg.UnloadContainer(); err != nil {
		return rt, err
	}

	return rt, nil
}

func (cfg workflowConfig) getRunscript() (string, error) {
	command := ""

	fimg, err := sif.LoadContainer(cfg.ApplicationContainer.Name+".sif", true)
	if err != nil {
		return command, err
	}

	descriptions, _, err := fimg.GetFromDescr(sif.Descriptor{
		Datatype: sif.DataDeffile,
	})
	if err != nil {
		return "", fmt.Errorf("Could not retrive container descriptions: %v", err)
	}

	for _, descriptor := range descriptions {
		defFile := descriptor.GetData(&fimg)
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
