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
			input:     "HelloWorld123",
			expOutput: "HelloWorld123",
		},
		{
			name:      "illegal characters",
			input:     "Hello/Wor)l!d_1&2*3@",
			expOutput: "HelloWorld123",
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
