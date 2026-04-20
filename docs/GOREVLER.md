# Gorevler

## PDF'deki Finale Gorevleri ile Esleme

### F01 Proje Kurulumu + Git Workflow + Docker Compose

- `secscan/` klasor yapisi olusturuldu.
- `docker-compose.yml` eklendi.
- `.github/workflows` dosyalari hazirlandi.

### F02 Backend API Iskeleti

- Gin tabanli `main.go`
- `scanner.Registry`
- `scanner.Service`
- API endpointleri

### F03 Port Scanner Modulu

- Temel TCP connect kontrolu eklendi.
- SSRF guard yazildi.

### F04 HTTP / Security Headers Analiz

- Beklenen baslik listesi ile analiz eklendi.

### F05 TLS / SSL Auditor

- HTTPS hedefleri icin TLS handshake denemesi eklendi.

### F06 Dir Fuzzer + XSS + SQLi

- Uc modul icin baslangic heuristikleri eklendi.
- Gercek tarama ve payload varyasyonlari sonraki iterasyona birakildi.

### F07 SBOM / CVE Analiz

- CVE modulu icin teknoloji algilama placeholder'i hazirlandi.
- OSV/NVD entegrasyonu sonraki adim olarak not edildi.

### F08 Frontend URL Input + SSE

- Ana sayfa formu
- Rapor sayfasi
- SSE event stream baglantisi

### F09 Dashboard + Radar Chart + PDF Export

- Dashboard alanlari kuruldu.
- Radar chart icin placeholder component eklendi.
- PDF export endpoint'i placeholder cikti donduruyor.

### F10 Deploy + Demo

- Henüz canli deploy yapilmadi.
- Dockerfile ve compose yapisi deploya temel hazirlik sagliyor.

## Kisa Yol Haritasi

1. Go ortaminda `go mod tidy` ile backend modullerini indir.
2. Frontend icin `npm install` ve `npm run dev` calistir.
3. Modulleri tek tek gercek test mantigi ile derinlestir.
4. Ekran goruntuleri alip README'ye ekle.

