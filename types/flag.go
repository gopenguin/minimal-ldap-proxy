package types

import (
	"fmt"
	"github.com/spf13/pflag"
	"strings"
	"sort"
)

var _ pflag.Value = &MapFlagValue{}

type MapFlagValue map[string]string

func NewMapFlag(name string, usage string) *pflag.Flag {
	return &pflag.Flag{
		Name:  name,
		Value: &MapFlagValue{},
		Usage: usage,
	}
}

func (v *MapFlagValue) String() string {
	var keys []string
	for key, _ := range *v {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var entries []string
	for _, key := range keys {
		entries = append(entries, key+"="+(*v)[key])
	}

	return strings.Join(entries, ",")
}

func (v *MapFlagValue) Set(value string) error {
	entries := strings.Split(value, ",")

	for _, entry := range entries {
		if entry == "" {
			continue
		}

		split := strings.Split(entry, "=")

		if len(split) != 2 {
			return fmt.Errorf("not a key value pair: %s", entry)
		}

		(*v)[split[0]] = split[1]
	}

	return nil
}

func (v *MapFlagValue) Type() string {
	return "map"
}

func (v *MapFlagValue) Get() interface{} {
	return (map[string]string)(*v)
}
