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

// LoadbalancerSpec defines the desired state of Loadbalancer
type LoadbalancerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	IP         string     `json:"ip"`
	AWSBackend AWSBackend `json:"aws,omitempty"`
}

// LoadbalancerStatus defines the observed state of Loadbalancer
type LoadbalancerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Backend BackendStatus `json:"backend"`
}

type BackendStatus struct {
	Phase     BackendPhase      `json:"phase"`
	Internal  bool              `json:"internal"`
	Endpoint  BackendEndpoint   `json:"endpoint"`
	Listeners []BackendListener `json:"listeners"`
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

// +kubebuilder:object:root=true

// Loadbalancer is the Schema for the loadbalancers API
type Loadbalancer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LoadbalancerSpec   `json:"spec,omitempty"`
	Status LoadbalancerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LoadbalancerList contains a list of Loadbalancer
type LoadbalancerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Loadbalancer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Loadbalancer{}, &LoadbalancerList{})
}
