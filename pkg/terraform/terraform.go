package terraform

import (
	"bytes"
	"os"
	"text/template"

	"github.com/takutakahashi/loadbalancer-controller/api/v1beta1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	ctrl "sigs.k8s.io/controller-runtime"
)

var BASE_DIR string = os.Getenv("PROJECT_ROOT")

type TerraformClient struct {
	clientset      *kubernetes.Clientset
	awsBackend     *v1beta1.AWSBackend
	tfvarsPath     string
	tfstateBackend string
}

func (t TerraformClient) createConfig() error {
	sc := t.clientset.CoreV1().ConfigMaps(t.awsBackend.Namespace)
	tfvars, err := t.genTfvars()
	if err != nil {
		return err
	}
	backendTf, err := t.genBackendTf()
	if err != nil {
		return err
	}
	sd := map[string]string{"tfvars": tfvars, "backend.tf": backendTf}
	configmap, err := sc.Get(t.awsBackend.Name, metav1.GetOptions{})
	if err != nil && apierrors.IsNotFound(err) {
		newConfigMap := corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      t.awsBackend.Name,
				Namespace: t.awsBackend.Namespace,
			},
			Data: sd}
		_, err = sc.Create(&newConfigMap)
	} else if err == nil {
		configmap.Data = sd
		_, err = sc.Update(configmap)
	}
	return err
}
func (t TerraformClient) createJob(ops string, force bool) error {
	job := t.buildJob(ops, force)
	_, err := t.clientset.BatchV1().Jobs(t.awsBackend.Namespace).Create(&job)
	return err
}

func (t TerraformClient) execute(ops string, force bool) error {
	err := t.createConfig()
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

func (t TerraformClient) genBackendTf() (string, error) {
	return t.genWithTpl(t.workDir() + "/backend.tf.tpl")
}

func (t TerraformClient) genTfvars() (string, error) {
	return t.genWithTpl(t.workDir() + "/template.tfvars.tpl")
}

func (t TerraformClient) genWithTpl(path string) (string, error) {
	awsBackend := t.awsBackend
	if awsBackend != nil {
		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return "", err
		}
		result := bytes.Buffer{}
		err = tmpl.Execute(&result, struct {
			B         *v1beta1.AWSBackend
			ServiceIn bool
		}{
			B:         awsBackend,
			ServiceIn: false,
		})
		if err != nil {
			return "", err
		}
		return result.String(), nil
	}
	return "", nil
}

func NewClientForAWSBackend(awsBackend v1beta1.AWSBackend) (TerraformClient, error) {
	return NewClient(v1beta1.Loadbalancer{Spec: v1beta1.LoadbalancerSpec{AWSBackend: awsBackend}})
}
func NewClient(lb v1beta1.Loadbalancer) (TerraformClient, error) {
	tc := TerraformClient{}
	config := ctrl.GetConfigOrDie()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return tc, err
	}
	tc.clientset = clientset
	tc.awsBackend = &lb.Spec.AWSBackend
	return tc, nil
}

func (t TerraformClient) workDir() string {
	if t.awsBackend != nil {
		return BASE_DIR + "/src/terraform/AWSBackend"
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
							Image:   "takutakahashi/loadbalancer-controller-toolkit",
							Command: cmd,
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      t.awsBackend.Name,
									MountPath: "/data",
								},
							},
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name:      "AWS_ACCESS_KEY_ID",
									ValueFrom: t.awsBackend.Spec.Credentials.AccesskeyID,
								},
								corev1.EnvVar{
									Name:      "AWS_SECRET_ACCESS_KEY",
									ValueFrom: t.awsBackend.Spec.Credentials.SecretAccesskey,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: t.awsBackend.Name,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: t.awsBackend.Name,
									},
								},
							},
						},
					},
				},
			},
		},
	}

}
