terraform {
  backend "s3" {
      bucket  = "{{ .B.Spec.BucketName }}"
      key     = "{{ .B.Namespace }}/{{ .B.Name }}.tfstate"
      region  = "{{ .B.Spec.Region }}"
    }
}
