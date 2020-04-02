package terraform

import (
	"os"
	"testing"

	"github.com/takutakahashi/loadbalancer-controller/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGenTfvars(t *testing.T) {
	cli, err := NewClient(lbMock())
	if err != nil {
		t.Fatal(err)
	}
	err = cli.genTfvars(os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
}
func TestInit(t *testing.T) {
	cli, err := NewClient(lbMock())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cli.tfvarsPath)
	err = cli.init()
	if err != nil {
		t.Fatal(err)
	}
}

func TestBuildJob(t *testing.T) {
	cli, err := NewClient(lbMock())
	if err != nil {
		t.Fatal(err)
	}
	job := cli.buildJob("plan", true)
	t.Log(job.Spec.Template.Spec.Containers[0].Command)
	t.Fatal(job.String())
}

func TestExecute(t *testing.T) {
	cli, err := NewClient(lbMock())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cli.tfvarsPath)
	o, err := cli.execute("plan", false)
	t.Log(string(o))
	if err != nil {
		t.Fatal(err)
	}
}

func awsBackendMock() v1beta1.AWSBackend {
	return v1beta1.AWSBackend{
		TypeMeta: metav1.TypeMeta{
			Kind: "AWSBackend",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-lb-test",
			Namespace: "default",
			Annotations: map[string]string{
				"hello": "world",
			},
		},
		Spec: v1beta1.AWSBackendSpec{
			Internal: false,
			Type:     v1beta1.TypeNetwork,
			VPC:      v1beta1.Identifier{ID: "aaa"},
			Subnets: []v1beta1.Identifier{
				v1beta1.Identifier{
					ID: "iii",
				},
			},
			Region: "ap-northeast-1",
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
	}
}
func lbMock() v1beta1.Loadbalancer {
	return v1beta1.Loadbalancer{
		Spec: v1beta1.LoadbalancerSpec{
			AWSBackend: awsBackendMock(),
		},
	}

}
