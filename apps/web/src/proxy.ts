import { NextResponse, type NextRequest } from "next/server";
import { updateSession } from "@/lib/supabase/middleware";

const protectedPrefixes = ["/member", "/admin"];

export async function proxy(request: NextRequest) {
  const { response, supabase } = updateSession(request);
  const pathname = request.nextUrl.pathname;

  if (!supabase) {
    return response;
  }

  let user = null;
  try {
    const result = await supabase.auth.getUser();
    user = result.data.user;
  } catch {
    // If Supabase is temporarily unavailable or misconfigured at the edge,
    // do not crash protected routes with a 500. Let the page-level auth gate render.
    return response;
  }

  const isProtected = protectedPrefixes.some((prefix) => pathname.startsWith(prefix));

  if (isProtected && !user) {
    const redirectUrl = request.nextUrl.clone();
    redirectUrl.pathname = "/auth/login";
    redirectUrl.searchParams.set("next", pathname);
    return NextResponse.redirect(redirectUrl);
  }

  if (pathname === "/auth/login" && user) {
    const redirectUrl = request.nextUrl.clone();
    redirectUrl.pathname = "/member";
    redirectUrl.search = "";
    return NextResponse.redirect(redirectUrl);
  }

  return response;
}

export const config = {
  matcher: ["/member/:path*", "/admin/:path*", "/auth/login"],
};
