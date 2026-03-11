import { CommunityComposer } from "@/components/community-composer";
import { ReportPostButton } from "@/components/report-post-button";
import { Shell } from "@/components/shell";
import { getCommunityPosts } from "@/lib/api";

export const metadata = {
  title: "看貨心得",
};

export default async function CommunityPage() {
  const { posts, source } = await getCommunityPosts();

  return (
    <Shell>
      <section className="glass rounded-[2rem] p-8">
        <div className="flex flex-wrap items-end justify-between gap-4">
          <div>
            <p className="label">Community</p>
            <h1 className="mt-4 font-display text-4xl font-semibold">看貨心得與風險分享</h1>
          </div>
          <div className="rounded-[1.5rem] bg-[#17342d] px-4 py-3 text-sm text-white">
            發文資格：Email 已驗證會員
          </div>
        </div>
        <div className="mt-8 grid gap-6 lg:grid-cols-[0.95fr_1.05fr]">
          <CommunityComposer />
          <div className="space-y-4">
            <p className="text-sm text-[color:var(--muted)]">目前資料來源：{source === "api" ? "即時 API" : "本地 seed fallback"}</p>
            {posts.map((post) => (
              <article key={post.id} className="rounded-[1.5rem] border border-black/8 bg-white/60 p-6">
                <div className="flex flex-wrap items-center justify-between gap-3">
                  <div>
                    <p className="text-sm text-[color:var(--muted)]">{post.office}</p>
                    <h2 className="mt-2 text-2xl font-semibold">{post.title}</h2>
                  </div>
                  <span className="rounded-full border border-black/10 px-3 py-1 text-sm">{post.author}</span>
                </div>
                <p className="mt-4 text-base leading-8 text-[color:var(--muted)]">{post.body}</p>
                <ReportPostButton postId={post.id} />
              </article>
            ))}
          </div>
        </div>
      </section>
    </Shell>
  );
}
