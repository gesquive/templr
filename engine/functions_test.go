package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookupIPv4Host(t *testing.T) {
	hosts, err := lookupIPv4Host("172.217.0.14")
	assert.NoError(t, err, "unexpected error")
	assert.Len(t, hosts, 1, "unexpected results")
}

func TestLookupIPv4Host_WithIPv6(t *testing.T) {
	hosts, err := lookupIPv4Host("2607:f8b0:4009:813::200e")
	assert.NoError(t, err, "unexpected error")
	assert.Len(t, hosts, 0, "unexpected results")
}

func TestLookupIPv6Host(t *testing.T) {
	hosts, err := lookupIPv6Host("2607:f8b0:4009:813::200e")
	assert.NoError(t, err, "unexpected error")
	assert.Len(t, hosts, 1, "unexpected results")
}

func TestLookupIPv6Host_WithIPv4(t *testing.T) {
	hosts, err := lookupIPv6Host("172.217.0.14")
	assert.NoError(t, err, "unexpected error")
	assert.Len(t, hosts, 0, "unexpected results")
}

func TestIsValidIPv4Addr(t *testing.T) {
	result := IsValidIPv4Addr("127.0.0.1")
	assert.True(t, result, "unexpected result")

	result = IsValidIPv4Addr("::1")
	assert.False(t, result, "unexpected result")

	result = IsValidIPv4Addr("rando")
	assert.False(t, result, "unexpected result")

	result = IsValidIPv4Addr("10.0.0.0/8")
	assert.False(t, result, "unexpected result")
}

func TestIsValidIPv6Addr(t *testing.T) {
	result := IsValidIPv6Addr("::1")
	assert.True(t, result, "unexpected result")

	result = IsValidIPv6Addr("127.0.0.1")
	assert.True(t, result, "unexpected result")

	result = IsValidIPv6Addr("rando")
	assert.False(t, result, "unexpected result")

	result = IsValidIPv6Addr("2001:db8::/32")
	assert.False(t, result, "unexpected result")
}

func TestIsValidIPv4CIDR(t *testing.T) {
	result := IsValidIPv4CIDR("10.0.0.0/8")
	assert.True(t, result, "unexpected result")

	result = IsValidIPv4CIDR("2001:db8::/32")
	assert.False(t, result, "unexpected result")

	result = IsValidIPv4CIDR("127.0.0.1")
	assert.False(t, result, "unexpected result")

	result = IsValidIPv4CIDR("::1")
	assert.False(t, result, "unexpected result")
}

func TestIsValidIPv6CIDR(t *testing.T) {
	result := IsValidIPv6CIDR("2001:db8::/32")
	assert.True(t, result, "unexpected result")

	result = IsValidIPv6CIDR("10.0.0.0/8")
	assert.True(t, result, "unexpected result")

	result = IsValidIPv6CIDR("::1")
	assert.False(t, result, "unexpected result")

	result = IsValidIPv6CIDR("127.0.0.1")
	assert.False(t, result, "unexpected result")
}
