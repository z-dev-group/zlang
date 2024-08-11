package cli

import (
	"strings"
	"testing"
)

func TestConvertZToC(t *testing.T) {
	tests := []struct {
		input   string
		expeced string
	}{
		{
			`let name="sevenpan"`,
			`char *name = "sevenpan";`,
		},
		{
			`let age = 18`,
			`int age = 18;`,
		},
		{
			`let age = 12; if (age > 18) {}`,
			``,
		},
	}
	for _, test := range tests {
		result := strings.Trim(ConvertZToC(test.input, false), "\n")

		if result != test.expeced {
			t.Fatalf("expected %s, got %s ", test.expeced, result)
		}
	}
}
