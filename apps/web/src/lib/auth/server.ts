import { createServerSupabaseClient } from "@/lib/supabase/server";

export type ServerAuthContext = {
  enabled: boolean;
  isAuthenticated: boolean;
  isAdmin: boolean;
  email: string | null;
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

async function getApiAuthContext(email: string) {
  try {
    const response = await fetch(`${getApiBaseUrl()}/v1/auth/context?email=${encodeURIComponent(email)}`, {
      headers: {
        Accept: "application/json",
      },
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
      data: { user },
    } = await supabase.auth.getUser();

    const email = user?.email?.toLowerCase() ?? null;
    if (!email) {
      return {
        enabled: true,
        isAuthenticated: false,
        isAdmin: false,
        email: null,
        role: "guest",
        capabilities: ["browse"],
        source: "supabase-session",
      };
    }

    const apiContext = await getApiAuthContext(email);
    if (apiContext) {
      return {
        enabled: true,
        isAuthenticated: true,
        isAdmin: apiContext.role == "admin",
        email,
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
      role: "guest",
      capabilities: ["browse"],
      source: "disabled",
    };
  }
}
