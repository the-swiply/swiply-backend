variable "profile-owner" {
  type = object({
    login    = string
    password = string
  })

  sensitive = true
}

variable "notification-owner" {
  type = object({
    login    = string
    password = string
  })

  sensitive = true
}

variable "chat-owner" {
  type = object({
    login    = string
    password = string
  })

  sensitive = true
}

variable "event-owner" {
  type = object({
    login    = string
    password = string
  })

  sensitive = true
}

variable "recommendation-owner" {
  type = object({
    login    = string
    password = string
  })

  sensitive = true
}

variable "randomcoffee-owner" {
  type = object({
    login    = string
    password = string
  })

  sensitive = true
}

variable "backoffice-owner" {
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
