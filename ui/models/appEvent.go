// Package models holds the ui model
package models

import (
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	tea "github.com/charmbracelet/bubbletea"
)

// AppEvent corresponds to an argocd application event
type AppEvent struct {
	Event v1alpha1.ApplicationWatchEvent
}

// ReceiveAppEvent wraps app events received from argocd
func ReceiveAppEvent(event v1alpha1.ApplicationWatchEvent) tea.Cmd {
	return func() tea.Msg {
		return AppEvent{Event: event}
	}
}
