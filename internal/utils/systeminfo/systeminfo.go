package systeminfo

import (
	"os"
	"os/user"
	"runtime"
	"strconv"

	"github.com/pkg/errors"
)

type Client interface {
	IsLinux() bool
	IsMacOS() bool
	GetSudoOriginalUser() (*UserInfo, error)
	IsSudo() bool
	GetSudoUsername() string
}

type client struct {
	user userHandler
}

type userHandler interface {
	Lookup(username string) (*user.User, error)
}

type userImpl struct {
}

type UserInfo struct {
	Id   int
	Name string
}

func New() Client {
	return newClientWithComponents(newUserHandler())
}

func newClientWithComponents(user userHandler) Client {
	return &client{
		user: user,
	}
}

func (c *client) IsLinux() bool {
	return runtime.GOOS == "linux"
}

func (c *client) IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

func (c *client) IsSudo() bool {
	return os.Getenv("SUDO_USER") != ""
}

func (c *client) GetSudoUsername() string {
	return os.Getenv("SUDO_USER")
}

func (c *client) GetSudoOriginalUser() (*UserInfo, error) {
	if !c.IsSudo() {
		return nil, errors.Errorf("Not running with sudo or SUDO_USER is not set")
	}

	sudoUsername := c.GetSudoUsername()
	info := &UserInfo{}
	originalUser, err := c.user.Lookup(sudoUsername)
	if err != nil {
		return nil, errors.Errorf("Unable to get information about user %s: %s", sudoUsername, err)
	}
	info.Name = originalUser.Username
	if userId, err := strconv.Atoi(originalUser.Uid); err != nil {
		return nil, errors.Errorf("Unable to convert user id to int: %s", err)
	} else {
		info.Id = userId
	}

	return info, nil
}

func newUserHandler() userHandler {
	return &userImpl{}
}

func (u *userImpl) Lookup(username string) (*user.User, error) {
	return user.Lookup(username)
}
