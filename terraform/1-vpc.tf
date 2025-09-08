resource "aws_vpc" "shop-ecommerce-vpc" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "shop-ecommerce-vpc"
  }
}
