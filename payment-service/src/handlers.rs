use actix_web::{web, HttpResponse, Responder};
use chrono::Utc;
use log::{info, error};
use uuid::Uuid;
use std::sync::Mutex;
use std::collections::HashMap;

use crate::models::{
    Payment, PaymentStatus, CreatePaymentRequest, RefundRequest, PaymentResponse
};

pub struct AppState {
    pub payments: Mutex<HashMap<String, Payment>>,
}
pub async fn create_payment(
    data: web::Data<AppState>,
    req: web::Json<CreatePaymentRequest>,
) -> impl Responder {
    let payment_id = Uuid::new_v4().to_string();
    let now = Utc::now();
    
    let payment = Payment {
        id: payment_id.clone(),
        order_id: req.order_id.clone(),
        amount: req.amount,
        currency: req.currency.clone(),
        status: PaymentStatus::Processing,
        payment_method: req.payment_method.clone(),
        created_at: now,
        updated_at: now,
        transaction_id: Some(format!("txn_{}", Uuid::new_v4().to_string())),
        metadata: req.metadata.clone(),
    };
    
    // In a real implementation, we would process the payment with a payment gateway
    // For now, we'll simulate a successful payment
    let mut payment_success = payment.clone();
    payment_success.status = PaymentStatus::Completed;
    
    // Store the payment in our in-memory database
    data.payments.lock().unwrap().insert(payment_id.clone(), payment_success.clone());
    
    info!("Created payment: {} for order: {}", payment_id, req.order_id);
    
    // Return the created payment
    HttpResponse::Created().json(PaymentResponse {
        payment: payment_success,
        success: true,
        message: Some("Payment processed successfully".to_string()),
    })
}

pub async fn get_payment(
    data: web::Data<AppState>,
    path: web::Path<String>,
) -> impl Responder {
    let payment_id = path.into_inner();
    
    // Get the payment from our in-memory database
    let payments = data.payments.lock().unwrap();
    match payments.get(&payment_id) {
        Some(payment) => {
            HttpResponse::Ok().json(PaymentResponse {
                payment: payment.clone(),
                success: true,
                message: None,
            })
        },
        None => {
            error!("Payment not found: {}", payment_id);
            HttpResponse::NotFound().json(serde_json::json!({
                "error": "Payment not found",
                "code": "payment_not_found"
            }))
        }
    }
}

pub async fn get_payments_by_order(
    data: web::Data<AppState>,
    path: web::Path<String>,
) -> impl Responder {
    let order_id = path.into_inner();
    
    let payments = data.payments.lock().unwrap();
    let order_payments: Vec<Payment> = payments
        .values()
        .filter(|p| p.order_id == order_id)
        .cloned()
        .collect();
    
    HttpResponse::Ok().json(order_payments)
}

pub async fn refund_payment(
    data: web::Data<AppState>,
    path: web::Path<String>,
    req: web::Json<RefundRequest>,
) -> impl Responder {
    let payment_id = path.into_inner();
    
    // Get the payment from our in-memory database
    let mut payments = data.payments.lock().unwrap();
    match payments.get_mut(&payment_id) {
        Some(payment) => {
            // Check if the payment can be refunded
            if payment.status != PaymentStatus::Completed {
                return HttpResponse::BadRequest().json(serde_json::json!({
                    "error": "Payment cannot be refunded",
                    "code": "invalid_refund_state"
                }));
            }
            
            // Update the payment status
            payment.status = PaymentStatus::Refunded;
            payment.updated_at = Utc::now();
            
            info!("Refunded payment: {} for order: {}", payment_id, payment.order_id);
            
            HttpResponse::Ok().json(PaymentResponse {
                payment: payment.clone(),
                success: true,
                message: Some(format!("Payment {} refunded successfully", payment_id)),
            })
        },
        None => {
            error!("Payment not found: {}", payment_id);
            HttpResponse::NotFound().json(serde_json::json!({
                "error": "Payment not found",
                "code": "payment_not_found"
            }))
        }
    }
}

pub async fn health_check() -> impl Responder {
    HttpResponse::Ok().json(serde_json::json!({
        "status": "ok",
        "service": "payment-service",
        "version": env!("CARGO_PKG_VERSION"),
        "timestamp": Utc::now().to_rfc3339()
    }))
}
