resource "penguin_tencentcloud_virtual_machine" "example" {
  name                 = "vm.example.com"
  zone                 = "ap-guangzhou-6"
  instance_type        = "SA2.MEDIUM2"
  security_group       = "sg-5hilszwp"
  system_image         = "img-7efla8nv"
  vpc_id               = "vpc-oahbq6lh"
  subnet_id            = "subnet-95tfs6am"
  system_disk_size_gib = 20

  bandwidth_limit_mbps = 20
  total_transfer_kb    = 1048576
}
