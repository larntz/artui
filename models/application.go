package models

import "strings"

// Application models an argocd application
type Application struct {
	Name       string
	Status     string
	LongStatus string
}

// Title returns the application title
func (i Application) Title() string { return i.Name }

// Description returns the application's current description
func (i Application) Description() string { return strings.Trim(i.Status, "\n") } // i.Status }

// FilterValue is what our list is filtered by
func (i Application) FilterValue() string { return i.Name }

// GetLongStatus get's the applications long status information
func (i Application) GetLongStatus() string { return i.LongStatus }
