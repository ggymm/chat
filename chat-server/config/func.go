package config

import (
	"path/filepath"
)

func AppLog(t ...string) string {
	fn := G.App.Name
	if len(t) > 0 {
		fn += "-" + t[0]
	}
	fn += ".log"

	// 判断服务类型
	switch G.App.Mode {
	case "debug":
		return filepath.Join(G.root, "temp", fn)
	case "release":
		return filepath.Join(G.root, "logs", fn)
	default:
		return filepath.Join(G.root, "logs", fn)
	}
}

func NodeId() int64 {
	return G.App.Id
}

func NodeMode() string {
	return G.App.Mode
}

func NodeHttp() string {
	return G.App.Http
}

func NodeSocket() string {
	return G.App.Socket
}
