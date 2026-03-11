import Link from "next/link";
import { RiskPanel } from "@/components/risk-panel";
import { Shell } from "@/components/shell";
import { TaxCalculator } from "@/components/tax-calculator";
import { getHomeData, historicalCoverageStart } from "@/lib/api";

export default async function Home() {
  const { auctions, posts, advisors, courses, source } = await getHomeData();

  return (
    <Shell>
      <section className="hero-grid">
        <div className="glass fade-up rounded-[2rem] p-8">
          <p className="label">Customs Auction Platform</p>
          <h1 className="mt-4 max-w-3xl font-display text-5xl font-semibold leading-tight tracking-tight">
            一個把海關標售公告、出價教學、推播提醒與顧問媒合同步收進來的 Web PWA。
          </h1>
          <p className="mt-5 max-w-2xl text-lg leading-8 text-[color:var(--muted)]">
            v1 聚焦四大海關公告整合、關鍵字推播、歷史行情、稅費試算、社群看貨心得與代標需求表單。
          </p>
          <div className="mt-8 flex flex-wrap gap-3">
            <Link className="rounded-full bg-[color:var(--accent)] px-5 py-3 font-medium text-white" href="/auctions">
              查看最新標售
            </Link>
            <Link className="rounded-full border border-black/10 px-5 py-3 font-medium" href="/member">
              設定推播與日曆
            </Link>
          </div>
          <div className="mt-10 grid gap-4 sm:grid-cols-3">
            <div className="rounded-[1.5rem] bg-white/65 p-4">
              <p className="text-3xl font-semibold">{auctions.length}</p>
              <p className="mt-1 text-sm text-[color:var(--muted)]">首頁載入標案數</p>
            </div>
            <div className="rounded-[1.5rem] bg-white/65 p-4">
              <p className="text-3xl font-semibold">30 分</p>
              <p className="mt-1 text-sm text-[color:var(--muted)]">預設輪詢間隔</p>
            </div>
            <div className="rounded-[1.5rem] bg-white/65 p-4">
              <p className="text-3xl font-semibold">{source === "api" ? "API" : "Seed"}</p>
              <p className="mt-1 text-sm text-[color:var(--muted)]">目前資料來源</p>
            </div>
          </div>
        </div>
        <RiskPanel />
      </section>

      <section className="mt-6 grid gap-6 lg:grid-cols-[1.2fr_0.8fr]">
        <div className="glass rounded-[2rem] p-6">
          <div className="flex items-end justify-between">
            <div>
              <p className="label">最新標售</p>
              <h2 className="mt-3 font-display text-3xl font-semibold">高價值批次與管制品項</h2>
            </div>
            <Link className="text-sm text-[color:var(--accent)]" href="/auctions">
              進入完整清單
            </Link>
          </div>
          <div className="mt-6 space-y-4">
            {auctions.map((auction) => (
              <Link
                key={auction.id}
                href={`/auctions/${auction.id}`}
                className="block rounded-[1.5rem] border border-black/8 bg-white/60 p-5 transition hover:-translate-y-0.5"
              >
                <div className="flex flex-wrap items-center justify-between gap-3">
                  <div>
                    <p className="text-sm text-[color:var(--muted)]">{auction.office} / {auction.category}</p>
                    <h3 className="mt-1 text-xl font-semibold">{auction.title}</h3>
                  </div>
                  <span className="rounded-full bg-[#f0dfc8] px-3 py-1 text-sm text-[#6d4c1a]">{auction.priceBand}</span>
                </div>
                <p className="mt-3 text-sm leading-7 text-[color:var(--muted)]">{auction.summary}</p>
              </Link>
            ))}
          </div>
        </div>
        <TaxCalculator />
      </section>

      <section className="mt-6 grid gap-6 md:grid-cols-2">
        <div className="glass rounded-[2rem] p-6">
          <p className="label">看貨心得</p>
          <div className="mt-4 space-y-4">
              {posts.map((post) => (
                <div key={post.id} className="rounded-[1.5rem] border border-black/8 bg-white/60 p-5">
                <p className="text-sm text-[color:var(--muted)]">{post.office} / {post.author}</p>
                <h3 className="mt-2 text-xl font-semibold">{post.title}</h3>
                <p className="mt-2 text-sm leading-7 text-[color:var(--muted)]">{post.body}</p>
              </div>
            ))}
          </div>
        </div>
        <div className="grid gap-6">
          <div className="glass rounded-[2rem] p-6">
            <p className="label">專業顧問</p>
            <div className="mt-4 space-y-4">
              {advisors.map((advisor) => (
                <div key={advisor.id} className="rounded-[1.5rem] border border-black/8 bg-white/60 p-5">
                  <div className="flex items-center justify-between gap-4">
                    <div>
                      <h3 className="text-xl font-semibold">{advisor.name}</h3>
                      <p className="text-sm text-[color:var(--muted)]">{advisor.specialty}</p>
                    </div>
                    <span className="rounded-full bg-[#163a30] px-3 py-1 text-sm text-white">{advisor.responseTime}</span>
                  </div>
                  <p className="mt-3 text-sm leading-7 text-[color:var(--muted)]">{advisor.note}</p>
                </div>
              ))}
            </div>
          </div>
          <div className="glass rounded-[2rem] p-6">
            <p className="label">付費內容</p>
            <div className="mt-4 space-y-4">
              {courses.map((course) => (
                <div key={course.id} className="rounded-[1.5rem] border border-black/8 bg-white/60 p-5">
                  <h3 className="text-lg font-semibold">{course.title}</h3>
                  <p className="mt-2 text-sm leading-7 text-[color:var(--muted)]">{course.summary}</p>
                </div>
              ))}
            </div>
            <p className="mt-4 text-sm text-[color:var(--muted)]">歷史成交資料起算日：{historicalCoverageStart}</p>
          </div>
        </div>
      </section>
    </Shell>
  );
}
