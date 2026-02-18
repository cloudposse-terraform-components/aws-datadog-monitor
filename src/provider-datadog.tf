module "datadog_configuration" {
  source  = "github.com/cloudposse-terraform-components/aws-datadog-credentials//src/modules/datadog_keys?ref=v1.535.13"
  enabled = true
  context = module.this.context
}

provider "datadog" {
  api_key  = module.datadog_configuration.datadog_api_key
  app_key  = module.datadog_configuration.datadog_app_key
  api_url  = module.datadog_configuration.datadog_api_url
  validate = local.enabled
}
