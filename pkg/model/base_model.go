package model

// Context key for configuration instance
const CTX_CONFIG_KEY = "ctx_config"

type ErrorImpl interface {
	WithCause(error) ErrorImpl
}
