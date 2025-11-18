use std::env;
use dotenv::dotenv;

#[derive(Debug, Clone)]
pub struct Config {
    pub server_host: String,
    pub server_port: u16,
    pub log_level: String,
    pub payment_gateway_url: String,
    pub payment_gateway_api_key: String,
}

impl Config {
    pub fn from_env() -> Self {
        dotenv().ok();

        let server_host = env::var("SERVER_HOST").unwrap_or_else(|_| "0.0.0.0".to_string());
        let server_port = env::var("SERVER_PORT")
            .unwrap_or_else(|_| "8085".to_string())
            .parse::<u16>()
            .expect("SERVER_PORT must be a valid port number");
        let log_level = env::var("LOG_LEVEL").unwrap_or_else(|_| "info".to_string());
        let payment_gateway_url = env::var("PAYMENT_GATEWAY_URL")
            .unwrap_or_else(|_| "https://api.example.com/payments".to_string());
        let payment_gateway_api_key = env::var("PAYMENT_GATEWAY_API_KEY")
            .unwrap_or_else(|_| "test_api_key".to_string());

        Self {
            server_host,
            server_port,
            log_level,
            payment_gateway_url,
            payment_gateway_api_key,
        }
    }

    pub fn server_address(&self) -> String {
        format!("{}:{}", self.server_host, self.server_port)
    }
}
