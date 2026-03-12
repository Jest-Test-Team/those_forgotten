use tonic::{Request, Response, Status, transport::Server};

pub mod platform {
    tonic::include_proto!("customs.platform.v1");
}

use platform::policy_service_server::{PolicyService, PolicyServiceServer};
use platform::{EvaluatePolicyRequest, EvaluatePolicyResponse};

#[derive(Default)]
struct PolicyGrpc;

#[tonic::async_trait]
impl PolicyService for PolicyGrpc {
    async fn evaluate_policy(
        &self,
        request: Request<EvaluatePolicyRequest>,
    ) -> Result<Response<EvaluatePolicyResponse>, Status> {
        let payload = request.into_inner();
        let Some(lot) = payload.lot else {
            return Err(Status::invalid_argument("lot is required"));
        };

        let haystack =
            format!("{} {} {}", lot.title, lot.category, lot.warnings.join(" ")).to_lowercase();

        let mut policy_labels = Vec::new();
        for (needle, label) in [
            ("菸", "regulated-tobacco"),
            ("酒", "regulated-alcohol"),
            ("醫療", "regulated-medical"),
            ("進口車", "regulated-vehicle"),
        ] {
            if haystack.contains(needle) {
                policy_labels.push(label.to_string());
            }
        }

        Ok(Response::new(EvaluatePolicyResponse {
            restricted: !policy_labels.is_empty(),
            policy_labels,
        }))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = std::env::var("POLICY_GRPC_ADDR")
        .unwrap_or_else(|_| "0.0.0.0:50052".to_string())
        .parse()?;

    Server::builder()
        .add_service(PolicyServiceServer::new(PolicyGrpc))
        .serve(addr)
        .await?;

    Ok(())
}
