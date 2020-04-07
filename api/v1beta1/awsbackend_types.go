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
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AWSBackendSpec defines the desired state of AWSBackend
type AWSBackendSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of AWSBackend. Edit AWSBackend_types.go to remove/update
	Internal    bool                  `json:"internal,omitempty"`
	Credentials AWSBackendCredentials `json:"credentials"`
	Type        AWSBackendType        `json:"type,omitempty"`
	VPC         Identifier            `json:"vpc,omitempty"`
	Region      string                `json:"region,omitempty"`
	Subnets     []Identifier          `json:"subnets,omitempty"`
	Listeners   []Listener            `json:"listeners"`
}

type AWSBackendCredentials struct {
	AccesskeyID     *corev1.EnvVarSource `json:"accessKeyID"`
	SecretAccesskey *corev1.EnvVarSource `json:"secretAccessKey"`
}

type Listener struct {
	Port          int                `json:"port"`
	Protocol      AWSBackendProtocol `json:"protocol"`
	DefaultAction AWSBackendAction   `json:"defaultAction"`
}

type AWSBackendAction struct {
	Type        AWSBackendActionType  `json:"type"`
	TargetGroup AWSBackendTargetGroup `json:"targetGroup"`
}

type AWSBackendTargetGroup struct {
	Port       int                  `json:"port"`
	Protocol   AWSBackendProtocol   `json:"protocol"`
	TargetType AWSBackendTargetType `json:"targetType"`
	Targets    []AWSBackendTarget   `json:"targets"`
}

type AWSBackendTarget struct {
	Destination AWSBackendDestination `json:"destination"`
	Port        int                   `json:"port"`
}

type AWSBackendDestination struct {
	InstanceID string `json:"instanceID,omitempty"`
	IP         string `json:"IP,omitempty"`
}

type AWSBackendActionType string

var (
	ActionTypeForward AWSBackendActionType = "forward"
)

type AWSBackendTargetType string

var (
	TargetTypeIP       AWSBackendTargetType = "ip"
	TargetTypeInstance AWSBackendTargetType = "instance"
)

type AWSBackendProtocol string

var (
	AWSBackendProtocolTCP AWSBackendProtocol = "TCP"
	AWSBackendProtocolUDP AWSBackendProtocol = "UDP"
)

type AWSBackendType string

var (
	TypeApplication AWSBackendType = "application"
	TypeNetwork     AWSBackendType = "network"
)

type Identifier struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

// AWSBackendStatus defines the observed state of AWSBackend
type AWSBackendStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Phase AWSBackendPhase `json:"phase"`
}

type AWSBackendPhase string

var (
	AWSBackendPhaseProvisioning AWSBackendPhase = "Provisioning"
	AWSBackendPhaseProvisioned  AWSBackendPhase = "Provisioned"
	AWSBackendPhaseReady        AWSBackendPhase = "Ready"
	AWSBackendPhaseDeleting     AWSBackendPhase = "Deleting"
	AWSBackendPhaseDeleted      AWSBackendPhase = "Deleted"
)

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

func (a AWSBackend) Yaml() string {
	d, err := yaml.Marshal(&a)
	if err != nil {
		return ""
	}
	return string(d)
}

func init() {
	SchemeBuilder.Register(&AWSBackend{}, &AWSBackendList{})
}
