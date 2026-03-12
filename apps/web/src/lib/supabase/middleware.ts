import { createServerClient } from "@supabase/ssr";
import { getSupabaseConfig } from "@/lib/supabase/config";
import { NextResponse, type NextRequest } from "next/server";

export function updateSession(request: NextRequest) {
  const config = getSupabaseConfig();
  if (!config) {
    return {
      response: NextResponse.next({
        request,
      }),
      supabase: null,
    };
  }

  const response = NextResponse.next({
    request,
  });

  const supabase = createServerClient(config.url, config.anonKey, {
    cookies: {
      getAll() {
        return request.cookies.getAll();
      },
      setAll(cookiesToSet) {
        cookiesToSet.forEach(({ name, value, options }) => {
          request.cookies.set(name, value);
          response.cookies.set(name, value, options);
        });
      },
    },
  });

  return { response, supabase };
}
