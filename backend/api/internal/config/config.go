package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf

	Mysql struct {
		DataSource string
	}

	Redis struct {
		Host string
		Type string
		Pass string
	}

	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	ExternalSync struct {
		Enabled  bool
		Database struct {
			Type     string
			Host     string
			Port     int
			User     string
			Password string
			Database string
			Schema   string
			Charset  string
		}
	}
}
