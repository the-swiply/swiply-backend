resource "yandex_iam_service_account" "k8s-sa" {
  folder_id = local.folder_id
  name      = local.sa_name
}

resource "yandex_resourcemanager_folder_iam_binding" "editor" {
  # Assign "editor" role to service account.
  folder_id = local.folder_id
  role      = "editor"
  members   = [
    "serviceAccount:${yandex_iam_service_account.k8s-sa.id}"
  ]
}

resource "yandex_resourcemanager_folder_iam_binding" "images-puller" {
  # Assign "container-registry.images.puller" role to service account.
  folder_id = local.folder_id
  role      = "container-registry.images.puller"
  members   = [
    "serviceAccount:${yandex_iam_service_account.k8s-sa.id}"
  ]
}
