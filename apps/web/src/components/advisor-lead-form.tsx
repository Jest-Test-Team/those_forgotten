"use client";

import { FormEvent, useState, useTransition } from "react";
import { createAdvisorLead } from "@/lib/browser-api";

type Props = {
  advisorId: string;
  advisorName: string;
};

export function AdvisorLeadForm({ advisorId, advisorName }: Props) {
  const [email, setEmail] = useState("");
  const [messageBody, setMessageBody] = useState("");
  const [statusMessage, setStatusMessage] = useState("");
  const [isPending, startTransition] = useTransition();

  function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    startTransition(async () => {
      const result = await createAdvisorLead({
        advisorId,
        name: "平台會員",
        email,
        message: messageBody,
        category: "代標需求",
      });

      setStatusMessage(result.message);
      if (result.ok) {
        setEmail("");
        setMessageBody("");
      }
    });
  }

  return (
    <div className="mt-5 rounded-[1.5rem] border border-black/8 bg-[#f8f1e5] p-4">
      <p className="text-sm font-medium">送出給 {advisorName} 的需求摘要</p>
      <form className="mt-3 space-y-3" onSubmit={onSubmit}>
        <input
          className="w-full rounded-2xl border border-black/10 bg-white/80 px-4 py-3 text-sm"
          placeholder="您的 Email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
        />
        <textarea
          className="min-h-24 w-full rounded-2xl border border-black/10 bg-white/80 px-4 py-3 text-sm"
          placeholder="描述標案、預算、需求與是否需驗車/資格協助"
          value={messageBody}
          onChange={(event) => setMessageBody(event.target.value)}
        />
        <button
          className="rounded-full bg-[color:var(--accent)] px-4 py-2 text-sm font-medium text-white"
          disabled={isPending}
          type="submit"
        >
          {isPending ? "送出中" : "提交需求"}
        </button>
      </form>
      {statusMessage ? <p className="mt-3 text-sm text-[color:var(--accent)]">{statusMessage}</p> : null}
    </div>
  );
}
