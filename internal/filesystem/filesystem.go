package filesystem

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/pkg/errors"
	"github.com/renegumroad/gum-cli/internal/log"
	"github.com/renegumroad/gum-cli/internal/systeminfo"
)

type Client interface {
	CurrentDir() (string, error)
	HomeDir() (string, error)
	Exists(path string) bool
	IsDir(path string) bool
	IsFile(path string) bool
	CopyFile(source, destination string) error
	IsSymlink(path string) bool
	RootDir() string
	ChownUser(path string, uid int) error
	Chown(path string, uid, gid int) error
	GetOwner(path string) (*UserInfo, error)
	EnsureNonSudoOwnership(path string) error
	MakeExecutable(path string) error
	IsExecutable(path string) bool
	CreateTempDir() (string, error)
	EqualFiles(source, destination string) (bool, error)
	WriteString(path, content string) error
	AppendString(path, content string) error
	MkdirAll(path string) error
}

type UserInfo struct {
	Id   int
	Name string
}

type client struct {
	sys systeminfo.Client
}

func New() Client {
	return newClientWithComponents(systeminfo.New())
}

func newClientWithComponents(sys systeminfo.Client) *client {
	return &client{
		sys: sys,
	}
}

func (c *client) CurrentDir() (string, error) {
	dir, err := os.Getwd()

	if err != nil {
		return "", errors.Errorf("Unable to get current directory: %s", err)
	}

	return dir, nil
}

func (c *client) HomeDir() (string, error) {
	dir, err := os.UserHomeDir()

	if err != nil {
		return "", errors.Errorf("Unable to get home directory: %s", err)
	}

	return dir, nil
}

func (c *client) Exists(path string) bool {
	_, err := os.Stat(path)

	return !errors.Is(err, os.ErrNotExist)
}

func (c *client) IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func (c *client) IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.Mode().IsRegular()
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
	return "/"
}

func (c *client) ChownUser(path string, uid int) error {
	return c.Chown(path, uid, -1)
}

func (c *client) Chown(path string, uid, gid int) error {
	if c.IsDir(path) {
		return filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			return os.Chown(name, uid, gid)
		})
	}
	return os.Chown(path, uid, gid)
}

func (c *client) GetOwner(path string) (*UserInfo, error) {
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

	userInfo := &UserInfo{}

	userInfo.Name = owner.Username
	if userId, err := strconv.Atoi(owner.Uid); err != nil {
		return nil, errors.Errorf("Unable to convert user id to int: %s", err)
	} else {
		userInfo.Id = userId
	}

	return userInfo, nil
}

func (c *client) EnsureNonSudoOwnership(path string) error {
	info, err := c.GetOwner(path)
	if err != nil {
		return errors.Errorf("Error getting %s owner information: %s", path, err)
	}

	if info.Id != 0 {
		log.Debugf("Won't update ownership of %s path since it is not owned by root. Owner: %s", path, info.Name)
		return nil
	}

	user, err := c.sys.GetSudoOriginalUser()
	if err != nil {
		return err
	}
	log.Debugf("Setting ownership of %s directory to %s", path, user.Name)

	return c.ChownUser(path, user.Id)
}

func (c *client) IsExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.Mode()&0111 != 0
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

func (c *client) CreateTempDir() (string, error) {
	return os.MkdirTemp("", "gum_*")
}

func (c *client) EqualFiles(source, destination string) (bool, error) {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return false, err
	}

	destinationInfo, err := os.Stat(destination)
	if err != nil {
		return false, err
	}

	if sourceInfo.Size() != destinationInfo.Size() {
		return false, nil
	}

	sourceContent, err := os.ReadFile(source)
	if err != nil {
		return false, err
	}

	destinationContent, err := os.ReadFile(destination)
	if err != nil {
		return false, err
	}

	return bytes.Equal(sourceContent, destinationContent), nil
}

func (c *client) AppendString(path, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

// WriteString writes the content to a file at the specified path.
// If the file already exists, it will be overwritten, otherwise it will be created.
func (c *client) WriteString(path, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

func (c *client) MkdirAll(path string) error {
	if c.Exists(path) && !c.IsDir(path) {
		return errors.Errorf("Path %s exists and is not a directory", path)
	}
	return os.MkdirAll(path, os.ModePerm)
}
