package test

// import (
// 	"app/models"
// 	"app/services"
// 	"testing"

// 	"github.com/stretchr/testify/mock"
// )

// type MockService struct {
// 	mock.Mock
// }

// func (m *MockService) ValidateUser(user models.User) []string {
// 	args := m.Called(user)
// 	return args.Get(0).([]string)
// }

// func (m *MockService) CheckUserExists(user models.User) (bool, error) {
// 	args := m.Called(user)
// 	return args.Bool(0), args.Error(1)
// }
// func (m *MockService) CreateUser(user models.User) error {
// 	args := m.Called(user)
// 	return args.Error(0)
// }
// func (m *MockService) HashPassword(user models.User) ([]byte, error) {
// 	args := m.Called(user)
// 	return args.Get(0).([]byte), args.Error(1)
// }

// func TestRegister(t *testing.T) {
// 	mockServices := new(MockService)
// 	user := &models.User{
// 		Username: "test",
// 		Email:    "a@b.com",
// 		Password: "password",
// 	}
// 	services.ValidateUser(*user) = mockServices.ValidateUser(*user)
// }
