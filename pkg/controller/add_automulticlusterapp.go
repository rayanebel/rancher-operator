package controller

import (
	"gitlab.thalesdigital.io/core-kube/rancher/rancher-operator/pkg/controller/automulticlusterapp"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, automulticlusterapp.Add)
}
