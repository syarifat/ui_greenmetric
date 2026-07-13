package tests

import (
	"github.com/goravel/framework/testing"

	"ui_greenmetric/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
