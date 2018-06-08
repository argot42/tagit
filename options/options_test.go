package options

import (
	"testing"
)

func TestNumeric (t *testing.T) {
	def := 0
	op := []string{ "hi", "how", "are", "you" }

	o := ChooseNumeric(def, op, "choose")

	t.Logf("option choosen: %d\n", o)
}
