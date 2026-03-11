"use client";

import { useState, useTransition } from "react";
import { reportCommunityPost } from "@/lib/browser-api";

type Props = {
  postId: string;
};

export function ReportPostButton({ postId }: Props) {
  const [message, setMessage] = useState("");
  const [isPending, startTransition] = useTransition();

  return (
    <div className="mt-4 flex flex-wrap items-center justify-between gap-3">
      <p className="text-sm text-[color:var(--warning)]">社群採先發後審，任何貼文都可檢舉並由管理員下架。</p>
      <button
        type="button"
        disabled={isPending}
        onClick={() => {
          startTransition(async () => {
            const result = await reportCommunityPost(postId, "疑似內容不實或描述不完整");
            setMessage(result.message);
          });
        }}
        className="rounded-full border border-black/10 bg-white px-4 py-2 text-sm font-medium"
      >
        {isPending ? "送出中" : "檢舉貼文"}
      </button>
      {message ? <p className="basis-full text-sm text-[color:var(--accent)]">{message}</p> : null}
    </div>
  );
}
