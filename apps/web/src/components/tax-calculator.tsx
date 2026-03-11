"use client";

import { useState } from "react";

export function TaxCalculator() {
  const [price, setPrice] = useState(50000);
  const vat = Math.round(price * 0.05);
  const fees = Math.round(price * 0.015);
  const total = price + vat + fees;

  return (
    <div className="glass rounded-[2rem] p-6">
      <p className="label">稅費試算器</p>
      <div className="mt-4 space-y-4">
        <label className="block text-sm text-[color:var(--muted)]">
          得標價
          <input
            className="mt-2 w-full rounded-2xl border border-black/10 bg-white/70 px-4 py-3 text-base outline-none"
            min={0}
            step={1000}
            type="number"
            value={price}
            onChange={(event) => setPrice(Number(event.target.value))}
          />
        </label>
        <div className="grid gap-3 rounded-[1.5rem] bg-[#17342d] p-4 text-[#f4ead9] sm:grid-cols-3">
          <div>
            <p className="text-sm text-white/70">營業稅 5%</p>
            <p className="mt-1 text-xl font-semibold">NT$ {vat.toLocaleString()}</p>
          </div>
          <div>
            <p className="text-sm text-white/70">規費估算 1.5%</p>
            <p className="mt-1 text-xl font-semibold">NT$ {fees.toLocaleString()}</p>
          </div>
          <div>
            <p className="text-sm text-white/70">預估總成本</p>
            <p className="mt-1 text-xl font-semibold text-[#f3b76d]">NT$ {total.toLocaleString()}</p>
          </div>
        </div>
      </div>
    </div>
  );
}
