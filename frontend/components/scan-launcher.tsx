"use client";

import { FormEvent, useState, useTransition } from "react";
import { useRouter } from "next/navigation";
import { startScan } from "../lib/api";

export function ScanLauncher() {
  const router = useRouter();
  const [target, setTarget] = useState("https://example.com");
  const [error, setError] = useState("");
  const [isPending, startTransition] = useTransition();

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const nextTarget = String(formData.get("target") || "").trim();
    setError("");

    try {
      const result = await startScan(nextTarget);
      startTransition(() => {
        router.push(`/scan/${result.scan_id}`);
      });
    } catch (submitError) {
      setError(submitError instanceof Error ? submitError.message : "Bilinmeyen hata");
    }
  }

  return (
    <form onSubmit={onSubmit} className="glass rounded-[32px] p-6 md:p-8">
      <div className="mb-4">
        <p className="mb-2 text-xs uppercase tracking-[0.35em] text-stone-500">Target URL</p>
        <input
          name="target"
          value={target}
          onChange={(event) => setTarget(event.target.value)}
          placeholder="https://app.example.com"
          className="w-full rounded-2xl border border-stone-300 bg-white/70 px-4 py-4 text-lg outline-none transition focus:border-orange-500"
        />
      </div>

      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <p className="max-w-xl text-sm text-stone-600">
          SecScan, SSRF guard ile özel ağ hedeflerini reddeder ve 7 güvenlik modülünü paralel çalıştırır.
        </p>
        <button
          type="submit"
          disabled={isPending}
          className="rounded-full bg-[var(--accent)] px-6 py-3 text-sm font-semibold text-white transition hover:bg-[var(--accent-dark)] disabled:opacity-60"
        >
          {isPending ? "Scan başlatılıyor..." : "Scan başlat"}
        </button>
      </div>

      {error ? <p className="mt-4 text-sm text-red-700">{error}</p> : null}
    </form>
  );
}
