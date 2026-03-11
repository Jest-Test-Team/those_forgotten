import Link from "next/link";
import { Shell } from "@/components/shell";
import { auctions } from "@/lib/site-data";

export const metadata = {
  title: "標售清單",
};

export default function AuctionsPage() {
  return (
    <Shell>
      <section className="glass rounded-[2rem] p-8">
        <p className="label">Auction Feed</p>
        <h1 className="mt-4 font-display text-4xl font-semibold">全台海關標售清單</h1>
        <p className="mt-4 max-w-2xl text-[color:var(--muted)]">
          列表頁保留官方連結、看貨日期、截標時間與風險標籤，供會員設定關鍵字推播與 ICS 日曆訂閱。
        </p>
        <div className="mt-8 space-y-4">
          {auctions.map((auction) => (
            <article key={auction.id} className="rounded-[1.5rem] border border-black/8 bg-white/60 p-5">
              <div className="flex flex-wrap items-center justify-between gap-3">
                <div>
                  <p className="text-sm text-[color:var(--muted)]">{auction.office} / {auction.category}</p>
                  <h2 className="mt-1 text-2xl font-semibold">{auction.title}</h2>
                </div>
                <span className="rounded-full bg-[#ecdac1] px-3 py-1 text-sm text-[#7b551d]">{auction.priceBand}</span>
              </div>
              <div className="mt-4 flex flex-wrap gap-2">
                {auction.warnings.map((warning) => (
                  <span key={warning} className="rounded-full bg-[#17342d] px-3 py-1 text-sm text-white">
                    {warning}
                  </span>
                ))}
              </div>
              <div className="mt-4 flex flex-wrap gap-3 text-sm text-[color:var(--muted)]">
                <p>看貨：{auction.viewingAt}</p>
                <p>截標：{auction.closingAt}</p>
              </div>
              <div className="mt-5 flex flex-wrap gap-3">
                <Link className="rounded-full bg-[color:var(--accent)] px-4 py-2 text-sm font-medium text-white" href={`/auctions/${auction.id}`}>
                  查看詳情
                </Link>
                <a className="rounded-full border border-black/10 px-4 py-2 text-sm font-medium" href={auction.officialUrl} target="_blank" rel="noreferrer">
                  官方公告
                </a>
              </div>
            </article>
          ))}
        </div>
      </section>
    </Shell>
  );
}
