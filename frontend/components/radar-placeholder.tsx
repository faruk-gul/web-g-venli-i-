type RadarPlaceholderProps = {
  modules: Array<{ name: string; grade: string }>;
};

export function RadarPlaceholder({ modules }: RadarPlaceholderProps) {
  return (
    <div className="glass rounded-[28px] p-6">
      <div className="mb-4 flex items-center justify-between">
        <div>
          <p className="text-xs uppercase tracking-[0.3em] text-stone-500">Dashboard</p>
          <h3 className="text-xl font-semibold">Radar chart placeholder</h3>
        </div>
        <span className="rounded-full bg-stone-900 px-3 py-1 text-xs text-white">F09-ready</span>
      </div>
      <div className="grid gap-3 md:grid-cols-2">
        {modules.map((module) => (
          <div key={module.name} className="rounded-2xl border border-stone-200 bg-white/60 px-4 py-3">
            <p className="text-sm font-semibold capitalize">{module.name}</p>
            <p className="text-sm text-stone-600">Skor: {module.grade}</p>
          </div>
        ))}
      </div>
    </div>
  );
}

