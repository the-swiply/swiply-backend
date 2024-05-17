resource "yandex_mdb_postgresql_cluster" "mdb" {
  folder_id           = local.folder_id
  name                = "postgres"
  environment         = "PRODUCTION"
  network_id          = yandex_vpc_network.k8s-network.id

  config {
    version = 15
    resources {
      resource_preset_id = "b2.medium"
      disk_type_id       = "network-ssd"
      disk_size          = 10
    }

    access {
      web_sql = true
    }
  }

  host {
    zone             = yandex_vpc_subnet.subnet-a.zone
    name             = "pg-host-1"
    subnet_id        = yandex_vpc_subnet.subnet-a.id
    assign_public_ip = true
  }
}

resource "yandex_mdb_postgresql_database" "profile" {
  cluster_id = yandex_mdb_postgresql_cluster.mdb.id
  name       = "profile"
  owner      = var.profile-owner.login
  depends_on = [
    yandex_mdb_postgresql_user.profile-owner
  ]
}

resource "yandex_mdb_postgresql_user" "profile-owner" {
  cluster_id = yandex_mdb_postgresql_cluster.mdb.id
  name       = var.profile-owner.login
  password   = var.profile-owner.password
  conn_limit = 20
}

resource "yandex_mdb_postgresql_database" "notification" {
  cluster_id = yandex_mdb_postgresql_cluster.mdb.id
  name       = "notification"
  owner      = var.notification-owner.login
  depends_on = [
    yandex_mdb_postgresql_user.notification-owner
  ]
}

resource "yandex_mdb_postgresql_user" "notification-owner" {
  cluster_id = yandex_mdb_postgresql_cluster.mdb.id
  name       = var.notification-owner.login
  password   = var.notification-owner.password
  conn_limit = 20
}

resource "yandex_mdb_postgresql_database" "chat" {
  cluster_id = yandex_mdb_postgresql_cluster.mdb.id
  name       = "chat"
  owner      = var.chat-owner.login
  depends_on = [
    yandex_mdb_postgresql_user.chat-owner
  ]
}

resource "yandex_mdb_postgresql_user" "chat-owner" {
  cluster_id = yandex_mdb_postgresql_cluster.mdb.id
  name       = var.chat-owner.login
  password   = var.chat-owner.password
  conn_limit = 20
}

resource "yandex_mdb_postgresql_database" "event" {
  cluster_id = yandex_mdb_postgresql_cluster.mdb.id
  name       = "event"
  owner      = var.event-owner.login
  depends_on = [
    yandex_mdb_postgresql_user.event-owner
  ]
}

resource "yandex_mdb_postgresql_user" "event-owner" {
  cluster_id = yandex_mdb_postgresql_cluster.mdb.id
  name       = var.event-owner.login
  password   = var.event-owner.password
  conn_limit = 20
}

resource "yandex_mdb_postgresql_database" "recommendation" {
  cluster_id = yandex_mdb_postgresql_cluster.mdb.id
  name       = "recommendation"
  owner      = var.recommendation-owner.login
  depends_on = [
    yandex_mdb_postgresql_user.recommendation-owner
  ]
}

resource "yandex_mdb_postgresql_user" "recommendation-owner" {
  cluster_id = yandex_mdb_postgresql_cluster.mdb.id
  name       = var.recommendation-owner.login
  password   = var.recommendation-owner.password
  conn_limit = 20
}

resource "yandex_mdb_postgresql_database" "randomcoffee" {
  cluster_id = yandex_mdb_postgresql_cluster.mdb.id
  name       = "randomcoffee"
  owner      = var.randomcoffee-owner.login
  depends_on = [
    yandex_mdb_postgresql_user.randomcoffee-owner
  ]
}

resource "yandex_mdb_postgresql_user" "randomcoffee-owner" {
  cluster_id = yandex_mdb_postgresql_cluster.mdb.id
  name       = var.randomcoffee-owner.login
  password   = var.randomcoffee-owner.password
  conn_limit = 20
}

resource "yandex_mdb_postgresql_database" "backoffice" {
  cluster_id = yandex_mdb_postgresql_cluster.mdb.id
  name       = "backoffice"
  owner      = var.backoffice-owner.login
  depends_on = [
    yandex_mdb_postgresql_user.backoffice-owner
  ]
}

resource "yandex_mdb_postgresql_user" "backoffice-owner" {
  cluster_id = yandex_mdb_postgresql_cluster.mdb.id
  name       = var.backoffice-owner.login
  password   = var.backoffice-owner.password
  conn_limit = 20
}