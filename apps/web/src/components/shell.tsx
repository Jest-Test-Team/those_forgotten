import { ReactNode } from "react";
import { SiteFooter } from "@/components/site-footer";
import { SiteHeader } from "@/components/site-header";

export function Shell({ children }: { children: ReactNode }) {
  return (
    <>
      <SiteHeader />
      <main className="shell py-8">{children}</main>
      <SiteFooter />
    </>
  );
}
