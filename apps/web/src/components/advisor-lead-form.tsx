"use client";

import { useState } from "react";

type Props = {
  advisorName: string;
};

export function AdvisorLeadForm({ advisorName }: Props) {
  const [submitted, setSubmitted] = useState(false);

  return (
    <div className="mt-5 rounded-[1.5rem] border border-black/8 bg-[#f8f1e5] p-4">
      <p className="text-sm font-medium">送出給 {advisorName} 的需求摘要</p>
      <form
        className="mt-3 space-y-3"
        onSubmit={(event) => {
          event.preventDefault();
          setSubmitted(true);
        }}
      >
        <input className="w-full rounded-2xl border border-black/10 bg-white/80 px-4 py-3 text-sm" placeholder="您的 Email" />
        <textarea className="min-h-24 w-full rounded-2xl border border-black/10 bg-white/80 px-4 py-3 text-sm" placeholder="描述標案、預算、需求與是否需驗車/資格協助" />
        <button className="rounded-full bg-[color:var(--accent)] px-4 py-2 text-sm font-medium text-white" type="submit">
          提交需求
        </button>
      </form>
      {submitted ? <p className="mt-3 text-sm text-[color:var(--accent)]">示意流程：需求單已建立，待管理後台指派。</p> : null}
    </div>
  );
}
