import { notFound } from "next/navigation";
import { Shell } from "@/components/shell";
import { articles } from "@/lib/site-data";

type Props = {
  params: Promise<{ slug: string }>;
};

export default async function KnowledgePage({ params }: Props) {
  const { slug } = await params;
  const article = articles.find((item) => item.slug === slug);

  if (!article) {
    notFound();
  }

  return (
    <Shell>
      <article className="glass rounded-[2rem] p-8">
        <p className="label">Knowledge Base</p>
        <h1 className="mt-4 font-display text-4xl font-semibold">{article.title}</h1>
        <p className="mt-4 max-w-2xl text-lg leading-8 text-[color:var(--muted)]">{article.summary}</p>
        <div className="mt-8 space-y-6">
          {article.sections.map((section) => (
            <section key={section.heading} className="rounded-[1.5rem] border border-black/8 bg-white/60 p-6">
              <h2 className="text-2xl font-semibold">{section.heading}</h2>
              <p className="mt-3 text-base leading-8 text-[color:var(--muted)]">{section.body}</p>
            </section>
          ))}
        </div>
      </article>
    </Shell>
  );
}
