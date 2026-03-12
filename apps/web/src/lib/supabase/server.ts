import { createServerClient } from "@supabase/ssr";
import { getSupabaseConfig } from "@/lib/supabase/config";
import { cookies } from "next/headers";

export async function createServerSupabaseClient() {
  const config = getSupabaseConfig();
  if (!config) {
    throw new Error("Missing or invalid Supabase server environment variables.");
  }

  const cookieStore = await cookies();

  return createServerClient(config.url, config.anonKey, {
    cookies: {
      getAll() {
        return cookieStore.getAll();
      },
      setAll(cookiesToSet) {
        try {
          cookiesToSet.forEach(({ name, value, options }) => {
            cookieStore.set(name, value, options);
          });
        } catch {
          // Server components may not be able to set cookies. Middleware/route handlers handle refresh.
        }
      },
    },
  });
}
