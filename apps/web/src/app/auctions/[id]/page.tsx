import { notFound } from "next/navigation";
import { RiskPanel } from "@/components/risk-panel";
import { Shell } from "@/components/shell";
import { getAuctionById, getAuctionHistory, historicalCoverageStart } from "@/lib/api";

type Props = {
  params: Promise<{ id: string }>;
};

export async function generateMetadata({ params }: Props) {
  const { id } = await params;
  const { auction } = await getAuctionById(id);

  if (!auction) {
    return { title: "標售不存在" };
  }

  return {
    title: auction.title,
    description: auction.summary,
  };
}

export default async function AuctionDetailPage({ params }: Props) {
  const { id } = await params;
  const [{ auction, source }, { rows: historyRows }] = await Promise.all([getAuctionById(id), getAuctionHistory(id)]);

  if (!auction) {
    notFound();
  }

  return (
    <Shell>
      <section className="grid gap-6 lg:grid-cols-[1.2fr_0.8fr]">
        <div className="glass rounded-[2rem] p-8">
          <p className="label">{auction.office}</p>
          <h1 className="mt-4 font-display text-4xl font-semibold">{auction.title}</h1>
          <p className="mt-4 max-w-2xl text-lg leading-8 text-[color:var(--muted)]">{auction.summary}</p>
          <div className="mt-6 flex flex-wrap gap-2">
            {auction.warnings.map((warning) => (
              <span key={warning} className="rounded-full bg-[#17342d] px-3 py-1 text-sm text-white">
                {warning}
              </span>
            ))}
          </div>
          <div className="mt-8 grid gap-4 sm:grid-cols-3">
            <div className="rounded-[1.5rem] bg-white/65 p-4">
              <p className="text-sm text-[color:var(--muted)]">看貨日期</p>
              <p className="mt-2 text-xl font-semibold">{auction.viewingAt}</p>
            </div>
            <div className="rounded-[1.5rem] bg-white/65 p-4">
              <p className="text-sm text-[color:var(--muted)]">投標截止</p>
              <p className="mt-2 text-xl font-semibold">{auction.closingAt}</p>
            </div>
            <div className="rounded-[1.5rem] bg-white/65 p-4">
              <p className="text-sm text-[color:var(--muted)]">參考區間</p>
              <p className="mt-2 text-xl font-semibold">{auction.priceBand}</p>
            </div>
          </div>
          <div className="mt-8 rounded-[1.5rem] border border-[#cb6f17]/30 bg-[#fff3e3] p-5 text-sm leading-7 text-[#7c4310]">
            <p className="font-semibold">官方標售提醒</p>
            <p className="mt-2">請以官方公告、附件與現場看貨為準。平台僅提供資訊整合與風險提示，不保證貨況與得標結果。</p>
          </div>
          <div className="mt-6 flex flex-wrap gap-3">
            <a className="rounded-full bg-[color:var(--accent)] px-5 py-3 text-sm font-medium text-white" href={auction.officialUrl} target="_blank" rel="noreferrer">
              前往官方公告
            </a>
            <a className="rounded-full border border-black/10 px-5 py-3 text-sm font-medium" href="/member">
              加入個人 ICS 日曆
            </a>
          </div>
          <div className="mt-8 rounded-[1.5rem] border border-black/8 bg-white/55 p-5">
            <p className="label">歷史成交資料</p>
            <p className="mt-3 text-sm leading-7 text-[color:var(--muted)]">
              成交資料庫目前僅涵蓋 {historicalCoverageStart} 之後的平台蒐集結果。v1 不回補既有歷史資料。
            </p>
            <div className="mt-4 space-y-2">
              {historyRows.length > 0 ? (
                historyRows.map((row) => (
                  <div key={row.id} className="rounded-2xl border border-black/8 bg-white/70 px-4 py-3 text-sm text-[color:var(--muted)]">
                    成交價 NT$ {row.finalPrice.toLocaleString()} / {new Date(row.recordedAt).toLocaleDateString("zh-TW")}
                  </div>
                ))
              ) : (
                <div className="rounded-2xl border border-black/8 bg-white/70 px-4 py-3 text-sm text-[color:var(--muted)]">
                  尚無歷史成交資料
                </div>
              )}
            </div>
            <p className="mt-4 text-sm text-[color:var(--muted)]">目前資料來源：{source === "api" ? "即時 API" : "本地 seed fallback"}</p>
          </div>
        </div>
        <RiskPanel />
      </section>
    </Shell>
  );
}
