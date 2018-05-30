package instruments

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/carldanley/scribe/src/compendium"
)

type FileInstrument struct {
	Spec     *compendium.TomeSpec
	Settings *FileInstrumentSettings
}

type FileInstrumentSettings struct {
	Path string
}

func (i *FileInstrument) ParseSettings(settings map[string]interface{}) {
	i.Settings = &FileInstrumentSettings{}

	if path, ok := settings["path"]; ok {
		i.Settings.Path = path.(string)
	} else {
		log.Println("A path must be specified for an instrument of type \"file\"")
		os.Exit(1)
	}
}

func (i *FileInstrument) PrepareWriteFolder(filePath string) {
	// attempt to `mkdir -p` the path
	if err := os.MkdirAll(path.Dir(filePath), 0700); err != nil {
		log.Println("Could not create the parent directory(ies) for:", path.Dir(filePath))
		os.Exit(1)
	}
}

func (i *FileInstrument) Prepare(settings map[string]interface{}) {
	// attempt to parse the settings
	i.ParseSettings(settings)

	// prepare the folder path (so writes can happen successfully)
	i.PrepareWriteFolder(i.Settings.Path)
}

func (i *FileInstrument) Write(secrets map[string]string) {
	log.Println("Writing to file:", i.Settings.Path, "...")

	// attempt to create a pointer to the File
	file, err := os.Create(i.Settings.Path)
	if err != nil {
		log.Println("Could not write to instrument file at:", i.Settings.Path)
		os.Exit(1)
	}

	// now, write all of the secrets to the file
	for key, value := range secrets {
		fmt.Fprintf(file, "%s=\"%s\"\n", key, value)
	}

	// close the file
	if file.Close() != nil {
		log.Println("Could not close instrument file at:", i.Settings.Path)
		os.Exit(1)
	}
}

func init() {
	RegisterInstrument("file", func(settings map[string]interface{}, tomeSpec *compendium.TomeSpec) Instrument {
		// create a new instrument
		instrument := &FileInstrument{
			Spec: tomeSpec,
		}

		// prepare the instrument with the settings
		instrument.Prepare(settings)

		// return the instrument
		return instrument
	})
}
