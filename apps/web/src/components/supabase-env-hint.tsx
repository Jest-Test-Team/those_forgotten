import { getSupabaseConfigStatus } from "@/lib/supabase/config";

type Props = {
  compact?: boolean;
};

export function SupabaseEnvHint({ compact = false }: Props) {
  const status = getSupabaseConfigStatus();

  if (status.enabled) {
    return null;
  }

  return (
    <div className="rounded-[1.25rem] border border-amber-200 bg-amber-50 px-4 py-4 text-sm text-amber-950">
      <p className="font-medium">目前未啟用 Supabase Auth。請先補齊以下 Vercel env：</p>
      <ul className={`mt-3 space-y-2 ${compact ? "text-xs" : ""}`}>
        {status.issues.map((issue) => (
          <li key={`${issue.envName}-${issue.kind}`}>
            <span className="font-semibold">{issue.envName}</span>:
            {" "}
            {issue.kind === "missing" ? "未設定" : "格式錯誤"}
            {" "}
            {issue.detail}
          </li>
        ))}
      </ul>
      {!compact ? (
        <p className="mt-3 text-xs text-amber-900/80">
          補完後重新部署，`/auth/login` 才會顯示可用的 Google 登入按鈕。
        </p>
      ) : null}
    </div>
  );
}
