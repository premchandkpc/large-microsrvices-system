use actix_web::{web, HttpResponse};
use std::sync::Arc;
use crate::models::{SearchRequest, VectorSearchRequest, HybridSearchRequest, SearchResponse, HealthResponse};
use crate::services::search::SearchService;

pub async fn health() -> HttpResponse {
    HttpResponse::Ok().json(HealthResponse {
        status: "UP".into(),
        service: "search-service".into(),
        version: "1.0.0".into(),
        qdrant: true,
        elasticsearch: true,
    })
}

pub async fn text_search(
    svc: web::Data<Arc<SearchService>>,
    body: web::Json<SearchRequest>,
) -> HttpResponse {
    match svc.fulltext_search(&body).await {
        Ok(results) => HttpResponse::Ok().json(results),
        Err(e) => {
            tracing::error!("Text search failed: {}", e);
            HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "search failed",
                "detail": e.to_string()
            }))
        }
    }
}

pub async fn vector_search(
    svc: web::Data<Arc<SearchService>>,
    body: web::Json<VectorSearchRequest>,
) -> HttpResponse {
    match svc.vector_search(&body).await {
        Ok(results) => HttpResponse::Ok().json(results),
        Err(e) => {
            tracing::error!("Vector search failed: {}", e);
            HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "vector search failed"
            }))
        }
    }
}

pub async fn hybrid_search(
    svc: web::Data<Arc<SearchService>>,
    body: web::Json<HybridSearchRequest>,
) -> HttpResponse {
    match svc.hybrid_search(&body).await {
        Ok(results) => HttpResponse::Ok().json(results),
        Err(e) => {
            tracing::error!("Hybrid search failed: {}", e);
            HttpResponse::InternalServerError().json(serde_json::json!({
                "error": "hybrid search failed"
            }))
        }
    }
}
