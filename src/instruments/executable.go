package instruments

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/carldanley/scribe/src/compendium"
)

type ExecutableInstrument struct {
	Spec           *compendium.TomeSpec
	Settings       *ExecutableInstrumentSettings
	CurrentCommand *exec.Cmd
	Secrets        *map[string]string
}

type ExecutableInstrumentSettings struct {
	Command            string
	CleanEnvironment   bool
	ExitOnCommandError bool
	RestartOnExit      bool
}

func (i *ExecutableInstrument) Prepare(settings map[string]interface{}) {
	i.Settings = &ExecutableInstrumentSettings{}

	if value, ok := settings["command"]; ok {
		i.Settings.Command = value.(string)
	} else {
		log.Println("A command must be specified for an instrument of type \"executable\"")
		os.Exit(1)
	}

	if value, ok := settings["cleanEnvironment"]; ok {
		i.Settings.CleanEnvironment = value.(bool)
	}

	if value, ok := settings["exitOnCommandError"]; ok {
		i.Settings.ExitOnCommandError = value.(bool)
	}

	if value, ok := settings["restartOnExit"]; ok {
		i.Settings.RestartOnExit = value.(bool)
	}
}

func (i *ExecutableInstrument) ExecuteCommand() {
	log.Println("Executing command:", i.Settings.Command, "...")

	args := strings.Split(i.Settings.Command, " ")
	i.CurrentCommand = exec.Command(args[0], args[1:]...)
	env := os.Environ()

	if i.Settings.CleanEnvironment == true {
		env = []string{}
	}

	for key, value := range *i.Secrets {
		env = append(env, fmt.Sprintf("%s=\"%s\"", key, value))
	}

	i.CurrentCommand.Env = env
	i.CurrentCommand.Stdout = os.Stdout
	i.CurrentCommand.Stderr = os.Stderr

	if err := i.CurrentCommand.Start(); err != nil {
		log.Println("Could not start executable:", err)
		os.Exit(1)
	}

	if err := i.CurrentCommand.Wait(); err != nil {
		log.Println("Executable finished with error:", err)

		if i.Settings.ExitOnCommandError == true {
			log.Println("Exiting...")
			os.Exit(1)
		}
	}

	i.CurrentCommand = nil
	if i.Settings.RestartOnExit == true {
		i.Write(*i.Secrets)
	}
}

func (i *ExecutableInstrument) Write(secrets map[string]string) {
	i.Secrets = &secrets

	if i.CurrentCommand != nil {
		i.CurrentCommand.Process.Signal(os.Interrupt)
		i.CurrentCommand.Process.Wait()
		return
	}

	go i.ExecuteCommand()
}

func init() {
	RegisterInstrument("executable", func(settings map[string]interface{}, tomeSpec *compendium.TomeSpec) Instrument {
		// create a new instrument
		instrument := &ExecutableInstrument{
			Spec: tomeSpec,
		}

		// prepare the instrument with the settings
		instrument.Prepare(settings)

		// return the instrument
		return instrument
	})
}
