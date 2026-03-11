import { disclaimers } from "@/lib/site-data";

export function RiskPanel() {
  return (
    <div className="glass rounded-[2rem] p-6">
      <p className="label">交易風險</p>
      <ul className="mt-4 space-y-3 text-sm leading-7 text-[color:var(--muted)]">
        {disclaimers.map((item) => (
          <li key={item} className="rounded-2xl border border-black/8 bg-white/60 px-4 py-3">
            {item}
          </li>
        ))}
      </ul>
    </div>
  );
}
