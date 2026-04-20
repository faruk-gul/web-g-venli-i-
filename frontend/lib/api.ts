export type ModuleResult = {
  name: string;
  status: string;
  grade: string;
  summary: string;
  findings?: string[];
  metadata?: Record<string, unknown>;
};

export type ScanRecord = {
  id: string;
  target_url: string;
  status: string;
  grade: string;
  progress: number;
  summary: string;
  recommendations: string[];
  error?: string;
  modules: ModuleResult[];
};

const API_BASE = process.env.NEXT_PUBLIC_API_BASE || "http://localhost:8080";

export async function startScan(targetURL: string) {
  const response = await fetch(`${API_BASE}/api/scan`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ target_url: targetURL })
  });

  if (!response.ok) {
    throw new Error("Scan baslatilamadi");
  }

  return response.json() as Promise<{ scan_id: string; status: string; target: string }>;
}

export async function getScan(id: string) {
  const response = await fetch(`${API_BASE}/api/scan/${id}`, {
    cache: "no-store"
  });

  if (!response.ok) {
    throw new Error("Scan bulunamadi");
  }

  return response.json() as Promise<ScanRecord>;
}

export function streamScan(id: string, onMessage: (event: unknown) => void) {
  const stream = new EventSource(`${API_BASE}/api/scan/${id}/stream`);
  stream.onmessage = (event) => {
    try {
      onMessage(JSON.parse(event.data));
    } catch {
      onMessage(event.data);
    }
  };
  return stream;
}

