package ari

import (
	"github.com/inconshreveable/log15"
	"os"
)

// Logger defaults to a discard handler (null output).
// If you wish to enable logging, you can set your own
// handler like so:
//
//	ari.Logger.SetHandler(log15.StderrHandler)
var Logger = log15.New()

func isNotErr(r *log15.Record) bool {
	return r.Lvl > log15.LvlError
}

func init() {
	// Null logger, by default
	Logger.SetHandler(log15.MultiHandler(
		log15.LvlFilterHandler(log15.LvlError, log15.StreamHandler(os.Stderr, &JsonFormatter{})),
		log15.FilterHandler(isNotErr, log15.StreamHandler(os.Stdout, &JsonFormatter{})),
	))
}
