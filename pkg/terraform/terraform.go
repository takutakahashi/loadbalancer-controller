package terraform

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"text/template"

	"github.com/google/uuid"
	"github.com/takutakahashi/loadbalancer-controller/api/v1beta1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var BASE_DIR string = os.Getenv("PROJECT_ROOT")

type TerraformClient struct {
	awsBackend     *v1beta1.AWSBackend
	tfvarsPath     string
	tfstateBackend string
}

func (t TerraformClient) init() error {
	prev, err := os.Getwd()
	if err != nil {
		return err
	}
	if err = os.Chdir(t.workDir()); err != nil {
		return err
	}
	defer os.Chdir(prev)
	init := exec.Command("terraform", "init")
	init.Env = os.Environ()
	err = init.Run()
	if err != nil {
		return err
	}
	return nil
}
func (t TerraformClient) createTfVarsSecret() error {
	return nil
}
func (t TerraformClient) createJob(ops string, force bool) error {
	_ = t.buildJob(ops, force)
	return nil
}

func (t TerraformClient) execute(ops string, force bool) error {
	err := t.createTfVarsSecret()
	if err != nil {
		return err
	}
	err = t.createJob(ops, force)
	if err != nil {
		return err
	}
	return nil
}

func (t TerraformClient) Apply() error {
	return t.execute("apply", true)
}

func (t TerraformClient) Destroy() error {
	return t.execute("destroy", true)
}

func genTfVarsPath() string {
	return fmt.Sprintf("/tmp/%s.tfvars", uuid.New().String())
}

func (t TerraformClient) genTfvars(w io.Writer) error {
	awsBackend := t.awsBackend
	if awsBackend != nil {
		tmpl, err := template.ParseFiles(t.workDir() + "/template.tfvars.tpl")
		if err != nil {
			return err
		}
		err = tmpl.Execute(w, struct {
			B         *v1beta1.AWSBackend
			ServiceIn bool
		}{
			B:         awsBackend,
			ServiceIn: false,
		})
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func NewClientForAWSBackend(awsBackend v1beta1.AWSBackend) (TerraformClient, error) {
	tc := TerraformClient{}
	tfvarsPath := genTfVarsPath()
	f, err := os.Create(tfvarsPath)
	if err != nil {
		return tc, err
	}
	defer f.Close()
	tc.awsBackend = &awsBackend
	err = tc.genTfvars(f)
	if err != nil {
		return TerraformClient{}, err
	}
	tc.tfvarsPath = tfvarsPath
	return tc, nil
}
func NewClient(lb v1beta1.Loadbalancer) (TerraformClient, error) {
	tc := TerraformClient{}
	tfvarsPath := genTfVarsPath()
	f, err := os.Create(tfvarsPath)
	if err != nil {
		return tc, err
	}
	defer f.Close()
	tc.awsBackend = &lb.Spec.AWSBackend
	err = tc.genTfvars(f)
	if err != nil {
		return TerraformClient{}, err
	}
	tc.tfvarsPath = tfvarsPath
	return tc, nil
}

func (t TerraformClient) workDir() string {
	if t.awsBackend != nil {
		return BASE_DIR + "/src/terraform/aws_backend"
	}
	return ""
}

func (t TerraformClient) buildJob(ops string, force bool) batchv1.Job {
	// secretName := ""
	om := metav1.ObjectMeta{
		Name:      t.awsBackend.Name,
		Namespace: t.awsBackend.Namespace,
	}
	cmd := []string{"/bin/terraform.sh", ops, t.awsBackend.Kind}
	if force {
		cmd = append(cmd, "force")
	}
	return batchv1.Job{
		ObjectMeta: om,
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Containers: []corev1.Container{
						corev1.Container{
							Name:    "tf",
							Image:   "takutakahashi/loadbalancer-controller",
							Command: cmd,
							//		EnvFrom: []corev1.EnvFromSource{
							//			corev1.EnvFromSource{
							//				SecretRef: &corev1.SecretEnvSource{
							//					LocalObjectReference: corev1.LocalObjectReference{
							//						Name: secretName,
							//					},
							//				},
							//			},
							//		},
						},
					},
				},
			},
		},
	}

}
