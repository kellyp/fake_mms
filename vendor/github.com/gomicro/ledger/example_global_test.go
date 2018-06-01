package ledger_test

import (
	"os"

	"github.com/gomicro/ledger"
)

var (
	log *ledger.Ledger
)

func init() {
	log = ledger.New(os.Stdout, ledger.DebugLevel)
}

func Example_global() {
	log.Debug("We've got a debug line here")
	log.Info("We've got a info line here")
	log.Warn("We've got a warn line here")
	log.Error("We've got a error line here")
	log.Fatal("We've got a fatal line here")
}
