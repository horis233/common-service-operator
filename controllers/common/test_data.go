//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package common

import (
	olmv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	apiv1alpha1 "github.com/IBM/operand-deployment-lifecycle-manager/api/v1alpha1"
	constant "github.com/IBM/operand-deployment-lifecycle-manager/controllers/constant"
)

// Return OperandRegistry obj
func OperandRegistryObj(name, namespace, subNamespace string) *apiv1alpha1.OperandRegistry {
	return &apiv1alpha1.OperandRegistry{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: apiv1alpha1.OperandRegistrySpec{
			Operators: []apiv1alpha1.Operator{
				{
					Name:            "etcd",
					Namespace:       subNamespace,
					SourceName:      "community-operators",
					SourceNamespace: "openshift-marketplace",
					PackageName:     "etcd",
					Channel:         "singlenamespace-alpha",
					Scope:           "public",
				},
				{
					Name:            "jenkins",
					Namespace:       subNamespace,
					SourceName:      "community-operators",
					SourceNamespace: "openshift-marketplace",
					PackageName:     "jenkins-operator",
					Channel:         "alpha",
					Scope:           "public",
				},
			},
		},
	}
}

// Return OperandConfig obj
func OperandConfigObj(name, namespace string) *apiv1alpha1.OperandConfig {
	return &apiv1alpha1.OperandConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: apiv1alpha1.OperandConfigSpec{
			Services: []apiv1alpha1.ConfigService{
				{
					Name: "etcd",
					Spec: map[string]runtime.RawExtension{
						"etcdCluster": {Raw: []byte(`{"size": 3}`)},
					},
				},
				{
					Name: "jenkins",
					Spec: map[string]runtime.RawExtension{
						"jenkins": {Raw: []byte(`{"service":{"port": 8081}}`)},
					},
				},
			},
		},
	}
}

// Return OperandRequest obj
func OperandRequestObj(registryName, registryNamespace, requestName, requestNamespace string) *apiv1alpha1.OperandRequest {
	return &apiv1alpha1.OperandRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:      requestName,
			Namespace: requestNamespace,
		},
		Spec: apiv1alpha1.OperandRequestSpec{
			Requests: []apiv1alpha1.Request{
				{
					Registry:          registryName,
					RegistryNamespace: registryNamespace,
					Operands: []apiv1alpha1.Operand{
						{
							Name: "etcd",
						},
						{
							Name: "jenkins",
							Bindings: map[string]apiv1alpha1.SecretConfigmap{
								"public": {
									Secret:    "secret3",
									Configmap: "cm3",
								},
							},
						},
					},
				},
			},
		},
	}
}

// Return OperandBindInfo obj
func OperandBindInfoObj(name, namespace, registryName, registryNamespace string) *apiv1alpha1.OperandBindInfo {
	return &apiv1alpha1.OperandBindInfo{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: apiv1alpha1.OperandBindInfoSpec{
			Operand:           "jenkins",
			Registry:          registryName,
			RegistryNamespace: registryNamespace,
			Bindings: map[string]apiv1alpha1.SecretConfigmap{
				"public": {
					Secret:    "secret1",
					Configmap: "cm1",
				},
				"private": {
					Secret:    "secret2",
					Configmap: "cm2",
				},
			},
		},
	}
}

func NamespaceObj(name string) *corev1.Namespace {
	return &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

func ConfigmapObj(name, namespace string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string]string{
			"test": name,
		},
	}
}

func SecretObj(name, namespace string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		StringData: map[string]string{
			"test": name,
		},
	}
}

func SubscriptionObj(name, namespace, csvVersion string) *olmv1alpha1.Subscription {
	labels := map[string]string{
		constant.OpreqLabel: "true",
	}
	return &olmv1alpha1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: &olmv1alpha1.SubscriptionSpec{
			Channel:                "alpha",
			Package:                name,
			CatalogSource:          "community-operators",
			CatalogSourceNamespace: "openshift-marketplace",
		},
		Status: olmv1alpha1.SubscriptionStatus{
			CurrentCSV:   name + "-csv.v" + csvVersion,
			InstalledCSV: name + "-csv.v" + csvVersion,
			Install: &olmv1alpha1.InstallPlanReference{
				APIVersion: "operators.coreos.com/v1alpha1",
				Kind:       "InstallPlan",
				Name:       name + "-install-plan",
				UID:        types.UID("install-plan-uid"),
			},
			InstallPlanRef: &corev1.ObjectReference{
				APIVersion: "operators.coreos.com/v1alpha1",
				Kind:       "InstallPlan",
				Name:       name + "-install-plan",
				Namespace:  namespace,
				UID:        types.UID("install-plan-uid"),
			},
		},
	}
}

// Return CSV obj
func CSVObj(name, namespace, example string) *olmv1alpha1.ClusterServiceVersion {
	return &olmv1alpha1.ClusterServiceVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Annotations: map[string]string{
				"alm-examples": example,
			},
		},
		Spec: olmv1alpha1.ClusterServiceVersionSpec{},
		Status: olmv1alpha1.ClusterServiceVersionStatus{
			Phase: olmv1alpha1.CSVPhaseSucceeded,
		},
	}
}

// Return InstallPlan obj
func InstallPlanObj(name, namespace string) *olmv1alpha1.InstallPlan {
	return &olmv1alpha1.InstallPlan{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: olmv1alpha1.InstallPlanSpec{},
		Status: olmv1alpha1.InstallPlanStatus{
			Phase: olmv1alpha1.InstallPlanPhaseComplete,
		},
	}
}
