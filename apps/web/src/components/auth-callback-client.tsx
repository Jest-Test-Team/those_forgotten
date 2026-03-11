"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { createBrowserSupabaseClient } from "@/lib/supabase/client";

export function AuthCallbackClient() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const code = searchParams.get("code");
  const next = searchParams.get("next") || "/member";
  const [message, setMessage] = useState("正在交換登入 session");

  useEffect(() => {
    if (!code) {
      return;
    }
    const authCode = code;

    async function run() {
      try {
        const supabase = createBrowserSupabaseClient();
        const { error } = await supabase.auth.exchangeCodeForSession(authCode);

        if (error) {
          setMessage(error.message);
          return;
        }

        router.replace(next);
      } catch {
        setMessage("尚未設定 Supabase 環境變數");
      }
    }

    run();
  }, [code, next, router]);

  if (!code) {
    return <h1 className="mt-4 font-display text-3xl font-semibold">缺少 OAuth code，無法完成登入</h1>;
  }

  return <h1 className="mt-4 font-display text-3xl font-semibold">{message}</h1>;
}
