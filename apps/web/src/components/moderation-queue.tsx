"use client";

import { useState, useTransition } from "react";
import { resolveCommunityReport } from "@/lib/browser-api";
import type { CommunityModerationReport } from "@/lib/api";

type Props = {
  initialReports: CommunityModerationReport[];
};

export function ModerationQueue({ initialReports }: Props) {
  const [reports, setReports] = useState(initialReports);
  const [message, setMessage] = useState("");
  const [activeId, setActiveId] = useState<string | null>(null);
  const [isPending, startTransition] = useTransition();

  function onResolve(reportId: string) {
    setActiveId(reportId);

    startTransition(async () => {
      const result = await resolveCommunityReport(reportId);
      setMessage(result.message);

      if (result.ok) {
        setReports((current) =>
          current.map((report) =>
            report.id === reportId
              ? {
                  ...report,
                  status: "resolved",
                }
              : report,
          ),
        );
      }

      setActiveId(null);
    });
  }

  return (
    <section className="rounded-[1.5rem] border border-black/8 bg-white/70 p-6">
      <div className="flex items-center justify-between gap-4">
        <div>
          <p className="label">Moderation Queue</p>
          <h2 className="mt-2 text-2xl font-semibold">社群檢舉待處理清單</h2>
        </div>
        <p className="rounded-full bg-amber-100 px-3 py-1 text-sm font-medium text-amber-900">
          {reports.filter((report) => report.status !== "resolved").length} pending
        </p>
      </div>
      <div className="mt-5 space-y-3">
        {reports.map((report) => (
          <article key={report.id} className="rounded-[1.25rem] border border-black/8 bg-white p-4">
            <div className="flex flex-wrap items-center justify-between gap-3">
              <div>
                <p className="text-sm text-[color:var(--muted)]">{report.office}</p>
                <h3 className="mt-1 text-lg font-semibold">{report.postTitle}</h3>
              </div>
              <p
                className={`rounded-full px-3 py-1 text-xs font-medium uppercase tracking-[0.2em] ${
                  report.status === "resolved" ? "bg-emerald-100 text-emerald-900" : "bg-rose-100 text-rose-900"
                }`}
              >
                {report.status}
              </p>
            </div>
            <p className="mt-3 text-sm leading-7 text-[color:var(--muted)]">{report.reason}</p>
            <div className="mt-3 flex flex-wrap items-center justify-between gap-3">
              <p className="text-xs text-[color:var(--muted)]">
                檢舉時間 {new Date(report.createdAt).toLocaleString("zh-TW", { hour12: false })}
              </p>
              <button
                type="button"
                disabled={isPending || report.status === "resolved"}
                onClick={() => onResolve(report.id)}
                className="rounded-full border border-black/10 bg-white px-4 py-2 text-sm font-medium disabled:cursor-not-allowed disabled:opacity-60"
              >
                {activeId === report.id && isPending ? "更新中" : report.status === "resolved" ? "已完成" : "標記完成"}
              </button>
            </div>
          </article>
        ))}
      </div>
      {message ? <p className="mt-4 text-sm text-[color:var(--accent)]">{message}</p> : null}
    </section>
  );
}
