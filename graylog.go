package gml

import (
	"crypto/tls"

	graylog "github.com/Devatoria/go-graylog"
)

// Graylog is a unifying interface for ...
type Graylog interface {
	Close() error
	Send(graylog.Message) error
}

// NewGraylog TODO
func NewGraylog(cfg *Config) (Graylog, error) {
	if cfg._isMock {
		return cfg._mock, cfg._mockErr
	}

	if cfg.UseTLS {
		return getGraylogTLS(cfg)
	}

	return getGraylog(cfg)
}

// getGraylogTLS MORE TODO
func getGraylogTLS(cfg *Config) (Graylog, error) {
	g, err := graylog.NewGraylogTLS(
		graylog.Endpoint{
			Transport: graylog.TCP,
			Address:   cfg.GraylogAddress,
			Port:      cfg.GraylogPort,
		},
		cfg.GraylogConnectionTimeout,
		&tls.Config{
			InsecureSkipVerify: cfg.InsecureSkipVerify,
		},
	)

	if err != nil {
		return nil, err
	}

	return g, nil
}

// getGraylog TODO
func getGraylog(cfg *Config) (Graylog, error) {
	g, err := graylog.NewGraylog(
		graylog.Endpoint{
			Transport: graylog.TCP,
			Address:   cfg.GraylogAddress,
			Port:      cfg.GraylogPort,
		},
	)

	if err != nil {
		return nil, err
	}

	return g, nil
}
