package app

import (
	"os"
)

type AppEnv interface {
	Getenv(string) string
	Setenv(string, string) error
}

type PrefixEnv struct {
	Prefix string
}

func (e *PrefixEnv) Getenv(key string) string {
	return os.Getenv(e.Prefix + key)
}

func (e *PrefixEnv) Setenv(key string, value string) error {
	return os.Setenv(key, value)
}
func NewPrefixEnv(prefix string) *PrefixEnv {
	return &PrefixEnv{
		Prefix: prefix,
	}
}

var OsEnv = NewPrefixEnv("")

type PresetEnv map[string]string

func (e PresetEnv) Getenv(key string) string {
	return e[key]
}

func (e PresetEnv) Setenv(key string, value string) error {
	e[key] = value
	return nil
}
