"use client";

import { useState } from "react";

export function CommunityComposer() {
  const [submitted, setSubmitted] = useState(false);

  return (
    <div className="glass rounded-[2rem] p-6">
      <p className="label">發佈心得</p>
      <form
        className="mt-4 space-y-3"
        onSubmit={(event) => {
          event.preventDefault();
          setSubmitted(true);
        }}
      >
        <input className="w-full rounded-2xl border border-black/10 bg-white/70 px-4 py-3" placeholder="貼文標題" />
        <textarea
          className="min-h-32 w-full rounded-2xl border border-black/10 bg-white/70 px-4 py-3"
          placeholder="描述實際看貨狀況、缺件、外觀磨損與估值觀察"
        />
        <div className="flex items-center justify-between gap-3">
          <p className="text-sm text-[color:var(--muted)]">先發後審，任何內容都可能被檢舉與下架。</p>
          <button className="rounded-full bg-[color:var(--accent)] px-5 py-3 text-sm font-medium text-white" type="submit">
            送出貼文
          </button>
        </div>
      </form>
      {submitted ? <p className="mt-4 text-sm text-[color:var(--accent)]">示意流程：貼文已建立並進入公開列表。</p> : null}
    </div>
  );
}
