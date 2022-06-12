package models

import (
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/larntz/artui/argo"
)

// GetApplicationMsg wraps a list of applications
type GetApplicationMsg struct {
	applications v1alpha1.ApplicationList
}

// GetApplications retrieves all applications from the argocd server
func GetApplications(m ArTUIModel) tea.Cmd {
	return func() tea.Msg {
		appList := argo.GetApplications(m.ArgoSessionRequest, m.APIClient)
		return GetApplicationMsg{applications: appList}
	}
}
