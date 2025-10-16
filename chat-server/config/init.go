package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

var G *Config

var (
	AppId     int64
	AppMode   string
	AppName   string
	AppHttp   string
	AppSocket string
)

var (
	ClusterNodes []string
)

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

func Init() {
	r := ""
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	base := filepath.Base(exe)
	if !strings.HasPrefix(exe, os.TempDir()) &&
		!strings.HasPrefix(base, "___") {
		r = filepath.Dir(exe)
	} else {
		_, filename, _, ok := runtime.Caller(0)
		if ok {
			r = filepath.Join(filepath.Dir(filename), "../")
		}
	}
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

	// 简化读取
	AppId = G.App.Id
	AppMode = G.App.Mode
	AppName = G.App.Name
	AppHttp = G.App.Http
	AppSocket = G.App.Socket
}

func Slog(t ...string) string {
	name := G.App.Name
	if len(t) > 0 {
		name += "-" + t[0]
	}
	name += ".log"

	// 判断类型
	switch G.App.Mode {
	case "debug":
		return filepath.Join(G.root, "temp", name)
	case "release":
		return filepath.Join(G.root, "logs", name)
	default:
		return filepath.Join(G.root, "logs", name)
	}
}
