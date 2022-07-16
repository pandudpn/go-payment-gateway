package pg

import (
	"github.com/pandudpn/go-pg/utils"
	"github.com/sirupsen/logrus"
)

type env string

const (
	Production env = "production"
	Staging    env = "staging"
)

type Config struct {
	// Set the environment gateway
	//
	// Default: Staging
	Environment env `json:"environment"`
	
	// When set to true, this will log your request, response or error to stdout
	// Use logrus as logging
	//
	// Default: true
	Logging bool `json:"logging"`
	
	// Format log when Logging is set to true
	//
	// Default: logrus.TextFormatter | [ts] - [level] - [message]
	LogFunc *logrus.Logger `json:"log_func"`
}

// DefaultConfig define all default value of configuration payment gateway
var DefaultConfig = &Config{
	Environment: Staging,
	Logging:     true,
	LogFunc:     utils.Logger(),
}
