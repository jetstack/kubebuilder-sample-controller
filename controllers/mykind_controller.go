/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mygroupv1beta1 "jetstack.io/example-controller/api/v1beta1"
)

// MyKindReconciler reconciles a MyKind object
type MyKindReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=mygroup.k8s.io,resources=mykinds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mygroup.k8s.io,resources=mykinds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mygroup.k8s.io,resources=mykinds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create

func (r *MyKindReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("mykind", req.NamespacedName)

	// your logic here
	log.Info("fetching MyKind resource")
	myKind := mygroupv1beta1.MyKind{}
	if err := r.Client.Get(ctx, req.NamespacedName, &myKind); err != nil {
		log.Error(err, "failed to get MyKind resource")
		// Ignore NotFound errors as they will be retried automatically if the
		// resource is created in future.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("checking if an existing Deployment exists for this resource")
	deployment := apps.Deployment{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: myKind.Namespace, Name: myKind.Name}, &deployment)
	if apierrors.IsNotFound(err) {
		log.Info("could not find existing Deployment for MyKind, creating one...")

		deployment = *buildDeployment(myKind)
		if err := r.Client.Create(ctx, &deployment); err != nil {
			log.Error(err, "failed to create Deployment resource")
			return ctrl.Result{}, err
		}

		log.Info("created Deployment resource for MyKind")
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, "failed to get Deployment for MyKind resource")
		return ctrl.Result{}, err
	}

	log.Info("existing Deployment resource already exists for MyKind, finished processing")

	return ctrl.Result{}, nil
}

func buildDeployment(myKind mygroupv1beta1.MyKind) *apps.Deployment {
	one := int32(1)
	deployment := apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            myKind.Name,
			Namespace:       myKind.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(&myKind, mygroupv1beta1.GroupVersion.WithKind("MyKind"))},
		},
		Spec: apps.DeploymentSpec{
			Replicas: &one,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"example-controller.jetstack.io/deployment-name": myKind.Name,
				},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"example-controller.jetstack.io/deployment-name": myKind.Name,
					},
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
						},
					},
				},
			},
		},
	}
	return &deployment
}

func (r *MyKindReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mygroupv1beta1.MyKind{}).
		Owns(&apps.Deployment{}).
		Complete(r)
}
