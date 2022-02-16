package models

import "strings"

// Application models an argocd application
type Application struct {
	Name       string
	Group      string
	Version    string
	Namespace  string
	Created    string
	LongStatus string
	Status     string
}

// Title returns the application title
func (i Application) Title() string { return i.Name }

// Description returns the application's current description
func (i Application) Description() string { return strings.Trim(i.Status, "\n") } // i.Status }

// FilterValue is what our list is filtered by
func (i Application) FilterValue() string { return i.Name }

// GetLongStatus get's the applications long status information
func (i Application) GetLongStatus() string { return i.LongStatus }

// ApplicationSpec holds the app's spec
type ApplicationSpec struct {
	DestinationNamespace string
	DesitnationServer    string
	Project              string
	SourcePath           string
	SourceRepo           string
	SourceRevision       string
	AutoPrune            string
	SelfHeal             string
}

// ApplicationStatus holds the app's status
type ApplicationStatus struct {
	HealthStatus string
	History      []AppHistory
}

// AppHistory holds the app's history
type AppHistory struct {
	ID            int
	DeployStarted string
	DeployedAt    string
	Revision      string
	Source        AppSource
}

// AppSource holds the app's source
type AppSource struct {
	Path           string
	RepoURL        string
	TargetRevision string
}
