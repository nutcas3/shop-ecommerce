# allows our VPC to be accessible from the internet
resource "aws_internet_gateway" "shop-ecommerce-igw" {
  vpc_id = aws_vpc.shop-ecommerce-vpc.id
}
