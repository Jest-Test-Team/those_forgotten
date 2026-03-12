use tonic::{Request, Response, Status, transport::Server};

pub mod platform {
    tonic::include_proto!("customs.platform.v1");
}

use platform::feed_service_server::{FeedService, FeedServiceServer};
use platform::{RenderCalendarRequest, RenderCalendarResponse};

#[derive(Default)]
struct FeedGrpc;

#[tonic::async_trait]
impl FeedService for FeedGrpc {
    async fn render_calendar(
        &self,
        request: Request<RenderCalendarRequest>,
    ) -> Result<Response<RenderCalendarResponse>, Status> {
        let payload = request.into_inner();
        if payload.token.trim().is_empty() {
            return Err(Status::invalid_argument("token is required"));
        }

        let mut lines = vec![
            "BEGIN:VCALENDAR".to_string(),
            "VERSION:2.0".to_string(),
            "PRODID:-//Those Forgotten//Feed Service//EN".to_string(),
        ];

        for lot in payload.lots {
            lines.push("BEGIN:VEVENT".to_string());
            lines.push(format!("UID:{}@those-forgotten", lot.auction_lot_id));
            lines.push(format!("SUMMARY:{} {}", lot.office, lot.title));
            lines.push(format!("DESCRIPTION:官方連結 {}", lot.official_url));
            lines.push("END:VEVENT".to_string());
        }

        lines.push("END:VCALENDAR".to_string());

        Ok(Response::new(RenderCalendarResponse {
            ics: lines.join("\r\n"),
        }))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = std::env::var("FEED_GRPC_ADDR")
        .unwrap_or_else(|_| "0.0.0.0:50054".to_string())
        .parse()?;

    Server::builder()
        .add_service(FeedServiceServer::new(FeedGrpc))
        .serve(addr)
        .await?;

    Ok(())
}
