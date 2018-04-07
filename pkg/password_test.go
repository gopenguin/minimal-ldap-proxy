package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerify(t *testing.T) {
	assert.True(t, Verify("test123", "$argon2i$v=19$m=32768,t=4,p=4$x2B68O4dFADMgr35a2JbPg$xab/2tyfCYcs4th0QZkDEJZk3rdZ2BSdOCkzy706ot8"))
	assert.False(t, Verify("test1234", "$argon2i$v=19$m=32768,t=4,p=4$x2B68O4dFADMgr35a2JbPg$xab/2tyfCYcs4th0QZkDEJZk3rdZ2BSdOCkzy706ot8"))

	assert.True(t, Verify("test123", "{SSHA}RrAeHR4zMHdNUfvtEibV9yTbtmMY7nF/"))
	assert.False(t, Verify("test124", "{SSHA}RrAeHR4zMHdNUfvtEibV9yTbtmMY7nF/"))

	assert.True(t, Verify("test123", "{SSHA256}ZJM6/dlkBV3P/LYLsNRmwgvF4iXwfndFkcpyYzi0u4pUwluKcBWpjF5SiS9fhc2Y"))
	assert.False(t, Verify("test124", "{SSHA256}ZJM6/dlkBV3P/LYLsNRmwgvF4iXwfndFkcpyYzi0u4pUwluKcBWpjF5SiS9fhc2Y"))

	assert.True(t, Verify("test123", "{SSHA512}QmNKY25YWxQ0V8mn3xtN5cV+cvcsNii2pfuUg34SgNYBR9Hl3bswKV6tffmeqTHjdXV26yS2Ogxe75lz32ZvPIFdD7H4P2N4NRnto3ek1bJSZGRNCdCJ5fXSu8Uomgoc"))
	assert.False(t, Verify("test124", "{SSHA512}QmNKY25YWxQ0V8mn3xtN5cV+cvcsNii2pfuUg34SgNYBR9Hl3bswKV6tffmeqTHjdXV26yS2Ogxe75lz32ZvPIFdD7H4P2N4NRnto3ek1bJSZGRNCdCJ5fXSu8Uomgoc"))
}
