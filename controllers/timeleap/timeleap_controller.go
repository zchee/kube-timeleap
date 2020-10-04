// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	timeleapv1alpha1 "github.com/zchee/kube-timeleap/apis/timeleap/v1alpha1"
)

// TimeLeapReconciler reconciles a TimeLeap object.
type TimeLeapReconciler struct {
	client.Client
	client.Reader
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// compile time check whether the TimeLeapReconciler implements ctrlreconcile.Reconciler interface.
var _ reconcile.Reconciler = (*TimeLeapReconciler)(nil)

// +kubebuilder:rbac:groups=timeleap.x-k8s.io,resources=timeleaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=timeleap.x-k8s.io,resources=timeleaps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=timeleap.x-k8s.io,resources=timeleaps/finalizers,verbs=update

// Reconcile implements a reconcile.Reconciler.
func (r *TimeLeapReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	_ = r.Log.WithValues("timeleap", req.NamespacedName)

	return reconcile.Result{}, nil
}

// SetupWithManager setups the Controller with manager.Manager.
func (r *TimeLeapReconciler) SetupWithManager(mgr manager.Manager) error {
	return builder.ControllerManagedBy(mgr).
		For(&timeleapv1alpha1.TimeLeap{}).
		Complete(r)
}
