package common

import (
	"os"
)

const resourceDir = "/workspace/resources/"

func LoadResource(filename string) (*os.File, error) {
	return os.Open(filename)
}
