use anyhow::Result;
use qdrant_client::client::QdrantClient as QdrantRawClient;
use qdrant_client::prelude::*;
use qdrant_client::qdrant::{
    SearchPoints, ScoredPoint, CreateCollectionBuilder, Distance, VectorParamsBuilder,
};
use crate::config::Config;

pub struct QdrantClient {
    client: QdrantRawClient,
    collection: String,
}

impl QdrantClient {
    pub async fn new(cfg: &Config) -> Result<Self> {
        let client = QdrantRawClient::from_url(&format!(
            "http://{}:{}", cfg.qdrant_host, cfg.qdrant_port
        ))
        .build()?;

        // Ensure collection exists
        let collections = client.list_collections().await?;
        let exists = collections.collections.iter()
            .any(|c| c.name == cfg.qdrant_collection);

        if !exists {
            client.create_collection(
                CreateCollectionBuilder::new(cfg.qdrant_collection.clone())
                    .vectors_config(VectorParamsBuilder::new(
                        cfg.embedding_dim as u64,
                        Distance::Cosine,
                    ))
            ).await?;
            tracing::info!("Created Qdrant collection: {}", cfg.qdrant_collection);
        }

        Ok(Self {
            client,
            collection: cfg.qdrant_collection.clone(),
        })
    }

    pub async fn search(
        &self,
        vector: Vec<f32>,
        top_k: usize,
        min_score: Option<f64>,
    ) -> Result<Vec<ScoredPoint>> {
        let results = self.client
            .search_points(&SearchPoints {
                collection_name: self.collection.clone(),
                vector,
                limit: top_k as u64,
                score_threshold: min_score.unwrap_or(0.0),
                with_payload: Some(true.into()),
                ..Default::default()
            })
            .await?;

        Ok(results.result)
    }

    pub async fn health(&self) -> bool {
        self.client.list_collections().await.is_ok()
    }
}
