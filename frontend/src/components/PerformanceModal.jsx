import { useState, useEffect, useCallback, useMemo } from 'react'
import { fetchUserPreferences, fetchRecommendations } from '../services/profile'
import { fetchAllTracks } from '../services/tracks'
import { SUBJECTS, SKILLS, getTrackStyle } from '../data/trackStyles'
import TrackDetailsModal from './TrackDetailsModal'

export default function PerformanceModal({ onClose }) {
  const [loading, setLoading] = useState(true)
  const [error, setError]     = useState('')
  const [prefs, setPrefs]     = useState(null)
  const [recs, setRecs]       = useState([])
  const [recsError, setRecsError] = useState('')
  const [tracks, setTracks]   = useState([])
  const [selectedTrack, setSelectedTrack] = useState(null)

  const handleClose = useCallback(() => onClose(), [onClose])

  useEffect(() => {
    const onKey = (e) => { if (e.key === 'Escape' && !selectedTrack) handleClose() }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [handleClose, selectedTrack])

  useEffect(() => {
    document.body.style.overflow = 'hidden'
    return () => { document.body.style.overflow = '' }
  }, [])

  useEffect(() => {
    let cancelled = false
    Promise.all([
      fetchUserPreferences().catch((e) => { throw e }),
      fetchAllTracks().catch(() => []),
      fetchRecommendations().catch((e) => ({ __err: e })),
    ])
      .then(([p, tr, rec]) => {
        if (cancelled) return
        setPrefs(p)
        setTracks(Array.isArray(tr) ? tr : [])
        if (rec && rec.__err) {
          setRecsError(rec.__err.message || 'Не удалось загрузить рекомендации')
        } else {
          setRecs(Array.isArray(rec?.recommendations) ? rec.recommendations : [])
        }
      })
      .catch((e) => !cancelled && setError(e.message || 'Не удалось загрузить данные'))
      .finally(() => !cancelled && setLoading(false))
    return () => { cancelled = true }
  }, [])

  const trackById = useMemo(() => {
    const map = {}
    tracks.forEach(t => { map[t.id] = t })
    return map
  }, [tracks])

  return (
    <div className={`modal-overlay${selectedTrack ? ' modal-overlay--no-bg' : ''}`} onClick={(e) => e.target === e.currentTarget && handleClose()}>
      <div className="modal modal--wide" role="dialog" aria-modal="true">
        <div className="modal-header">
          <h2>Успеваемость и рекомендации</h2>
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
                value: prefs?.grades?.[s.key] ?? 0,
              }))}
            />
            <Chart
              title="Навыки"
              max={10}
              color="#7C3AED"
              items={SKILLS.map(s => ({
                label: s.label,
                value: prefs?.skills?.[s.key] ?? 0,
              }))}
            />

            <section className="chart-section">
              <h3 className="chart-title">Рекомендованные треки</h3>
              {recsError ? (
                <div className="recs-empty">{recsError}</div>
              ) : recs.length === 0 ? (
                <div className="recs-empty">Рекомендации пока недоступны.</div>
              ) : (
                <ul className="recs-list">
                  {recs.map((r) => {
                    const track = trackById[r.track_id]
                    const style = getTrackStyle(track?.type)
                    const pct = Math.round((r.score ?? 0) * 100)
                    return (
                      <li key={r.track_id}>
                        <button
                          type="button"
                          className="rec-item"
                          onClick={() => track && setSelectedTrack(track)}
                          disabled={!track}
                          title={track ? 'Открыть описание трека' : 'Трек недоступен'}
                        >
                          <span className="rec-rank">#{r.rank}</span>
                          <span className="rec-icon" style={{ background: style.color }}>
                            <span dangerouslySetInnerHTML={{ __html: style.icon(style.shapeColor) }} />
                          </span>
                          <span className="rec-main">
                            <span className="rec-name">{r.track_name || track?.name || 'Трек'}</span>
                            <span className="rec-category">{style.label}</span>
                          </span>
                          <span className="rec-score">
                            <span className="rec-score-bar">
                              <span className="rec-score-fill" style={{ width: `${pct}%` }} />
                            </span>
                            <span className="rec-score-val">{pct}%</span>
                          </span>
                        </button>
                      </li>
                    )
                  })}
                </ul>
              )}
            </section>
          </div>
        )}
      </div>

      {selectedTrack && (
        <TrackDetailsModal track={selectedTrack} onClose={() => setSelectedTrack(null)} />
      )}
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
