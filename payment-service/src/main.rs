mod config;
mod models;
mod handlers;

use actix_web::{web, App, HttpServer, middleware};
use std::sync::Mutex;
use std::collections::HashMap;
use log::info;

use crate::config::Config;
use crate::handlers::{AppState, create_payment, get_payment, get_payments_by_order, refund_payment, health_check};

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let config = Config::from_env();
    
    env_logger::init_from_env(env_logger::Env::new().default_filter_or(&config.log_level));
    
    info!("Starting payment service on {}", config.server_address());
    
    let app_state = web::Data::new(AppState {
        payments: Mutex::new(HashMap::new()),
    });
    
    HttpServer::new(move || {
        App::new()
            .app_data(app_state.clone())
            .wrap(middleware::Logger::default())
            .service(
                web::scope("/api/payments")
                    .route("", web::post().to(create_payment))
                    .route("/{id}", web::get().to(get_payment))
                    .route("/{id}/refund", web::post().to(refund_payment))
                    .route("/order/{order_id}", web::get().to(get_payments_by_order))
            )
            .route("/health", web::get().to(health_check))
    })
    .bind(config.server_address())?
    .run()
    .await
}
