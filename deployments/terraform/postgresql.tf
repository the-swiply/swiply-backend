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

resource "yandex_mdb_postgresql_database" "srv" {
  cluster_id = yandex_mdb_postgresql_cluster.mdb.id
  name       = "srv"
  owner      = var.postgres-admin.login
  depends_on = [
    yandex_mdb_postgresql_user.postgres-admin
  ]
}

resource "yandex_mdb_postgresql_user" "postgres-admin" {
  cluster_id = yandex_mdb_postgresql_cluster.mdb.id
  name       = var.postgres-admin.login
  password   = var.postgres-admin.password
}