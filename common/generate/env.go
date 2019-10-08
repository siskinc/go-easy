package generate

import (
	"os"
)

func GetGoFile() string {
	return os.Getenv("GOFILE")
}

func GetGoArch() string {
	return os.Getenv("GOARCH")
}

func GetGoOS() string {
	return os.Getenv("GOOS")
}

func GetGoLine() string {
	return os.Getenv("GOLINE")
}

func GetGoPackage() string {
	return os.Getenv("GOPACKAGE")
}
