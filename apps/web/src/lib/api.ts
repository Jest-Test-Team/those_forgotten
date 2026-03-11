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
  category?: string;
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

export type AuctionHistory = {
  id: string;
  finalPrice: number;
  recordedAt: string;
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
  source: DataSource;
};

const defaultHistory: Record<string, AuctionHistory[]> = {
  "lot-camera-001": [{ id: "result-001", finalPrice: 58000, recordedAt: "2026-03-10T12:00:00Z" }],
  "lot-car-002": [{ id: "result-002", finalPrice: 193000, recordedAt: "2026-03-09T12:00:00Z" }],
};

function getApiBaseUrl() {
  return process.env.API_BASE_URL || process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";
}

async function safeFetch<T>(path: string): Promise<T | null> {
  try {
    const response = await fetch(`${getApiBaseUrl()}${path}`, {
      headers: { Accept: "application/json" },
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
  const seedMatch = seedAuctions.find((item) => item.id === record.id);

  return {
    id: record.id,
    title: record.title,
    office: record.customsOffice,
    category: record.category ?? "未分類",
    closingAt: record.closingAt,
    viewingAt: "詳見官方公告",
    priceBand: "依公告與現場貨況估值",
    warnings: record.disclaimers,
    officialUrl: seedMatch?.officialUrl ?? "/member",
    summary: `${record.customsOffice} / ${record.category ?? "未分類"} 標案，請搭配官方文件與現場看貨判斷風險。`,
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

export async function getAdminDashboardData(): Promise<AdminDashboardData> {
  const [auctionResult, postResult, advisorResult, courseResult] = await Promise.all([
    getAuctions(),
    getCommunityPosts(),
    getAdvisors(),
    getCourses(),
  ]);

  return {
    liveAuctionCount: auctionResult.auctions.length,
    communityPostCount: postResult.posts.length,
    advisorCount: advisorResult.advisors.length,
    courseCount: courseResult.courses.length,
    source:
      auctionResult.source === "api" ||
      postResult.source === "api" ||
      advisorResult.source === "api" ||
      courseResult.source === "api"
        ? "api"
        : "seed",
  };
}

export function getKnowledgeArticle(slug: string) {
  return articles.find((item) => item.slug === slug) ?? null;
}

export { historicalCoverageStart };
