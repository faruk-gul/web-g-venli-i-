# Mimari

## Genel Bakis

SecScan iki katmandan olusur:

- `backend/`: Tarama is motoru, API ve olay yayini
- `frontend/`: Kullanici arayuzu, scan baslatma ve rapor goruntuleme

## Backend Akisi

1. `POST /api/scan` istegi gelir.
2. `ValidateTarget` SSRF guard kontrolu yapar.
3. `Service.Start` kaydi olusturur ve async taramayi baslatir.
4. `Registry`, secilen modulleri resolve eder.
5. Her scanner goroutine icinde calisir.
6. Sonuclar in-memory store'a yazilir.
7. SSE subscriber'lara progress event gonderilir.

## Frontend Akisi

1. Ana sayfadaki form scan'i baslatir.
2. Backend `scan_id` dondurur.
3. Kullanici `/scan/[id]` ekranina yonlenir.
4. Sayfa ilk sonucu fetch eder.
5. SSE event geldikce veri tekrar cekilir.
6. Moduller, tavsiyeler ve genel not ekranda gosterilir.

## Dizinler

- `backend/internal/api`: HTTP route ve handler'lar
- `backend/internal/scanner`: Domain modelleri, servis ve moduller
- `frontend/app`: App Router sayfalari
- `frontend/components`: UI bilesenleri
- `frontend/lib`: API yardimcilari

## Tasarim Notlari

- In-memory store, gelistirme hizini artirmak icin secildi.
- Gercek uretim surumunde Redis/PostgreSQL eklenebilir.
- PDF export su an placeholder uretir; ileride gercek bir PDF kutuphanesi ile degistirilmeli.

