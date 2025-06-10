package cfg

import (
	"fmt"
	zfg "github.com/chaindead/zerocfg"
	zfgEnv "github.com/chaindead/zerocfg/env"
	zfgYaml "github.com/chaindead/zerocfg/yaml"
)

var (
	configYamlPath = zfg.Str("cfg_yaml_path", "./config.yaml", "CFGYAMLPATH", zfg.Alias("c"))
)

func Init() error {
	if err := zfg.Parse(zfgEnv.New(), zfgYaml.New(configYamlPath)); err != nil {
		return err
	}
	fmt.Println("starting with config\n", zfg.Show())
	return nil
}
