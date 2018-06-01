package ledger_test

import (
	"github.com/gomicro/ledger"
)

var lvl = ledger.InfoLevel

// ExampleExported demonstrates the default exported logger
func Example_exported() {
	ledger.Debug("We've got a debug line here")
	ledger.Info("We've got a info line here")
	ledger.Warn("We've got a warn line here")
	ledger.Error("We've got a error line here")
	ledger.Fatal("We've got a fatal line here")
}
