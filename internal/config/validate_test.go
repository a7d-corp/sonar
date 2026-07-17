package config

import (
	"testing"
)

func TestValidateGlobalConfig(t *testing.T) {
	testCases := []struct {
		input   *Globals
		name    string
		wantErr bool
	}{
		{
			name: "test valid name",
			input: &Globals{
				Name:      "valid-name",
				Labels:    make(map[string]string),
				Namespace: "default",
			},
			wantErr: false,
		},
		{
			name: "test over-length name",
			input: &Globals{
				Name:      "iexoochohNgoughie1wu8Ahyoof7de0iequiuzoo0mic8siasdc",
				Labels:    make(map[string]string),
				Namespace: "default",
			},
			wantErr: true,
		},
		{
			name: "test name with invalid characters",
			input: &Globals{
				Name:      "invalid_name",
				Labels:    make(map[string]string),
				Namespace: "default",
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			errs := ValidateGlobalConfig(testCase.input)
			gotErr := false
			if errs != nil {
				gotErr = true
			}
			if gotErr != testCase.wantErr {
				t.Errorf("error expected: %t, got %t, %v", testCase.wantErr, gotErr, errs)
			}
		})
	}
}
