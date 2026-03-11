import { AuthStatus } from "@/components/auth-status";
import { Shell } from "@/components/shell";

export const metadata = {
  title: "登入",
};

export default function LoginPage() {
  return (
    <Shell>
      <section className="glass rounded-[2rem] p-8">
        <p className="label">Supabase Auth</p>
        <h1 className="mt-4 font-display text-4xl font-semibold">登入與 Session 啟用</h1>
        <p className="mt-4 max-w-2xl text-base leading-8 text-[color:var(--muted)]">
          這一版先接上 Google OAuth browser flow。完成登入後會導回會員中心；下一階段再補 server-side middleware、
          RBAC 與真正的 protected routes。
        </p>
        <div className="mt-8 rounded-[1.5rem] border border-black/8 bg-white/60 p-6">
          <AuthStatus />
        </div>
      </section>
    </Shell>
  );
}
