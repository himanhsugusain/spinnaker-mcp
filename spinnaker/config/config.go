// Package config for configuration google iap auth
package config

import (
	"context"
	"gopkg.in/yaml.v3"
	"os"
)

type AuthConfig struct {
	OAuthClientID         string `json:"oauthClientId" yaml:"oauthClientId"`
	ServiceAccountKeyPath string `json:"serviceAccountKeyPath" yaml:"serviceAccountKeyPath"`
}

type Gate struct {
	Endpoint     string `json:"endpoint" yaml:"endpoint"`
	RetryTimeout int    `json:"retryTimeout,omitempty" yaml:"retryTimeout,omitempty"`
}

type Config struct {
	Gate *Gate       `json:"gate" yaml:"gate"`
	Auth *AuthConfig `json:"auth" yaml:"auth"`
}

func NewConfig(ctx context.Context) (*Config, error) {
	cfg := Config{}
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}
