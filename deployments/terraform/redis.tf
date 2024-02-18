resource "yandex_mdb_redis_cluster" "mdb" {
  folder_id        = local.folder_id
  name             = "redis"
  environment      = "PRODUCTION"
  network_id       = yandex_vpc_network.k8s-network.id
  persistence_mode = "ON"
  tls_enabled      = true

  config {
    password  = var.redis-admin-password
    version   = "7.2"
    databases = 100
  }

  resources {
    resource_preset_id = "b3-c1-m4"
    disk_size          = 24
    disk_type_id       = "network-ssd"
  }

  host {
    zone             = yandex_vpc_subnet.subnet-a.zone
    subnet_id        = yandex_vpc_subnet.subnet-a.id
    assign_public_ip = true
  }
}
