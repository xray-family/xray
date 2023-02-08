package helper

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectString(t *testing.T) {
	assert.Equal(t, "1", SelectString(true, "1", "0"))
	assert.Equal(t, "0", SelectString(false, "1", "0"))
}

func TestSplit(t *testing.T) {
	assert.ElementsMatch(t, []string{"api", "v1"}, Split("/api/v1", "/"))
	assert.ElementsMatch(t, []string{"api", "v1"}, Split("/api/v1/", "/"))
	assert.ElementsMatch(t, []string{"ming", "hong", "hu"}, Split("ming, hong, hu", ","))
	assert.ElementsMatch(t, []string{"ming", "hong", "hu"}, Split(",ming, hong, hu, ", ","))
	assert.ElementsMatch(t, []string{"ming", "hong", "hu"}, Split("\nming, hong, hu\n", ","))
}

func TestMustHash(t *testing.T) {
	var data = []byte("123")
	assert.Equal(t, "202cb962ac59075b964b07152d234b70", MustHash(data, md5.New(), hex.EncodeToString))
}
