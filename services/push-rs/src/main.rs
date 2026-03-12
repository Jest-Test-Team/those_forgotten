use tonic::{Request, Response, Status, transport::Server};

pub mod platform {
    tonic::include_proto!("customs.platform.v1");
}

use platform::push_service_server::{PushService, PushServiceServer};
use platform::{PrepareNotificationRequest, PrepareNotificationResponse};

#[derive(Default)]
struct PushGrpc;

#[tonic::async_trait]
impl PushService for PushGrpc {
    async fn prepare_notification(
        &self,
        request: Request<PrepareNotificationRequest>,
    ) -> Result<Response<PrepareNotificationResponse>, Status> {
        let payload = request.into_inner();
        let Some(lot) = payload.lot else {
            return Err(Status::invalid_argument("lot is required"));
        };

        Ok(Response::new(PrepareNotificationResponse {
            title: format!("{} 新標售通知", lot.office),
            body: format!("{} 命中關鍵字「{}」", lot.title, payload.keyword),
            click_url: lot.official_url,
        }))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = std::env::var("PUSH_GRPC_ADDR")
        .unwrap_or_else(|_| "0.0.0.0:50053".to_string())
        .parse()?;

    Server::builder()
        .add_service(PushServiceServer::new(PushGrpc))
        .serve(addr)
        .await?;

    Ok(())
}
