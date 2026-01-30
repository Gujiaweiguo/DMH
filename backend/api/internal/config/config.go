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

	RateLimit struct {
		PosterGenerate struct {
			MaxRequests    int    `json:",default=5"`
			WindowDuration int    `json:",default=60"`
			Storage        string `json:",default=memory"`
		}
		Default struct {
			MaxRequests    int    `json:",default=100"`
			WindowDuration int    `json:",default=60"`
			Storage        string `json:",default=memory"`
		}
	}

	WeChatPay struct {
		AppID           string `json:",optional"`
		MchID           string `json:",optional"`
		APIKey          string `json:",optional"`
		APIKeyV3        string `json:",optional"`
		APIClientCert   string `json:",optional"`
		APIClientKey    string `json:",optional"`
		NotifyURL       string `json:",optional"`
		RefundNotifyURL string `json:",optional"`
		Sandbox         bool   `json:",default=true"`
		CacheTTL        int    `json:",default=7200"`
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
