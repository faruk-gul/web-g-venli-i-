import Link from "next/link";
import { ScanLauncher } from "../components/scan-launcher";

const modules = [
  "ports",
  "headers",
  "tls",
  "fuzz",
  "xss",
  "sqli",
  "cve"
];

export default function HomePage() {
  return (
    <main className="px-4 py-8 md:px-8 md:py-12">
      <div className="shell">
        <section className="mb-8 rounded-[36px] border border-[var(--line)] bg-[rgba(255,252,247,0.66)] p-8 shadow-[0_24px_70px_rgba(48,30,17,0.12)] md:p-12">
          <p className="mb-4 text-xs uppercase tracking-[0.45em] text-stone-500">Seri 3 Final Project</p>
          <div className="grid gap-8 md:grid-cols-[1.2fr_0.8fr] md:items-end">
            <div>
              <h1 className="max-w-3xl text-4xl font-semibold leading-tight md:text-6xl">
                SecScan ile bir URL'i tara, güvenlik sinyallerini tek ekranda toparla.
              </h1>
              <p className="mt-5 max-w-2xl text-lg leading-8 text-stone-700">
                Go + Gin backend, Next.js dashboard, SSE canlı ilerleme ve SSRF guard ile hazırlanmış proje iskeleti.
              </p>
            </div>
            <div className="glass rounded-[32px] p-6">
              <p className="text-sm uppercase tracking-[0.3em] text-stone-500">Modül seti</p>
              <div className="mt-4 flex flex-wrap gap-2">
                {modules.map((module) => (
                  <span key={module} className="rounded-full border border-stone-300 px-3 py-2 text-sm capitalize">
                    {module}
                  </span>
                ))}
              </div>
              <Link href="#launch" className="mt-6 inline-block text-sm font-semibold text-[var(--accent)]">
                Scan başlat bölümüne git
              </Link>
            </div>
          </div>
        </section>

        <section id="launch" className="mb-8">
          <ScanLauncher />
        </section>

        <section className="grid gap-6 md:grid-cols-3">
          <article className="glass rounded-[28px] p-6">
            <p className="text-xs uppercase tracking-[0.3em] text-stone-500">F03-F07</p>
            <h2 className="mt-2 text-2xl font-semibold">Backend çekirdeği</h2>
            <p className="mt-3 text-sm leading-7 text-stone-700">
              Scanner registry, in-memory scan store, SSE event yayını ve SSRF koruması ile kuruldu.
            </p>
          </article>

          <article className="glass rounded-[28px] p-6">
            <p className="text-xs uppercase tracking-[0.3em] text-stone-500">F08-F09</p>
            <h2 className="mt-2 text-2xl font-semibold">Canlı rapor akışı</h2>
            <p className="mt-3 text-sm leading-7 text-stone-700">
              Kullanıcı scan başlatır, rapor ekranında ilerlemeyi izler ve özet tavsiyeleri görür.
            </p>
          </article>

          <article className="glass rounded-[28px] p-6">
            <p className="text-xs uppercase tracking-[0.3em] text-stone-500">Teslim odaklı</p>
            <h2 className="mt-2 text-2xl font-semibold">README + teknik döküman</h2>
            <p className="mt-3 text-sm leading-7 text-stone-700">
              Proje yapısı, yapılacaklar, mimari kararlar ve geliştirme adımları ayrıca yazıldı.
            </p>
          </article>
        </section>
      </div>
    </main>
  );
}

