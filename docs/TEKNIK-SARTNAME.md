# Teknik Sartname

## Amac

SecScan, kullanicidan alinan bir hedef URL uzerinde temel web guvenligi sinyallerini analiz eden ve sonucunu not, bulgu ve oneriler ile raporlayan bir platformdur.

## Fonksiyonel Gereksinimler

1. Sistem kullanicidan bir URL alir.
2. Tarama asenkron baslar ve bir `scan_id` dondurur.
3. Durum bilgisi JSON olarak okunabilir.
4. SSE ile canli progress akisi saglanir.
5. PDF rapor cikisi sunulur.
6. Tarama oncesi SSRF guard uygulanir.
7. Asagidaki moduller desteklenir:
   - ports
   - headers
   - tls
   - fuzz
   - xss
   - sqli
   - cve

## Teknik Gereksinimler

- Backend dili: Go 1.22+
- HTTP framework: Gin
- Frontend: Next.js 14 App Router
- Stil: TailwindCSS
- Container: Docker Compose
- CI: GitHub Actions
- Security scan: Semgrep + Trivy

## Guvenlik Gereksinimleri

- Private IP bloklari taranmamali.
- Sadece `http` ve `https` kabul edilmeli.
- Tarama loglari hassas bilgi sızdirmamali.
- Gelecekte rate limiting ve auth katmani eklenmeli.

## Basari Kriterleri

- `POST /api/scan` isler.
- `GET /api/scan/:id` sonuc dondurur.
- `GET /api/scan/:id/stream` frontend'e olay gonderir.
- Frontend scan ekranina yonlenir ve ilerlemeyi gosterir.

