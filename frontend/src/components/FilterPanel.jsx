import { useEffect, useRef } from 'react'
import { TRACK_STYLES, DIFFICULTY_LABELS, LEARNING_STYLES, PROFESSIONAL_GOALS } from '../data/trackStyles'

const TYPES = Object.values(TRACK_STYLES).map(s => ({ value: s.key, label: s.label }))

function Section({ title, children }) {
  return (
    <div className="fp-section">
      <h4 className="fp-section-title">{title}</h4>
      {children}
    </div>
  )
}

function CheckRow({ checked, onChange, label }) {
  return (
    <label className="fp-check-row">
      <input type="checkbox" checked={checked} onChange={onChange} />
      <span>{label}</span>
    </label>
  )
}

export default function FilterPanel({ filters, onChange, onClose }) {
  const ref = useRef(null)

  useEffect(() => {
    function onPointerDown(e) {
      if (ref.current && !ref.current.contains(e.target)) onClose()
    }
    document.addEventListener('pointerdown', onPointerDown)
    return () => document.removeEventListener('pointerdown', onPointerDown)
  }, [onClose])

  function toggleSet(key, value) {
    onChange(prev => {
      const next = new Set(prev[key])
      next.has(value) ? next.delete(value) : next.add(value)
      return { ...prev, [key]: next }
    })
  }

  function toggleBool(key, value) {
    onChange(prev => ({ ...prev, [key]: prev[key] === value ? null : value }))
  }

  return (
    <div className="fp-panel" ref={ref}>
      <Section title="Сложность">
        {[1, 2, 3, 4, 5].map(v => (
          <CheckRow
            key={v}
            checked={filters.difficulty.has(v)}
            onChange={() => toggleSet('difficulty', v)}
            label={`${v} — ${DIFFICULTY_LABELS[v]}`}
          />
        ))}
      </Section>

      <Section title="Направление">
        {TYPES.map(t => (
          <CheckRow
            key={t.value}
            checked={filters.types.has(t.value)}
            onChange={() => toggleSet('types', t.value)}
            label={t.label}
          />
        ))}
      </Section>

      <Section title="Стиль обучения">
        {LEARNING_STYLES.map(s => (
          <CheckRow
            key={s.value}
            checked={filters.learningStyles.has(s.value)}
            onChange={() => toggleSet('learningStyles', s.value)}
            label={s.label}
          />
        ))}
      </Section>

      <Section title="Сертификаты">
        <CheckRow
          checked={filters.certificates === 1}
          onChange={() => toggleBool('certificates', 1)}
          label="Есть сертификаты"
        />
      </Section>

      <Section title="Профессиональные цели">
        {PROFESSIONAL_GOALS.map(g => (
          <CheckRow
            key={g.value}
            checked={filters.goals.has(g.value)}
            onChange={() => toggleSet('goals', g.value)}
            label={g.label}
          />
        ))}
      </Section>
    </div>
  )
}
