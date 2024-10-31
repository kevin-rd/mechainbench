package core

import (
	"golang.org/x/net/context"
	"mechainbench/app/config"
)

type Context struct {
	context.Context

	AppConfig *config.Config
}
