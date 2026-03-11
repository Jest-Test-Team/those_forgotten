import Link from "next/link";
import { Shell } from "@/components/shell";

export default function NotFound() {
  return (
    <Shell>
      <section className="glass rounded-[2rem] p-8 text-center">
        <p className="label">404</p>
        <h1 className="mt-4 font-display text-4xl font-semibold">找不到這個頁面</h1>
        <p className="mt-4 text-[color:var(--muted)]">請回到標售清單或首頁重新查詢。</p>
        <Link className="mt-6 inline-flex rounded-full bg-[color:var(--accent)] px-5 py-3 text-white" href="/">
          回首頁
        </Link>
      </section>
    </Shell>
  );
}
