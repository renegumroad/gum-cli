package systeminfo

import "runtime"

func IsLinux() bool {
	return runtime.GOOS == "linux"
}

func IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

func IsWindows() bool {
  return runtime.GOOS == "windows"
}
