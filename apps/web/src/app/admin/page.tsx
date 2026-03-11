import { Shell } from "@/components/shell";
import { getAdminDashboardData } from "@/lib/api";

const queues = [
  "爬蟲健康度與手動重抓",
  "社群檢舉與下架處理",
  "文章 / 課程發布",
  "顧問需求單與客服查詢",
  "廣告版位與排程管理",
];

export const metadata = {
  title: "管理後台",
};

export default async function AdminPage() {
  const dashboard = await getAdminDashboardData();

  return (
    <Shell>
      <section className="glass rounded-[2rem] p-8">
        <p className="label">Admin Backoffice</p>
        <h1 className="mt-4 font-display text-4xl font-semibold">營運與審核總覽</h1>
        <div className="mt-6 grid gap-4 md:grid-cols-5">
          <div className="rounded-[1.5rem] bg-white/65 p-4">
            <p className="text-sm text-[color:var(--muted)]">Live Auctions</p>
            <p className="mt-2 text-3xl font-semibold">{dashboard.liveAuctionCount}</p>
          </div>
          <div className="rounded-[1.5rem] bg-white/65 p-4">
            <p className="text-sm text-[color:var(--muted)]">Community Posts</p>
            <p className="mt-2 text-3xl font-semibold">{dashboard.communityPostCount}</p>
          </div>
          <div className="rounded-[1.5rem] bg-white/65 p-4">
            <p className="text-sm text-[color:var(--muted)]">Advisors</p>
            <p className="mt-2 text-3xl font-semibold">{dashboard.advisorCount}</p>
          </div>
          <div className="rounded-[1.5rem] bg-white/65 p-4">
            <p className="text-sm text-[color:var(--muted)]">Data Source</p>
            <p className="mt-2 text-3xl font-semibold">{dashboard.source === "api" ? "API" : "Seed"}</p>
          </div>
          <div className="rounded-[1.5rem] bg-white/65 p-4">
            <p className="text-sm text-[color:var(--muted)]">API Repository</p>
            <p className="mt-2 text-3xl font-semibold">{dashboard.backendHealth.repository}</p>
            <p className="mt-1 text-xs text-[color:var(--muted)]">{dashboard.backendHealth.status}</p>
          </div>
        </div>
        <div className="mt-8 grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {queues.map((item, index) => (
            <div key={item} className="rounded-[1.5rem] border border-black/8 bg-white/60 p-5">
              <p className="text-sm text-[color:var(--muted)]">Queue {index + 1}</p>
              <h2 className="mt-2 text-2xl font-semibold">{item}</h2>
              <p className="mt-3 text-sm leading-7 text-[color:var(--muted)]">
                需保留操作審計記錄，並在高風險品項頁面明確顯示法規與免責警示。
              </p>
            </div>
          ))}
        </div>
        <div className="mt-8 grid gap-4 lg:grid-cols-[1.3fr_0.7fr]">
          <section className="rounded-[1.5rem] border border-black/8 bg-white/70 p-6">
            <div className="flex items-center justify-between gap-4">
              <div>
                <p className="label">Moderation Queue</p>
                <h2 className="mt-2 text-2xl font-semibold">社群檢舉待處理清單</h2>
              </div>
              <p className="rounded-full bg-amber-100 px-3 py-1 text-sm font-medium text-amber-900">
                {dashboard.moderationReports.length} pending
              </p>
            </div>
            <div className="mt-5 space-y-3">
              {dashboard.moderationReports.map((report) => (
                <article key={report.id} className="rounded-[1.25rem] border border-black/8 bg-white p-4">
                  <div className="flex flex-wrap items-center justify-between gap-3">
                    <div>
                      <p className="text-sm text-[color:var(--muted)]">{report.office}</p>
                      <h3 className="mt-1 text-lg font-semibold">{report.postTitle}</h3>
                    </div>
                    <p className="rounded-full bg-rose-100 px-3 py-1 text-xs font-medium uppercase tracking-[0.2em] text-rose-900">
                      {report.status}
                    </p>
                  </div>
                  <p className="mt-3 text-sm leading-7 text-[color:var(--muted)]">{report.reason}</p>
                  <p className="mt-3 text-xs text-[color:var(--muted)]">
                    檢舉時間 {new Date(report.createdAt).toLocaleString("zh-TW", { hour12: false })}
                  </p>
                </article>
              ))}
            </div>
          </section>
          <div className="space-y-4">
            <section className="rounded-[1.5rem] border border-black/8 bg-white/70 p-6">
              <div className="flex items-center justify-between gap-4">
                <div>
                  <p className="label">Advisor Queue</p>
                  <h2 className="mt-2 text-2xl font-semibold">顧問需求單</h2>
                </div>
                <p className="rounded-full bg-sky-100 px-3 py-1 text-sm font-medium text-sky-900">
                  {dashboard.advisorLeads.length} open
                </p>
              </div>
              <div className="mt-5 space-y-3">
                {dashboard.advisorLeads.map((lead) => (
                  <article key={lead.id} className="rounded-[1.25rem] border border-black/8 bg-white p-4">
                    <p className="text-sm text-[color:var(--muted)]">{lead.category}</p>
                    <h3 className="mt-1 text-lg font-semibold">
                      {lead.name} {"->"} {lead.advisorName}
                    </h3>
                    <p className="mt-3 text-sm leading-7 text-[color:var(--muted)]">{lead.message}</p>
                    <p className="mt-3 text-xs text-[color:var(--muted)]">
                      {lead.email} | {new Date(lead.createdAt).toLocaleString("zh-TW", { hour12: false })}
                    </p>
                  </article>
                ))}
              </div>
            </section>
            <section className="rounded-[1.5rem] border border-black/8 bg-white/70 p-6">
              <p className="label">Operations Notes</p>
              <h2 className="mt-2 text-2xl font-semibold">審核規則</h2>
              <ul className="mt-4 space-y-3 text-sm leading-7 text-[color:var(--muted)]">
                <li>所有檢舉需保留操作審計記錄，避免誤刪與後續爭議。</li>
                <li>高風險品項頁面必須同步顯示官方連結、免責聲明與特殊資格標籤。</li>
                <li>若內容涉及實拍圖，優先確認是否有誤導標題或缺少看貨依據。</li>
              </ul>
            </section>
          </div>
        </div>
      </section>
    </Shell>
  );
}
