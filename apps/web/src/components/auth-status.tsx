"use client";

import Link from "next/link";
import { useEffect, useState, useTransition } from "react";
import { createBrowserSupabaseClient } from "@/lib/supabase/client";
import type { SupabaseClient } from "@supabase/supabase-js";

type AuthState = {
  email: string | null;
  ready: boolean;
  enabled: boolean;
};

export function AuthStatus() {
  const [client] = useState<SupabaseClient | null>(() => {
    try {
      return createBrowserSupabaseClient();
    } catch {
      return null;
    }
  });
  const [authState, setAuthState] = useState<AuthState>(() => ({
    email: null,
    ready: client === null,
    enabled: client !== null,
  }));
  const [message, setMessage] = useState("");
  const [isPending, startTransition] = useTransition();

  useEffect(() => {
    if (!client) {
      return undefined;
    }

    client.auth.getUser().then(({ data }) => {
      setAuthState({
        email: data.user?.email ?? null,
        ready: true,
        enabled: true,
      });
    });

    const {
      data: { subscription },
    } = client.auth.onAuthStateChange((_event, session) => {
      setAuthState({
        email: session?.user?.email ?? null,
        ready: true,
        enabled: true,
      });
    });

    return () => subscription.unsubscribe();
  }, [client]);

  function signIn() {
    if (!client) {
      setMessage("尚未設定 Supabase 環境變數");
      return;
    }

    startTransition(async () => {
      const callbackPath = process.env.NEXT_PUBLIC_SUPABASE_AUTH_REDIRECT || "/auth/callback";
      const redirectTo = `${window.location.origin}${callbackPath}?next=/member`;
      const { data, error } = await client.auth.signInWithOAuth({
        provider: "google",
        options: {
          redirectTo,
        },
      });

      if (error) {
        setMessage(error.message);
        return;
      }

      if (data.url) {
        window.location.href = data.url;
      }
    });
  }

  function signOut() {
    if (!client) {
      setMessage("登出流程目前不可用");
      return;
    }

    startTransition(async () => {
      const { error } = await client.auth.signOut();
      if (error) {
        setMessage(error.message);
        return;
      }
      setMessage("已登出");
    });
  }

  if (!authState.ready) {
    return <p className="text-sm text-[color:var(--muted)]">驗證狀態讀取中</p>;
  }

  if (!authState.enabled) {
    return (
      <Link href="/auth/login" className="rounded-full border border-black/10 px-4 py-2 text-sm">
        設定登入
      </Link>
    );
  }

  return (
    <div className="flex items-center gap-3">
      {authState.email ? <p className="hidden text-sm text-[color:var(--muted)] lg:block">{authState.email}</p> : null}
      <button
        type="button"
        onClick={authState.email ? signOut : signIn}
        disabled={isPending}
        className="rounded-full border border-black/10 bg-white/70 px-4 py-2 text-sm"
      >
        {isPending ? "處理中" : authState.email ? "登出" : "Google 登入"}
      </button>
      {message ? <p className="hidden text-xs text-[color:var(--accent)] xl:block">{message}</p> : null}
    </div>
  );
}
