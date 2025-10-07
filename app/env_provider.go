package app

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

type EnvProvider interface {
	MustGet(key string) string
	Get(key string) (string, bool)

	String(key string, def string) string
	Int(key string, def int) int
}

type OsEnvProvider struct {
}

func (e *OsEnvProvider) Int(key string, def int) int {

	strVal, has := e.Get(key)

	if !has {
		return def
	}

	v, parseErr := strconv.ParseInt(strVal, 10, 64)

	if parseErr != nil {
		return def
	}

	if v > math.MaxInt32 || v < math.MinInt32 {
		return def
	} else {
		return int(v)
	}
}

func (e *OsEnvProvider) String(key string, def string) string {

	strVal, has := e.Get(key)

	if !has {
		return def
	} else {

		if strVal != "" {
			return strVal
		} else {
			return def
		}
	}
}

func (e *OsEnvProvider) Get(key string) (string, bool) {
	return os.LookupEnv(key)
}

func (e *OsEnvProvider) MustGet(key string) string {
	val, ok := e.Get(key)
	if !ok {
		panic(fmt.Sprintf("env var `%s` not found", key))
	} else {
		return val
	}
}

func NewEnvProvider() EnvProvider { return &OsEnvProvider{} }
