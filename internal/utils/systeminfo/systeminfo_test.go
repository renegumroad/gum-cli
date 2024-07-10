package systeminfo

import (
	"os"
	"os/user"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type systeminfoSuite struct {
	suite.Suite
}

// Mock for the os/user.Lookup function
type UserLookupMock struct {
	mock.Mock
}

func (m *UserLookupMock) Lookup(username string) (*user.User, error) {
	args := m.Called(username)
	return args.Get(0).(*user.User), args.Error(1)
}

type SystemInfoSuite struct {
	suite.Suite
}

func (s *systeminfoSuite) TestGetSudoOriginalUserWithSudo() {
	// Set SUDO_USER environment variable
	os.Setenv("SUDO_USER", "testuser")
	defer os.Unsetenv("SUDO_USER")

	// Mock user.Lookup to return a specific user
	expectedUser := &user.User{Username: "testuser", Uid: "1001"}
	userMock := &UserLookupMock{}
	userMock.On("Lookup", "testuser").Return(expectedUser, nil)

	client := newClientWithComponents(userMock)
	userInfo, err := client.GetSudoOriginalUser()

	s.Require().NoError(err)
	s.Require().NotNil(userInfo)
	s.Require().Equal("testuser", userInfo.Name)
	s.Require().Equal(1001, userInfo.Id)

	userMock.AssertExpectations(s.T())
}

func (s *systeminfoSuite) TestGetSudoOriginalUserWithoutSudo() {
	// Ensure SUDO_USER environment variable is not set
	os.Unsetenv("SUDO_USER")

	client := New()
	userInfo, err := client.GetSudoOriginalUser()

	s.Require().Error(err)
	s.Require().Nil(userInfo)
	s.Require().EqualError(err, "Not running with sudo or SUDO_USER is not set")
}

func TestSystemInfoSuite(t *testing.T) {
	suite.Run(t, new(systeminfoSuite))
}
