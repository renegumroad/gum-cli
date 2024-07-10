package assets

import "embed"

//go:embed *
var assets embed.FS

func GetAsset(name string) ([]byte, error) {
	return assets.ReadFile(name)
}
