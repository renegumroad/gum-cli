package filesystem

import (
	"io"
	"os"

	"github.com/renehernandez/gum-cli/internal/utils/systeminfo"
)

func CurrentDir() (string, error) {
	return os.Getwd()
}

func Exists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func CopyFile(source, destination string) error {
	if IsSymlink(source) {
		return copySymlink(source, destination)
	}

	return copyFile(source, destination)
}

func IsSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeSymlink != 0
}

func copyFile(source, destination string) error {
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

func copySymlink(source, destination string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}

	if Exists(destination) {
		if err := os.Remove(destination); err != nil {
			return err
		}
	}

	return os.Symlink(link, destination)
}

func RootDir() string {
	if systeminfo.IsWindows() {
		return os.Getenv("SystemDrive") + "\\"
	}

	return "/"
}
