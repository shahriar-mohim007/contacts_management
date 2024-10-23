package mocks

import (
	"contacts/repository"
	"context"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*repository.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*repository.User), args.Error(1)
}

func (m *MockRepository) CreateUser(ctx context.Context, user *repository.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) ActivateUserByID(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) GetAllContacts(ctx context.Context, userID uuid.UUID) ([]repository.Contact, error) {
	args := m.Called(ctx, userID)

	// Ensure args.Get(0) is not nil and is the correct type before returning
	if contacts, ok := args.Get(0).([]repository.Contact); ok {
		return contacts, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) CreateContact(ctx context.Context, contact *repository.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *MockRepository) GetContactByID(ctx context.Context, contactID uuid.UUID) (*repository.ContactWithUserResponse, error) {
	args := m.Called(ctx, contactID)
	return args.Get(0).(*repository.ContactWithUserResponse), args.Error(1)
}

func (m *MockRepository) PatchContact(ctx context.Context, contactID uuid.UUID, contact *repository.Contact) error {
	args := m.Called(ctx, contactID, contact)
	return args.Error(0)
}

func (m *MockRepository) DeleteContactByID(ctx context.Context, contactID uuid.UUID) error {
	args := m.Called(ctx, contactID)
	return args.Error(0)
}
