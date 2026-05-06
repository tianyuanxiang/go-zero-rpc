package config

import (
	"path/filepath"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestLoadSysConfig(t *testing.T) {
	var c Config
	path := filepath.Join("..", "..", "etc", "sys.yaml")

	if err := conf.Load(path, &c); err != nil {
		t.Fatalf("load sys config: %v", err)
	}
}
