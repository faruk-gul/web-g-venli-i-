# SecScan

SecScan, Seri 3 Web Guvenligi finali icin hazirlanmis full-stack bir guvenlik tarama projesidir. Kullanici bir URL girer, backend tarafinda 7 guvenlik modulu paralel calisir, frontend tarafinda SSE ile canli ilerleme ve rapor ekrani gorulur.

## Ogrenci Bilgileri

- Ogrenci: Faruk Gul
- Okul No: 24080410022
- Ders: BMU1208 Web Tabanli Programlama
- Proje: Seri 3 Web Guvenligi Final Projesi

## Stack

- Backend: Go 1.22 + Gin
- Frontend: Next.js 14 + TypeScript + TailwindCSS
- Orkestrasyon: Docker Compose
- Guvenlik otomasyonu: GitHub Actions, Semgrep, Trivy

## Proje Amaci

Bu proje, kullanicinin verdigi bir hedef URL uzerinde temel guvenlik kontrolleri yapmayi amaclar. Sistem, farkli guvenlik modullerini paralel sekilde calistirir ve sonuclari kullaniciya rapor halinde sunar.

## Hazir Ozellikler

- `POST /api/scan` ile yeni scan baslatma
- `GET /api/scan/:id` ile scan sonucu goruntuleme
- `GET /api/scan/:id/stream` ile SSE canli ilerleme akisi
- `GET /api/scan/:id/report.pdf` ile PDF rapor indirme
- `GET /health` ile servis durum kontrolu
- SSRF guard ile private ve local IP bloklarini engelleme
- 7 modullu scanner yapisi
- Frontend ana sayfa, scan baslatma formu ve rapor ekrani
- Docker ve GitHub Actions dosyalari

## Kullanilan Moduller

- `ports` : Port tarama
- `headers` : Security header analizi
- `tls` : TLS/SSL kontrolu
- `fuzz` : Dizin ve path tarama
- `xss` : XSS kontrolu
- `sqli` : SQL Injection kontrolu
- `cve` : Teknoloji ve CVE analizi

## Proje Yapisi

```text
secscan/
|-- backend/
|   |-- main.go
|   \-- internal/
|       |-- api/
|       \-- scanner/
|-- frontend/
|   |-- app/
|   |-- components/
|   \-- lib/
|-- docs/
|   |-- TEKNIK-SARTNAME.md
|   |-- MIMARI.md
|   \-- GOREVLER.md
|-- docker-compose.yml
\-- .github/workflows/
