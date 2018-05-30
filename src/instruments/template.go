package instruments

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/carldanley/scribe/src/compendium"
)

type TemplateInstrument struct {
	Template *template.Template
	Settings *TemplateInstrumentSettings
}

type TemplateInstrumentSettings struct {
	TemplatePath string
	OutputPath   string
}

func (i *TemplateInstrument) Prepare(settings map[string]interface{}) {
	// attempt to parse the settings
	i.PrepareSettings(settings)

	// prepare the folder path (so writes can happen successfully)
	i.PrepareWriteFolder(i.Settings.OutputPath)
}

func (i *TemplateInstrument) PrepareSettings(settings map[string]interface{}) {
	i.Settings = &TemplateInstrumentSettings{}

	// Make sure a templatePath was specified
	if value, ok := settings["templatePath"]; ok {
		i.Settings.TemplatePath = value.(string)
	}

	// Make sure a outputPath was specified
	if value, ok := settings["outputPath"]; ok {
		i.Settings.OutputPath = value.(string)
	}

	// Attempt to read the template
	fileContents, err := ioutil.ReadFile(i.Settings.TemplatePath)
	if err != nil {
		log.Println("Could not read template file:", i.Settings.TemplatePath, "...")
		os.Exit(1)
	}

	// Attempt to read the template
	i.Template = template.Must(template.New("output").Funcs(sprig.TxtFuncMap()).Parse(string(fileContents)))
}

func (i *TemplateInstrument) PrepareWriteFolder(filePath string) {
	// attempt to `mkdir -p` the path
	if err := os.MkdirAll(path.Dir(filePath), 0700); err != nil {
		log.Println("Could not create the parent directory(ies) for:", path.Dir(filePath))
		os.Exit(1)
	}
}

func (i *TemplateInstrument) Write(secrets map[string]string) {
	log.Println("Writing to file:", i.Settings.OutputPath, "...")

	// create a new file
	file, err := os.Create(i.Settings.OutputPath)
	if err != nil {
		log.Println("Could not write to instrument file at:", i.Settings.OutputPath)
		os.Exit(1)
	}

	// make sure we close the file when we're done with it
	defer file.Close()

	// setup a writer for the template
	writer := bufio.NewWriter(file)

	// generate the template
	i.Template.Execute(writer, secrets)

	// flush the file
	writer.Flush()
}

func init() {
	RegisterInstrument("template", func(settings map[string]interface{}, tomeSpec *compendium.TomeSpec) Instrument {
		// create a new file instrument
		instrument := &TemplateInstrument{}

		// prepare the instrument with the settings
		instrument.Prepare(settings)

		// return the instrument
		return instrument
	})
}
