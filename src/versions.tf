terraform {
  required_version = ">= 1.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.9.0, < 7.0.0"
    }
    datadog = {
      source  = "datadog/datadog"
      version = ">= 4.0.0"
    }
  }
}
