use anyhow::Result;
use elasticsearch::{
    Elasticsearch, SearchParts,
    http::transport::Transport,
};
use serde_json::{json, Value};
use crate::config::Config;

pub struct ElasticClient {
    client: Elasticsearch,
    index: String,
}

impl ElasticClient {
    pub async fn new(cfg: &Config) -> Result<Self> {
        let transport = Transport::single_node(&cfg.elasticsearch_url)?;
        let client = Elasticsearch::new(transport);

        // Ensure index exists
        let exists = client.indices()
            .exists(elasticsearch::indices::IndicesExistsParts::Index(&[&cfg.elasticsearch_index]))
            .send().await?
            .status_code() == 200;

        if !exists {
            client.indices()
                .create(elasticsearch::indices::IndicesCreateParts::Index(&cfg.elasticsearch_index))
                .body(json!({
                    "settings": {
                        "number_of_shards": 2,
                        "number_of_replicas": 1
                    },
                    "mappings": {
                        "properties": {
                            "document_id": {"type": "keyword"},
                            "user_id": {"type": "keyword"},
                            "content": {"type": "text"},
                            "created_at": {"type": "date"}
                        }
                    }
                }))
                .send().await?;
            tracing::info!("Created ES index: {}", cfg.elasticsearch_index);
        }

        Ok(Self {
            client,
            index: cfg.elasticsearch_index.clone(),
        })
    }

    pub async fn search(
        &self,
        query: &str,
        from: u32,
        size: u32,
    ) -> Result<(Vec<Value>, u64, f64)> {
        let response = self.client
            .search(SearchParts::Index(&[&self.index]))
            .body(json!({
                "query": {
                    "multi_match": {
                        "query": query,
                        "fields": ["content^2", "document_id"],
                        "fuzziness": "AUTO"
                    }
                },
                "from": from,
                "size": size,
                "highlight": {
                    "fields": {
                        "content": {}
                    }
                }
            }))
            .send().await?;

        let response_body: Value = response.json().await?;
        let took = response_body["took"].as_f64().unwrap_or(0.0);
        let total = response_body["hits"]["total"]["value"].as_u64().unwrap_or(0);
        let hits = response_body["hits"]["hits"]
            .as_array()
            .cloned()
            .unwrap_or_default();

        Ok((hits, total, took))
    }

    pub async fn health(&self) -> bool {
        self.client.ping().send().await.is_ok()
    }
}
