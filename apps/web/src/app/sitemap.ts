import type { MetadataRoute } from "next";
import { articles, auctions } from "@/lib/site-data";

export default function sitemap(): MetadataRoute.Sitemap {
  const base = "https://customs-auction.local";

  return [
    "",
    "/auctions",
    "/community",
    "/advisors",
    "/member",
    "/admin",
    ...auctions.map((auction) => `/auctions/${auction.id}`),
    ...articles.map((article) => `/knowledge/${article.slug}`),
  ].map((path) => ({
    url: `${base}${path}`,
    lastModified: new Date(),
  }));
}
