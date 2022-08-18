package main

import (
	pluginapi "github.com/apptainer/apptainer/pkg/plugin"
	clicallback "github.com/apptainer/apptainer/pkg/plugin/callback/cli"
)

// nolint
var Plugin = pluginapi.Plugin{
	Manifest: pluginapi.Manifest{
		Name:        "TRIC",
		Author:      "Dominic Kennedy (email: dominicmkennedy@gmail.com)",
		Version:     "0.1.0",
		Description: "Traceability and Reproducibility through Individual Containerisation",
	},
	Callbacks: []pluginapi.Callback{
		(clicallback.Command)(callbackRegisterCmd),
	},
	Install: nil,
}
