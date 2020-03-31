package client

import (
	"os"
	"testing"

	"github.com/takutakahashi/loadbalancer-controller/api/v1beta1"
)

func TestGenTfvars(t *testing.T) {
	lb := v1beta1.Loadbalancer{
		Spec: v1beta1.LoadbalancerSpec{
			AWSBackend: v1beta1.AWSBackend{
				Spec: v1beta1.AWSBackendSpec{
					Internal: false,
					Type:     v1beta1.TypeNetwork,
					VPC:      v1beta1.Identifier{ID: "aaa"},
					Subnets: []v1beta1.Identifier{
						v1beta1.Identifier{
							ID: "iii",
						},
					},
					Region: "ap-northeast1",
					Listeners: []v1beta1.Listener{
						v1beta1.Listener{
							Port:     8080,
							Protocol: v1beta1.AWSBackendProtocolTCP,
							DefaultAction: v1beta1.AWSBackendAction{
								Type: v1beta1.ActionTypeForward,
								TargetGroup: v1beta1.AWSBackendTargetGroup{
									Port:       8080,
									Protocol:   v1beta1.AWSBackendProtocolTCP,
									TargetType: v1beta1.TargetTypeIP,
									Targets: []v1beta1.AWSBackendTarget{
										v1beta1.AWSBackendTarget{
											Destination: v1beta1.AWSBackendDestination{
												IP: "10.0.0.3",
											},
											Port: 8080,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	cli, err := NewClient(lb)
	if err != nil {
		t.Fatal(err)
	}
	err = cli.genTfvars(os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
}
