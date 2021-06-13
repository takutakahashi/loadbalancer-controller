package terraform

import (
	"testing"

	"fmt"
	"github.com/takutakahashi/loadbalancer-controller/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func clientsetMock() (*envtest.Environment, *kubernetes.Clientset, error) {
	testEnv := &envtest.Environment{}
	restConfig, err := testEnv.Start()
	if err != nil {
		return testEnv, nil, err
	}
	c, e := kubernetes.NewForConfig(restConfig)
	return testEnv, c, e
}

func TestBuildJob(t *testing.T) {
	cli, err := NewClient(lbMock())
	if err != nil {
		t.Fatal(err)
	}
	job := cli.buildJob("plan", false)
	expected := []string{"/bin/terraform.sh", "plan", "AWSBackend"}
	if len(job.Spec.Template.Spec.InitContainers[1].Command) != len(expected) {
		t.Fatalf("expected: %v, actual: %v", expected, job.Spec.Template.Spec.Containers[0].Command)
	}
	for i, s := range job.Spec.Template.Spec.InitContainers[1].Command {
		if s != expected[i] {
			t.Fatalf("expected: %v, actual: %v", expected, job.Spec.Template.Spec.Containers[0].Command)
		}
	}
}

func TestExecute(t *testing.T) {
	lb := lbMock()
	cli, err := NewClient(lb)
	if err != nil {
		t.Fatal(err)
	}
	testEnv, clientset, err := clientsetMock()
	defer testEnv.Stop()
	if err != nil {
		t.Fatal(err)
	}

	cli.clientset = clientset
	err = cli.execute("plan", false, false)
	if err != nil {
		t.Fatal(err)
	}

	job, err := clientset.BatchV1().Jobs(lb.Namespace).Get(fmt.Sprintf("%s-%s", lb.Spec.AWSBackend.Name, lb.Spec.AWSBackend.ResourceVersion), metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"/bin/terraform.sh", "plan", "AWSBackend"}
	if len(job.Spec.Template.Spec.InitContainers[1].Command) != len(expected) {
		t.Fatalf("expected: %v, actual: %v", expected, job.Spec.Template.Spec.InitContainers[1].Command)
	}
	for i, s := range job.Spec.Template.Spec.InitContainers[1].Command {
		if s != expected[i] {
			t.Fatalf("expected: %v, actual: %v", expected, job.Spec.Template.Spec.InitContainers[1].Command)
		}
	}
	cm, err := clientset.CoreV1().ConfigMaps(cli.awsBackend.Namespace).Get(cli.awsBackend.Name, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := cm.Data["tfvars"]; !ok {
		t.Fatal(cm)
	}
	if _, ok := cm.Data["backend.tf"]; !ok {
		t.Fatal(cm)
	}
}

func TestApply(t *testing.T) {
	lb := lbMock()
	cli, err := NewClient(lb)
	if err != nil {
		t.Fatal(err)
	}
	testEnv, clientset, err := clientsetMock()
	defer testEnv.Stop()
	if err != nil {
		t.Fatal(err)
	}
	cli.clientset = clientset
	err = cli.Apply()
	job, err := clientset.BatchV1().Jobs(lb.Namespace).Get(fmt.Sprintf("%s-%s", lb.Spec.AWSBackend.Name, lb.Spec.AWSBackend.ResourceVersion), metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{"/bin/terraform.sh", "apply", "AWSBackend", "true"}
	if len(job.Spec.Template.Spec.InitContainers[1].Command) != len(expected) {
		t.Fatalf("expected: %v, actual: %v", expected, job.Spec.Template.Spec.InitContainers[1].Command)
	}
	for i, s := range job.Spec.Template.Spec.InitContainers[1].Command {
		if s != expected[i] {
			t.Fatalf("expected: %v, actual: %v", expected, job.Spec.Template.Spec.InitContainers[1].Command)
		}
	}
	cm, err := clientset.CoreV1().ConfigMaps(cli.awsBackend.Namespace).Get(cli.awsBackend.Name, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := cm.Data["tfvars"]; !ok {
		t.Fatal(cm)
	}
}

func TestGenTemplates(t *testing.T) {
	lb := lbMock()
	cli, err := NewClient(lb)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := cli.genTfvars()
	if err != nil {
		t.Log(expected)
		t.Fatal(err)
	}
	expected, err = cli.genBackendTf()
	if err != nil {
		t.Log(expected)
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
			ResourceVersion: "12345",
		},
		Spec: v1beta1.AWSBackendSpec{
			Internal: false,
			Type:     v1beta1.TypeNetwork,
			VPC:      v1beta1.Identifier{ID: "vpc-082f7dcbe447d7ba7"},
			Subnets: []v1beta1.Identifier{
				v1beta1.Identifier{
					ID: "subnet-04ddc1d62069e344c",
				},
			},
			Region: "ap-northeast-1",
			Listeners: []v1beta1.Listener{
				v1beta1.Listener{
					Port:     8080,
					Protocol: v1beta1.BackendProtocolTCP,
					DefaultAction: v1beta1.AWSBackendAction{
						Type: v1beta1.ActionTypeForward,
						TargetGroup: v1beta1.AWSBackendTargetGroup{
							Port:       8080,
							Protocol:   v1beta1.BackendProtocolTCP,
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
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-lb",
			Namespace: "default",
		},
		Spec: v1beta1.LoadbalancerSpec{
			AWSBackend: awsBackendMock(),
		},
	}

}
