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

// BackendSpec defines the desired state of Backend
type BackendSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Backend. Edit Backend_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// BackendStatus defines the observed state of Backend

// +kubebuilder:object:root=true

// Backend is the Schema for the backends API
type Backend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackendSpec   `json:"spec,omitempty"`
	Status BackendStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BackendList contains a list of Backend
type BackendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Backend `json:"items"`
}
type BackendStatus struct {

	// +optional
	Phase BackendPhase `json:"phase,omitempty"`
	// +optional
	Internal bool `json:"internal,omitempty"`
	// +optional
	Endpoint BackendEndpoint `json:"endpoint,omitempty"`
	// +optional
	Listeners []BackendListener `json:"listeners,omitempty"`
}

type BackendListener struct {
	Protocol BackendProtocol `json:"protocol"`
	Port     int             `json:"port"`
}

type BackendEndpoint struct {
	IP  string `json:"IP"`
	DNS string `json:"DNS"`
}
type BackendPhase string

type BackendProtocol string

func (b BackendProtocol) String() string {
	switch b {
	case BackendProtocolTCP:
		return "TCP"
	case BackendProtocolUDP:
		return "UDP"
	default:
		return ""
	}
}

var (
	BackendProtocolTCP BackendProtocol = "TCP"
	BackendProtocolUDP BackendProtocol = "UDP"
)
var (
	BackendPhaseProvisioning BackendPhase = "Provisioning"
	BackendPhaseProvisioned  BackendPhase = "Provisioned"
	BackendPhaseReady        BackendPhase = "Ready"
	BackendPhaseDeleting     BackendPhase = "Deleting"
	BackendPhaseDeleted      BackendPhase = "Deleted"
)

func init() {
	SchemeBuilder.Register(&Backend{}, &BackendList{})
}
