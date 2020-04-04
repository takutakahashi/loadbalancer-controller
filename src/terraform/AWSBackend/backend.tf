terraform {
  backend "s3" {
      bucket  = "takutakahashi-tfstate"
      key     = "lb-name"
      region  = "ap-northeast-1"
      profile = "default"
    }
}
