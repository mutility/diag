package testdiag_test

import (
	"testing"

	"github.com/mutility/diag"
	"github.com/mutility/diag/testdiag"
)

func TestExample(t *testing.T) {
	td := testdiag.Interface(t)

	diag.MaskValue(td, "haha")

	diag.Print(td, "ha")                               // logs "ha"
	diag.Printf(td, "haha")                            // logs "***"
	diag.Group(td, "hahaha", func(td diag.Interface) { // logs: "***ha:"
		diag.Warning(td, "hahaha") // logs "  ***ha"
	})
	diag.Print(td, "hahahahaha") // logs "******ha"
}
