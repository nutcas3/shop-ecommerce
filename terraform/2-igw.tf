resource "aws_internet_gateway" "shop-ecommerce-igw" {
  vpc_id = aws_vpc.shop-ecommerce-vpc.id
}
