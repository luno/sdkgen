package clientgen

import (
	"net/http"
	"testing"
)

func TestExampleRequest(t *testing.T) {
	tests := []struct {
		name     string
		endpoint Endpoint
		want     string
	}{
		{
			name: "GET request with no required fields",
			endpoint: Endpoint{
				Method: http.MethodGet,
				Path:   "/api/test",
				Request: Type{
					Kind:        KindStruct,
					StructProps: &StructProps{},
				},
			},
			want: "curl https://api.luno.com/api/test",
		},
		{
			name: "POST request with no required fields",
			endpoint: Endpoint{
				Method: http.MethodPost,
				Path:   "/api/test",
				Request: Type{
					Kind:        KindStruct,
					StructProps: &StructProps{},
				},
			},
			want: "curl -X POST \\\n     https://api.luno.com/api/test",
		},
		{
			name: "GET request requiring auth",
			endpoint: Endpoint{
				Method:       http.MethodGet,
				Path:         "/api/test",
				RequiresAuth: true,
				Request: Type{
					Kind:        KindStruct,
					StructProps: &StructProps{},
				},
			},
			want: "curl -u api_key_id:api_key_secret \\\n     https://api.luno.com/api/test",
		},
		{
			name: "GET request with required scalar field",
			endpoint: Endpoint{
				Method: http.MethodGet,
				Path:   "/api/test",
				Request: Type{
					Kind: KindStruct,
					StructProps: &StructProps{
						Properties: []Property{
							{
								Name:     "pair",
								Type:     Type{Kind: KindString},
								Required: true,
								Example:  "XBTZAR",
							},
						},
					},
				},
			},
			want: "curl -F 'pair=XBTZAR' \\\n     https://api.luno.com/api/test",
		},
		{
			name: "GET request with path parameter substitution",
			endpoint: Endpoint{
				Method: http.MethodGet,
				Path:   "/api/test/{id}",
				Request: Type{
					Kind: KindStruct,
					StructProps: &StructProps{
						Properties: []Property{
							{
								Name:     "id",
								Type:     Type{Kind: KindString},
								Required: true,
								Example:  "abc123",
							},
						},
					},
				},
			},
			want: "curl https://api.luno.com/api/test/abc123",
		},
		{
			name: "GET request with required array field (uses SplitSeq)",
			endpoint: Endpoint{
				Method: http.MethodGet,
				Path:   "/api/test",
				Request: Type{
					Kind: KindStruct,
					StructProps: &StructProps{
						Properties: []Property{
							{
								Name: "ids",
								Type: Type{
									Kind:       KindArray,
									ArrayProps: &ArrayProps{Type: Type{Kind: KindString}},
								},
								Required: true,
								Example:  "a,b,c",
							},
						},
					},
				},
			},
			want: "curl -F 'ids=a' \\\n     -F 'ids=b' \\\n     -F 'ids=c' \\\n     https://api.luno.com/api/test",
		},
		{
			name: "optional fields are skipped",
			endpoint: Endpoint{
				Method: http.MethodGet,
				Path:   "/api/test",
				Request: Type{
					Kind: KindStruct,
					StructProps: &StructProps{
						Properties: []Property{
							{
								Name:     "optional_field",
								Type:     Type{Kind: KindString},
								Required: false,
								Example:  "some_value",
							},
						},
					},
				},
			},
			want: "curl https://api.luno.com/api/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.endpoint.ExampleRequest()
			if got != tt.want {
				t.Errorf("ExampleRequest() =\n%q\nwant:\n%q", got, tt.want)
			}
		})
	}
}

func TestResponseStructure(t *testing.T) {
	tests := []struct {
		name     string
		endpoint Endpoint
		want     string
	}{
		{
			name: "empty response returns empty string",
			endpoint: Endpoint{
				Response: Type{Kind: KindUnknown},
			},
			want: "",
		},
		{
			name: "struct response with string and integer fields",
			endpoint: Endpoint{
				Response: Type{
					Kind: KindStruct,
					StructProps: &StructProps{
						Properties: []Property{
							{Name: "name", Type: Type{Kind: KindString}},
							{Name: "count", Type: Type{Kind: KindInteger}},
						},
					},
				},
			},
			want: "{\n  \"count\": 0,\n  \"name\": \"\"\n}",
		},
		{
			name: "struct response with decimal and timestamp fields",
			endpoint: Endpoint{
				Response: Type{
					Kind: KindStruct,
					StructProps: &StructProps{
						Properties: []Property{
							{Name: "amount", Type: Type{Kind: KindDecimal}},
							{Name: "timestamp", Type: Type{Kind: KindTimestamp}},
						},
					},
				},
			},
			want: "{\n  \"amount\": \"0.0\",\n  \"timestamp\": 0\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.endpoint.ResponseStructure()
			if got != tt.want {
				t.Errorf("ResponseStructure() =\n%q\nwant:\n%q", got, tt.want)
			}
		})
	}
}
