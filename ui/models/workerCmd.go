// Package models holds the ui model
package models

import (
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

// WorkerCommand differenciates possible commands
type WorkerCommand int

// Possible Commands
const (
	Refresh WorkerCommand = iota
	HardRefresh
	Sync
)

// WorkerCmd corresponds to an argocd application event
type WorkerCmd struct {
	Cmd WorkerCommand
	App v1alpha1.Application
}
