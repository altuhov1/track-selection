import { useState, useEffect, useMemo, useCallback, useRef } from 'react'
import TrackCard from '../components/TrackCard'
import TrackDetailsModal from '../components/TrackDetailsModal'
import TrackFormModal from '../components/TrackFormModal'
import PerformanceModal from '../components/PerformanceModal'
import TrackPickerModal from '../components/TrackPickerModal'
import FilterPanel from '../components/FilterPanel'
import { fetchAllTracks, deleteTrack } from '../services/tracks'
import { fetchProfileCompletion } from '../services/profile'
import { selectTrack, getSelectedTracks, unselectTrack } from '../services/selection'
import { TRACK_STYLES, DIFFICULTY_LABELS, LEARNING_STYLES, PROFESSIONAL_GOALS } from '../data/trackStyles'

const EMPTY_FILTERS = {
  difficulty:     new Set(),
  types:          new Set(),
  learningStyles: new Set(),
  certificates:   null,
  goals:          new Set(),
}

function countActiveFilters(f) {
  return f.difficulty.size + f.types.size + f.learningStyles.size +
    (f.certificates !== null ? 1 : 0) + f.goals.size
}

export default function Home({ user, onOpenProfile, onOpenLogin }) {
  const isAdmin   = user?.role === 'admin'
  const isStudent = user?.role === 'student'

  const [tracks, setTracks]   = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError]     = useState('')

  const [search, setSearch]       = useState('')
  const [filters, setFilters]     = useState(EMPTY_FILTERS)
  const [filterOpen, setFilterOpen] = useState(false)
  const filterBtnRef = useRef(null)

  const [detailsTrack, setDetailsTrack]       = useState(null)
  const [performanceOpen, setPerformanceOpen] = useState(false)
  const [performanceChecking, setPerformanceChecking] = useState(false)
  const [formTrack, setFormTrack]             = useState(null)

  const [selectedTrack, setSelectedTrack] = useState(null)
  const [pickerOpen, setPickerOpen]       = useState(false)
  const [selectingId, setSelectingId]     = useState(null)
  const [pickerError, setPickerError]     = useState('')

  const [searchStuck, setSearchStuck] = useState(false)
  const sentinelRef = useRef(null)

  useEffect(() => {
    const el = sentinelRef.current
    if (!el) return
    const headerH = parseInt(getComputedStyle(document.documentElement).getPropertyValue('--header-h')) || 64
    const observer = new IntersectionObserver(
      ([entry]) => setSearchStuck(!entry.isIntersecting),
      { threshold: 0, rootMargin: `-${headerH}px 0px 0px 0px` }
    )
    observer.observe(el)
    return () => observer.disconnect()
  }, [])

  const loadTracks = useCallback(() => {
    setLoading(true)
    fetchAllTracks()
      .then(list => setTracks(Array.isArray(list) ? list : []))
      .catch(e => setError(e.message || 'Не удалось загрузить треки'))
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => { loadTracks() }, [loadTracks])

  useEffect(() => {
    if (!isStudent) return
    getSelectedTracks()
      .then(data => setSelectedTrack((data?.tracks || [])[0] || null))
      .catch(() => {})
  }, [isStudent])

  const filtered = useMemo(() => {
    return tracks.filter(t => {
      if (search && !t.name?.toLowerCase().includes(search.toLowerCase())) return false
      if (filters.difficulty.size && !filters.difficulty.has(t.difficulty)) return false
      if (filters.types.size && !filters.types.has(t.type)) return false
      if (filters.learningStyles.size && !filters.learningStyles.has(t.learning_style)) return false
      if (filters.certificates !== null && t.has_certificates !== filters.certificates) return false
      if (filters.goals.size && !t.professional_goals?.some(g => filters.goals.has(g))) return false
      return true
    })
  }, [tracks, search, filters])

  const activeCount = countActiveFilters(filters)

  function removeFilter(key, value) {
    setFilters(prev => {
      if (prev[key] instanceof Set) {
        const next = new Set(prev[key])
        next.delete(value)
        return { ...prev, [key]: next }
      }
      return { ...prev, [key]: null }
    })
  }

  function resetFilters() { setFilters(EMPTY_FILTERS) }

  async function handlePerformanceClick() {
    if (!user) { onOpenLogin?.(); return }
    setPerformanceChecking(true)
    try {
      const { is_complete } = await fetchProfileCompletion()
      if (is_complete) setPerformanceOpen(true)
      else onOpenProfile?.()
    } catch {
      onOpenProfile?.()
    } finally {
      setPerformanceChecking(false)
    }
  }

  async function handlePickTrack(track) {
    if (selectingId) return
    if (selectedTrack?.id === track.id) return
    setPickerError('')
    setSelectingId(track.id)
    try {
      if (selectedTrack) await unselectTrack(selectedTrack.id)
      await selectTrack(track.id)
      setSelectedTrack(track)
      setPickerOpen(false)
    } catch (e) {
      setPickerError(e.message || 'Не удалось выбрать трек')
    } finally {
      setSelectingId(null)
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

  // Build active filter chips for sidebar
  const activeChips = []
  filters.difficulty.forEach(v =>
    activeChips.push({ key: 'difficulty', value: v, label: `Сложность: ${DIFFICULTY_LABELS[v]}` }))
  filters.types.forEach(v =>
    activeChips.push({ key: 'types', value: v, label: `${TRACK_STYLES[v]?.label || v}` }))
  filters.learningStyles.forEach(v =>
    activeChips.push({ key: 'learningStyles', value: v, label: LEARNING_STYLES.find(s => s.value === v)?.label }))
  if (filters.certificates !== null)
    activeChips.push({ key: 'certificates', value: filters.certificates, label: 'Есть сертификаты' })
  filters.goals.forEach(v =>
    activeChips.push({ key: 'goals', value: v, label: PROFESSIONAL_GOALS.find(g => g.value === v)?.label }))

  return (
    <main style={{ flex: 1 }}>
      <section className="hero">
        <div className="container">
          <h1>Каталог треков обучения</h1>
          <p>Выберите направление, которое подходит именно вам</p>
        </div>
      </section>

      {isStudent && (
        <div className="action-bar">
          <div className="container">
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

              <button
                className={`action-card action-card--square${selectedTrack ? ' action-card--selected' : ''}`}
                onClick={() => setPickerOpen(true)}
              >
                <div className="action-card-visual action-card-visual--pick">
                  <svg width="72" height="72" viewBox="0 0 72 72" fill="none" aria-hidden="true">
                    <rect x="10" y="10" width="52" height="52" rx="14" fill="#fff" opacity=".4"/>
                    <path d="M20 38 L30 48 L52 24" stroke="#1F2937" strokeWidth="4" strokeLinecap="round" strokeLinejoin="round"/>
                  </svg>
                </div>
                <div className="action-card-body">
                  <span className="action-card-eyebrow">Решение</span>
                  <h3 className="action-card-title">Выбор трека</h3>
                  <p className="action-card-meta">
                    {selectedTrack ? selectedTrack.name : 'Не выбран'}
                  </p>
                </div>
              </button>
            </section>
          </div>
        </div>
      )}

      <div ref={sentinelRef} className="search-sticky-sentinel" />
      <div className={`search-sticky-bar${searchStuck ? ' search-sticky-bar--stuck' : ''}`}>
        <div className="container">
          <div className={`search-filter-row${!isAdmin ? ' search-filter-row--grid' : ''}`}>
            <div className="search-box search-box--full">
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

            <div className="filter-btn-wrap" ref={filterBtnRef}>
              <button
                className={`btn filter-btn filter-btn--full${activeCount > 0 ? ' filter-btn--active' : ''}`}
                onClick={() => setFilterOpen(o => !o)}
              >
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <line x1="4" y1="6" x2="20" y2="6"/>
                  <line x1="8" y1="12" x2="16" y2="12"/>
                  <line x1="11" y1="18" x2="13" y2="18"/>
                </svg>
                Фильтры
                {activeCount > 0 && <span className="filter-badge">{activeCount}</span>}
              </button>

              {filterOpen && (
                <FilterPanel
                  filters={filters}
                  onChange={setFilters}
                  onClose={() => setFilterOpen(false)}
                />
              )}
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
          <div className="layout">
            <div className="cards-grid">
              {error ? (
                <div className="empty-state"><p>{error}</p></div>
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

            <aside className="sidebar">
              <p className="results-count" style={{ marginTop: 0, paddingTop: 0, borderTop: 'none' }}>
                {loading ? 'Загрузка…' : `Найдено: ${filtered.length} ${declension(filtered.length)}`}
              </p>

              {activeChips.length > 0 && (
                <div className="active-filters">
                  <div className="active-filters-header">
                    <span>Активные фильтры</span>
                    <button className="active-filters-reset" onClick={resetFilters}>Сбросить</button>
                  </div>
                  <div className="active-chips">
                    {activeChips.map((chip, i) => (
                      <button
                        key={i}
                        className="active-chip"
                        onClick={() => removeFilter(chip.key, chip.value)}
                      >
                        {chip.label}
                        <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round">
                          <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
                        </svg>
                      </button>
                    ))}
                  </div>
                </div>
              )}
            </aside>
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
      {pickerOpen && (
        <TrackPickerModal
          tracks={tracks}
          selectedTrack={selectedTrack}
          selectingId={selectingId}
          error={pickerError}
          onPick={handlePickTrack}
          onClose={() => { setPickerOpen(false); setPickerError('') }}
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
