package main

import (
	"os"

	"github.com/multiprocessio/datastation/runner"
)

const VERSION = "development"
const APP_NAME = "DataStation Runner (Go)"

func main() {
	runner.Verbose = true
	runner.Logln(APP_NAME + " " + VERSION)
	projectId := ""
	panelId := ""
	panelMetaOut := ""
	settingsFile := ""

	args := os.Args
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "--dsproj" {
			projectId = args[i+1]
			i++
			continue
		}

		if args[i] == "--evalPanel" {
			panelId = args[i+1]
			i++
			continue
		}

		if args[i] == "--metaFile" {
			panelMetaOut = args[i+1]
			i++
			continue
		}

		if args[i] == "--settingsFile" {
			settingsFile = args[i+1]
			i++
			continue
		}
	}

	if projectId == "" {
		runner.Fatalln("No project id given.")
	}

	if panelId == "" {
		runner.Fatalln("No panel id given.")
	}

	if panelMetaOut == "" {
		runner.Fatalln("No panel meta out given.")
	}

	if settingsFile == "" {
		settingsFile = runner.SettingsFileDefaultLocation
	}

	settings, err := runner.LoadSettings(settingsFile)
	if err != nil {
		runner.Logln("Could not load settings, assuming defaults.")
		settings = runner.DefaultSettings
	}

	ec := runner.NewEvalContext(*settings, runner.DEFAULT_FS_BASE)

	err = ec.Eval(projectId, panelId)
	if err != nil {
		runner.Logln("Failed to eval: %s", err)

		if _, ok := err.(*runner.DSError); !ok {
			err = runner.Edse(err)
			err.(*runner.DSError).Stack = "Unknown"
		}

		err := runner.WriteJSONFile(panelMetaOut, map[string]interface{}{
			"exception": err,
		})
		if err != nil {
			runner.Fatalln("Could not write panel meta out: %s", err)
		}

		// Explicitly don't fail here so that the parent can read the exception from disk
	}
}
