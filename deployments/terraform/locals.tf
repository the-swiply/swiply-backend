locals {
  zone_a_v4_cidr_blocks = "10.1.0.0/16" # Set the CIDR block for subnet in the ru-central1-a availability zone.
  folder_id             = "b1grti8d755eddfv0jcp"            # Set your cloud folder ID.
  k8s_version           = "1.27"            # Set the Kubernetes version.
  sa_name               = "cluster-swiply"            # Set a service account name. It must be unique in a cloud.
}
