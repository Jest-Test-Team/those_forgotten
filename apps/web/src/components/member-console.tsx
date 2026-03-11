"use client";

import { useMemo, useState, useTransition } from "react";
import { createKeywordSubscription, createWebPushSubscription } from "@/lib/browser-api";

const keywordOptions = ["相機", "進口車", "名牌包", "原木"];

export function MemberConsole() {
  const [keywords, setKeywords] = useState<string[]>(["相機", "進口車"]);
  const [draft, setDraft] = useState("");
  const [pushEnabled, setPushEnabled] = useState(false);
  const [message, setMessage] = useState("");
  const [isPending, startTransition] = useTransition();

  const icsUrl = useMemo(
    () => "https://customs-auction.local/v1/auctions/calendar.ics?token=demo-member-token",
    [],
  );

  function addKeyword() {
    const value = draft.trim();
    if (!value || keywords.includes(value)) {
      return;
    }

    startTransition(async () => {
      const result = await createKeywordSubscription(value);
      setMessage(result.message);

      if (result.ok) {
        setKeywords((current) => [...current, value]);
        setDraft("");
      }
    });
  }

  function toggleKeyword(keyword: string) {
    setKeywords((current) =>
      current.includes(keyword) ? current.filter((item) => item !== keyword) : [...current, keyword],
    );
  }

  function togglePush() {
    if (pushEnabled) {
      setPushEnabled(false);
      setMessage("瀏覽器通知已停用");
      return;
    }

    startTransition(async () => {
      const result = await createWebPushSubscription({
        endpoint: "https://push.example.dev/subscriptions/demo-member",
        keys: {
          p256dh: "demo-p256dh-key",
          auth: "demo-auth-key",
        },
      });

      setMessage(result.message);
      if (result.ok) {
        setPushEnabled(true);
      }
    });
  }

  return (
    <div className="grid gap-6">
      <div className="glass rounded-[2rem] p-6">
        <p className="label">推播訂閱</p>
        <div className="mt-4 flex flex-wrap gap-3">
          {keywordOptions.map((keyword) => {
            const active = keywords.includes(keyword);
            return (
              <button
                key={keyword}
                type="button"
                onClick={() => toggleKeyword(keyword)}
                className={`rounded-full px-4 py-2 text-sm transition ${
                  active ? "bg-[#17342d] text-white" : "border border-black/10 bg-white/65"
                }`}
              >
                {keyword}
              </button>
            );
          })}
        </div>
        <div className="mt-4 flex gap-3">
          <input
            className="flex-1 rounded-full border border-black/10 bg-white/70 px-4 py-3 text-sm outline-none"
            placeholder="新增自訂關鍵字"
            value={draft}
            onChange={(event) => setDraft(event.target.value)}
          />
          <button
            type="button"
            onClick={addKeyword}
            disabled={isPending}
            className="rounded-full bg-[color:var(--accent)] px-5 py-3 text-sm font-medium text-white"
          >
            {isPending ? "送出中" : "新增"}
          </button>
        </div>
        <p className="mt-4 text-sm text-[color:var(--muted)]">
          目前訂閱：{keywords.join("、")}。付費會員匹配新案時即時接收 Web Push。
        </p>
        {message ? <p className="mt-3 text-sm text-[color:var(--accent)]">{message}</p> : null}
      </div>

      <div className="glass rounded-[2rem] p-6">
        <p className="label">Web Push</p>
        <div className="mt-4 flex items-center justify-between rounded-[1.5rem] border border-black/8 bg-white/60 p-4">
          <div>
            <p className="font-semibold">瀏覽器通知</p>
            <p className="mt-1 text-sm text-[color:var(--muted)]">啟用後可接收關鍵字推播與截標提醒。</p>
          </div>
          <button
            type="button"
            onClick={togglePush}
            disabled={isPending}
            className={`rounded-full px-4 py-2 text-sm font-medium ${pushEnabled ? "bg-[#17342d] text-white" : "border border-black/10"}`}
          >
            {isPending ? "處理中" : pushEnabled ? "已啟用" : "啟用通知"}
          </button>
        </div>
      </div>

      <div className="glass rounded-[2rem] p-6">
        <p className="label">ICS 日曆</p>
        <div className="mt-4 rounded-[1.5rem] border border-black/8 bg-white/60 p-4">
          <p className="font-semibold">個人化日曆 feed</p>
          <code className="mt-3 block overflow-x-auto rounded-2xl bg-[#17342d] px-4 py-3 text-sm text-white">
            {icsUrl}
          </code>
          <div className="mt-4 flex flex-wrap gap-3 text-sm">
            <a className="rounded-full bg-[color:var(--accent)] px-4 py-2 text-white" href={icsUrl}>
              複製連結
            </a>
            <a className="rounded-full border border-black/10 px-4 py-2" href={`https://calendar.google.com/calendar/r?cid=${encodeURIComponent(icsUrl)}`}>
              加入 Google
            </a>
          </div>
        </div>
      </div>
    </div>
  );
}
