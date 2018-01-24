package types

import (
	"strings"
	"fmt"
	"github.com/spf13/pflag"
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
	var entries []string
	for key, value := range *v {
		entries = append(entries, key+":"+value)
	}

	return strings.Join(entries, ",")
}

func (v *MapFlagValue) Set(value string) error {
	entries := strings.Split(value, ",")

	for _, entry := range entries {
		if entry == "" {
			continue
		}

		split := strings.Split(entry, ":")

		if len(split) != 2 {
			return fmt.Errorf("not a key value pair: %s", entry)
		}

		(*v)[split[0]] = split[1]
	}

	return nil
}

func (v MapFlagValue) Type() string {
	return "map[string]string"
}

func (v MapFlagValue) Get() interface{} {
	return (map[string]string)(v)
}
