package common

import "sigs.k8s.io/controller-runtime/pkg/client"

type Entity interface {
	GetObject() client.Object
	IsReady() bool
	IsReconciled() bool
	SetError(err error)
	SetReady(ready bool)
	SetReconciled(reconciled bool)
}
