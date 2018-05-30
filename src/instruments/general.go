package instruments

import (
	"log"
	"os"

	"github.com/carldanley/scribe/src/compendium"
)

var registeredInstruments map[string]GetInstrument

func RegisterInstrument(name string, callback GetInstrument) {
	// make sure registered tomes has been initialized
	if registeredInstruments == nil {
		registeredInstruments = map[string]GetInstrument{}
	}

	// make sure the instrument hasn't been registered already
	if _, ok := registeredInstruments[name]; ok {
		log.Println("Instrument \"" + name + "\" has already been registered")
		os.Exit(1)
	}

	// register the callback that will return this instrument
	registeredInstruments[name] = callback
}

func CreateInstrument(settings map[string]interface{}, tomeSpec *compendium.TomeSpec) *Instrument {
	var instrument Instrument

	// if this instrument type was registered, create a new instance of it
	if createInstrument, ok := registeredInstruments[settings["type"].(string)]; ok {
		instrument = createInstrument(settings, tomeSpec)
	} else {
		log.Println("Instrument \"" + settings["type"].(string) + "\" was not previously registered")
		os.Exit(1)
	}

	// return a pointer to the newly created instrument
	return &instrument
}
