resource "penguin_tencentcloud_elastic_ip" "example" {
  region               = "ap-guangzhou"
  bandwidth_limit_mbps = 20
  address_name         = "prod-eip-01"
}
