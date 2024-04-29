output "aws-access-key-id" {
  value = yandex_iam_service_account_static_access_key.sa-static-key.access_key
}

output "aws-access-secret-key" {
  value     = yandex_iam_service_account_static_access_key.sa-static-key.secret_key
  sensitive = true
}
