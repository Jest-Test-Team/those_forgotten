import { Suspense } from "react";
import { AuthCallbackClient } from "@/components/auth-callback-client";

export default function AuthCallbackPage() {
  return (
    <main className="shell py-24">
      <div className="glass rounded-[2rem] p-8">
        <p className="label">OAuth Callback</p>
        <Suspense fallback={<h1 className="mt-4 font-display text-3xl font-semibold">正在初始化登入流程</h1>}>
          <AuthCallbackClient />
        </Suspense>
      </div>
    </main>
  );
}
