import { AuthStatus } from "@/components/auth-status";
import { Shell } from "@/components/shell";

export const metadata = {
  title: "登入",
};

type Props = {
  searchParams?: Promise<{
    next?: string;
    error?: string;
  }>;
};

export default async function LoginPage({ searchParams }: Props) {
  const params = (await searchParams) ?? {};
  const nextPath = params.next || "/member";

  return (
    <Shell>
      <section className="glass rounded-[2rem] p-8">
        <p className="label">Supabase Auth</p>
        <h1 className="mt-4 font-display text-4xl font-semibold">登入與 Session 啟用</h1>
        <p className="mt-4 max-w-2xl text-base leading-8 text-[color:var(--muted)]">
          這一版已接上 Supabase SSR callback 與 middleware。完成登入後會回到原本目標頁；下一階段再補 admin role、
          member role 與真正的 RBAC。
        </p>
        {params.error ? (
          <p className="mt-4 text-sm text-rose-700">登入流程失敗：{params.error}</p>
        ) : null}
        <div className="mt-8 rounded-[1.5rem] border border-black/8 bg-white/60 p-6">
          <AuthStatus nextPath={nextPath} />
        </div>
      </section>
    </Shell>
  );
}
