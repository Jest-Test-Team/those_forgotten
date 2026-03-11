import { createServerSupabaseClient } from "@/lib/supabase/server";

export type ServerAuthContext = {
  enabled: boolean;
  isAuthenticated: boolean;
  isAdmin: boolean;
  email: string | null;
};

function parseAdminEmails() {
  return (process.env.ADMIN_EMAILS || "")
    .split(",")
    .map((value) => value.trim().toLowerCase())
    .filter(Boolean);
}

export async function getServerAuthContext(): Promise<ServerAuthContext> {
  try {
    const supabase = await createServerSupabaseClient();
    const {
      data: { user },
    } = await supabase.auth.getUser();

    const email = user?.email?.toLowerCase() ?? null;
    const adminEmails = parseAdminEmails();

    return {
      enabled: true,
      isAuthenticated: Boolean(user),
      isAdmin: Boolean(email && adminEmails.includes(email)),
      email,
    };
  } catch {
    return {
      enabled: false,
      isAuthenticated: false,
      isAdmin: false,
      email: null,
    };
  }
}
