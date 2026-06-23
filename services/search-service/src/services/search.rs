use anyhow::Result;
use std::collections::HashMap;
use std::sync::Arc;
use qdrant_client::qdrant::ScoredPoint;
use serde_json::Value;
use crate::config::Config;
use crate::db::qdrant::QdrantClient;
use crate::db::elastic::ElasticClient;
use crate::models::{
    SearchRequest, VectorSearchRequest, HybridSearchRequest,
    SearchResult, SearchResponse,
};

pub struct SearchService {
    qdrant: QdrantClient,
    elastic: ElasticClient,
    cfg: Config,
}

impl SearchService {
    pub fn new(qdrant: QdrantClient, elastic: ElasticClient, cfg: &Config) -> Self {
        Self {
            qdrant,
            elastic,
            cfg: cfg.clone(),
        }
    }

    pub async fn fulltext_search(&self, req: &SearchRequest) -> Result<SearchResponse> {
        let page = req.page.unwrap_or(1);
        let page_size = req.page_size.unwrap_or(20);
        let from = (page - 1) * page_size;

        let start = std::time::Instant::now();
        let (hits, total, es_took) = self.elastic.search(&req.query, from, page_size).await?;
        let took_ms = start.elapsed().as_secs_f64() * 1000.0;

        let results = hits.iter().map(|hit| {
            let source = &hit["_source"];
            let highlight = hit["highlight"]["content"]
                .as_array()
                .and_then(|h| h.first())
                .and_then(|h| h.as_str())
                .unwrap_or("");

            SearchResult {
                document_id: source["document_id"].as_str().unwrap_or("").to_string(),
                score: hit["_score"].as_f64().unwrap_or(0.0),
                fields: HashMap::new(),
                highlights: vec![highlight.to_string()],
                snippet: source["content"].as_str().unwrap_or("").chars().take(200).collect(),
            }
        }).collect();

        Ok(SearchResponse {
            results,
            total_hits: total,
            page,
            page_size,
            took_ms,
        })
    }

    pub async fn vector_search(&self, req: &VectorSearchRequest) -> Result<SearchResponse> {
        let top_k = req.top_k.unwrap_or(self.cfg.top_k_default);

        let start = std::time::Instant::now();
        let points = self.qdrant.search(
            req.query_vector.clone(),
            top_k,
            req.min_score.or(Some(self.cfg.min_score)),
        ).await?;
        let took_ms = start.elapsed().as_secs_f64() * 1000.0;

        let results = points.iter().map(|p| {
            let payload = p.payload.as_ref().map(|m| {
                m.iter().map(|(k, v)| {
                    (k.clone(), format!("{:?}", v))
                }).collect::<HashMap<_, _>>()
            }).unwrap_or_default();

            SearchResult {
                document_id: payload.get("document_id").cloned().unwrap_or_default(),
                score: p.score as f64,
                fields: payload,
                highlights: vec![],
                snippet: String::new(),
            }
        }).collect();

        Ok(SearchResponse {
            results,
            total_hits: results.len() as u64,
            page: 1,
            page_size: top_k as u32,
            took_ms,
        })
    }

    pub async fn hybrid_search(&self, req: &HybridSearchRequest) -> Result<SearchResponse> {
        let alpha = req.alpha.unwrap_or(0.5);
        let top_k = req.top_k.unwrap_or(self.cfg.top_k_default);

        let start = std::time::Instant::now();

        // Run both searches in parallel
        let (es_result, qdrant_result) = tokio::join!(
            self.elastic.search(&req.query, 0, top_k as u32),
            self.qdrant.search(req.query_vector.clone(), top_k, None)
        );

        let (es_hits, _, _) = es_result?;
        let qdrant_points = qdrant_result?;

        // Normalize and fuse results
        let mut fused: HashMap<String, (f64, f64)> = HashMap::new();

        for (i, hit) in es_hits.iter().enumerate() {
            let doc_id = hit["_source"]["document_id"].as_str().unwrap_or("");
            let score = (1.0 - (i as f64 / es_hits.len() as f64)) * (1.0 - alpha);
            fused.entry(doc_id.to_string())
                .and_modify(|(s, _)| *s += score)
                .or_insert((score, 0.0));
        }

        for (i, point) in qdrant_points.iter().enumerate() {
            let doc_id = point.payload.as_ref()
                .and_then(|m| m.get("document_id"))
                .map(|v| format!("{:?}", v))
                .unwrap_or_default();
            let score = (1.0 - (i as f64 / qdrant_points.len() as f64)) * alpha;
            fused.entry(doc_id)
                .and_modify(|(s, _)| *s += score)
                .or_insert((score, 0.0));
        }

        let mut results: Vec<SearchResult> = fused.into_iter()
            .map(|(doc_id, (score, _))| SearchResult {
                document_id: doc_id,
                score,
                fields: HashMap::new(),
                highlights: vec![],
                snippet: String::new(),
            })
            .collect();

        results.sort_by(|a, b| b.score.partial_cmp(&a.score).unwrap_or(std::cmp::Ordering::Equal));
        results.truncate(top_k);

        let took_ms = start.elapsed().as_secs_f64() * 1000.0;

        Ok(SearchResponse {
            total_hits: results.len() as u64,
            results,
            page: 1,
            page_size: top_k as u32,
            took_ms,
        })
    }
}
