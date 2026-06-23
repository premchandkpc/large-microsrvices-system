use std::sync::Arc;
use actix_web::{web, App, HttpServer, middleware::Logger};
use tracing_subscriber::EnvFilter;
use search_service::config::Config;
use search_service::api::handlers;
use search_service::services::search::SearchService;
use search_service::db::qdrant::QdrantClient;
use search_service::db::elastic::ElasticClient;

#[actix_web::main]
async fn main() -> anyhow::Result<()> {
    // Initialize tracing
    tracing_subscriber::fmt()
        .with_env_filter(
            EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| EnvFilter::new("info"))
        )
        .json()
        .init();

    let cfg = Config::from_env()?;
    tracing::info!("Starting search-service on port {}", cfg.port);

    let qdrant = QdrantClient::new(&cfg).await?;
    let elastic = ElasticClient::new(&cfg).await?;
    let search_svc = Arc::new(SearchService::new(qdrant, elastic, &cfg));

    let bind_addr = format!("0.0.0.0:{}", cfg.port);

    HttpServer::new(move || {
        App::new()
            .wrap(Logger::default())
            .app_data(web::Data::from(search_svc.clone()))
            .service(
                web::scope("/api/v1/search")
                    .route("/text", web::post().to(handlers::text_search))
                    .route("/vector", web::post().to(handlers::vector_search))
                    .route("/hybrid", web::post().to(handlers::hybrid_search))
            )
            .route("/health", web::get().to(handlers::health))
    })
    .bind(&bind_addr)?
    .run()
    .await?;

    Ok(())
}
