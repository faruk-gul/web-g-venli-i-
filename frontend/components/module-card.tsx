import { ModuleResult } from "../lib/api";

const gradeStyle: Record<string, string> = {
  "A+": "bg-emerald-100 text-emerald-900",
  A: "bg-emerald-50 text-emerald-800",
  B: "bg-amber-100 text-amber-900",
  C: "bg-orange-100 text-orange-900",
  D: "bg-red-100 text-red-900",
  F: "bg-red-200 text-red-950"
};

export function ModuleCard({ module }: { module: ModuleResult }) {
  return (
    <article className="glass rounded-[28px] p-5">
      <div className="mb-4 flex items-center justify-between gap-3">
        <div>
          <p className="text-xs uppercase tracking-[0.3em] text-stone-500">{module.name}</p>
          <h3 className="text-xl font-semibold capitalize">{module.name} modülü</h3>
        </div>
        <span className={`rounded-full px-3 py-1 text-sm font-semibold ${gradeStyle[module.grade] || "bg-stone-200 text-stone-800"}`}>
          {module.grade}
        </span>
      </div>
      <p className="mb-4 text-sm text-stone-700">{module.summary}</p>
      {module.findings?.length ? (
        <ul className="space-y-2 text-sm text-stone-700">
          {module.findings.map((finding) => (
            <li key={finding}>• {finding}</li>
          ))}
        </ul>
      ) : null}
    </article>
  );
}

