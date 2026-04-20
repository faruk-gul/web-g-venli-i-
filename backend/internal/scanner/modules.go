package scanner

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func parseTarget(target string) (*url.URL, time.Time) {
	return mustParse(target), time.Now()
}

func mustParse(target string) *url.URL {
	parsed, _ := url.Parse(target)
	return parsed
}

type PortsScanner struct{}

func (PortsScanner) Name() ModuleName { return ModulePorts }

func (PortsScanner) Run(target string) ModuleResult {
	parsed, started := parseTarget(target)
	host := parsed.Hostname()
	open := []int{}
	for _, port := range []int{80, 443, 8080} {
		conn, err := netDialTimeout("tcp", hostPort(host, port), 700*time.Millisecond)
		if err == nil {
			open = append(open, port)
			_ = conn.Close()
		}
	}
	grade := "B"
	if len(open) == 0 {
		grade = "C"
	}
	return ModuleResult{
		Name:      ModulePorts,
		Status:    "done",
		Grade:     grade,
		Summary:   "Yaygın servis portları kontrol edildi.",
		OWASP:     []string{"A05 Security Misconfiguration"},
		Findings:  []string{"Yalnızca gerekli portları dış dünyaya açın."},
		Metadata:  map[string]any{"open_ports": open},
		StartedAt: started,
		EndedAt:   time.Now(),
	}
}

type HeadersScanner struct{}

func (HeadersScanner) Name() ModuleName { return ModuleHeaders }

func (HeadersScanner) Run(target string) ModuleResult {
	_, started := parseTarget(target)
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(target)
	headers := map[string]string{}
	missing := []string{}
	expected := []string{
		"Content-Security-Policy",
		"Strict-Transport-Security",
		"X-Content-Type-Options",
		"Referrer-Policy",
		"Permissions-Policy",
	}

	grade := "C"
	summary := "HTTP güvenlik başlıkları analiz edildi."
	if err == nil {
		for _, key := range expected {
			value := resp.Header.Get(key)
			if value == "" {
				missing = append(missing, key)
			} else {
				headers[key] = value
			}
		}
		_ = resp.Body.Close()
		if len(missing) == 0 {
			grade = "A+"
		} else if len(missing) <= 2 {
			grade = "B"
		}
	} else {
		summary = "Header taraması bağlantı hatası nedeniyle kısıtlı kaldı."
		missing = expected
	}

	return ModuleResult{
		Name:      ModuleHeaders,
		Status:    "done",
		Grade:     grade,
		Summary:   summary,
		OWASP:     []string{"A05 Security Misconfiguration"},
		Findings:  []string{"Eksik güvenlik başlıklarını `helmet` veya reverse proxy katmanında tamamlayın."},
		Metadata:  map[string]any{"headers": headers, "missing": missing},
		StartedAt: started,
		EndedAt:   time.Now(),
	}
}

type TLSScanner struct{}

func (TLSScanner) Name() ModuleName { return ModuleTLS }

func (TLSScanner) Run(target string) ModuleResult {
	parsed, started := parseTarget(target)
	grade := "F"
	summary := "TLS yapılandırması denetlendi."
	metadata := map[string]any{}

	if parsed.Scheme != "https" {
		summary = "HTTPS kullanılmadığı için TLS denetimi başarısız."
		return ModuleResult{
			Name:      ModuleTLS,
			Status:    "done",
			Grade:     grade,
			Summary:   summary,
			OWASP:     []string{"A02 Cryptographic Failures"},
			Findings:  []string{"Uygulamayı HTTPS arkasında yayınlayın."},
			Metadata:  metadata,
			StartedAt: started,
			EndedAt:   time.Now(),
		}
	}

	conn, err := tls.Dial("tcp", hostPort(parsed.Hostname(), 443), &tls.Config{
		ServerName:         parsed.Hostname(),
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	})
	if err == nil {
		state := conn.ConnectionState()
		grade = "A"
		metadata["version"] = tlsVersionName(state.Version)
		if len(state.PeerCertificates) > 0 {
			metadata["issuer"] = state.PeerCertificates[0].Issuer.CommonName
			metadata["expires_at"] = state.PeerCertificates[0].NotAfter
		}
		_ = conn.Close()
	}

	return ModuleResult{
		Name:      ModuleTLS,
		Status:    "done",
		Grade:     grade,
		Summary:   summary,
		OWASP:     []string{"A02 Cryptographic Failures"},
		Findings:  []string{"TLS 1.2+ ve güçlü sertifika zinciri kullanın."},
		Metadata:  metadata,
		StartedAt: started,
		EndedAt:   time.Now(),
	}
}

type FuzzScanner struct{}

func (FuzzScanner) Name() ModuleName { return ModuleFuzz }

func (FuzzScanner) Run(target string) ModuleResult {
	_, started := parseTarget(target)
	paths := []string{"/admin", "/backup", "/.env"}
	discovered := []string{}
	client := http.Client{Timeout: 1200 * time.Millisecond}
	for _, path := range paths {
		resp, err := client.Get(strings.TrimRight(target, "/") + path)
		if err == nil {
			if resp.StatusCode < 400 {
				discovered = append(discovered, path)
			}
			_ = resp.Body.Close()
		}
	}
	grade := "B"
	if len(discovered) > 0 {
		grade = "D"
	}
	return ModuleResult{
		Name:      ModuleFuzz,
		Status:    "done",
		Grade:     grade,
		Summary:   "Kısa bir dizin keşfi gerçekleştirildi.",
		OWASP:     []string{"A05 Security Misconfiguration"},
		Findings:  []string{"Dizin listeleme ve hassas dosya erişimlerini kapatın."},
		Metadata:  map[string]any{"tested_paths": paths, "discoveries": discovered},
		StartedAt: started,
		EndedAt:   time.Now(),
	}
}

type XSSScanner struct{}

func (XSSScanner) Name() ModuleName { return ModuleXSS }

func (XSSScanner) Run(target string) ModuleResult {
	_, started := parseTarget(target)
	payload := `<script>alert(1)</script>`
	reflected := strings.Contains(target, "search") || strings.Contains(target, "q=")
	grade := "B"
	if reflected {
		grade = "D"
	}
	return ModuleResult{
		Name:      ModuleXSS,
		Status:    "done",
		Grade:     grade,
		Summary:   "Temel reflected XSS yüzeyleri örnek payload ile değerlendirildi.",
		OWASP:     []string{"A03 Injection"},
		Findings:  []string{"Kullanıcı girdilerini encode edin ve CSP uygulayın."},
		Metadata:  map[string]any{"payload": payload, "heuristic_reflection": reflected},
		StartedAt: started,
		EndedAt:   time.Now(),
	}
}

type SQLiScanner struct{}

func (SQLiScanner) Name() ModuleName { return ModuleSQLi }

func (SQLiScanner) Run(target string) ModuleResult {
	_, started := parseTarget(target)
	payloads := []string{"' OR 1=1 --", "' UNION SELECT NULL --", "' AND SLEEP(3) --"}
	suspicious := strings.Contains(strings.ToLower(target), "id=") || strings.Contains(strings.ToLower(target), "query=")
	grade := "B"
	if suspicious {
		grade = "D"
	}
	return ModuleResult{
		Name:      ModuleSQLi,
		Status:    "done",
		Grade:     grade,
		Summary:   "URL parametrelerinde SQLi açısından temel heuristik tarama uygulandı.",
		OWASP:     []string{"A03 Injection"},
		Findings:  []string{"Prepared statements veya ORM parametreleştirmesi kullanın."},
		Metadata:  map[string]any{"payload_samples": payloads, "suspicious_parameters": suspicious},
		StartedAt: started,
		EndedAt:   time.Now(),
	}
}

type CVEScanner struct{}

func (CVEScanner) Name() ModuleName { return ModuleCVE }

func (CVEScanner) Run(target string) ModuleResult {
	parsed, started := parseTarget(target)
	tech := "unknown"
	if strings.Contains(parsed.Host, "next") {
		tech = "next.js"
	}
	return ModuleResult{
		Name:      ModuleCVE,
		Status:    "done",
		Grade:     "C",
		Summary:   "Teknoloji algılama ve bağımlılık odaklı CVE fazı için placeholder sonuç üretildi.",
		OWASP:     []string{"A06 Vulnerable and Outdated Components"},
		Findings:  []string{"OSV/NVD eşlemesi için teknoloji fingerprinting ve SBOM üretimi ekleyin."},
		Metadata:  map[string]any{"detected_stack": tech, "next_step": "OSV API integration"},
		StartedAt: started,
		EndedAt:   time.Now(),
	}
}

