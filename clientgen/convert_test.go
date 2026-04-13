package clientgen

import (
	"encoding/json"
	"testing"

	openapi "github.com/go-openapi/spec"
	"github.com/google/go-cmp/cmp"
	"github.com/luno/sdkgen/clientgen/testspecs"
)

func TestConvertSpec(t *testing.T) {
	tcs := []struct {
		name      string
		inputSpec string
		expAPI    API
	}{
		{
			name:      "standard spec",
			inputSpec: testspecs.TestSpec,
			expAPI: API{
				Description: "Spec to test conversion",
				Sections: []Section{
					{
						Name:        "Test1",
						Description: "Test1 Endpoint",
						Endpoints: []Endpoint{
							{
								Method:      "POST",
								Path:        "/test1",
								Name:        "Test1",
								Description: []string{"Test1 endpoint"},
								Request: Type{
									Kind: KindStruct,
									Name: "Test1Request",
									StructProps: &StructProps{
										Properties: []Property{
											{
												Name:        "field1",
												Description: []string{"Field1 of test endpoint"},
												Type: Type{
													Kind: KindString,
												},
												Example:  "Field1 Example",
												Required: true,
											},
											{
												Name:        "field2",
												Description: []string{"Field2 of test endpoint"},
												Type: Type{
													Kind: KindInteger,
												},
												Example: "Field2 Example",
											},
										},
									},
								},
								Response: Type{
									Kind: KindStruct,
									Name: "Test1Response",
									StructProps: &StructProps{
										[]Property{
											{
												Name:        "field1",
												Description: []string{"field1 desc"},
												Type: Type{
													Kind: KindDecimal,
												},
											},
											{
												Name:        "field2",
												Description: []string{"Unix timestamp in milliseconds"},
												Type: Type{
													Kind: KindTimestamp,
												},
											},
											{
												Name:        "field3",
												Description: []string{"Unix timestamp in milliseconds string"},
												Type: Type{
													Kind: KindTimestamp,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			var spec *openapi.Swagger
			if err := json.Unmarshal([]byte(tc.inputSpec), &spec); err != nil {
				t.Errorf("unmarshal failed: %s", err)
			}

			api := ConvertSpec(spec)

			if !cmp.Equal(tc.expAPI, api) {
				t.Errorf("expected api differs:\n%s", cmp.Diff(tc.expAPI, api))
			}
		})
	}
}

func TestGetPackageName(t *testing.T) {
	tests := []struct {
		name       string
		extName    string
		extensions map[string]any
		want       string
	}{
		{
			name:       "nil x-go-package returns empty string",
			extName:    "MyType",
			extensions: map[string]any{},
			want:       "",
		},
		{
			name:       "package name present, no x-go-name override",
			extName:    "MyType",
			extensions: map[string]any{"x-go-package": "mypkg"},
			want:       "mypkg.MyType",
		},
		{
			name:       "package name present with x-go-name override",
			extName:    "MyType",
			extensions: map[string]any{"x-go-package": "mypkg", "x-go-name": "OverriddenType"},
			want:       "mypkg.OverriddenType",
		},
		{
			name:       "different package name",
			extName:    "Order",
			extensions: map[string]any{"x-go-package": "trade"},
			want:       "trade.Order",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getPackageName(tt.extName, tt.extensions)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
