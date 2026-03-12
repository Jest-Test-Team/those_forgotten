import { createBrowserSupabaseClient } from "@/lib/supabase/client";

type MutationResult = {
  ok: boolean;
  message: string;
};

let browserSupabaseClient: ReturnType<typeof createBrowserSupabaseClient> | null = null;

function getApiBaseUrl() {
  return process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";
}

function getSupabaseClient() {
  if (!browserSupabaseClient) {
    browserSupabaseClient = createBrowserSupabaseClient();
  }

  return browserSupabaseClient;
}

async function getAuthHeaders(includeJson: boolean): Promise<HeadersInit> {
  const headers: HeadersInit = {};
  if (includeJson) {
    headers["Content-Type"] = "application/json";
  }

  try {
    const client = getSupabaseClient();
    const {
      data: { session },
    } = await client.auth.getSession();

    if (session?.access_token) {
      headers.Authorization = `Bearer ${session.access_token}`;
    }
  } catch {
    // Keep unauthenticated fallback behavior for local/demo flows.
  }

  return headers;
}

async function postJson(path: string, payload: unknown, fallbackMessage: string): Promise<MutationResult> {
  try {
    const response = await fetch(`${getApiBaseUrl()}${path}`, {
      method: "POST",
      headers: await getAuthHeaders(true),
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      const body = (await response.json().catch(() => null)) as { error?: string } | null;
      return {
        ok: false,
        message: body?.error ?? "請求失敗",
      };
    }

    return {
      ok: true,
      message: fallbackMessage,
    };
  } catch {
    return {
      ok: true,
      message: `${fallbackMessage}（目前使用 fallback 流程，待 API 上線後改為真實送出）`,
    };
  }
}

async function deleteRequest(path: string, fallbackMessage: string): Promise<MutationResult> {
  try {
    const response = await fetch(`${getApiBaseUrl()}${path}`, {
      method: "DELETE",
      headers: await getAuthHeaders(false),
    });

    if (!response.ok) {
      const body = (await response.json().catch(() => null)) as { error?: string } | null;
      return {
        ok: false,
        message: body?.error ?? "請求失敗",
      };
    }

    return {
      ok: true,
      message: fallbackMessage,
    };
  } catch {
    return {
      ok: true,
      message: `${fallbackMessage}（目前使用 fallback 流程，待 API 上線後改為真實送出）`,
    };
  }
}

export function createKeywordSubscription(keyword: string) {
  return postJson("/v1/keyword-subscriptions", { keyword }, `關鍵字「${keyword}」已加入訂閱`);
}

export function deleteKeywordSubscription(id: string, keyword: string) {
  return deleteRequest(`/v1/keyword-subscriptions/${id}`, `關鍵字「${keyword}」已移除`);
}

export function createCommunityPost(input: { title: string; body: string; office: string }) {
  return postJson("/v1/community/posts", input, "貼文已送出並進入公開列表");
}

export function reportCommunityPost(postId: string, reason: string) {
  return postJson(`/v1/community/posts/${postId}/report`, { reason }, "檢舉已送出，管理員將進一步審核");
}

export function resolveCommunityReport(reportId: string) {
  return postJson(`/v1/admin/community-reports/${reportId}/resolve`, {}, "檢舉狀態已更新為 resolved");
}

export function createWebPushSubscription(input: {
  endpoint: string;
  keys: {
    p256dh: string;
    auth: string;
  };
}) {
  return postJson("/v1/web-push-subscriptions", input, "瀏覽器通知已啟用");
}

export function createAdvisorLead(input: {
  advisorId: string;
  name: string;
  email: string;
  message: string;
  category?: string;
}) {
  return postJson(
    "/v1/advisor-leads",
    {
      advisor_id: input.advisorId,
      name: input.name,
      email: input.email,
      message: input.message,
      category: input.category,
    },
    "需求單已建立，待管理後台指派",
  );
}
