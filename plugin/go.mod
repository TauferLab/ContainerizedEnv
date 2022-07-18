module workflow_creation

go 1.13

require (
	github.com/satori/go.uuid v1.2.1-0.20180404165556-75cca531ea76 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/sylabs/sif v1.2.3
	github.com/sylabs/singularity v0.0.0
)

replace github.com/sylabs/singularity => ./singularity_source
