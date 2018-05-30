package instruments

import (
	"fmt"
	"log"

	"github.com/carldanley/scribe/src/compendium"
)

type CatInstrument struct{}

func (i *CatInstrument) Prepare(settings map[string]interface{}) {
	// do nothing
}

func (i *CatInstrument) Write(secrets map[string]string) {
	log.Println("Output:")

	for key, value := range secrets {
		fmt.Printf("%s=\"%s\"\n", key, value)
	}
}

func init() {
	RegisterInstrument("cat", func(settings map[string]interface{}, tomeSpec *compendium.TomeSpec) Instrument {
		// create a new instrument
		instrument := &CatInstrument{}

		// prepare the instrument with the settings
		instrument.Prepare(settings)

		// return the instrument
		return instrument
	})
}
