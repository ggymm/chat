package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

var G *Config

type Config struct {
	root string

	App struct {
		Id   int64  `toml:"id"`
		Mode string `toml:"mode"`
		Name string `toml:"name"`

		Http   string `toml:"http"`
		Socket string `toml:"socket"`
	}

	Cluster struct {
		Nodes []string `toml:"nodes"`
	}
}

func root() string {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	base := filepath.Base(exe)
	if !strings.HasPrefix(exe, os.TempDir()) &&
		!strings.HasPrefix(base, "___") {
		return filepath.Dir(exe)
	} else {
		_, filename, _, ok := runtime.Caller(0)
		if ok {
			return filepath.Join(filepath.Dir(filename), "../")
		}
	}
	return ""
}

func Init() {
	r := root()
	p := filepath.Join(r, "config.toml")

	// 读取日志
	buf, err := os.ReadFile(p)
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal(buf, &G)
	if err != nil {
		panic(err)
	}

	G.root = r
}
