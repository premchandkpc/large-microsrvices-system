use actix_web::{test, App, web};
use search_service::config::Config;
use search_service::handlers;

#[actix_web::test]
async fn test_health_endpoint() {
    let cfg = Config {
        port: 8086,
        environment: "test".to_string(),
        qdrant_host: "localhost".to_string(),
        qdrant_port: 6333,
        elasticsearch_url: "http://localhost:9200".to_string(),
        kafka_brokers: vec![],
        redis_addr: "localhost:6379".to_string(),
        embedding_dim: 1536,
        top_k: 10,
        min_score: 0.7,
    };

    let app = test::init_service(
        App::new()
            .app_data(web::Data::new(cfg))
            .route("/health", web::get().to(handlers::health))
    ).await;

    let req = test::TestRequest::get().uri("/health").to_request();
    let resp = test::call_service(&app, req).await;
    assert!(resp.status().is_success());
}
