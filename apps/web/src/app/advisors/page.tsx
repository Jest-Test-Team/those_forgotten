import { Shell } from "@/components/shell";
import { advisors } from "@/lib/site-data";

export const metadata = {
  title: "顧問媒合",
};

export default function AdvisorsPage() {
  return (
    <Shell>
      <section className="glass rounded-[2rem] p-8">
        <p className="label">Advisor Directory</p>
        <h1 className="mt-4 font-display text-4xl font-semibold">專業代標與顧問媒合</h1>
        <p className="mt-4 max-w-2xl text-lg leading-8 text-[color:var(--muted)]">
          v1 僅提供顧問名錄與需求表單，不處理站內私訊與金流。顧問資料為自我聲明，請自行驗證專業背景。
        </p>
        <div className="mt-8 grid gap-4 md:grid-cols-2">
          {advisors.map((advisor) => (
            <article key={advisor.id} className="rounded-[1.5rem] border border-black/8 bg-white/60 p-6">
              <p className="text-sm text-[color:var(--muted)]">{advisor.specialty}</p>
              <h2 className="mt-2 text-2xl font-semibold">{advisor.name}</h2>
              <p className="mt-3 text-base leading-8 text-[color:var(--muted)]">{advisor.note}</p>
              <div className="mt-5 flex items-center justify-between">
                <span className="text-sm text-[color:var(--accent)]">{advisor.responseTime}</span>
                <button className="rounded-full bg-[color:var(--accent)] px-4 py-2 text-sm font-medium text-white">
                  提交需求
                </button>
              </div>
            </article>
          ))}
        </div>
      </section>
    </Shell>
  );
}
