package terraform

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"text/template"

	"github.com/google/uuid"
	"github.com/takutakahashi/loadbalancer-controller/api/v1beta1"
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

func (t TerraformClient) execute(c string, force bool) ([]byte, error) {
	prev, err := os.Getwd()
	if err != nil {
		return []byte{}, err
	}
	if err = os.Chdir(t.workDir()); err != nil {
		return []byte{}, err
	}
	defer os.Chdir(prev)
	err = t.init()
	var cmd *exec.Cmd
	if force {
		cmd = exec.Command("terraform", c, "-auto-approve", "-var-file", t.tfvarsPath)
	} else {
		cmd = exec.Command("terraform", c, "-var-file", t.tfvarsPath)
	}
	cmd.Env = os.Environ()
	return cmd.Output()
}

func (t TerraformClient) Apply() error {
	_, err := t.execute("apply", true)
	return err
}

func (t TerraformClient) Destroy() error {
	_, err := t.execute("destroy", true)
	return err
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
