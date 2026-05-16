package subscription

import "testing"

func TestIsValidListInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   ListInput
		wantErr bool
	}{
		{
			name:  "valid input",
			input: ListInput{Limit: 50, Offset: 10},
		},
		{
			name:    "zero limit",
			input:   ListInput{Limit: 0, Offset: 0},
			wantErr: true,
		},
		{
			name:    "limit greater than max",
			input:   ListInput{Limit: 101, Offset: 0},
			wantErr: true,
		},
		{
			name:    "negative offset",
			input:   ListInput{Limit: 10, Offset: -1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := isValidListInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("isValidListInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
