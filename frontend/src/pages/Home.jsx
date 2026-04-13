import { useState, useMemo } from 'react'
import { TRACKS } from '../data/tracks'
import TrackCard from '../components/TrackCard'

const CATEGORIES = [
  { id: 'all',    label: 'Все треки' },
  { id: 'dev',    label: 'Разработка' },
  { id: 'data',   label: 'Данные' },
  { id: 'design', label: 'Дизайн' },
  { id: 'devops', label: 'DevOps' },
]

export default function Home() {
  const [activeCategory, setActiveCategory] = useState('all')
  const [search, setSearch]                 = useState('')
  const [filterLevels, setFilterLevels]     = useState(new Set())
  const [filterCosts, setFilterCosts]       = useState(new Set())

  const filtered = useMemo(() => {
    return TRACKS.filter(t => {
      if (activeCategory !== 'all' && t.category !== activeCategory) return false
      if (search && !t.title.toLowerCase().includes(search.toLowerCase())) return false
      if (filterLevels.size && !filterLevels.has(t.level)) return false
      if (filterCosts.has('free') && !t.isFree)  return false
      if (filterCosts.has('paid') &&  t.isFree)  return false
      return true
    })
  }, [activeCategory, search, filterLevels, filterCosts])

  function toggleLevel(value) {
    setFilterLevels(prev => {
      const next = new Set(prev)
      next.has(value) ? next.delete(value) : next.add(value)
      return next
    })
  }

  function toggleCost(value) {
    setFilterCosts(prev => {
      const next = new Set(prev)
      next.has(value) ? next.delete(value) : next.add(value)
      return next
    })
  }

  return (
    <main style={{ flex: 1 }}>
      {/* Hero */}
      <section className="hero">
        <div className="container">
          <h1>Каталог треков обучения</h1>
          <p>Почти везде есть бесплатная часть</p>
        </div>
      </section>

      {/* Tabs bar */}
      <div className="tabs-bar">
        <div className="container">
          <div className="tabs-bar-inner">
            <div className="tabs" role="tablist">
              {CATEGORIES.map(cat => (
                <button
                  key={cat.id}
                  className={`tab${activeCategory === cat.id ? ' active' : ''}`}
                  onClick={() => setActiveCategory(cat.id)}
                  role="tab"
                  aria-selected={activeCategory === cat.id}
                >
                  {cat.label}
                </button>
              ))}
            </div>

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
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="content-area">
        <div className="container">
          <div className="layout">

            {/* Sidebar */}
            <aside className="sidebar">
              <div className="filter-group">
                <h3>Уровень</h3>
                {[
                  { value: 'beginner', label: 'С нуля' },
                  { value: 'middle',   label: 'Базовый' },
                  { value: 'advanced', label: 'Продвинутый' },
                ].map(({ value, label }) => (
                  <label key={value} className="filter-label">
                    <input
                      type="checkbox"
                      checked={filterLevels.has(value)}
                      onChange={() => toggleLevel(value)}
                    />
                    {label}
                  </label>
                ))}
              </div>

              <div className="filter-group">
                <h3>Стоимость</h3>
                {[
                  { value: 'free', label: 'Бесплатно' },
                  { value: 'paid', label: 'Платно' },
                ].map(({ value, label }) => (
                  <label key={value} className="filter-label">
                    <input
                      type="checkbox"
                      checked={filterCosts.has(value)}
                      onChange={() => toggleCost(value)}
                    />
                    {label}
                  </label>
                ))}
              </div>

              <p className="results-count">Найдено: {filtered.length} {declension(filtered.length)}</p>
            </aside>

            {/* Cards */}
            <div className="cards-grid">
              {filtered.length === 0 ? (
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
                filtered.map(track => <TrackCard key={track.id} track={track} />)
              )}
            </div>

          </div>
        </div>
      </div>
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
