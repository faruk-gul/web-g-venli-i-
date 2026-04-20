"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { getScan, ScanRecord, streamScan } from "../../../lib/api";
import { ModuleCard } from "../../../components/module-card";
import { RadarPlaceholder } from "../../../components/radar-placeholder";

type PageProps = {
  params: {
    id: string;
  };
};

export default function ScanDetailPage({ params }: PageProps) {
  const [record, setRecord] = useState<ScanRecord | null>(null);
  const [events, setEvents] = useState<string[]>([]);
  const [error, setError] = useState("");

  useEffect(() => {
    let active = true;
    let eventSource: EventSource | null = null;

    async function load() {
      try {
        const initial = await getScan(params.id);
        if (!active) return;
        setRecord(initial);

        eventSource = streamScan(params.id, async (event: any) => {
          if (!active) return;
          if (event?.message) {
            setEvents((current) => [event.message, ...current].slice(0, 8));
          }
          const next = await getScan(params.id);
          if (active) {
            setRecord(next);
          }
        });
      } catch (loadError) {
        setError(loadError instanceof Error ? loadError.message : "Yükleme hatası");
      }
    }

    void load();

    return () => {
      active = false;
      eventSource?.close();
    };
  }, [params.id]);

  if (error) {
    return <main className="shell px-4 py-10 text-red-700">{error}</main>;
  }

  if (!record) {
    return <main className="shell px-4 py-10">Rapor yükleniyor...</main>;
  }

  return (
    <main className="px-4 py-8 md:px-8 md:py-12">
      <div className="shell">
        <div className="mb-6 flex items-center justify-between gap-4">
          <div>
            <p className="text-xs uppercase tracking-[0.35em] text-stone-500">Scan Report</p>
            <h1 className="mt-2 text-3xl font-semibold md:text-5xl">{record.target_url}</h1>
          </div>
          <Link href="/" className="rounded-full border border-stone-300 px-4 py-2 text-sm">
            Yeni scan
          </Link>
        </div>

        <section className="mb-8 grid gap-6 md:grid-cols-[1.2fr_0.8fr]">
          <article className="glass rounded-[30px] p-6 md:p-8">
            <div className="mb-5 flex items-center justify-between">
              <div>
                <p className="text-xs uppercase tracking-[0.3em] text-stone-500">Genel skor</p>
                <h2 className="mt-2 text-4xl font-semibold">{record.grade}</h2>
              </div>
              <div className="text-right">
                <p className="text-xs uppercase tracking-[0.3em] text-stone-500">İlerleme</p>
                <p className="mt-2 text-3xl font-semibold">{record.progress}%</p>
              </div>
            </div>
            <div className="mb-4 h-3 rounded-full bg-stone-200">
              <div className="h-3 rounded-full bg-[var(--accent)] transition-all" style={{ width: `${record.progress}%` }} />
            </div>
            <p className="text-sm leading-7 text-stone-700">{record.summary}</p>
            {record.error ? <p className="mt-4 text-sm text-red-700">{record.error}</p> : null}
          </article>

          <article className="glass rounded-[30px] p-6 md:p-8">
            <p className="text-xs uppercase tracking-[0.3em] text-stone-500">Canlı olaylar</p>
            <div className="mt-4 space-y-3">
              {events.length ? (
                events.map((item, index) => (
                  <div key={`${item}-${index}`} className="rounded-2xl border border-stone-200 bg-white/70 px-4 py-3 text-sm">
                    {item}
                  </div>
                ))
              ) : (
                <div className="rounded-2xl border border-dashed border-stone-300 px-4 py-3 text-sm text-stone-600">
                  Olay akışı bekleniyor.
                </div>
              )}
            </div>
            <a
              href={`${process.env.NEXT_PUBLIC_API_BASE || "http://localhost:8080"}/api/scan/${record.id}/report.pdf`}
              className="mt-6 inline-block text-sm font-semibold text-[var(--accent)]"
            >
              PDF raporunu indir
            </a>
          </article>
        </section>

        <section className="mb-8">
          <RadarPlaceholder modules={record.modules.map((module) => ({ name: module.name, grade: module.grade }))} />
        </section>

        <section className="mb-8">
          <div className="mb-4 flex items-center justify-between">
            <div>
              <p className="text-xs uppercase tracking-[0.3em] text-stone-500">Tavsiyeler</p>
              <h2 className="text-2xl font-semibold">Öncelikli aksiyonlar</h2>
            </div>
          </div>
          <div className="glass rounded-[30px] p-6">
            {record.recommendations.length ? (
              <ul className="space-y-3 text-sm text-stone-700">
                {record.recommendations.map((item) => (
                  <li key={item}>• {item}</li>
                ))}
              </ul>
            ) : (
              <p className="text-sm text-stone-600">Henüz tavsiye üretilmedi.</p>
            )}
          </div>
        </section>

        <section className="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
          {record.modules.map((module) => (
            <ModuleCard key={module.name} module={module} />
          ))}
        </section>
      </div>
    </main>
  );
}

