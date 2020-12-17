package game

import "testing"

func TestStatusFromString(t *testing.T) {
	var tests = []struct {
		str            string
		expectedStatus Status
		expectedError  error
	}{
		{"Prototype", Prototype, nil},
		{"Signed", Signed, nil},
		{"Published", Published, nil},
		{"Archived", Archived, nil},
		{"Not a status", "", InvalidStatus{}},
	}

	for _, tt := range tests {
		actual, err := StatusFromString(tt.str)
		if tt.expectedError != nil {
			if _, ok := err.(InvalidStatus); !ok {
				t.Errorf("Expected error on invalid status, got none")
			}
		}

		if actual != tt.expectedStatus {
			t.Errorf("String '%s' did not produce expected status. Got '%s'", tt.str, actual)
		}
	}
}
