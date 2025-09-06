use serde::{Deserialize, Serialize};
use chrono::{DateTime, Utc};
use std::collections::HashMap;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Order {
    pub id: String,
    pub user_id: String,
    pub items: Vec<OrderItem>,
    pub total: f64,
    pub status: OrderStatus,
    pub shipping_address: Address,
    pub billing_address: Address,
    pub payment_id: Option<String>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct OrderItem {
    pub product_id: String,
    pub name: String,
    pub price: f64,
    pub quantity: i32,
    pub subtotal: f64,
}

#[derive(Debug, Serialize, Deserialize, Clone, PartialEq)]
pub enum OrderStatus {
    #[serde(rename = "pending")]
    Pending,
    #[serde(rename = "processing")]
    Processing,
    #[serde(rename = "shipped")]
    Shipped,
    #[serde(rename = "delivered")]
    Delivered,
    #[serde(rename = "cancelled")]
    Cancelled,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Address {
    pub street: String,
    pub city: String,
    pub state: String,
    pub postal_code: String,
    pub country: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CreateOrderRequest {
    pub user_id: String,
    pub items: Vec<OrderItemRequest>,
    pub shipping_address: Address,
    pub billing_address: Address,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct OrderItemRequest {
    pub product_id: String,
    pub quantity: i32,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct UpdateOrderStatusRequest {
    pub status: OrderStatus,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct OrderResponse {
    pub order: Order,
    pub success: bool,
    pub message: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ErrorResponse {
    pub error: String,
    pub code: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ProductResponse {
    pub id: String,
    pub name: String,
    pub price: f64,
    pub description: String,
    pub category: String,
    pub image_url: String,
    pub stock: i32,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ReserveInventoryRequest {
    pub order_id: String,
    pub items: HashMap<String, i32>, // product_id -> quantity
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CreatePaymentRequest {
    pub order_id: String,
    pub amount: f64,
    pub currency: String,
    pub payment_method: String,
}
