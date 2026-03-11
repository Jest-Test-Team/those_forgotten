import { MemberConsole } from "@/components/member-console";
import { Shell } from "@/components/shell";
import { getAuctions, getCourses, historicalCoverageStart } from "@/lib/api";

export const metadata = {
  title: "會員中心",
};

export default async function MemberPage() {
  const [{ auctions, source: auctionSource }, { courses, source: courseSource }] = await Promise.all([
    getAuctions(),
    getCourses(),
  ]);

  return (
    <Shell>
      <div className="grid gap-6 lg:grid-cols-[0.95fr_1.05fr]">
        <section className="glass rounded-[2rem] p-8">
          <p className="label">Member Console</p>
          <h1 className="mt-4 font-display text-4xl font-semibold">關鍵字推播、日曆整合與訂閱權限</h1>
          <div className="mt-6 space-y-4 text-sm leading-7 text-[color:var(--muted)]">
            <div className="rounded-[1.5rem] border border-black/8 bg-white/60 p-5">
              <p className="font-semibold text-[color:var(--foreground)]">登入方式</p>
              <p className="mt-2">Google OAuth 與 Email magic link。會員發文門檻為 Email 已驗證。</p>
            </div>
            <div className="rounded-[1.5rem] border border-black/8 bg-white/60 p-5">
              <p className="font-semibold text-[color:var(--foreground)]">個人日曆</p>
              <p className="mt-2">以簽名 ICS feed 提供持續同步，並搭配 Google / Apple / Outlook 一鍵加入。</p>
            </div>
            <div className="rounded-[1.5rem] border border-black/8 bg-white/60 p-5">
              <p className="font-semibold text-[color:var(--foreground)]">歷史資料涵蓋</p>
              <p className="mt-2">成交資料起算日：{historicalCoverageStart}</p>
            </div>
            <div className="rounded-[1.5rem] border border-black/8 bg-white/60 p-5">
              <p className="font-semibold text-[color:var(--foreground)]">資料來源</p>
              <p className="mt-2">標案：{auctionSource === "api" ? "API" : "Seed"} / 課程：{courseSource === "api" ? "API" : "Seed"}</p>
            </div>
          </div>
        </section>
        <section className="grid gap-6">
          <MemberConsole />
          <div className="glass rounded-[2rem] p-6">
            <p className="label">已追蹤標案</p>
            <div className="mt-4 space-y-3">
              {auctions.slice(0, 2).map((auction) => (
                <div key={auction.id} className="rounded-[1.5rem] border border-black/8 bg-white/60 p-4">
                  <p className="font-semibold">{auction.title}</p>
                  <p className="mt-1 text-sm text-[color:var(--muted)]">{auction.viewingAt} 看貨 / {auction.closingAt} 截標</p>
                </div>
              ))}
            </div>
          </div>
          <div className="glass rounded-[2rem] p-6">
            <p className="label">已購買內容</p>
            <div className="mt-4 space-y-3">
              {courses.map((course) => (
                <div key={course.id} className="rounded-[1.5rem] border border-black/8 bg-white/60 p-4">
                  <p className="font-semibold">{course.title}</p>
                  <p className="mt-1 text-sm text-[color:var(--muted)]">{course.summary}</p>
                </div>
              ))}
            </div>
          </div>
        </section>
      </div>
    </Shell>
  );
}
