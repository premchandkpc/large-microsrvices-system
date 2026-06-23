use config::{ConfigError, Environment};
use serde::Deserialize;

#[derive(Debug, Deserialize, Clone)]
pub struct Config {
    pub port: u16,

    pub qdrant_host: String,
    pub qdrant_port: u16,
    pub qdrant_collection: String,

    pub elasticsearch_url: String,
    pub elasticsearch_index: String,

    pub redis_addr: String,
    pub redis_password: Option<String>,

    pub kafka_brokers: String,

    pub embedding_dim: usize,
    pub top_k_default: usize,
    pub min_score: f64,
}

impl Config {
    pub fn from_env() -> Result<Self, ConfigError> {
        let mut cfg = config::Config::builder()
            .set_default("port", 8086)?
            .set_default("qdrant_host", "localhost")?
            .set_default("qdrant_port", 6333)?
            .set_default("qdrant_collection", "documents")?
            .set_default("elasticsearch_url", "http://localhost:9200")?
            .set_default("elasticsearch_index", "documents")?
            .set_default("redis_addr", "localhost:6379")?
            .set_default("redis_password", None::<String>)?
            .set_default("kafka_brokers", "localhost:9092")?
            .set_default("embedding_dim", 1536)?
            .set_default("top_k_default", 10)?
            .set_default("min_score", 0.7)?
            .add_source(Environment::with_prefix("SS"));

        cfg.build()?.try_deserialize()
    }
}
