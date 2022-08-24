package models

import (
	"errors"
	"log"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/charmbracelet/lipgloss"
)

// Find correct application
func getApplication(m ArTUIModel) (v1alpha1.Application, error) {
	log.Printf("getApplications: len(m.Applications.Items) = %d", len(m.Applications.Items))

	if m.List.SelectedItem() != nil {
		log.Printf("getApplication: SelectedItem = %v", m.List.SelectedItem().FilterValue())
		for _, v := range m.Applications.Items {
			if v.Name == m.List.SelectedItem().FilterValue() {
				log.Printf("getApplication: want=%s, got=%s", m.List.SelectedItem().FilterValue(), v.Name)
				return v, nil
			}
		}
	} else {
		log.Printf("getApplciation: m.List.SelectedItem() == nil")
	}
	log.Printf("getApplication: failed to find application")
	return v1alpha1.Application{}, errors.New("failed-to-find-app")
}

// Return windowHeight - (header + footer)
func getContentHeight(m ArTUIModel) int {
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())

	return m.WindowHeight - (headerHeight + footerHeight + 1)
}
