import { createServerSupabaseClient } from "@/lib/supabase/server";

export type ServerAuthContext = {
  enabled: boolean;
  isAuthenticated: boolean;
  isAdmin: boolean;
  email: string | null;
  accessToken: string | null;
  role: string;
  capabilities: string[];
  source: string;
};

function parseAdminEmails() {
  return (process.env.ADMIN_EMAILS || "")
    .split(",")
    .map((value) => value.trim().toLowerCase())
    .filter(Boolean);
}

function getApiBaseUrl() {
  return process.env.API_BASE_URL || process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";
}

async function getApiAuthContext(email: string, accessToken: string | null) {
  try {
    const headers: HeadersInit = {
      Accept: "application/json",
    };
    if (accessToken) {
      headers.Authorization = `Bearer ${accessToken}`;
    }

    const path = accessToken ? "/v1/auth/context" : `/v1/auth/context?email=${encodeURIComponent(email)}`;
    const response = await fetch(`${getApiBaseUrl()}${path}`, {
      headers,
      cache: "no-store",
    });

    if (!response.ok) {
      return null;
    }

    const payload = (await response.json()) as {
      data?: {
        role?: string;
        capabilities?: string[];
        source?: string;
      };
    };

    if (!payload.data?.role) {
      return null;
    }

    return {
      role: payload.data.role,
      capabilities: payload.data.capabilities ?? [],
      source: payload.data.source ?? "api",
    };
  } catch {
    return null;
  }
}

export async function getServerAuthContext(): Promise<ServerAuthContext> {
  try {
    const supabase = await createServerSupabaseClient();
    const {
      data: { session },
    } = await supabase.auth.getSession();
    const {
      data: { user },
    } = await supabase.auth.getUser();

    const email = user?.email?.toLowerCase() ?? null;
    const accessToken = session?.access_token ?? null;
    if (!email) {
      return {
        enabled: true,
        isAuthenticated: false,
        isAdmin: false,
        email: null,
        accessToken: null,
        role: "guest",
        capabilities: ["browse"],
        source: "supabase-session",
      };
    }

    const apiContext = await getApiAuthContext(email, accessToken);
    if (apiContext) {
      return {
        enabled: true,
        isAuthenticated: true,
        isAdmin: apiContext.role == "admin",
        email,
        accessToken,
        role: apiContext.role,
        capabilities: apiContext.capabilities,
        source: apiContext.source,
      };
    }

    const adminEmails = parseAdminEmails();
    const isAdmin = adminEmails.includes(email);

    return {
      enabled: true,
      isAuthenticated: true,
      isAdmin: isAdmin,
      email,
      accessToken,
      role: isAdmin ? "admin" : "member",
      capabilities: isAdmin ? ["browse", "member", "admin"] : ["browse", "member"],
      source: "web-fallback-allowlist",
    };
  } catch {
    return {
      enabled: false,
      isAuthenticated: false,
      isAdmin: false,
      email: null,
      accessToken: null,
      role: "guest",
      capabilities: ["browse"],
      source: "disabled",
    };
  }
}
