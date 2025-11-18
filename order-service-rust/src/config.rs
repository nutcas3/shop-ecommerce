use std::env;
use dotenv::dotenv;

#[derive(Debug, Clone)]
pub struct Config {
    pub server_host: String,
    pub server_port: u16,
    pub log_level: String,
    pub product_service_url: String,
    pub inventory_service_url: String,
    pub payment_service_url: String,
}

impl Config {
    pub fn from_env() -> Self {
        dotenv().ok();

        let server_host = env::var("SERVER_HOST").unwrap_or_else(|_| "0.0.0.0".to_string());
        let server_port = env::var("SERVER_PORT")
            .unwrap_or_else(|_| "8087".to_string())
            .parse::<u16>()
            .expect("SERVER_PORT must be a valid port number");
        let log_level = env::var("LOG_LEVEL").unwrap_or_else(|_| "info".to_string());
        let product_service_url = env::var("PRODUCT_SERVICE_URL")
            .unwrap_or_else(|_| "http://product-service:8082".to_string());
        let inventory_service_url = env::var("INVENTORY_SERVICE_URL")
            .unwrap_or_else(|_| "http://inventory-service:8086".to_string());
        let payment_service_url = env::var("PAYMENT_SERVICE_URL")
            .unwrap_or_else(|_| "http://payment-service:8085".to_string());

        Self {
            server_host,
            server_port,
            log_level,
            product_service_url,
            inventory_service_url,
            payment_service_url,
        }
    }

    pub fn server_address(&self) -> String {
        format!("{}:{}", self.server_host, self.server_port)
    }
}
