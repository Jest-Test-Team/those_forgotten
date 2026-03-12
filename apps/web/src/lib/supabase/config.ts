type SupabaseConfig = {
  url: string;
  anonKey: string;
};

export type SupabaseConfigIssue =
  | {
      envName: "NEXT_PUBLIC_SUPABASE_URL";
      kind: "missing" | "invalid";
      detail: string;
    }
  | {
      envName: "NEXT_PUBLIC_SUPABASE_ANON_KEY";
      kind: "missing";
      detail: string;
    };

export type SupabaseConfigStatus = {
  config: SupabaseConfig | null;
  enabled: boolean;
  issues: SupabaseConfigIssue[];
};

function isHttpUrl(value: string) {
  try {
    const parsed = new URL(value);
    return parsed.protocol === "http:" || parsed.protocol === "https:";
  } catch {
    return false;
  }
}

export function getSupabaseConfigStatus(): SupabaseConfigStatus {
  const url = process.env.NEXT_PUBLIC_SUPABASE_URL?.trim() ?? "";
  const anonKey = process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY?.trim() ?? "";
  const issues: SupabaseConfigIssue[] = [];

  if (!url) {
    issues.push({
      envName: "NEXT_PUBLIC_SUPABASE_URL",
      kind: "missing",
      detail: "缺少 Supabase 專案 URL。",
    });
  } else if (!isHttpUrl(url)) {
    issues.push({
      envName: "NEXT_PUBLIC_SUPABASE_URL",
      kind: "invalid",
      detail: "必須是合法的 http 或 https URL。",
    });
  }

  if (!anonKey) {
    issues.push({
      envName: "NEXT_PUBLIC_SUPABASE_ANON_KEY",
      kind: "missing",
      detail: "缺少 Supabase anon key。",
    });
  }

  if (issues.length > 0) {
    return {
      config: null,
      enabled: false,
      issues,
    };
  }

  return {
    config: {
      url,
      anonKey,
    },
    enabled: true,
    issues: [],
  };
}

export function getSupabaseConfig(): SupabaseConfig | null {
  return getSupabaseConfigStatus().config;
}
