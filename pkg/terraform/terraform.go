package terraform

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/takutakahashi/loadbalancer-controller/api/v1beta1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
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

func (t TerraformClient) genTerraformFiles() (map[string]string, error) {

	tfvars, err := t.genTfvars()
	if err != nil {
		return nil, err
	}
	backendTf, err := t.genBackendTf()
	if err != nil {
		return nil, err
	}
	lbtf, err := t.genWithTpl(t.workDir() + "/lb.tf.tpl")
	if err != nil {
		return nil, err
	}
	vartf, err := t.genWithTpl(t.workDir() + "/var.tf.tpl")
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"tfvars":     tfvars,
		"backend.tf": backendTf,
		"lb.tf":      lbtf,
		"var.tf":     vartf,
	}, nil
}

func (t TerraformClient) createConfig() error {
	sc := t.clientset.CoreV1().ConfigMaps(t.awsBackend.Namespace)
	sd, err := t.genTerraformFiles()
	if err != nil {
		return err
	}
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
func (t TerraformClient) ensureJob(ops string, force bool) error {
	job := t.buildJob(ops, force)
	c := t.clientset.BatchV1().Jobs(t.awsBackend.Namespace)
	if _, err := c.Get(job.Name, metav1.GetOptions{}); err == nil {
		return nil
	} else if !apierrors.IsNotFound(err) {
		return err
	}
	_, err := c.Create(&job)
	return err
}

func (t TerraformClient) execute(ops string, force bool, watch bool) error {
	err := t.createConfig()
	if err != nil {
		return err
	}
	processing, err := t.isProcessing()
	if err != nil {
		return err
	}
	if processing {
		return errors.New("before task is processing")
	}
	err = t.ensureJob(ops, force)
	if err != nil {
		return err
	}
	if !watch {
		return nil
	} else {
		err = t.watchCompleteOrError()
		if err != nil {
			return err
		}

		return t.cleanup()
	}
}

func (t TerraformClient) cleanup() error {
	name := fmt.Sprintf("%s-%s", t.awsBackend.Name, t.awsBackend.ResourceVersion)
	jc := t.clientset.BatchV1().Jobs(t.awsBackend.Namespace)
	job, err := jc.Get(name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	pc := t.clientset.CoreV1().Pods(t.awsBackend.Namespace)
	pods, err := pc.List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, pod := range pods.Items {
		if pod.Labels["controller-uid"] == job.Labels["controller-uid"] {
			return pc.Delete(pod.Name, &metav1.DeleteOptions{})
		}
	}
	return errors.New("failed to cleanup")
}

func (t TerraformClient) GetStatus() (v1beta1.BackendStatus, error) {
	endpoint, err := t.GetEndpointStatus()
	if err != nil {
		return *t.awsBackend.Status.DeepCopy(), err
	}
	if endpoint.IP == "" {
		addr, _ := net.ResolveIPAddr("ip", endpoint.DNS)
		if addr != nil {
			endpoint.IP = addr.String()
		}
	}
	var listeners []v1beta1.BackendListener
	for _, l := range t.awsBackend.Spec.Listeners {
		bl := v1beta1.BackendListener{
			Protocol: l.Protocol,
			Port:     l.Port,
		}
		listeners = append(listeners, bl)
	}

	status := t.awsBackend.Status.DeepCopy()
	status.Endpoint = endpoint
	status.Listeners = listeners
	status.Internal = t.awsBackend.Spec.Internal
	return *status, nil
}

func (t TerraformClient) getDNSRegex() string {
	return fmt.Sprintf("%s.*.elb.*.amazonaws.com", t.awsBackend.Name)
}

func (t TerraformClient) GetEndpointStatus() (v1beta1.BackendEndpoint, error) {
	cm, err := t.clientset.CoreV1().ConfigMaps(t.awsBackend.Namespace).Get(t.awsBackend.Name, metav1.GetOptions{})
	if err != nil {
		return v1beta1.BackendEndpoint{}, err
	}
	r := regexp.MustCompile(t.getDNSRegex())
	matches := r.FindStringSubmatch(cm.Data["show"])
	if len(matches) > 0 {
		return v1beta1.BackendEndpoint{
			DNS: matches[0],
		}, nil
	} else {
		return v1beta1.BackendEndpoint{}, errors.New("no dns record found in show")
	}
}

func (t TerraformClient) updateReport(job *batchv1.Job) error {
	logs := map[string]io.ReadCloser{}
	c := t.clientset.CoreV1().Pods(t.awsBackend.Namespace)
	pods, err := c.List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, pod := range pods.Items {
		if pod.Labels["controller-uid"] == job.Labels["controller-uid"] {
			logs["show"], err = c.GetLogs(pod.Name, &v1.PodLogOptions{Container: "show"}).Stream()
			if err != nil {
				return nil
			}
			logs["plan"], err = c.GetLogs(pod.Name, &v1.PodLogOptions{Container: "plan"}).Stream()
			if err != nil {
				return nil
			}
			break
		}
	}
	cmclient := t.clientset.CoreV1().ConfigMaps(t.awsBackend.Namespace)
	cm, err := cmclient.Get(t.awsBackend.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	for k, v := range logs {
		defer v.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(v)
		cm.Data[k] = buf.String()
	}
	_, err = cmclient.Update(cm)
	return err
}

func (t TerraformClient) isProcessing() (bool, error) {
	c := t.clientset.BatchV1().Jobs(t.awsBackend.Namespace)
	opt := metav1.ListOptions{}
	jobs, err := c.List(opt)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	for _, job := range jobs.Items {
		if !strings.Contains(job.Name, t.awsBackend.Name) {
			continue
		}
		if job.Status.CompletionTime != nil {
			return false, nil
		}

		if job.Status.Failed > 0 {
			return false, nil
		}
	}
	return false, nil
}

func (t TerraformClient) watchCompleteOrError() error {
	name := fmt.Sprintf("%s-%s", t.awsBackend.Name, t.awsBackend.ResourceVersion)
	namespace := t.awsBackend.Namespace
	c := t.clientset.BatchV1().Jobs(namespace)
	opt := metav1.GetOptions{}
	job, err := c.Get(name, opt)
	if err != nil {
		return err
	}

	// job is completed
	if job.Status.CompletionTime != nil {
		t.updateReport(job)
		return c.Delete(name, &metav1.DeleteOptions{})
	}

	// job is failed
	if job.Status.Failed > 0 {
		t.updateReport(job)
		c.Delete(name, &metav1.DeleteOptions{})
		return errors.New("Job errored")
	}

	// job is still processing
	return errors.New("Job is processing")
}

func (t TerraformClient) Apply() error {
	return t.execute("apply", true, true)
}

func (t TerraformClient) Destroy() error {
	return t.execute("destroy", true, true)
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
		tmpl, err := template.New(filepath.Base(path)).Funcs(sprig.TxtFuncMap()).ParseFiles(path)
		if err != nil {
			return "", err
		}
		result := bytes.Buffer{}
		var length int
		if len(awsBackend.Name) < 32 {
			length = len(awsBackend.Name)
		} else {
			length = 31
		}
		err = tmpl.Execute(&result, struct {
			Name      string
			B         *v1beta1.AWSBackend
			ServiceIn bool
		}{
			Name:      string([]rune(awsBackend.Name[:length])),
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
		Name:      fmt.Sprintf("%s-%s", t.awsBackend.Name, t.awsBackend.ResourceVersion),
		Namespace: t.awsBackend.Namespace,
	}
	cmd := []string{"/bin/terraform.sh", ops, t.awsBackend.Kind}
	if force {
		cmd = append(cmd, "true")
	}
	backOffLimit := int32(0)
	terminationGracePeriodSeconds := int64(300)
	vm := []corev1.VolumeMount{
		{
			Name:      t.awsBackend.Name,
			MountPath: "/data",
		},
	}
	env := []corev1.EnvVar{

		{
			Name:      "AWS_ACCESS_KEY_ID",
			ValueFrom: t.awsBackend.Spec.Credentials.AccesskeyID,
		},
		{
			Name:      "AWS_SECRET_ACCESS_KEY",
			ValueFrom: t.awsBackend.Spec.Credentials.SecretAccesskey,
		},
	}
	version := os.Getenv("DAEMON_VERSION")
	if version == "" {
		version = "latest"
	}
	image := fmt.Sprintf("%s:%s", "ghcr.io/takutakahashi/loadbalancer-controller-daemon", version)
	return batchv1.Job{
		ObjectMeta: om,
		Spec: batchv1.JobSpec{
			BackoffLimit: &backOffLimit,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy:                 corev1.RestartPolicyNever,
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
					InitContainers: []corev1.Container{
						{
							Name:         "plan",
							Image:        image,
							Command:      []string{"/bin/terraform.sh", "plan", t.awsBackend.Kind},
							VolumeMounts: vm,
							Env:          env,
						},
						{
							Name:         "tf",
							Image:        image,
							Command:      cmd,
							VolumeMounts: vm,
							Env:          env,
						},
					},
					Containers: []corev1.Container{
						{
							Name:         "show",
							Image:        image,
							Command:      []string{"/bin/terraform.sh", "show", t.awsBackend.Kind},
							VolumeMounts: vm,
							Env:          env,
						},
					},
					Volumes: []corev1.Volume{
						{
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
