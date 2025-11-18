use serde::{Deserialize, Serialize};
use chrono::{DateTime, Utc};

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Payment {
    pub id: String,
    pub order_id: String,
    pub amount: f64,
    pub currency: String,
    pub status: PaymentStatus,
    pub payment_method: PaymentMethod,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub transaction_id: Option<String>,
    pub metadata: Option<PaymentMetadata>,
}

#[derive(Debug, Serialize, Deserialize, Clone, PartialEq)]
pub enum PaymentStatus {
    #[serde(rename = "pending")]
    Pending,
    #[serde(rename = "processing")]
    Processing,
    #[serde(rename = "completed")]
    Completed,
    #[serde(rename = "failed")]
    Failed,
    #[serde(rename = "refunded")]
    Refunded,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum PaymentMethod {
    #[serde(rename = "credit_card")]
    CreditCard,
    #[serde(rename = "paypal")]
    PayPal,
    #[serde(rename = "bank_transfer")]
    BankTransfer,
    #[serde(rename = "crypto")]
    Crypto,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct PaymentMetadata {
    pub card_last_four: Option<String>,
    pub card_brand: Option<String>,
    pub customer_email: Option<String>,
    pub customer_name: Option<String>,
    pub billing_address: Option<BillingAddress>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BillingAddress {
    pub street: String,
    pub city: String,
    pub state: String,
    pub postal_code: String,
    pub country: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CreatePaymentRequest {
    pub order_id: String,
    pub amount: f64,
    pub currency: String,
    pub payment_method: PaymentMethod,
    pub metadata: Option<PaymentMetadata>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct RefundRequest {
    pub amount: Option<f64>,
    pub reason: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct PaymentResponse {
    pub payment: Payment,
    pub success: bool,
    pub message: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ErrorResponse {
    pub error: String,
    pub code: String,
}
