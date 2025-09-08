# Create a VPC with a CIDR block of 10.0.0.0/16 and approximately 2^16 IP addresses
resource "aws_vpc" "shop-ecommerce-vpc" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "shop-ecommerce-vpc"
  }
}
