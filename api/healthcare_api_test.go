package api

import (
	"reflect"
	"testing"

	"github.com/imaware/fhir-operator/mocks"
	"google.golang.org/api/healthcare/v1"
)

const PROJECT_ID = "imaware-test"
const LOCATION = "us-central1"
const FHIR_ID = "test-sdcsdcsdc"
const DATASET_ID = "demo-dataset"
const FHIR_VERSION = "R4"
const FAKE_REPRESENTATION = "{}"
const RESOURCE_TYPE = "Observation"
const RESOURCE_ID = "123"

func Test_build_fhir_store_create_call(t *testing.T) {
	_, err := BuildFHIRStoreCreateCall(PROJECT_ID, LOCATION, DATASET_ID, FHIR_VERSION, FHIR_ID)
	if err != nil {
		t.Errorf("Failed to build FHIR store create call due to %v", err)
	}
}

func Test_build_dataset_create_call(t *testing.T) {
	_, err := BuildDatasetCreateCall(PROJECT_ID, LOCATION, DATASET_ID)
	if err != nil {
		t.Errorf("Failed to build Dataset create call due to %v", err)
	}
}

func Test_fhir_build_get_call(t *testing.T) {
	_, err := BuildFHIRStoreGetCall(PROJECT_ID, LOCATION, DATASET_ID, FHIR_ID)
	if err != nil {
		t.Errorf("Failed to build fhir get call due to %v", err)
	}
}

func Test_dataset_build_get_call(t *testing.T) {
	_, err := BuildDatasetGetCall(PROJECT_ID, LOCATION, DATASET_ID)
	if err != nil {
		t.Errorf("Failed to build dataset get call due to %v", err)
	}
}

func Test_fhir_build_delete_call(t *testing.T) {
	_, err := BuildFHIRStoreDeleteCall(PROJECT_ID, LOCATION, DATASET_ID, FHIR_ID)
	if err != nil {
		t.Errorf("Failed to build fhir delete call due to %v", err)
	}
}

func Test_build_fhir_store_resource_update_call(t *testing.T) {
	_, err := BuildFHIRStoreResourceUpdateCall(PROJECT_ID, LOCATION, DATASET_ID, FHIR_ID, FAKE_REPRESENTATION, RESOURCE_TYPE, RESOURCE_ID)
	if err != nil {
		t.Errorf("Failed to build fhir resource update call %v", err)
	}
}

func Test_build_fhir_store_resource_get_call(t *testing.T) {
	_, err := BuildFHIRStoreResourceGetCall(PROJECT_ID, LOCATION, DATASET_ID, FHIR_ID, RESOURCE_TYPE, RESOURCE_ID)
	if err != nil {
		t.Errorf("Failed to build fhir resource get call %v", err)
	}
}

func Test_build_fhir_store_resource_delete_call(t *testing.T) {
	_, err := BuildFHIRStoreResourceDeleteCall(PROJECT_ID, LOCATION, DATASET_ID, FHIR_ID, RESOURCE_TYPE, RESOURCE_ID)
	if err != nil {
		t.Errorf("Failed to build fhir resource delete call %v", err)
	}
}

func Test_build_fhir_store_get_iam_get_call(t *testing.T) {
	_, err := BuildFhirStoreGetIAMPolicyRequest(PROJECT_ID, LOCATION, DATASET_ID, FHIR_ID)
	if err != nil {
		t.Errorf("Failed to build fhir store iam policy get call %v", err)
	}
}

func Test_build_fhir_store_create_or_update_iam_get_call(t *testing.T) {
	fhirStorePolicy := &healthcare.Policy{}
	_, err := BuildFhirStoreUpdateOrCreateIAMPolicyRequest(PROJECT_ID, LOCATION, DATASET_ID, FHIR_ID, fhirStorePolicy)
	if err != nil {
		t.Errorf("Failed to build fhir store iam policy create or update call %v", err)
	}
}

func Test_fhir_create_call(t *testing.T) {

	mockClientCall := &mocks.MockFhirCreateCall{}
	_, err := CreateFHIRStore(mockClientCall)
	if err != nil {
		t.Errorf("Failed to call mock FHIR create due to %v", err)
	}
}

func Test_dataset_create_call(t *testing.T) {

	mockClientCall := &mocks.MockDatasetCreateCall{}
	_, err := CreateDataset(mockClientCall)
	if err != nil {
		t.Errorf("Failed to call mock dataset create due to %v", err)
	}
}

func Test_fhir_get_call(t *testing.T) {
	mockGetCall := &mocks.MockFhirGetCall{}
	_, err := GetFHIRStore(mockGetCall)
	if err != nil {
		t.Errorf("Failed to call mock FHIR get due to %v", err)
	}
}

func Test_fhir_iam_policy_get_call(t *testing.T) {
	mockFhirIAMPolicyGetCall := &mocks.MockFhirGetIAMPolicyCall{}
	_, err := GetFHIRStoreIAMPolicy(mockFhirIAMPolicyGetCall)
	if err != nil {
		t.Errorf("Failed to call mock FHIR get IAM Policy due to %v", err)
	}
}

func Test_dataset_get_call(t *testing.T) {
	mockGetCall := &mocks.MockDatastoreGetCall{}
	_, err := GetDataset(mockGetCall)
	if err != nil {
		t.Errorf("Failed to call mock Dataset get due to %v", err)
	}
}

func Test_fhir_delete_call(t *testing.T) {
	mockDeleteCall := &mocks.MockFhirDeleteCall{}
	_, err := DeleteFHIRStore(mockDeleteCall)
	if err != nil {
		t.Errorf("Failed to call mock FHIR delete due to %v", err)
	}
}

func Test_fhir_resource_delete_call(t *testing.T) {
	mockDeleteCall := &mocks.MockFhirResourceDeleteCall{}
	_, err := DeleteFHIRResource(mockDeleteCall)
	if err != nil {
		t.Errorf("Failed to call mock FHIR resource delete due to %v", err)
	}
}

func Test_fhir_resource_update_call(t *testing.T) {
	mockUpdateCall := &mocks.MockFhirResourceUpdateCall{}
	_, err := UpdateFHIRResource(mockUpdateCall)
	if err != nil {
		t.Errorf("Failed to call mock FHIR resource update due to %v", err)
	}
}

func Test_fhir_resource_get_call(t *testing.T) {
	mockGetCall := &mocks.MockFhirResourceGetCall{}
	_, err := GetFHIRResource(mockGetCall)
	if err != nil {
		t.Errorf("Failed to call mock FHIR resource get due to %v", err)
	}
}

func TestExportFHIRStore(t *testing.T) {

	type args struct {
		fhirStoreExportCall FHIRStoreClientExportCall
	}
	tests := []struct {
		name    string
		args    args
		want    *healthcare.Operation
		wantErr bool
	}{
		{
			name: "Call ExportFHIRStore()",
			args: args{
				fhirStoreExportCall: &mocks.MockFhirExportCall{},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExportFHIRStore(tt.args.fhirStoreExportCall)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExportFHIRStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExportFHIRStore() = %v, want %v", got, tt.want)
			}
		})
	}
}
