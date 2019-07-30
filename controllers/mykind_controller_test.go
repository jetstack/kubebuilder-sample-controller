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
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apps "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mygroupv1beta1 "jetstack.io/example-controller/api/v1beta1"
)

var _ = Context("Inside of a new namespace", func() {
	ctx := context.TODO()
	ns := SetupTest(ctx)

	Describe("when no existing resources exist", func() {

		It("should create a new Deployment resource with the specified name and one replica if none is provided", func() {
			myKind := &mygroupv1beta1.MyKind{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testresource",
					Namespace: ns.Name,
				},
				Spec: mygroupv1beta1.MyKindSpec{
					DeploymentName: "deployment-name",
				},
			}

			err := k8sClient.Create(ctx, myKind)
			Expect(err).NotTo(HaveOccurred(), "failed to create test MyKind resource")

			deployment := &apps.Deployment{}
			Eventually(
				getResourceFunc(ctx, client.ObjectKey{Name: "deployment-name", Namespace: myKind.Namespace}, deployment),
				time.Second*5, time.Millisecond*500).Should(BeNil())

			Expect(*deployment.Spec.Replicas).To(Equal(int32(1)))
		})

		It("should create a new Deployment resource with the specified name and two replicas if two is specified", func() {
			myKind := &mygroupv1beta1.MyKind{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testresource",
					Namespace: ns.Name,
				},
				Spec: mygroupv1beta1.MyKindSpec{
					DeploymentName: "deployment-name",
					Replicas:       pointer.Int32Ptr(2),
				},
			}

			err := k8sClient.Create(ctx, myKind)
			Expect(err).NotTo(HaveOccurred(), "failed to create test MyKind resource")

			deployment := &apps.Deployment{}
			Eventually(
				getResourceFunc(ctx, client.ObjectKey{Name: "deployment-name", Namespace: myKind.Namespace}, deployment),
				time.Second*5, time.Millisecond*500).Should(BeNil())

			Expect(*deployment.Spec.Replicas).To(Equal(int32(2)))
		})

		It("should allow updating the replicas count after creating a MyKind resource", func() {
			deploymentObjectKey := client.ObjectKey{
				Name:      "deployment-name",
				Namespace: ns.Name,
			}
			myKindObjectKey := client.ObjectKey{
				Name:      "testresource",
				Namespace: ns.Name,
			}
			myKind := &mygroupv1beta1.MyKind{
				ObjectMeta: metav1.ObjectMeta{
					Name:      myKindObjectKey.Name,
					Namespace: myKindObjectKey.Namespace,
				},
				Spec: mygroupv1beta1.MyKindSpec{
					DeploymentName: deploymentObjectKey.Name,
				},
			}

			err := k8sClient.Create(ctx, myKind)
			Expect(err).NotTo(HaveOccurred(), "failed to create test MyKind resource")

			deployment := &apps.Deployment{}
			Eventually(
				getResourceFunc(ctx, deploymentObjectKey, deployment),
				time.Second*5, time.Millisecond*500).Should(BeNil(), "deployment resource should exist")

			Expect(*deployment.Spec.Replicas).To(Equal(int32(1)), "replica count should be equal to 1")

			err = k8sClient.Get(ctx, myKindObjectKey, myKind)
			Expect(err).NotTo(HaveOccurred(), "failed to retrieve MyKind resource")

			myKind.Spec.Replicas = pointer.Int32Ptr(2)
			err = k8sClient.Update(ctx, myKind)
			Expect(err).NotTo(HaveOccurred(), "failed to Update MyKind resource")

			Eventually(getDeploymentReplicasFunc(ctx, deploymentObjectKey)).
				Should(Equal(int32(2)), "expected Deployment resource to be scale to 2 replicas")
		})

		It("should clean up an old Deployment resource if the deploymentName is changed", func() {
			deploymentObjectKey := client.ObjectKey{
				Name:      "deployment-name",
				Namespace: ns.Name,
			}
			newDeploymentObjectKey := client.ObjectKey{
				Name:      "new-deployment",
				Namespace: ns.Name,
			}
			myKindObjectKey := client.ObjectKey{
				Name:      "testresource",
				Namespace: ns.Name,
			}
			myKind := &mygroupv1beta1.MyKind{
				ObjectMeta: metav1.ObjectMeta{
					Name:      myKindObjectKey.Name,
					Namespace: myKindObjectKey.Namespace,
				},
				Spec: mygroupv1beta1.MyKindSpec{
					DeploymentName: deploymentObjectKey.Name,
				},
			}

			err := k8sClient.Create(ctx, myKind)
			Expect(err).NotTo(HaveOccurred(), "failed to create test MyKind resource")

			deployment := &apps.Deployment{}
			Eventually(
				getResourceFunc(ctx, deploymentObjectKey, deployment),
				time.Second*5, time.Millisecond*500).Should(BeNil(), "deployment resource should exist")

			err = k8sClient.Get(ctx, myKindObjectKey, myKind)
			Expect(err).NotTo(HaveOccurred(), "failed to retrieve MyKind resource")

			myKind.Spec.DeploymentName = newDeploymentObjectKey.Name
			err = k8sClient.Update(ctx, myKind)
			Expect(err).NotTo(HaveOccurred(), "failed to Update MyKind resource")

			Eventually(
				getResourceFunc(ctx, deploymentObjectKey, deployment),
				time.Second*5, time.Millisecond*500).ShouldNot(BeNil(), "old deployment resource should be deleted")

			Eventually(
				getResourceFunc(ctx, newDeploymentObjectKey, deployment),
				time.Second*5, time.Millisecond*500).Should(BeNil(), "new deployment resource should be created")
		})
	})
})

func getResourceFunc(ctx context.Context, key client.ObjectKey, obj runtime.Object) func() error {
	return func() error {
		return k8sClient.Get(ctx, key, obj)
	}
}

func getDeploymentReplicasFunc(ctx context.Context, key client.ObjectKey) func() int32 {
	return func() int32 {
		depl := &apps.Deployment{}
		err := k8sClient.Get(ctx, key, depl)
		Expect(err).NotTo(HaveOccurred(), "failed to get Deployment resource")

		return *depl.Spec.Replicas
	}
}
