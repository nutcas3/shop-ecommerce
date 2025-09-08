resource "aws_eip" "shop-ecommerce-nat" {
  
}

resource "aws_nat_gateway" "nat" {
  allocation_id = aws_eip.shop-ecommerce-nat.id
  subnet_id     = aws_subnet.public-us-east-1a.id 
  depends_on = [aws_internet_gateway.shop-ecommerce-igw]
}
