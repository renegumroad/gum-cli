package yaml

import (
	"os"
	"reflect"

	"github.com/pkg/errors"
	"github.com/renegumroad/gum-cli/internal/filesystem"
	lib "gopkg.in/yaml.v3"
)

type Client interface {
	Read(path string, out interface{}) error
	Load(data []byte, out interface{}) error
}

type client struct {
	fs filesystem.Client
}

func New() Client {
	return &client{
		fs: filesystem.New(),
	}
}

func (c *client) Read(path string, out interface{}) error {
	if !c.fs.Exists(path) {
		return errors.Errorf("File does not exist: %s", path)
	}

	if !c.fs.IsFile(path) {
		return errors.Errorf("Path is not a file: %s", path)
	}

	bytes, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	return c.Load(bytes, out)
}

func (c *client) Load(data []byte, out interface{}) error {
	if data == nil {
		return errors.Errorf("data argument is nil")
	}

	if out == nil {
		return errors.Errorf("out argument is nil")
	}

	err := lib.Unmarshal(data, out)
	if err != nil {
		return errors.Errorf("Failed to unmarshal yaml: %s", err)
	}

	if out == nil || reflect.ValueOf(out).Elem().IsZero() {
		return errors.Errorf("Unmarshalled yaml does not contain expected data: %v", out)
	}

	return nil
}
