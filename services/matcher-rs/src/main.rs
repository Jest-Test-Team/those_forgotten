use tonic::{Request, Response, Status, transport::Server};

pub mod platform {
    tonic::include_proto!("customs.platform.v1");
}

use platform::matcher_service_server::{MatcherService, MatcherServiceServer};
use platform::{MatchAuctionRequest, MatchAuctionResponse};

#[derive(Default)]
struct MatcherGrpc;

#[tonic::async_trait]
impl MatcherService for MatcherGrpc {
    async fn match_auction(
        &self,
        request: Request<MatchAuctionRequest>,
    ) -> Result<Response<MatchAuctionResponse>, Status> {
        let payload = request.into_inner();
        let Some(lot) = payload.lot else {
            return Err(Status::invalid_argument("lot is required"));
        };

        let haystack = format!(
            "{} {} {} {}",
            lot.title,
            lot.category,
            lot.office,
            lot.warnings.join(" ")
        )
        .to_lowercase();

        let matched_keywords = payload
            .keywords
            .into_iter()
            .filter(|keyword| haystack.contains(&keyword.to_lowercase()))
            .collect::<Vec<_>>();

        let mut derived_tags = Vec::new();
        if haystack.contains("相機") || haystack.contains("鏡頭") {
            derived_tags.push("3c".to_string());
        }
        if haystack.contains("進口車") || haystack.contains("驗車") {
            derived_tags.push("vehicle".to_string());
        }

        Ok(Response::new(MatchAuctionResponse {
            matched_keywords,
            derived_tags,
        }))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = std::env::var("MATCHER_GRPC_ADDR")
        .unwrap_or_else(|_| "0.0.0.0:50051".to_string())
        .parse()?;

    Server::builder()
        .add_service(MatcherServiceServer::new(MatcherGrpc))
        .serve(addr)
        .await?;

    Ok(())
}
