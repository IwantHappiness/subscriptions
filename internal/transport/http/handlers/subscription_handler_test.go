package handlers

import (
	"net/http/httptest"
	"testing"
)

func TestGetListInputFromRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		target     string
		wantLimit  int
		wantOffset int
		wantErr    bool
	}{
		{
			name:       "defaults",
			target:     "/api/v1/subscriptions",
			wantLimit:  defaultListLimit,
			wantOffset: defaultListOffset,
		},
		{
			name:       "custom pagination",
			target:     "/api/v1/subscriptions?limit=25&offset=50",
			wantLimit:  25,
			wantOffset: 50,
		},
		{
			name:    "invalid limit",
			target:  "/api/v1/subscriptions?limit=bad",
			wantErr: true,
		},
		{
			name:    "invalid offset",
			target:  "/api/v1/subscriptions?offset=-1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("GET", tt.target, nil)
			got, err := getListInputFromRequest(req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("getListInputFromRequest() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if got.Limit != tt.wantLimit || got.Offset != tt.wantOffset {
				t.Fatalf("getListInputFromRequest() = %+v, want limit %d offset %d", got, tt.wantLimit, tt.wantOffset)
			}
		})
	}
}
