import { Shell } from "@/components/shell";

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

export default function AdminPage() {
  return (
    <Shell>
      <section className="glass rounded-[2rem] p-8">
        <p className="label">Admin Backoffice</p>
        <h1 className="mt-4 font-display text-4xl font-semibold">營運與審核總覽</h1>
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
      </section>
    </Shell>
  );
}
