package config

import (
	"testing"

	"github.com/go-test/deep"
)

func TestValidateGlobalConfig(t *testing.T) {
	testCases := []struct {
		input           *Globals
		output          *Globals
		name            string
		validateGlobals bool
		wantErr         bool
	}{
		// Test name validation
		{
			name: "test valid name",
			input: &Globals{
				Name: "valid-name",
				Labels: map[string]string{
					"owner":      "sonar",
					"created-by": "sonar",
				},
				Namespace: "default",
			},
			output: &Globals{
				Name:     "valid-name",
				FullName: "sonar-valid-name",
				Labels: map[string]string{
					"name":       "valid-name",
					"owner":      "sonar",
					"created-by": "sonar",
				},
				Namespace: "default",
			},
			validateGlobals: true,
			wantErr:         false,
		},
		{
			name: "test exact-length name",
			input: &Globals{
				Name: "duhuddmhpmitxtevvixngwmnbwjmqbrdqquuiknfqxfnrwybdh",
				Labels: map[string]string{
					"owner":      "sonar",
					"created-by": "sonar",
				},
				Namespace: "default",
			},
			output: &Globals{
				Name:     "duhuddmhpmitxtevvixngwmnbwjmqbrdqquuiknfqxfnrwybdh",
				FullName: "sonar-duhuddmhpmitxtevvixngwmnbwjmqbrdqquuiknfqxfnrwybdh",
				Labels: map[string]string{
					"name":       "duhuddmhpmitxtevvixngwmnbwjmqbrdqquuiknfqxfnrwybdh",
					"owner":      "sonar",
					"created-by": "sonar",
				},
				Namespace: "default",
			},
			validateGlobals: true,
			wantErr:         false,
		},
		{
			name: "test over-length name",
			input: &Globals{
				Name: "duhuddmhpmitxtevvixngwmnbwjmqbrdqquuiknfqxfnrwybdh1",
				Labels: map[string]string{
					"owner":      "sonar",
					"created-by": "sonar",
				},
				Namespace: "default",
			},
			output: &Globals{
				Name: "duhuddmhpmitxtevvixngwmnbwjmqbrdqquuiknfqxfnrwybdh1",
				Labels: map[string]string{
					"owner":      "sonar",
					"created-by": "sonar",
				},
				Namespace: "default",
			},
			validateGlobals: false,
			wantErr:         true,
		},
		{
			name: "test name with invalid characters",
			input: &Globals{
				Name: "invalid_name",
				Labels: map[string]string{
					"owner":      "sonar",
					"created-by": "sonar",
				},
				Namespace: "default",
			},
			output: &Globals{
				Name: "invalid_name",
				Labels: map[string]string{
					"owner":      "sonar",
					"created-by": "sonar",
				},
				Namespace: "default",
			},
			validateGlobals: false,
			wantErr:         true,
		},
		{
			name: "test no provided name",
			input: &Globals{
				Name: "",
				Labels: map[string]string{
					"owner":      "sonar",
					"created-by": "sonar",
				},
				Namespace: "default",
			},
			output: &Globals{
				Name:     "sonar",
				FullName: "sonar",
				Labels: map[string]string{
					"name":       "sonar",
					"owner":      "sonar",
					"created-by": "sonar",
				},
				Namespace: "default",
			},
			validateGlobals: true,
			wantErr:         false,
		},
		// Test namespace validation
		{
			name: "test valid namespace",
			input: &Globals{
				Name: "valid-name",
				Labels: map[string]string{
					"owner":      "sonar",
					"created-by": "sonar",
				},
				Namespace: "valid-namespace",
			},
			output: &Globals{
				Name: "valid-name",
				Labels: map[string]string{
					"name":       "valid-name",
					"created-by": "sonar",
				},
				Namespace: "valid-namespace",
			},
			validateGlobals: false,
			wantErr:         false,
		},
		{
			name: "test namespace with invalid characters",
			input: &Globals{
				Name: "valid-name",
				Labels: map[string]string{
					"owner":      "sonar",
					"created-by": "sonar",
				},
				Namespace: "invalid_namespace",
			},
			output: &Globals{
				Name: "valid-name",
				Labels: map[string]string{
					"owner":      "sonar",
					"created-by": "sonar",
				},
				Namespace: "invalid_namespace",
			},
			validateGlobals: false,
			wantErr:         true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			errs := ValidateGlobalConfig(testCase.input)
			gotErr := errs != nil

			if gotErr != testCase.wantErr {
				t.Errorf("error expected: %t, got %t, %v", testCase.wantErr, gotErr, errs)
			}

			if testCase.validateGlobals {
				if diff := deep.Equal(*testCase.input, *testCase.output); diff != nil {
					t.Error(diff)
				}
			}
		})
	}
}
