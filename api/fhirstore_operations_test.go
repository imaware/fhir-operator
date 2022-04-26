package api

import (
	"reflect"
	"testing"

	"github.com/imaware/fhir-operator/api/v1alpha1"
	"github.com/imaware/fhir-operator/mocks"
	"google.golang.org/api/healthcare/v1"
)

var (
	mockDatasetBadGetCall           = &mocks.MockDatasetGetCallBadRequest{}
	mockDatasetGoodGetCall          = &mocks.MockDatastoreGetCall{}
	mockFhirGoodCreateCall          = &mocks.MockFhirCreateCall{}
	mockFhirBadCreateCall           = &mocks.MockFhirCreateCallBadRequest{}
	mockFhirBadGetCall              = &mocks.MockFhirGetCallBadRequest{}
	mockFhirBadDeleteCall           = &mocks.MockFhirDeleteCallBadRequest{}
	mockFhirGoodDeleteCall          = &mocks.MockFhirDeleteCall{}
	mockFhirIAMPolicyGoodCreateCall = &mocks.MockFhirCreateOrUpdateIAMPolicyCall{}
	mockFhirIAMPolicyBadCreateCall  = &mocks.MockFhirCreateOrUpdateIAMPolicyCallBadRequest{}
	mockFhirIAMPolicyGoodGetCall    = &mocks.MockFhirGetIAMPolicyCall{}
	mockFhirIAMPolicyBadGetCall     = &mocks.MockFhirGetIAMPolicyCallBadRequest{}
)

func Test_readandor_create_failed_get_dataset_request(t *testing.T) {
	var fhirStore = &v1alpha1.FhirStore{}
	err := ReadAndOrCreateFHIRStore(mockDatasetBadGetCall, nil, nil, fhirStore)
	if err == nil {
		t.Error("returned no error, wanted an error")
	}
}

func Test_readandor_create_create_fhirstore_fhir_request(t *testing.T) {
	var fhirStore = &v1alpha1.FhirStore{}
	err := ReadAndOrCreateFHIRStore(mockDatasetGoodGetCall, mockFhirBadGetCall, mockFhirGoodCreateCall, fhirStore)
	if err != nil {
		t.Errorf("returned error %v, but wanted no error", err.Error())
	}
}

func Test_readandor_create_failed_create_fhirstore_fhir_request(t *testing.T) {
	var fhirStore = &v1alpha1.FhirStore{}
	err := ReadAndOrCreateFHIRStore(mockDatasetGoodGetCall, mockFhirBadGetCall, mockFhirBadCreateCall, fhirStore)
	if err == nil {
		t.Error("returned no error, wanted an error")
	}
}

func Test_readandor_delete_failed_delete_fhirstor_fhir_request(t *testing.T) {
	var fhirStore = &v1alpha1.FhirStore{}
	err := ReadAndOrDeleteFHIRStore(mockFhirGoodCreateCall, mockFhirBadDeleteCall, fhirStore)
	if err == nil {
		t.Error("no error returned and wanted an error")
	}
}

func Test_readandor_delete_delete_fhirstor_fhir_request(t *testing.T) {
	var fhirStore = &v1alpha1.FhirStore{}
	err := ReadAndOrDeleteFHIRStore(mockFhirGoodCreateCall, mockFhirGoodDeleteCall, fhirStore)
	if err != nil {
		t.Errorf("error returned %v and wanted no error", err.Error())
	}
}

func Test_read_fhirstore_iam_policy(t *testing.T) {
	var fhirStore = &v1alpha1.FhirStore{}
	policy, err := ReadFHIRStoreIAMPolicy(mockFhirIAMPolicyGoodGetCall, fhirStore)
	if err != nil {
		t.Errorf("expected no error, but got error %v", err.Error())
	}
	if policy == nil {
		t.Error("expected a policy but got none")
	}
}

func Test_read_fhirstore_failed_iam_policy(t *testing.T) {
	var fhirStore = &v1alpha1.FhirStore{}
	_, err := ReadFHIRStoreIAMPolicy(mockFhirIAMPolicyBadGetCall, fhirStore)
	if err == nil {
		t.Error("expected an error, but got no error")
	}
	if fhirStore.Status.Status != FAILED {
		t.Errorf("expected fhirStore.Status to be %v but got %v", FAILED, fhirStore.Status.Status)
	}
}

func Test_creatorupdate_fhirstore_iam_policy(t *testing.T) {
	var fhirStore = &v1alpha1.FhirStore{}
	err := CreateOrUpdateFHIRStoreIAMPolicy(mockFhirIAMPolicyGoodCreateCall, fhirStore)
	if err != nil {
		t.Errorf("expected no error, but got error %v", err.Error())
	}

}

func Test_creatorupdate_fhirstore_failed_iam_policy(t *testing.T) {
	var fhirStore = &v1alpha1.FhirStore{}
	err := CreateOrUpdateFHIRStoreIAMPolicy(mockFhirIAMPolicyBadCreateCall, fhirStore)
	if err == nil {
		t.Errorf("expected an error, but got none")
	}
	if fhirStore.Status.Status != FAILED {
		t.Errorf("expected fhirStore.Status to be %v but got %v", FAILED, fhirStore.Status.Status)
	}
}

func Test_generate_IAM_policy(t *testing.T) {
	roleAdmin := "roles/Administrator"
	user2SA := "serviceAccount:user2@imaware-test.iam.gserviceaccount.com"
	user3SA := "serviceAccount:user3@imaware-test.iam.gserviceaccount.com"
	expectedGeneratedPolicy := []*healthcare.Binding{{
		Role:    roleAdmin,
		Members: []string{user2SA, user3SA},
	}}

	newAuth := make(map[string]v1alpha1.FhirStoreSpecAuth)
	newAuth[roleAdmin] = v1alpha1.FhirStoreSpecAuth{
		Members: []string{user2SA, user3SA},
	}
	generatedPolicyBindings := GenerateIAMPolicyBindings(newAuth)
	if len(generatedPolicyBindings) != len(expectedGeneratedPolicy) {
		t.Errorf("expected lenght of geneted policy to be %d but got length of %d", len(expectedGeneratedPolicy), len(generatedPolicyBindings))
	}
	if generatedPolicyBindings[0].Role != roleAdmin {
		t.Errorf("expected generted policy role to be %v but got %v", roleAdmin, generatedPolicyBindings[0].Role)
	}
	if !reflect.DeepEqual(generatedPolicyBindings[0].Members, newAuth[roleAdmin].Members) {
		t.Errorf("expected generated members to be %v but got %v", generatedPolicyBindings[0].Members, newAuth[roleAdmin].Members)
	}

}

func Test_generate_streaming_configs(t *testing.T) {
	id1 := "id1"
	id2 := "id2"
	expectedLen := 2
	fhirstoreBigquerryConfig1 := v1alpha1.FhirStoreSpecOptionsBigquery{Id: id1}
	fhirstoreBigquerryConfig2 := v1alpha1.FhirStoreSpecOptionsBigquery{Id: id2}
	fhirstoreBigquerryConfigs := []v1alpha1.FhirStoreSpecOptionsBigquery{fhirstoreBigquerryConfig1, fhirstoreBigquerryConfig2}
	streamingConfigs := GenerateFhirStoreBigQueryConfigs(fhirstoreBigquerryConfigs)
	if len(streamingConfigs) != expectedLen {
		t.Errorf("expected length of streaming configs to be %d but got length of %d", expectedLen, len(streamingConfigs))
	}

}
