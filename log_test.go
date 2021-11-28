package syncwd

import (
	"github.com/fwhezfwhez/errorx"
	"testing"
)

func TestErrorf(t *testing.T) {
	Errorf("%s \n", errorx.NewFromStringf("%s", "hehe"))
	Printf("%s \n", errorx.NewFromStringf("%s", "hehe"))
}
