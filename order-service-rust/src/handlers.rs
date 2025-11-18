use actix_web::{web, HttpResponse, Responder};
use chrono::Utc;
use log::{info, error};
use uuid::Uuid;
use std::sync::Mutex;
use std::collections::HashMap;
use reqwest::Client;

use crate::models::{
    Order, OrderItem, OrderStatus, CreateOrderRequest, UpdateOrderStatusRequest,
    OrderResponse, ProductResponse, ReserveInventoryRequest, CreatePaymentRequest
};
use crate::config::Config;

pub struct AppState {
    pub orders: Mutex<HashMap<String, Order>>,
    pub http_client: Client,
    pub config: Config,
}

pub async fn create_order(
    data: web::Data<AppState>,
    req: web::Json<CreateOrderRequest>,
) -> impl Responder {
    let order_id = Uuid::new_v4().to_string();
    let now = Utc::now();
    
    let mut order_items = Vec::new();
    let mut total = 0.0;
    let mut product_quantities = HashMap::new();
    
    for item in &req.items {
        // In a real implementation, we would fetch product details from the product service
        // For now, we'll simulate this with a mock response
        let product_url = format!(
            "{}/api/products/{}", 
            data.config.product_service_url, 
            item.product_id
        );
        
        match data.http_client.get(&product_url).send().await {
            Ok(response) => {
                if response.status().is_success() {
                    match response.json::<ProductResponse>().await {
                        Ok(product) => {
                            let subtotal = product.price * item.quantity as f64;
                            order_items.push(OrderItem {
                                product_id: item.product_id.clone(),
                                name: product.name,
                                price: product.price,
                                quantity: item.quantity,
                                subtotal,
                            });
                            total += subtotal;
                            product_quantities.insert(item.product_id.clone(), item.quantity);
                        },
                        Err(e) => {
                            error!("Failed to parse product response: {}", e);
                            return HttpResponse::InternalServerError().json(serde_json::json!({
                                "error": "Failed to fetch product details",
                                "code": "product_fetch_error"
                            }));
                        }
                    }
                } else {
                    error!("Product service returned error: {}", response.status());
                    return HttpResponse::BadRequest().json(serde_json::json!({
                        "error": format!("Product not found: {}", item.product_id),
                        "code": "product_not_found"
                    }));
                }
            },
            Err(e) => {
                error!("Failed to connect to product service: {}", e);
                return HttpResponse::InternalServerError().json(serde_json::json!({
                    "error": "Failed to connect to product service",
                    "code": "service_unavailable"
                }));
            }
        }
    }
    
    let inventory_url = format!("{}/api/inventory/reserve", data.config.inventory_service_url);
    let reserve_request = ReserveInventoryRequest {
        order_id: order_id.clone(),
        items: product_quantities,
    };
    
    match data.http_client.post(&inventory_url)
        .json(&reserve_request)
        .send()
        .await 
    {
        Ok(response) => {
            if !response.status().is_success() {
                error!("Inventory service returned error: {}", response.status());
                return HttpResponse::BadRequest().json(serde_json::json!({
                    "error": "Failed to reserve inventory",
                    "code": "inventory_reservation_failed"
                }));
            }
        },
        Err(e) => {
            error!("Failed to connect to inventory service: {}", e);
            return HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "Failed to connect to inventory service",
                "code": "service_unavailable"
            }));
        }
    }
    
    let order = Order {
        id: order_id.clone(),
        user_id: req.user_id.clone(),
        items: order_items,
        total,
        status: OrderStatus::Pending,
        shipping_address: req.shipping_address.clone(),
        billing_address: req.billing_address.clone(),
        payment_id: None,
        created_at: now,
        updated_at: now,
    };
    
    data.orders.lock().unwrap().insert(order_id.clone(), order.clone());
    
    info!("Created order: {} for user: {}", order_id, req.user_id);
    
    HttpResponse::Created().json(OrderResponse {
        order,
        success: true,
        message: Some("Order created successfully".to_string()),
    })
}

pub async fn get_order(
    data: web::Data<AppState>,
    path: web::Path<String>,
) -> impl Responder {
    let order_id = path.into_inner();
    
    let orders = data.orders.lock().unwrap();
    match orders.get(&order_id) {
        Some(order) => {
            HttpResponse::Ok().json(OrderResponse {
                order: order.clone(),
                success: true,
                message: None,
            })
        },
        None => {
            error!("Order not found: {}", order_id);
            HttpResponse::NotFound().json(serde_json::json!({
                "error": "Order not found",
                "code": "order_not_found"
            }))
        }
    }
}

pub async fn get_user_orders(
    data: web::Data<AppState>,
    path: web::Path<String>,
) -> impl Responder {
    let user_id = path.into_inner();
    
    let orders = data.orders.lock().unwrap();
    let user_orders: Vec<Order> = orders
        .values()
        .filter(|o| o.user_id == user_id)
        .cloned()
        .collect();
    
    HttpResponse::Ok().json(user_orders)
}

pub async fn update_order_status(
    data: web::Data<AppState>,
    path: web::Path<String>,
    req: web::Json<UpdateOrderStatusRequest>,
) -> impl Responder {
    let order_id = path.into_inner();
    
    let mut orders = data.orders.lock().unwrap();
    match orders.get_mut(&order_id) {
        Some(order) => {
            order.status = req.status.clone();
            order.updated_at = Utc::now();
            
            if req.status == OrderStatus::Cancelled {
                let inventory_url = format!(
                    "{}/api/inventory/reservation/{}/release", 
                    data.config.inventory_service_url, 
                    order_id
                );
                
                match data.http_client.post(&inventory_url).send().await {
                    Ok(_) => {
                        info!("Released inventory for cancelled order: {}", order_id);
                    },
                    Err(e) => {
                        error!("Failed to release inventory: {}", e);
                    }
                }
            }
            
            info!("Updated order status: {} to {:?}", order_id, req.status);
            
            HttpResponse::Ok().json(OrderResponse {
                order: order.clone(),
                success: true,
                message: Some(format!("Order status updated to {:?}", req.status)),
            })
        },
        None => {
            error!("Order not found: {}", order_id);
            HttpResponse::NotFound().json(serde_json::json!({
                "error": "Order not found",
                "code": "order_not_found"
            }))
        }
    }
}

pub async fn process_payment(
    data: web::Data<AppState>,
    path: web::Path<String>,
) -> impl Responder {
    let order_id = path.into_inner();
    
    // Get the order from our in-memory database
    let mut orders = data.orders.lock().unwrap();
    let order = match orders.get(&order_id) {
        Some(o) => o.clone(),
        None => {
            error!("Order not found: {}", order_id);
            return HttpResponse::NotFound().json(serde_json::json!({
                "error": "Order not found",
                "code": "order_not_found"
            }));
        }
    };
    
    // Check if payment has already been processed
    if order.payment_id.is_some() {
        return HttpResponse::BadRequest().json(serde_json::json!({
            "error": "Payment has already been processed for this order",
            "code": "payment_already_processed"
        }));
    }
    
    // Process payment
    let payment_url = format!("{}/api/payments", data.config.payment_service_url);
    let payment_request = CreatePaymentRequest {
        order_id: order_id.clone(),
        amount: order.total,
        currency: "USD".to_string(),
        payment_method: "credit_card".to_string(),
    };
    
    match data.http_client.post(&payment_url)
        .json(&payment_request)
        .send()
        .await 
    {
        Ok(response) => {
            if response.status().is_success() {
                match response.json::<serde_json::Value>().await {
                    Ok(payment_response) => {
                        // Extract payment ID from response
                        let payment_id = payment_response["payment"]["id"].as_str()
                            .unwrap_or("unknown")
                            .to_string();
                        
                        // Update order with payment ID and change status to Processing
                        if let Some(order) = orders.get_mut(&order_id) {
                            order.payment_id = Some(payment_id.clone());
                            order.status = OrderStatus::Processing;
                            order.updated_at = Utc::now();
                            
                            // Confirm inventory reservation
                            let inventory_url = format!(
                                "{}/api/inventory/reservation/{}/confirm", 
                                data.config.inventory_service_url, 
                                order_id
                            );
                            
                            match data.http_client.post(&inventory_url).send().await {
                                Ok(_) => {
                                    info!("Confirmed inventory reservation for order: {}", order_id);
                                },
                                Err(e) => {
                                    error!("Failed to confirm inventory reservation: {}", e);
                                    // Continue with order processing even if confirmation fails
                                }
                            }
                            
                            info!("Processed payment for order: {}, payment ID: {}", order_id, payment_id);
                            
                            return HttpResponse::Ok().json(OrderResponse {
                                order: order.clone(),
                                success: true,
                                message: Some("Payment processed successfully".to_string()),
                            });
                        }
                    },
                    Err(e) => {
                        error!("Failed to parse payment response: {}", e);
                    }
                }
            }
            
            error!("Payment service returned error: {}", response.status());
            HttpResponse::BadRequest().json(serde_json::json!({
                "error": "Failed to process payment",
                "code": "payment_processing_failed"
            }))
        },
        Err(e) => {
            error!("Failed to connect to payment service: {}", e);
            HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "Failed to connect to payment service",
                "code": "service_unavailable"
            }))
        }
    }
}

// Health check endpoint
pub async fn health_check() -> impl Responder {
    HttpResponse::Ok().json(serde_json::json!({
        "status": "ok",
        "service": "order-service",
        "version": env!("CARGO_PKG_VERSION"),
        "timestamp": Utc::now().to_rfc3339()
    }))
}
