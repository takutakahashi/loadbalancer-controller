terraform {
  backend "s3" {
      bucket  = "loadbalancer-controller-tfstate"
      key     = "{{ .B.Namespace }}/{{ .B.Name }}.tfstate"
      region  = "{{ .B.Spec.Region }}"
    }
}
