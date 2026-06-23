use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SearchRequest {
    pub query: String,
    pub index_name: Option<String>,
    pub tenant_id: Option<String>,
    pub page: Option<u32>,
    pub page_size: Option<u32>,
    pub filters: Option<Vec<String>>,
    pub sort_by: Option<String>,
    pub sort_order: Option<String>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct VectorSearchRequest {
    pub query_vector: Vec<f32>,
    pub index_name: Option<String>,
    pub tenant_id: Option<String>,
    pub top_k: Option<usize>,
    pub min_score: Option<f64>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct HybridSearchRequest {
    pub query: String,
    pub query_vector: Vec<f32>,
    pub index_name: Option<String>,
    pub tenant_id: Option<String>,
    pub top_k: Option<usize>,
    pub alpha: Option<f64>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SearchResult {
    pub document_id: String,
    pub score: f64,
    pub fields: HashMap<String, String>,
    pub highlights: Vec<String>,
    pub snippet: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct SearchResponse {
    pub results: Vec<SearchResult>,
    pub total_hits: u64,
    pub page: u32,
    pub page_size: u32,
    pub took_ms: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct HealthResponse {
    pub status: String,
    pub service: String,
    pub version: String,
    pub qdrant: bool,
    pub elasticsearch: bool,
}
