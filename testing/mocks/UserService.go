// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	database "go-rest-api-auth/internal/database"

	mock "github.com/stretchr/testify/mock"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: user
func (_m *UserService) CreateUser(user database.UserDTO) (database.UserDTO, error) {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 database.UserDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(database.UserDTO) (database.UserDTO, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(database.UserDTO) database.UserDTO); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(database.UserDTO)
	}

	if rf, ok := ret.Get(1).(func(database.UserDTO) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUser provides a mock function with given fields: userID
func (_m *UserService) DeleteUser(userID int) error {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetALlUsers provides a mock function with given fields:
func (_m *UserService) GetALlUsers() ([]database.UserDTO, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetALlUsers")
	}

	var r0 []database.UserDTO
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]database.UserDTO, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []database.UserDTO); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]database.UserDTO)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserById provides a mock function with given fields: userID
func (_m *UserService) GetUserById(userID int) (database.UserDTO, error) {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for GetUserById")
	}

	var r0 database.UserDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (database.UserDTO, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(int) database.UserDTO); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Get(0).(database.UserDTO)
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByName provides a mock function with given fields: username
func (_m *UserService) GetUserByName(username string) (database.UserDTO, error) {
	ret := _m.Called(username)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByName")
	}

	var r0 database.UserDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (database.UserDTO, error)); ok {
		return rf(username)
	}
	if rf, ok := ret.Get(0).(func(string) database.UserDTO); ok {
		r0 = rf(username)
	} else {
		r0 = ret.Get(0).(database.UserDTO)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: user
func (_m *UserService) UpdateUser(user database.UserDTO) error {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(database.UserDTO) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUserService creates a new instance of UserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserService(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserService {
	mock := &UserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
