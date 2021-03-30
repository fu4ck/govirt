package runtime

import (
	"github.com/spf13/viper"
	"strings"
)

type Cfg struct {
	Name    string
	Address string
}

type GOVirt struct {
	Cfg *Cfg
}

func NewGoVirt(cfg *Cfg) *GOVirt {
	return &GOVirt{
		Cfg: cfg,
	}
}

func InitConfig(cfgname string) (*Cfg, error) {
	if cfgname != "" {
		viper.SetConfigFile(cfgname)
	} else {
		viper.AddConfigPath("conf")
		viper.SetConfigName("govirt")
	}

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Cfg{
		Address: viper.GetString("listen_address"),
		Name:    viper.GetString("name"),
	}
	return cfg, nil
}
