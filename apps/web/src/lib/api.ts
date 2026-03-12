import {
  advisors as seedAdvisors,
  articles,
  Auction,
  auctions as seedAuctions,
  CommunityPost,
  courses as seedCourses,
  historicalCoverageStart,
  posts as seedPosts,
} from "@/lib/site-data";

type AuctionApiRecord = {
  id: string;
  title: string;
  customsOffice: string;
  closingAt: string;
  viewingAt: string;
  category?: string;
  officialUrl: string;
  summary: string;
  disclaimers: string[];
};

type AuctionHistoryApiRecord = {
  id: string;
  finalPrice: number;
  recordedAt: string;
};

type CommunityPostApiRecord = {
  id: string;
  title: string;
  body: string;
  office: string;
  author: string;
};

type AdvisorApiRecord = {
  id: string;
  name: string;
  specialty: string;
  description: string;
};

type CourseApiRecord = {
  id: string;
  title: string;
  description: string;
};

type JsonEnvelope<T> = {
  data: T;
};

type DataSource = "api" | "seed";

type HealthResponse = {
  status: string;
  repository: string;
};

type KeywordSubscriptionApiRecord = {
  id: string;
  keyword: string;
};

type CommunityReportApiRecord = {
  id: string;
  postId: string;
  postTitle?: string;
  office?: string;
  reason: string;
  status: string;
  createdAt: string;
};

type AdvisorLeadApiRecord = {
  id: string;
  advisorId: string;
  advisorName?: string;
  name: string;
  email: string;
  message: string;
  category?: string;
  createdAt?: string;
};

type CrawlerStatusApiRecord = {
  office: string;
  status: string;
  lastRunAt: string;
  nextRunAt: string;
  lastChecksum: string;
  lastRowCount: number;
  triggerSource: string;
};

export type AuctionHistory = {
  id: string;
  finalPrice: number;
  recordedAt: string;
};

export type CommunityModerationReport = {
  id: string;
  postId: string;
  postTitle: string;
  office: string;
  reason: string;
  status: string;
  createdAt: string;
};

export type AdminAdvisorLead = {
  id: string;
  advisorId: string;
  advisorName: string;
  name: string;
  email: string;
  message: string;
  category: string;
  createdAt: string;
};

export type AdminCrawlerStatus = {
  office: string;
  status: string;
  lastRunAt: string;
  nextRunAt: string;
  lastChecksum: string;
  lastRowCount: number;
  triggerSource: string;
};

export type HomeData = {
  auctions: Auction[];
  posts: CommunityPost[];
  advisors: typeof seedAdvisors;
  courses: typeof seedCourses;
  source: DataSource;
};

export type AdminDashboardData = {
  liveAuctionCount: number;
  communityPostCount: number;
  advisorCount: number;
  courseCount: number;
  moderationReports: CommunityModerationReport[];
  advisorLeads: AdminAdvisorLead[];
  crawlerStatuses: AdminCrawlerStatus[];
  source: DataSource;
  backendHealth: HealthResponse;
};

export type KeywordSubscription = {
  id: string;
  keyword: string;
};

const defaultHistory: Record<string, AuctionHistory[]> = {
  "lot-camera-001": [{ id: "result-001", finalPrice: 58000, recordedAt: "2026-03-10T12:00:00Z" }],
  "lot-car-002": [{ id: "result-002", finalPrice: 193000, recordedAt: "2026-03-09T12:00:00Z" }],
};

function getApiBaseUrl() {
  return process.env.API_BASE_URL || process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";
}

async function safeFetch<T>(path: string, options?: { actorEmail?: string }): Promise<T | null> {
  try {
    const headers: HeadersInit = { Accept: "application/json" };
    if (options?.actorEmail) {
      headers["X-Actor-Email"] = options.actorEmail;
    }

    const response = await fetch(`${getApiBaseUrl()}${path}`, {
      headers,
      next: { revalidate: 300 },
    });

    if (!response.ok) {
      return null;
    }

    return (await response.json()) as T;
  } catch {
    return null;
  }
}

function normalizeAuction(record: AuctionApiRecord): Auction {
  return {
    id: record.id,
    title: record.title,
    office: record.customsOffice,
    category: record.category ?? "未分類",
    closingAt: record.closingAt,
    viewingAt: record.viewingAt,
    priceBand: "依公告與現場貨況估值",
    warnings: record.disclaimers,
    officialUrl: record.officialUrl,
    summary: record.summary,
  };
}

export async function getAuctions(): Promise<{ auctions: Auction[]; source: DataSource }> {
  const response = await safeFetch<JsonEnvelope<AuctionApiRecord[]>>("/v1/auctions");
  if (!response?.data?.length) {
    return { auctions: seedAuctions, source: "seed" };
  }

  return {
    auctions: response.data.map(normalizeAuction),
    source: "api",
  };
}

export async function getAuctionById(id: string): Promise<{ auction: Auction | null; source: DataSource }> {
  const response = await safeFetch<JsonEnvelope<AuctionApiRecord>>(`/v1/auctions/${id}`);
  if (!response?.data) {
    return {
      auction: seedAuctions.find((item) => item.id === id) ?? null,
      source: "seed",
    };
  }

  return {
    auction: normalizeAuction(response.data),
    source: "api",
  };
}

export async function getAuctionHistory(id: string): Promise<{ rows: AuctionHistory[]; source: DataSource }> {
  const response = await safeFetch<JsonEnvelope<AuctionHistoryApiRecord[]>>(`/v1/auctions/${id}/history`);
  if (!response?.data) {
    return {
      rows: defaultHistory[id] ?? [],
      source: "seed",
    };
  }

  return {
    rows: response.data.map((row) => ({
      id: row.id,
      finalPrice: row.finalPrice,
      recordedAt: row.recordedAt,
    })),
    source: "api",
  };
}

export async function getCommunityPosts(): Promise<{ posts: CommunityPost[]; source: DataSource }> {
  const response = await safeFetch<JsonEnvelope<CommunityPostApiRecord[]>>("/v1/community/posts");
  if (!response?.data?.length) {
    return { posts: seedPosts, source: "seed" };
  }

  return {
    posts: response.data.map((post) => ({
      id: post.id,
      title: post.title,
      office: post.office,
      author: post.author,
      body: post.body,
    })),
    source: "api",
  };
}

export async function getAdvisors(): Promise<{ advisors: typeof seedAdvisors; source: DataSource }> {
  const response = await safeFetch<JsonEnvelope<AdvisorApiRecord[]>>("/v1/advisors");
  if (!response?.data?.length) {
    return { advisors: seedAdvisors, source: "seed" };
  }

  return {
    advisors: response.data.map((advisor, index) => ({
      id: advisor.id,
      name: advisor.name,
      specialty: advisor.specialty,
      note: advisor.description,
      responseTime: seedAdvisors[index]?.responseTime ?? "待設定回覆 SLA",
    })),
    source: "api",
  };
}

export async function getCourses(): Promise<{ courses: typeof seedCourses; source: DataSource }> {
  const response = await safeFetch<JsonEnvelope<CourseApiRecord[]>>("/v1/courses");
  if (!response?.data?.length) {
    return { courses: seedCourses, source: "seed" };
  }

  return {
    courses: response.data.map((course) => ({
      id: course.id,
      title: course.title,
      summary: course.description,
    })),
    source: "api",
  };
}

export async function getKeywordSubscriptions(): Promise<{ subscriptions: KeywordSubscription[]; source: DataSource }> {
  const response = await safeFetch<JsonEnvelope<KeywordSubscriptionApiRecord[]>>("/v1/keyword-subscriptions");
  if (!response?.data?.length) {
    return {
      subscriptions: [
        { id: "seed-sub-1", keyword: "相機" },
        { id: "seed-sub-2", keyword: "進口車" },
      ],
      source: "seed",
    };
  }

  return {
    subscriptions: response.data.map((subscription) => ({
      id: subscription.id,
      keyword: subscription.keyword,
    })),
    source: "api",
  };
}

export async function getHomeData(): Promise<HomeData> {
  const [auctionResult, postResult, advisorResult, courseResult] = await Promise.all([
    getAuctions(),
    getCommunityPosts(),
    getAdvisors(),
    getCourses(),
  ]);

  const source: DataSource =
    auctionResult.source === "api" ||
    postResult.source === "api" ||
    advisorResult.source === "api" ||
    courseResult.source === "api"
      ? "api"
      : "seed";

  return {
    auctions: auctionResult.auctions,
    posts: postResult.posts,
    advisors: advisorResult.advisors,
    courses: courseResult.courses,
    source,
  };
}

export async function getAdminDashboardData(actorEmail?: string): Promise<AdminDashboardData> {
  const [auctionResult, postResult, advisorResult, courseResult, moderationResult, advisorLeadResult, crawlerResult] = await Promise.all([
    getAuctions(),
    getCommunityPosts(),
    getAdvisors(),
    getCourses(),
    safeFetch<JsonEnvelope<CommunityReportApiRecord[]>>("/v1/admin/community-reports", { actorEmail }),
    safeFetch<JsonEnvelope<AdvisorLeadApiRecord[]>>("/v1/admin/advisor-leads", { actorEmail }),
    safeFetch<JsonEnvelope<CrawlerStatusApiRecord[]>>("/v1/admin/crawler-status", { actorEmail }),
  ]);
  const health = await safeFetch<HealthResponse>("/healthz");

  return {
    liveAuctionCount: auctionResult.auctions.length,
    communityPostCount: postResult.posts.length,
    advisorCount: advisorResult.advisors.length,
    courseCount: courseResult.courses.length,
    moderationReports:
      moderationResult?.data?.length
        ? moderationResult.data.map((report) => ({
            id: report.id,
            postId: report.postId,
            postTitle: report.postTitle ?? "待查貼文",
            office: report.office ?? "未標記關別",
            reason: report.reason,
            status: report.status,
            createdAt: report.createdAt,
          }))
        : [
            {
              id: "seed-report-001",
              postId: "post-001",
              postTitle: "臺北關相機批次看貨紀錄",
              office: "臺北關",
              reason: "缺少看貨照片佐證",
              status: "pending",
              createdAt: "2026-03-11T10:00:00Z",
            },
          ],
    advisorLeads:
      advisorLeadResult?.data?.length
        ? advisorLeadResult.data.map((lead) => ({
            id: lead.id,
            advisorId: lead.advisorId,
            advisorName: lead.advisorName ?? "待指派顧問",
            name: lead.name,
            email: lead.email,
            message: lead.message,
            category: lead.category ?? "一般諮詢",
            createdAt: lead.createdAt ?? "2026-03-11T10:30:00Z",
          }))
        : [
            {
              id: "seed-lead-001",
              advisorId: "advisor-001",
              advisorName: "王顧問",
              name: "示例會員",
              email: "member@example.com",
              message: "需要協助驗車與提領安排。",
              category: "進口車驗車",
              createdAt: "2026-03-11T10:30:00Z",
            },
          ],
    crawlerStatuses:
      crawlerResult?.data?.length
        ? crawlerResult.data.map((crawler) => ({
            office: crawler.office,
            status: crawler.status,
            lastRunAt: crawler.lastRunAt,
            nextRunAt: crawler.nextRunAt,
            lastChecksum: crawler.lastChecksum,
            lastRowCount: crawler.lastRowCount,
            triggerSource: crawler.triggerSource,
          }))
        : [
            {
              office: "基隆關",
              status: "healthy",
              lastRunAt: "2026-03-11T10:18:00Z",
              nextRunAt: "2026-03-11T10:48:00Z",
              lastChecksum: "keelung-demo",
              lastRowCount: 4,
              triggerSource: "schedule",
            },
            {
              office: "臺中關",
              status: "warning",
              lastRunAt: "2026-03-11T09:57:00Z",
              nextRunAt: "2026-03-11T10:27:00Z",
              lastChecksum: "taichung-demo",
              lastRowCount: 0,
              triggerSource: "retry",
            },
          ],
    source:
      auctionResult.source === "api" ||
      postResult.source === "api" ||
      advisorResult.source === "api" ||
      courseResult.source === "api"
        ? "api"
        : "seed",
    backendHealth: health ?? {
      status: "degraded",
      repository: "unreachable",
    },
  };
}

export function getKnowledgeArticle(slug: string) {
  return articles.find((item) => item.slug === slug) ?? null;
}

export { historicalCoverageStart };
