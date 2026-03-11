import { AuthStatus } from "@/components/auth-status";
import Link from "next/link";

const navItems = [
  { href: "/auctions", label: "標售清單" },
  { href: "/knowledge/bid-form-guide", label: "新手指南" },
  { href: "/community", label: "看貨心得" },
  { href: "/advisors", label: "顧問媒合" },
  { href: "/member", label: "會員中心" },
  { href: "/admin", label: "管理後台" },
];

export function SiteHeader() {
  return (
    <header className="shell pt-6">
      <div className="glass fade-up flex items-center justify-between rounded-full px-5 py-3">
        <Link href="/" className="font-display text-xl font-semibold tracking-tight">
          海關標售雷達
        </Link>
        <nav className="hidden gap-4 text-sm text-[color:var(--muted)] md:flex">
          {navItems.map((item) => (
            <Link key={item.href} href={item.href} className="transition hover:text-[color:var(--accent-strong)]">
              {item.label}
            </Link>
          ))}
        </nav>
        <div className="hidden md:block">
          <AuthStatus />
        </div>
      </div>
    </header>
  );
}
