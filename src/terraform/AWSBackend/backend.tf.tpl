terraform {
  backend "s3" {
      bucket  = "takutakahashi-tfstate"
      key     = "{{ .B.Namespace }}/{{ .B.Name }}.tfstate"
      region  = "{{ .B.Spec.Resion }}"
    }
}
