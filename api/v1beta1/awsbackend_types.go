/*


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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AWSBackendSpec defines the desired state of AWSBackend
type AWSBackendSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of AWSBackend. Edit AWSBackend_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// AWSBackendStatus defines the observed state of AWSBackend
type AWSBackendStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// AWSBackend is the Schema for the awsbackends API
type AWSBackend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSBackendSpec   `json:"spec,omitempty"`
	Status AWSBackendStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AWSBackendList contains a list of AWSBackend
type AWSBackendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSBackend `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSBackend{}, &AWSBackendList{})
}
