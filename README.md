# SecScan

SecScan, Seri 3 Web Guvenligi finali icin hazirlanmis full-stack bir guvenlik tarama projesidir. Kullanici bir URL girer, backend tarafinda 7 guvenlik modulu paralel calisir, frontend tarafinda SSE ile canli ilerleme ve rapor ekrani gorulur.

## Stack

- Backend: Go 1.22 + Gin
- Frontend: Next.js 14 + TypeScript + TailwindCSS
- Orkestrasyon: Docker Compose
- Guvenlik otomasyonu: GitHub Actions, Semgrep, Trivy

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
```

## Hazir Ozellikler

- `/health`, `POST /api/scan`, `GET /api/scan/:id`, `GET /api/scan/:id/stream`, `GET /api/scan/:id/report.pdf`
- In-memory scan store ve async tarama akisi
- SSRF guard: private/link-local/loopback CIDR bloklari reddedilir
- 7 scanner modulu icin temel registry ve sonuc modeli
- Frontend ana sayfa, scan baslatma akisi ve rapor ekrani
- PDF raporu icin indirilebilir placeholder cikti
- CI ve security workflow dosyalari

## Calistirma

### Backend

```bash
cd backend
go mod tidy
go run .
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

### Docker Compose

```bash
docker compose up --build
```

## API Ornekleri

### Scan baslat

```bash
curl -X POST http://localhost:8080/api/scan \
  -H "Content-Type: application/json" \
  -d '{"target_url":"https://example.com"}'
```

### Scan sonucu

```bash
curl http://localhost:8080/api/scan/scan-0001
```

## Sonraki Adimlar

1. Backend tarafinda modulleri gercek tarama mantigi ile derinlestir.
2. Radar chart ve gercek PDF export katmanini ekle.
3. NVD veya OSV baglantisini tamamlayip CVE sonucunu zenginlestir.
4. Testler ve production deploy akisini ekle.

## Not

Bu iskelet, PDF'deki F01-F09 gereksinimlerini hizli baslangic icin somutlastirir. Go araci bu ortamda kurulu olmadigi icin backend derlemesi burada dogrulanamadi; dosya yapisi ve kaynak kod bu varsayimla hazirlandi.

