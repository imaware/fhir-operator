package mocks

import (
	"context"
	"errors"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MockK8sClientCreatePassed struct {
}

type MockK8sClientCreateAlreadyExists struct {
}

type MockK8sClientBad struct {
}

func (m *MockK8sClientCreatePassed) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	return nil
}

func (m *MockK8sClientCreatePassed) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	return nil
}

func (m *MockK8sClientCreatePassed) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	return nil
}

func (m *MockK8sClientCreatePassed) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	return nil
}

func (m *MockK8sClientCreatePassed) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	return nil
}

func (m *MockK8sClientCreateAlreadyExists) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	return k8serrors.NewAlreadyExists(schema.GroupResource{}, "")
}

func (m *MockK8sClientCreateAlreadyExists) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	return nil
}

func (m *MockK8sClientCreateAlreadyExists) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	return nil
}

func (m *MockK8sClientCreateAlreadyExists) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	return nil
}

func (m *MockK8sClientCreateAlreadyExists) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	return nil
}

func (m *MockK8sClientBad) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	return errors.New("")
}

func (m *MockK8sClientBad) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOptions) error {
	return errors.New("")
}

func (m *MockK8sClientBad) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	return errors.New("")
}

func (m *MockK8sClientBad) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	return errors.New("")
}

func (m *MockK8sClientBad) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	return errors.New("")
}
