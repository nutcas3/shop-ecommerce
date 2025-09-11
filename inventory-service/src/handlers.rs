use actix_web::{web, HttpResponse, Responder};
use chrono::{Duration, Utc};
use log::{info, error};
use uuid::Uuid;
use std::sync::Mutex;
use std::collections::HashMap;

use crate::models::{
    InventoryItem, InventoryStatus, InventoryReservation, ReservationStatus,
    UpdateInventoryRequest, ReserveInventoryRequest, ReleaseInventoryRequest,
    InventoryResponse, ReservationResponse
};

pub struct AppState {
    pub inventory: Mutex<HashMap<String, InventoryItem>>,
    pub reservations: Mutex<HashMap<String, InventoryReservation>>,
}

pub async fn get_inventory_item(
    data: web::Data<AppState>,
    path: web::Path<String>,
) -> impl Responder {
    let product_id = path.into_inner();
    
    let inventory = data.inventory.lock().unwrap();
    match inventory.get(&product_id) {
        Some(item) => {
            HttpResponse::Ok().json(InventoryResponse {
                item: item.clone(),
                success: true,
                message: None,
            })
        },
        None => {
            error!("Inventory item not found: {}", product_id);
            HttpResponse::NotFound().json(serde_json::json!({
                "error": "Inventory item not found",
                "code": "inventory_not_found"
            }))
        }
    }
}

pub async fn get_all_inventory(
    data: web::Data<AppState>,
) -> impl Responder {
    let inventory = data.inventory.lock().unwrap();
    let items: Vec<InventoryItem> = inventory.values().cloned().collect();
    
    HttpResponse::Ok().json(items)
}

pub async fn create_inventory_item(
    data: web::Data<AppState>,
    path: web::Path<String>,
    req: web::Json<UpdateInventoryRequest>,
) -> impl Responder {
    let product_id = path.into_inner();
    let now = Utc::now();
    
    let mut inventory = data.inventory.lock().unwrap();
    if inventory.contains_key(&product_id) {
        return HttpResponse::BadRequest().json(serde_json::json!({
            "error": "Inventory item already exists",
            "code": "inventory_already_exists"
        }));
    }
    
    let status = if req.quantity > 10 {
        InventoryStatus::InStock
    } else if req.quantity > 0 {
        InventoryStatus::LowStock
    } else {
        InventoryStatus::OutOfStock
    };
    
    let item = InventoryItem {
        id: product_id.clone(),
        product_id: product_id.clone(),
        quantity: req.quantity,
        reserved: 0,
        available: req.quantity,
        location: req.location.clone(),
        status,
        created_at: now,
        updated_at: now,
    };
    
    inventory.insert(product_id.clone(), item.clone());
    
    info!("Created inventory item for product: {}", product_id);
    
    HttpResponse::Created().json(InventoryResponse {
        item,
        success: true,
        message: Some("Inventory item created successfully".to_string()),
    })
}

pub async fn update_inventory_item(
    data: web::Data<AppState>,
    path: web::Path<String>,
    req: web::Json<UpdateInventoryRequest>,
) -> impl Responder {
    let product_id = path.into_inner();
    
    let mut inventory = data.inventory.lock().unwrap();
    match inventory.get_mut(&product_id) {
        Some(item) => {
            item.quantity = req.quantity;
            item.available = req.quantity - item.reserved;
            item.location = req.location.clone();
            item.updated_at = Utc::now();
            
            item.status = if item.available > 10 {
                InventoryStatus::InStock
            } else if item.available > 0 {
                InventoryStatus::LowStock
            } else {
                InventoryStatus::OutOfStock
            };
            
            info!("Updated inventory item for product: {}", product_id);
            
            HttpResponse::Ok().json(InventoryResponse {
                item: item.clone(),
                success: true,
                message: Some("Inventory item updated successfully".to_string()),
            })
        },
        None => {
            error!("Inventory item not found: {}", product_id);
            HttpResponse::NotFound().json(serde_json::json!({
                "error": "Inventory item not found",
                "code": "inventory_not_found"
            }))
        }
    }
}

pub async fn reserve_inventory(
    data: web::Data<AppState>,
    req: web::Json<ReserveInventoryRequest>,
) -> impl Responder {
    let reservation_id = Uuid::new_v4().to_string();
    let now = Utc::now();
    
    let mut inventory = data.inventory.lock().unwrap();
    for (product_id, quantity) in &req.items {
        match inventory.get(product_id) {
            Some(item) => {
                if item.available < *quantity {
                    return HttpResponse::BadRequest().json(serde_json::json!({
                        "error": format!("Not enough inventory for product: {}", product_id),
                        "code": "insufficient_inventory"
                    }));
                }
            },
            None => {
                return HttpResponse::BadRequest().json(serde_json::json!({
                    "error": format!("Product not found in inventory: {}", product_id),
                    "code": "product_not_found"
                }));
            }
        }
    }
    
    for (product_id, quantity) in &req.items {
        if let Some(item) = inventory.get_mut(product_id) {
            item.reserved += *quantity;
            item.available -= *quantity;
            item.updated_at = now;
            
            item.status = if item.available > 10 {
                InventoryStatus::InStock
            } else if item.available > 0 {
                InventoryStatus::LowStock
            } else {
                InventoryStatus::OutOfStock
            };
        }
    }
    
    let reservation = InventoryReservation {
        id: reservation_id.clone(),
        order_id: req.order_id.clone(),
        items: req.items.clone(),
        status: ReservationStatus::Pending,
        created_at: now,
        expires_at: Some(now + Duration::minutes(30)), // Reservation expires in 30 minutes
    };
    
    let mut reservations = data.reservations.lock().unwrap();
    reservations.insert(reservation_id.clone(), reservation.clone());
    
    info!("Created inventory reservation: {} for order: {}", reservation_id, req.order_id);
    
    HttpResponse::Created().json(ReservationResponse {
        reservation,
        success: true,
        message: Some("Inventory reserved successfully".to_string()),
    })
}

pub async fn confirm_reservation(
    data: web::Data<AppState>,
    path: web::Path<String>,
) -> impl Responder {
    let reservation_id = path.into_inner();
    
    let mut reservations = data.reservations.lock().unwrap();
    match reservations.get_mut(&reservation_id) {
        Some(reservation) => {
            reservation.status = ReservationStatus::Confirmed;
            reservation.expires_at = None; // Remove expiration
            
            info!("Confirmed inventory reservation: {}", reservation_id);
            
            HttpResponse::Ok().json(ReservationResponse {
                reservation: reservation.clone(),
                success: true,
                message: Some("Reservation confirmed successfully".to_string()),
            })
        },
        None => {
            error!("Reservation not found: {}", reservation_id);
            HttpResponse::NotFound().json(serde_json::json!({
                "error": "Reservation not found",
                "code": "reservation_not_found"
            }))
        }
    }
}

pub async fn release_reservation(
    data: web::Data<AppState>,
    path: web::Path<String>,
) -> impl Responder {
    let reservation_id = path.into_inner();
    
    // Get the reservation
    let mut reservations = data.reservations.lock().unwrap();
    let reservation = match reservations.get_mut(&reservation_id) {
        Some(r) => r.clone(),
        None => {
            error!("Reservation not found: {}", reservation_id);
            return HttpResponse::NotFound().json(serde_json::json!({
                "error": "Reservation not found",
                "code": "reservation_not_found"
            }));
        }
    };
    
    if let Some(r) = reservations.get_mut(&reservation_id) {
        r.status = ReservationStatus::Released;
    }
    
    let mut inventory = data.inventory.lock().unwrap();
    let now = Utc::now();
    
    for (product_id, quantity) in &reservation.items {
        if let Some(item) = inventory.get_mut(product_id) {
            item.reserved -= *quantity;
            item.available += *quantity;
            item.updated_at = now;
            
            item.status = if item.available > 10 {
                InventoryStatus::InStock
            } else if item.available > 0 {
                InventoryStatus::LowStock
            } else {
                InventoryStatus::OutOfStock
            };
        }
    }
    
    info!("Released inventory reservation: {}", reservation_id);
    
    HttpResponse::Ok().json(ReservationResponse {
        reservation,
        success: true,
        message: Some("Reservation released successfully".to_string()),
    })
}

pub async fn health_check() -> impl Responder {
    HttpResponse::Ok().json(serde_json::json!({
        "status": "ok",
        "service": "inventory-service",
        "version": env!("CARGO_PKG_VERSION"),
        "timestamp": Utc::now().to_rfc3339()
    }))
}
