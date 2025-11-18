mod config;
mod models;
mod handlers;

use actix_web::{web, App, HttpServer, middleware};
use std::sync::Mutex;
use std::collections::HashMap;
use log::info;
use reqwest::Client;

use crate::config::Config;
use crate::handlers::{
    AppState, create_order, get_order, get_user_orders, 
    update_order_status, process_payment, health_check
};

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let config = Config::from_env();
    
    env_logger::init_from_env(env_logger::Env::new().default_filter_or(&config.log_level));
    
    info!("Starting order service on {}", config.server_address());
    
    let http_client = Client::new();
    
    let app_state = web::Data::new(AppState {
        orders: Mutex::new(HashMap::new()),
        http_client,
        config: config.clone(),
    });
    
    HttpServer::new(move || {
        App::new()
            .app_data(app_state.clone())
            .wrap(middleware::Logger::default())
            .service(
                web::scope("/api/orders")
                    .route("", web::post().to(create_order))
                    .route("/{id}", web::get().to(get_order))
                    .route("/{id}/status", web::put().to(update_order_status))
                    .route("/{id}/payment", web::post().to(process_payment))
                    .route("/user/{user_id}", web::get().to(get_user_orders))
            )
            .route("/health", web::get().to(health_check))
    })
    .bind(config.server_address())?
    .run()
    .await
}
