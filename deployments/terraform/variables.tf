variable "postgres-admin" {
  type = object({
    login    = string
    password = string
  })

  sensitive = true
}