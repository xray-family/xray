package constant

import "testing"

func TestUError_Error(t *testing.T) {
	println(ErrHeaderTooLarge.Error())
}
