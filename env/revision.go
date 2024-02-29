package env

import (
	_ "embed"
)

//go:embed revision.txt
var revision string

func Revision() string {
	return revision
}
