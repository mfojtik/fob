package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/mfojtik/fob/pkg/plugin"
	control_plane "github.com/mfojtik/fob/pkg/plugin/control-plane"
)

type defaultPluginManager struct {
	plugins []plugin.Plugin
}

func (d *defaultPluginManager) Add(plugin plugin.Plugin) {
	d.plugins = append(d.plugins, plugin)
}

func (d *defaultPluginManager) Run(ctx context.Context, options plugin.PluginOptions) error {
	var errors []error
	for _, p := range d.plugins {
		if err := p.Execute(ctx, options); err != nil {
			errors = append(errors, fmt.Errorf("%T: %v", p, err))
		}
	}
	if len(errors) == 0 {
		return nil
	}
	var errMessage []string
	for _, err := range errors {
		errMessage = append(errMessage, err.Error())
	}
	return fmt.Errorf("ERRORS:\n%s", strings.Join(errMessage, "\n"))
}

var plugins = &defaultPluginManager{}

func main() {
	// plugins.Add(&sample.SamplePlugin{})
	// plugins.Add(&pods.PodStatusPlugin{})

	plugins.Add(&control_plane.ControlPlanePlugin{})

	if len(os.Args) != 2 {
		fmt.Printf("usage: %s JOB_URL\n", path.Base(os.Args[0]))
		os.Exit(255)
	}

	options := plugin.PluginOptions{
		JobUrl: os.Args[1],
		Output: os.Stdout,
	}

	if err := plugins.Run(context.TODO(), options); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
