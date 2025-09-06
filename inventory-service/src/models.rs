use serde::{Deserialize, Serialize};
use chrono::{DateTime, Utc};
use std::collections::HashMap;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct InventoryItem {
    pub id: String,
    pub product_id: String,
    pub quantity: i32,
    pub reserved: i32,
    pub available: i32,
    pub location: Option<String>,
    pub status: InventoryStatus,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Serialize, Deserialize, Clone, PartialEq)]
pub enum InventoryStatus {
    #[serde(rename = "in_stock")]
    InStock,
    #[serde(rename = "low_stock")]
    LowStock,
    #[serde(rename = "out_of_stock")]
    OutOfStock,
    #[serde(rename = "discontinued")]
    Discontinued,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct InventoryReservation {
    pub id: String,
    pub order_id: String,
    pub items: HashMap<String, i32>, // product_id -> quantity
    pub status: ReservationStatus,
    pub created_at: DateTime<Utc>,
    pub expires_at: Option<DateTime<Utc>>,
}

#[derive(Debug, Serialize, Deserialize, Clone, PartialEq)]
pub enum ReservationStatus {
    #[serde(rename = "pending")]
    Pending,
    #[serde(rename = "confirmed")]
    Confirmed,
    #[serde(rename = "released")]
    Released,
    #[serde(rename = "failed")]
    Failed,
}

// Request/Response DTOs
#[derive(Debug, Serialize, Deserialize)]
pub struct UpdateInventoryRequest {
    pub quantity: i32,
    pub location: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ReserveInventoryRequest {
    pub order_id: String,
    pub items: HashMap<String, i32>, // product_id -> quantity
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ReleaseInventoryRequest {
    pub order_id: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct InventoryResponse {
    pub item: InventoryItem,
    pub success: bool,
    pub message: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ReservationResponse {
    pub reservation: InventoryReservation,
    pub success: bool,
    pub message: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ErrorResponse {
    pub error: String,
    pub code: String,
}
