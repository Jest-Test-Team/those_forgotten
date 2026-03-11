"use client";

import { FormEvent, useState, useTransition } from "react";
import { createCommunityPost } from "@/lib/browser-api";

export function CommunityComposer() {
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const [message, setMessage] = useState("");
  const [isPending, startTransition] = useTransition();

  function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    startTransition(async () => {
      const result = await createCommunityPost({
        title,
        body,
        office: "會員看貨分享",
      });

      setMessage(result.message);
      if (result.ok) {
        setTitle("");
        setBody("");
      }
    });
  }

  return (
    <div className="glass rounded-[2rem] p-6">
      <p className="label">發佈心得</p>
      <form className="mt-4 space-y-3" onSubmit={onSubmit}>
        <input
          className="w-full rounded-2xl border border-black/10 bg-white/70 px-4 py-3"
          placeholder="貼文標題"
          value={title}
          onChange={(event) => setTitle(event.target.value)}
        />
        <textarea
          className="min-h-32 w-full rounded-2xl border border-black/10 bg-white/70 px-4 py-3"
          placeholder="描述實際看貨狀況、缺件、外觀磨損與估值觀察"
          value={body}
          onChange={(event) => setBody(event.target.value)}
        />
        <div className="flex items-center justify-between gap-3">
          <p className="text-sm text-[color:var(--muted)]">先發後審，任何內容都可能被檢舉與下架。</p>
          <button
            className="rounded-full bg-[color:var(--accent)] px-5 py-3 text-sm font-medium text-white"
            disabled={isPending}
            type="submit"
          >
            {isPending ? "送出中" : "送出貼文"}
          </button>
        </div>
      </form>
      {message ? <p className="mt-4 text-sm text-[color:var(--accent)]">{message}</p> : null}
    </div>
  );
}
