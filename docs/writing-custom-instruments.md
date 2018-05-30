# Writing Custom Instruments

## Overview

If `scribe` does not support an instrument that accomplishes something you may be trying to achieve, it is possible to register your own instrumentation with `scribe` to carry out custom logic. Creating custom instrumentation is meant to be easy, painless and straight to the point.

## Interface Design

The [interface for an instrument](../src/instruments/types.go) is very simple:

```go
type Instrument interface {
	Prepare(settings map[string]interface{})
	Write(secrets map[string]string)
}
```

## Writing a Custom Instrument

Here's what a very basic instrument looks like:

```go
package instruments

import (
	"fmt"
	"log"

	"github.com/carldanley/scribe/src/compendium"
)

type MyCustomInstrument struct{}

func (i *MyCustomInstrument) Prepare(settings map[string]interface{}) {
	// do anything you want to prep
}

func (i *MyCustomInstrument) Write(secrets map[string]string) {
	// direct the composited secrets to wherever you want
}

func init() {
	RegisterInstrument("my-custom-instrument", func(settings map[string]interface{}, tomeSpec *compendium.TomeSpec) Instrument {
		// create a new file instrument
		instrument := &MyCustomInstrument{}

		// prepare the instrument with the settings
		instrument.Prepare(settings)

		// return the instrument
		return instrument
	})
}
```

The rest is **very** straight-forward.
