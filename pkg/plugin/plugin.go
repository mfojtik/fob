package plugin

import (
	"context"
	"fmt"
	"io"
)

type PluginOptions struct {
	JobUrl string
	Output io.Writer
}

func (o PluginOptions) Printf(message string, v ...interface{}) {
	if _, err := fmt.Fprintf(o.Output, message+"\n", v...); err != nil {
		panic(err)
	}
}

type PluginManager interface {
	Add(Plugin)
}

type Plugin interface {
	Register(PluginManager) error
	Execute(context.Context, PluginOptions) error
}
