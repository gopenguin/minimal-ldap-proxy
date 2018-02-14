package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMapFlag(t *testing.T) {
	const (
		name  = "testFlag"
		usage = "This is a test flag"
	)

	flag := NewMapFlag(name, usage)

	assert.Equal(t, flag.Name, name)
	assert.Equal(t, flag.Usage, usage)
	assert.NotNil(t, flag.Value)
}

func TestMapFlagValue_Set(t *testing.T) {
	tests := []struct {
		desc        string
		stringValue string
		expected    MapFlagValue
		errorString string
	}{
		{
			desc:        "empty value",
			stringValue: "",
			expected:    MapFlagValue{},
		},
		{
			desc:        "single value",
			stringValue: "key:value",
			expected:    MapFlagValue{"key": "value"},
		},
		{
			desc:        "multiple values",
			stringValue: "k1:v1,k2:v2,k3:v3",
			expected:    MapFlagValue{"k1": "v1", "k2": "v2", "k3": "v3"},
		},
		{
			desc:        "not a proper key value pair",
			stringValue: "asdf",
			errorString: "not a key value pair: asdf",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			value := MapFlagValue{}
			err := value.Set(test.stringValue)

			if test.errorString != "" {
				assert.Error(t, err, test.errorString)
				assert.Equal(t, value, MapFlagValue{})
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.expected, value)
			}
		})
	}
}

func TestMapFlagValue_Get(t *testing.T) {
	tests := []struct {
		desc   string
		string string
		value  MapFlagValue
	}{
		{
			desc:  "empty value",
			value: MapFlagValue{},
		},
		{
			desc:  "single value",
			value: MapFlagValue{"key": "value"},
		},
		{
			desc:  "multiple values",
			value: MapFlagValue{"k1": "v1", "k2": "v2", "k3": "v3"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			ret := test.value.Get()

			assert.Equal(t, (map[string]string)(test.value), ret)
		})
	}
}

func TestMapFlagValue_String(t *testing.T) {
	tests := []struct {
		desc   string
		string string
		value  MapFlagValue
	}{
		{
			desc:   "empty value",
			string: "",
			value:  MapFlagValue{},
		},
		{
			desc:   "single value",
			string: "key:value",
			value:  MapFlagValue{"key": "value"},
		},
		{
			desc:   "multiple values",
			string: "k1:v1,k2:v2,k3:v3",
			value:  MapFlagValue{"k1": "v1", "k2": "v2", "k3": "v3"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			ret := test.value.String()

			assert.EqualValues(t, test.string, ret)
		})
	}
}

func TestMapFlagValue_Type(t *testing.T) {
	assert.Equal(t, "map[string]string", MapFlagValue{}.Type())
}
