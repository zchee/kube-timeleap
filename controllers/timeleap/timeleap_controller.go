// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	timeleapv1alpha1 "github.com/zchee/kube-timeleap/apis/timeleap/v1alpha1"
)

// TimeLeapReconciler reconciles a TimeLeap object
type TimeLeapReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=timeleap.x-k8s.io,resources=timeleaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=timeleap.x-k8s.io,resources=timeleaps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=timeleap.x-k8s.io,resources=timeleaps/finalizers,verbs=update

func (r *TimeLeapReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("timeleap", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *TimeLeapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&timeleapv1alpha1.TimeLeap{}).
		Complete(r)
}
