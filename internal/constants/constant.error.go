package constants

import "errors"

var (
	// config
	ErrLoadingConfig = errors.New("failed to load config file")
	ErrParsingConfig = errors.New("failed to parse env to config struct")
	ErrEmptyVar      = errors.New("required variabel environment is empty")
)
