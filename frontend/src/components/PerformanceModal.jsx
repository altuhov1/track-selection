import { useState, useEffect, useCallback } from 'react'
import { fetchUserPreferences } from '../services/profile'
import { SUBJECTS, SKILLS } from '../data/trackStyles'

export default function PerformanceModal({ onClose }) {
  const [loading, setLoading] = useState(true)
  const [error, setError]     = useState('')
  const [data, setData]       = useState(null)

  const handleClose = useCallback(() => onClose(), [onClose])

  useEffect(() => {
    const onKey = (e) => { if (e.key === 'Escape') handleClose() }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [handleClose])

  useEffect(() => {
    document.body.style.overflow = 'hidden'
    return () => { document.body.style.overflow = '' }
  }, [])

  useEffect(() => {
    fetchUserPreferences()
      .then(setData)
      .catch((e) => setError(e.message || 'Не удалось загрузить данные'))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div className="modal-overlay" onClick={(e) => e.target === e.currentTarget && handleClose()}>
      <div className="modal modal--wide" role="dialog" aria-modal="true">
        <div className="modal-header">
          <h2>Успеваемость</h2>
          <button className="modal-close" onClick={handleClose} aria-label="Закрыть">×</button>
        </div>

        {loading ? (
          <div className="profile-loading">Загрузка…</div>
        ) : error ? (
          <div className="form-error">{error}</div>
        ) : (
          <div className="charts">
            <Chart
              title="Оценки"
              max={5}
              color="#2563EB"
              items={SUBJECTS.map(s => ({
                label: s.label,
                value: data?.grades?.[s.key] ?? 0,
              }))}
            />
            <Chart
              title="Навыки"
              max={10}
              color="#7C3AED"
              items={SKILLS.map(s => ({
                label: s.label,
                value: data?.skills?.[s.key] ?? 0,
              }))}
            />
          </div>
        )}
      </div>
    </div>
  )
}

function Chart({ title, items, max, color }) {
  return (
    <section className="chart-section">
      <h3 className="chart-title">{title}</h3>
      <div className="chart" style={{ '--chart-color': color }}>
        <div className="chart-axis">
          {Array.from({ length: max + 1 }).map((_, i) => (
            <span key={i} className="chart-axis-tick">{max - i}</span>
          ))}
        </div>
        <div className="chart-bars">
          {items.map((it, i) => {
            const pct = Math.max(0, Math.min(1, it.value / max)) * 100
            return (
              <div key={i} className="chart-col">
                <div className="chart-bar-wrap">
                  <div
                    className="chart-bar"
                    style={{ height: `${pct}%` }}
                    title={`${it.label}: ${it.value}`}
                  >
                    <span className="chart-value">{it.value}</span>
                  </div>
                </div>
                <span className="chart-label">{it.label}</span>
              </div>
            )
          })}
        </div>
      </div>
    </section>
  )
}
