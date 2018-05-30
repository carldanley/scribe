package instruments

import "github.com/carldanley/scribe/src/compendium"

type Instrument interface {
	Prepare(settings map[string]interface{})
	Write(secrets map[string]string)
}

type GetInstrument func(settings map[string]interface{}, tomeSpec *compendium.TomeSpec) Instrument
