package filesystem

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/pkg/errors"
	"github.com/renehernandez/gum-cli/internal/utils/systeminfo"
)

type Client interface {
	CurrentDir() (string, error)
	Exists(path string) bool
	CopyFile(source, destination string) error
	IsSymlink(path string) bool
	RootDir() string
	ChownUserRecursively(path string, uid int) error
	ChownRecursively(path string, uid, gid int) error
	GetOwner(path string) (*userInfo, error)
	MakeExecutable(path string) error
}

type userInfo struct {
	Id   int
	Name string
}

type client struct {
	sys systeminfo.Client
}

func NewClient() Client {
	return &client{
		sys: systeminfo.NewClient(),
	}
}

func (c *client) CurrentDir() (string, error) {
	return os.Getwd()
}

func (c *client) Exists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func (c *client) CopyFile(source, destination string) error {
	if c.IsSymlink(source) {
		return c.copySymlink(source, destination)
	}

	return c.copyFile(source, destination)
}

func (c *client) IsSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeSymlink != 0
}

func (c *client) copyFile(source, destination string) error {
	content, err := os.Open(source)
	if err != nil {
		return err
	}
	defer content.Close()

	dest, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, content)
	return err
}

func (c *client) copySymlink(source, destination string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}

	if c.Exists(destination) {
		if err := os.Remove(destination); err != nil {
			return err
		}
	}

	return os.Symlink(link, destination)
}

func (c *client) RootDir() string {
	if c.sys.IsWindows() {
		return os.Getenv("SystemDrive") + "\\"
	}

	return "/"
}

func (c *client) ChownUserRecursively(path string, uid int) error {
	return c.ChownRecursively(path, uid, -1)
}

func (c *client) ChownRecursively(path string, uid, gid int) error {
	err := filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		return os.Chown(name, uid, gid)
	})

	return err
}

func (c *client) GetOwner(path string) (*userInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// On Unix systems, the Sys method returns a *syscall.Stat_t object
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, errors.Errorf("not a syscall.Stat_t")
	}

	// Lookup the user based on UID
	owner, err := user.LookupId(fmt.Sprint(stat.Uid))
	if err != nil {
		return nil, err
	}

	userInfo := &userInfo{}

	userInfo.Name = owner.Username
	if userId, err := strconv.Atoi(owner.Uid); err != nil {
		return nil, errors.Errorf("Unable to convert user id to int: %s", err)
	} else {
		userInfo.Id = userId
	}

	return userInfo, nil
}

func (c *client) MakeExecutable(path string) error {
	// Get the current permissions of the file
	fileInfo, err := os.Stat(path)
	if err != nil {
		return errors.Errorf("Failed to get %s file info while trying to make it executable : %s", path, err)
	}

	// Calculate the new permissions: add executable bits for user, group, and others
	newPermissions := fileInfo.Mode() | 0111

	// Change the file mode
	if err := os.Chmod(path, newPermissions); err != nil {
		return errors.Errorf("Failed to change file mode of %s to %o : %s", path, newPermissions, err)
	}

	return nil
}
