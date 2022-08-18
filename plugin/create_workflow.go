package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/apptainer/sif/v2/pkg/sif"
	uuid "github.com/satori/go.uuid"
)

func (cfg workflowConfig) createWorkflow() error {
	// application container
	fmt.Fprintf(os.Stdout, "Building: application container %s\n", cfg.ApplicationContainer.Name)
	if err := cfg.ApplicationContainer.buildAppContainer(); err != nil {
		return fmt.Errorf("error creating application container %s: %v", cfg.ApplicationContainer.Name, err)
	}
	fmt.Fprintf(os.Stdout, "Completed: application container %s\n", cfg.ApplicationContainer.Name)

	// input containers
	for i, inputContainer := range cfg.InputContainer {
		fmt.Fprintf(os.Stdout, "Building: input container %d: %s\n", i+1, inputContainer.Name)
		if err := inputContainer.createInputContainer(); err != nil {
			return fmt.Errorf("error creating input container %d: %s: %v", i, inputContainer.Name, err)
		}
		fmt.Fprintf(os.Stdout, "Completed: input container %d: %s\n", i+1, inputContainer.Name)
	}

	// output container
	fmt.Fprintf(os.Stdout, "Building: output container %s\n", cfg.OutputContainer.Name)
	if err := cfg.OutputContainer.createOutputContainer(); err != nil {
		return fmt.Errorf("error creating output container %s: %v", cfg.ApplicationContainer.Name, err)
	}
	fmt.Fprintf(os.Stdout, "Completed: output container %s\n", cfg.OutputContainer.Name)

	fmt.Fprintln(os.Stdout, "workflow has been set up")

	return nil
}

func (cfg containerConfig) buildAppContainer() error {
	if err := exec.Command(
		"apptainer",
		"build",
		"--fakeroot",
		"--force",
		cfg.Name+".sif",
		cfg.InPath,
	).Run(); err != nil {
		return fmt.Errorf("error building application container %v", err)
	}

	if err := addStaticMetadata(cfg.Name, false); err != nil {
		return fmt.Errorf("error adding static metadata to application container: %v", err)
	}

	return nil
}

func (cfg containerConfig) createInputContainer() error {
	inputDir := "input_dir/" + cfg.Name

	if err := exec.Command(
		"mkdir",
		"-p",
		inputDir,
	).Run(); err != nil {
		return fmt.Errorf("error creating tmp input directory: %v", err)
	}

	if cfg.InPath != "" {
		if err := exec.Command(
			"cp",
			"-r",
			cfg.InPath,
			inputDir,
		).Run(); err != nil {
			return fmt.Errorf("error copying input files: %v", err)
		}
	} else {
		inputDir = "input_dir/"
	}

	inputContainerFS, err := os.Create(cfg.Name)
	if err != nil {
		return fmt.Errorf("error creating input file: %v", err)
	}
	defer inputContainerFS.Close()

	if err := inputContainerFS.Truncate(cfg.Size); err != nil {
		return fmt.Errorf("error truncating the input file: %v", err)
	}

	if err := exec.Command(
		"mkfs.ext3",
		"-d",
		inputDir,
		cfg.Name,
	).Run(); err != nil {
		return fmt.Errorf("error creating input filesystem: %v", err)
	}

	if err := exec.Command(
		"rm",
		"-r",
		"input_dir/",
	).Run(); err != nil {
		return fmt.Errorf("error removing input files: %v", err)
	}

	if err := exec.Command(
		"tune2fs",
		"-m",
		"0",
		cfg.Name,
	).Run(); err != nil {
		return fmt.Errorf("error tuneing input filesystem: %v", err)
	}

	inputPartition, err := os.ReadFile(cfg.Name)
	if err != nil {
		return err
	}

	inputSifDesc, err := sif.NewDescriptorInput(sif.DataPartition, bytes.NewReader(inputPartition), sif.OptObjectName(cfg.Name), sif.OptPartitionMetadata(sif.FsExt3, sif.PartData, "amd64"))
	if err != nil {
		return err
	}

	inputContainerUUID := uuid.NewV4()

	sif.CreateContainerAtPath(cfg.Name+".sif", sif.OptCreateWithID(inputContainerUUID.String()), sif.OptCreateWithDescriptors(inputSifDesc))

	if err := addStaticMetadata(cfg.Name, true); err != nil {
		return fmt.Errorf("error adding static metadata to input container: %v", err)
	}

	if err := exec.Command(
		"rm",
		cfg.Name,
	).Run(); err != nil {
		return fmt.Errorf("error removing input filesystem: %v", err)
	}

	return nil
}

func (cfg containerConfig) createOutputContainer() error {
	output_dir := "output_dir/" + cfg.Name

	outputContainerFS, err := os.Create(cfg.Name)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputContainerFS.Close()

	if err := outputContainerFS.Truncate(cfg.Size); err != nil {
		return fmt.Errorf("error truncating the output file: %v", err)
	}

	if err := exec.Command(
		"mkdir",
		"-p",
		output_dir,
	).Run(); err != nil {
		return fmt.Errorf("error creating tmp output directory: %v", err)
	}

	if err := exec.Command(
		"mkfs.ext3",
		"-d",
		"output_dir",
		cfg.Name,
	).Run(); err != nil {
		return fmt.Errorf("error creating output filesystem: %v", err)
	}

	if err := exec.Command(
		"rm",
		"-r",
		"output_dir",
	).Run(); err != nil {
		return fmt.Errorf("error removing tmp output directory: %v", err)
	}

	if err := exec.Command(
		"tune2fs",
		"-m",
		"0",
		cfg.Name,
	).Run(); err != nil {
		return fmt.Errorf("error tuneing output filesystem: %v", err)
	}

	outputPartition, err := os.ReadFile(cfg.Name)
	if err != nil {
		return err
	}

	outputSifDesc, err := sif.NewDescriptorInput(sif.DataPartition, bytes.NewReader(outputPartition), sif.OptObjectName(cfg.Name), sif.OptPartitionMetadata(sif.FsExt3, sif.PartData, "amd64"))
	if err != nil {
		return err
	}

	outputContainerUUID := uuid.NewV4()

	sif.CreateContainerAtPath(cfg.Name+".sif", sif.OptCreateWithID(outputContainerUUID.String()), sif.OptCreateWithDescriptors(outputSifDesc))

	if err := exec.Command(
		"rm",
		cfg.Name,
	).Run(); err != nil {
		return fmt.Errorf("error removing output filesystem: %v", err)
	}

	return nil
}

func addStaticMetadata(name string, isInputContainer bool) error {
	path := name + ".sif"

	containerImg, err := sif.LoadContainerFromPath(path, sif.OptLoadWithFlag(os.O_RDWR))
	if err != nil {
		return err
	}

	containerUuid, err := uuid.FromString(containerImg.ID())
	if err != nil {
		return err
	}

	metadata := containerMetadata{
		UUID:             containerUuid,
		Name:             name,
		CreationTime:     containerImg.CreatedAt(),
		ExecutionCommand: "no operation",
		RecordTrail:      nil,
	}

	if isInputContainer {
		metadata.RecordTrail = &recordTrail{
			OutputContainer: &struct {
				Name string
				UUID uuid.UUID
			}{
				Name: name,
				UUID: containerUuid,
			},
			InputContainers:      nil,
			ApplicationContainer: nil,
		}
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	applicationSifMetadata, err := sif.NewDescriptorInput(sif.DataGenericJSON, bytes.NewReader(metadataJSON), sif.OptObjectName("metadata"))
	if err != nil {
		return err
	}

	if err := containerImg.AddObject(applicationSifMetadata); err != nil {
		return err
	}

	if err := containerImg.UnloadContainer(); err != nil {
		return err
	}

	return nil
}
