package conf

import (
	"v2ray.com/core/app/admin"
)

type AdminConfig struct {
	Addr      string   `json:"addr"`
	ContextPath      string   `json:"contextPath"`
	PublicPath      string   `json:"publicPath"`
}

func (c *AdminConfig) Build() (*admin.Config, error) {
	if c.Addr == "" {
		return nil, newError("admin addr can't be empty.")
	}


	return &admin.Config{
		Addr:     c.Addr,
		ContextPath: c.ContextPath,
		PublicPath: c.PublicPath,
	}, nil
}
