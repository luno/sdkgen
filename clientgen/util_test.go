package clientgen

import "testing"

func TestSanitiseVariableName(t *testing.T) {
	tcs := []struct {
		name      string
		input     string
		expOutput string
	}{
		{
			name:      "no change",
			input:     "Hello_World123",
			expOutput: "Hello_World123",
		},
		{
			name:      "illegal characters",
			input:     "Hello_/Wor)l!d1&2*3@",
			expOutput: "Hello_World123",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			output := SanitiseVariableName(tc.input)
			if tc.expOutput != output {
				t.Fatalf("output %s didn't match expected output %s", output, tc.expOutput)
			}
		})
	}
}
