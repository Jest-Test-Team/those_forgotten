"use client";

import Link from "next/link";
import { ReactNode, useEffect, useState } from "react";
import { createBrowserSupabaseClient } from "@/lib/supabase/client";

type Props = {
  children: ReactNode;
  title: string;
  description: string;
};

type GateState = {
  ready: boolean;
  signedIn: boolean;
  enabled: boolean;
};

export function AuthGate({ children, title, description }: Props) {
  const [client] = useState(() => {
    try {
      return createBrowserSupabaseClient();
    } catch {
      return null;
    }
  });
  const [gateState, setGateState] = useState<GateState>(() => ({
    ready: client === null,
    signedIn: false,
    enabled: client !== null,
  }));

  useEffect(() => {
    if (!client) {
      return undefined;
    }

    client.auth.getUser().then(({ data }) => {
      setGateState({
        ready: true,
        signedIn: Boolean(data.user),
        enabled: true,
      });
    });

    const {
      data: { subscription },
    } = client.auth.onAuthStateChange((_event, session) => {
      setGateState({
        ready: true,
        signedIn: Boolean(session?.user),
        enabled: true,
      });
    });

    return () => subscription.unsubscribe();
  }, [client]);

  if (!gateState.ready) {
    return (
      <section className="glass rounded-[2rem] p-8">
        <p className="label">Auth Check</p>
        <h1 className="mt-4 font-display text-3xl font-semibold">驗證狀態讀取中</h1>
      </section>
    );
  }

  if (!gateState.enabled) {
    return (
      <section className="glass rounded-[2rem] p-8">
        <p className="label">Supabase Auth</p>
        <h1 className="mt-4 font-display text-3xl font-semibold">尚未設定登入環境</h1>
        <p className="mt-4 text-base leading-8 text-[color:var(--muted)]">
          目前缺少 Supabase URL 或 anon key，請先完成 web 環境設定。
        </p>
      </section>
    );
  }

  if (!gateState.signedIn) {
    return (
      <section className="glass rounded-[2rem] p-8">
        <p className="label">Protected Surface</p>
        <h1 className="mt-4 font-display text-3xl font-semibold">{title}</h1>
        <p className="mt-4 text-base leading-8 text-[color:var(--muted)]">{description}</p>
        <div className="mt-6">
          <Link href="/auth/login" className="rounded-full bg-[color:var(--accent)] px-5 py-3 text-sm font-medium text-white">
            前往登入
          </Link>
        </div>
      </section>
    );
  }

  return <>{children}</>;
}
