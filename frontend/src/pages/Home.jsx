import { useState, useEffect, useMemo, useCallback } from 'react'
import TrackCard from '../components/TrackCard'
import TrackDetailsModal from '../components/TrackDetailsModal'
import TrackFormModal from '../components/TrackFormModal'
import PerformanceModal from '../components/PerformanceModal'
import { fetchAllTracks, deleteTrack } from '../services/tracks'
import { fetchProfileCompletion } from '../services/profile'
import { DIFFICULTY_LABELS } from '../data/trackStyles'

export default function Home({ user, onOpenProfile, onOpenLogin }) {
  const isAdmin   = user?.role === 'admin'
  const isStudent = user?.role === 'student'

  const [tracks, setTracks]     = useState([])
  const [loading, setLoading]   = useState(true)
  const [error, setError]       = useState('')

  const [search, setSearch]     = useState('')
  const [filterDiff, setFilterDiff] = useState(new Set())

  const [detailsTrack, setDetailsTrack]   = useState(null)
  const [performanceOpen, setPerformanceOpen] = useState(false)
  const [performanceChecking, setPerformanceChecking] = useState(false)
  const [formTrack, setFormTrack]         = useState(null) // track | 'new' | null

  const loadTracks = useCallback(() => {
    setLoading(true)
    fetchAllTracks()
      .then(list => setTracks(Array.isArray(list) ? list : []))
      .catch(e => setError(e.message || 'Не удалось загрузить треки'))
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => { loadTracks() }, [loadTracks])

  const filtered = useMemo(() => {
    return tracks.filter(t => {
      if (search && !t.name?.toLowerCase().includes(search.toLowerCase())) return false
      if (filterDiff.size && !filterDiff.has(t.difficulty)) return false
      return true
    })
  }, [tracks, search, filterDiff])

  function toggleDiff(v) {
    setFilterDiff(prev => {
      const next = new Set(prev)
      next.has(v) ? next.delete(v) : next.add(v)
      return next
    })
  }

  async function handlePerformanceClick() {
    if (!user) { onOpenLogin?.(); return }
    setPerformanceChecking(true)
    try {
      const { is_complete } = await fetchProfileCompletion()
      if (is_complete) setPerformanceOpen(true)
      else onOpenProfile?.()
    } catch {
      // если профиль не инициализирован — откроем анкету
      onOpenProfile?.()
    } finally {
      setPerformanceChecking(false)
    }
  }

  async function handleDelete(track) {
    if (!confirm(`Удалить трек «${track.name}»?`)) return
    try {
      await deleteTrack(track.id)
      loadTracks()
    } catch (e) {
      alert(e.message || 'Не удалось удалить трек')
    }
  }

  return (
    <main style={{ flex: 1 }}>
      <section className="hero">
        <div className="container">
          <h1>Каталог треков обучения</h1>
          <p>Выберите направление, которое подходит именно вам</p>
        </div>
      </section>

      <div className="toolbar-bar">
        <div className="container">
          <div className="toolbar-inner">
            <div className="search-box">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
                <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
              </svg>
              <input
                type="search"
                placeholder="Поиск трека"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                aria-label="Поиск трека"
              />
            </div>

            {isAdmin && (
              <button className="btn btn-primary" onClick={() => setFormTrack('new')}>
                + Новый трек
              </button>
            )}
          </div>
        </div>
      </div>

      <div className="content-area">
        <div className="container">

          {isStudent && (
            <section className="top-block">
              <button
                className="action-card action-card--wide"
                onClick={handlePerformanceClick}
                disabled={performanceChecking}
              >
                <div className="action-card-visual action-card-visual--perf">
                  <svg width="56" height="56" viewBox="0 0 64 64" fill="none" aria-hidden="true">
                    <rect x="8" y="8" width="48" height="48" rx="12" fill="#fff" opacity=".45"/>
                    <rect x="14" y="36" width="6" height="16" rx="2" fill="#4B5563"/>
                    <rect x="24" y="28" width="6" height="24" rx="2" fill="#374151"/>
                    <rect x="34" y="20" width="6" height="32" rx="2" fill="#1F2937"/>
                    <path d="M14 22 L22 18 L30 22 L42 14 L50 18" stroke="#1F2937" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" fill="none"/>
                    <circle cx="22" cy="18" r="2.5" fill="#1F2937"/>
                    <circle cx="42" cy="14" r="2.5" fill="#1F2937"/>
                  </svg>
                </div>
                <div className="action-card-body">
                  <span className="action-card-eyebrow">Для вас</span>
                  <h3 className="action-card-title">Успеваемость и рекомендации</h3>
                  <p className="action-card-meta">
                    Посмотрите свой прогресс и подходящие треки на основе оценок и навыков
                  </p>
                </div>
              </button>

              <button className="action-card action-card--square" disabled title="Скоро">
                <div className="action-card-visual action-card-visual--pick">
                  <svg width="72" height="72" viewBox="0 0 72 72" fill="none" aria-hidden="true">
                    <rect x="10" y="10" width="52" height="52" rx="14" fill="#fff" opacity=".4"/>
                    <path d="M20 38 L30 48 L52 24" stroke="#4B5563" strokeWidth="4" strokeLinecap="round" strokeLinejoin="round"/>
                  </svg>
                </div>
                <div className="action-card-body">
                  <span className="action-card-eyebrow">Решение</span>
                  <h3 className="action-card-title">Выбор трека</h3>
                  <p className="action-card-meta">Скоро</p>
                </div>
              </button>
            </section>
          )}

          <div className="layout">
            <aside className="sidebar">
              <div className="filter-group">
                <h3>Сложность</h3>
                {[1, 2, 3, 4, 5].map((v) => (
                  <label key={v} className="filter-label">
                    <input
                      type="checkbox"
                      checked={filterDiff.has(v)}
                      onChange={() => toggleDiff(v)}
                    />
                    <span className="filter-label-text">
                      <span className="difficulty-stars small" aria-hidden="true">
                        {'★'.repeat(v)}{'☆'.repeat(5 - v)}
                      </span>
                      <span>{DIFFICULTY_LABELS[v]}</span>
                    </span>
                  </label>
                ))}
              </div>

              <p className="results-count">
                {loading ? 'Загрузка…' : `Найдено: ${filtered.length} ${declension(filtered.length)}`}
              </p>
            </aside>

            <div className="cards-grid">
              {error ? (
                <div className="empty-state">
                  <p>{error}</p>
                </div>
              ) : loading ? (
                <div className="empty-state"><p>Загрузка треков…</p></div>
              ) : filtered.length === 0 ? (
                <div className="empty-state">
                  <svg width="64" height="64" viewBox="0 0 64 64" fill="none">
                    <circle cx="32" cy="32" r="30" stroke="#9CA3AF" strokeWidth="2"/>
                    <path d="M22 38 Q32 26 42 38" stroke="#9CA3AF" strokeWidth="2" strokeLinecap="round"/>
                    <circle cx="24" cy="26" r="2.5" fill="#9CA3AF"/>
                    <circle cx="40" cy="26" r="2.5" fill="#9CA3AF"/>
                  </svg>
                  <p>Треки не найдены. Попробуйте изменить фильтры.</p>
                </div>
              ) : (
                filtered.map(track => (
                  <TrackCard
                    key={track.id}
                    track={track}
                    onOpen={setDetailsTrack}
                    isAdmin={isAdmin}
                    onEdit={(t) => setFormTrack(t)}
                    onDelete={handleDelete}
                  />
                ))
              )}
            </div>
          </div>
        </div>
      </div>

      {detailsTrack && (
        <TrackDetailsModal track={detailsTrack} onClose={() => setDetailsTrack(null)} />
      )}
      {performanceOpen && (
        <PerformanceModal onClose={() => setPerformanceOpen(false)} />
      )}
      {formTrack && (
        <TrackFormModal
          track={formTrack === 'new' ? null : formTrack}
          onClose={() => setFormTrack(null)}
          onSaved={loadTracks}
        />
      )}
    </main>
  )
}

function declension(n) {
  const abs = Math.abs(n)
  if (abs % 100 >= 11 && abs % 100 <= 14) return 'треков'
  switch (abs % 10) {
    case 1: return 'трек'
    case 2: case 3: case 4: return 'трека'
    default: return 'треков'
  }
}
