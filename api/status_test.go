package api

import "testing"

func TestFHIRStoreExportFailedStatus(t *testing.T) {
	type args struct {
		errorString string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "FHIRStoreExportFailedStatus",
			args: args{"errorString"},
			want: "Failed to export FHIR store due to [ errorString ]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FHIRStoreExportFailedStatus(tt.args.errorString); got != tt.want {
				t.Errorf("FHIRStoreExportFailedStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
