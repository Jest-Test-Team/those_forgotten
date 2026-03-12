import { createBrowserClient } from "@supabase/ssr";
import { getSupabaseConfig } from "@/lib/supabase/config";

export function createBrowserSupabaseClient() {
  const config = getSupabaseConfig();
  if (!config) {
    throw new Error("Missing or invalid Supabase browser environment variables.");
  }

  return createBrowserClient(config.url, config.anonKey);
}
