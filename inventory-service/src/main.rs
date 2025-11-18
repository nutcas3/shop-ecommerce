mod config;
mod models;
mod handlers;

use actix_web::{web, App, HttpServer, middleware};
use std::sync::Mutex;
use std::collections::HashMap;
use log::info;

use crate::config::Config;
use crate::handlers::{
    AppState, get_inventory_item, get_all_inventory, create_inventory_item, 
    update_inventory_item, reserve_inventory, confirm_reservation, 
    release_reservation, health_check
};

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let config = Config::from_env();
    
    env_logger::init_from_env(env_logger::Env::new().default_filter_or(&config.log_level));
    
    info!("Starting inventory service on {}", config.server_address());
    
    let app_state = web::Data::new(AppState {
        inventory: Mutex::new(HashMap::new()),
        reservations: Mutex::new(HashMap::new()),
    });
    
    HttpServer::new(move || {
        App::new()
            .app_data(app_state.clone())
            .wrap(middleware::Logger::default())
            .service(
                web::scope("/api/inventory")
                    .route("", web::get().to(get_all_inventory))
                    .route("/{product_id}", web::get().to(get_inventory_item))
                    .route("/{product_id}", web::post().to(create_inventory_item))
                    .route("/{product_id}", web::put().to(update_inventory_item))
                    .route("/reserve", web::post().to(reserve_inventory))
                    .route("/reservation/{reservation_id}/confirm", web::post().to(confirm_reservation))
                    .route("/reservation/{reservation_id}/release", web::post().to(release_reservation))
            )
            .route("/health", web::get().to(health_check))
    })
    .bind(config.server_address())?
    .run()
    .await
}
