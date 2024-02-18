variable "postgres-admin" {
  type = object({
    login    = string
    password = string
  })

  sensitive = true
}

variable "redis-admin-password" {
  type = string

  sensitive = true
}
