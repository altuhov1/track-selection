import { useState, useEffect, useCallback } from 'react'
import { fetchUserPreferences, updateUserPreferences } from '../services/profile'
import {
  SKILLS,
  LEARNING_STYLES,
  PROFESSIONAL_GOALS,
} from '../data/trackStyles'

const EMPTY_SKILLS = Object.fromEntries(SKILLS.map(s => [s.key, 5]))

export default function ProfileModal({ user, onClose, onSaved }) {
  const isAdmin = user?.role === 'admin'

  const [loading, setLoading] = useState(!isAdmin)
  const [saving, setSaving]   = useState(false)
  const [error, setError]     = useState('')

  const [skills, setSkills]                 = useState(EMPTY_SKILLS)
  const [learningStyle, setLearningStyle]   = useState(3)
  const [certificates, setCertificates]     = useState(0)
  const [goals, setGoals]                   = useState(new Set())

  const handleClose = useCallback(() => { if (!saving) onClose() }, [saving, onClose])

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
    if (isAdmin) return
    fetchUserPreferences()
      .then((prefs) => {
        if (prefs.skills)        setSkills({ ...EMPTY_SKILLS, ...prefs.skills })
        if (prefs.learning_style) setLearningStyle(prefs.learning_style)
        if (prefs.certificates != null) setCertificates(prefs.certificates)
        if (Array.isArray(prefs.professional_goals)) setGoals(new Set(prefs.professional_goals))
      })
      .catch((e) => {
        if (e.status !== 404) setError(e.message || 'Не удалось загрузить данные')
      })
      .finally(() => setLoading(false))
  }, [isAdmin])

  function setSkill(key, value) {
    setSkills(s => ({ ...s, [key]: value }))
  }
  function toggleGoal(v) {
    setGoals(prev => {
      const next = new Set(prev)
      next.has(v) ? next.delete(v) : next.add(v)
      return next
    })
  }

  async function handleSave(e) {
    e.preventDefault()
    setError('')
    setSaving(true)
    try {
      await updateUserPreferences({
        skills,
        learning_style: learningStyle,
        certificates,
        professional_goals: Array.from(goals),
      })
      onSaved?.()
      onClose()
    } catch (err) {
      setError(err.message || 'Не удалось сохранить')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="modal-overlay" onClick={(e) => e.target === e.currentTarget && handleClose()}>
      <div className="modal modal--wide" role="dialog" aria-modal="true">
        <div className="modal-header">
          <h2>Профиль</h2>
          <button className="modal-close" onClick={handleClose} aria-label="Закрыть">×</button>
        </div>

        <div className="profile-user">
          <div className="profile-user-row">
            <span className="profile-user-label">Имя</span>
            <span className="profile-user-value">{user.first_name} {user.last_name}</span>
          </div>
          <div className="profile-user-row">
            <span className="profile-user-label">Email</span>
            <span className="profile-user-value">{user.email}</span>
          </div>
          <div className="profile-user-row">
            <span className="profile-user-label">Роль</span>
            <span className="profile-user-value">{isAdmin ? 'Администратор' : 'Студент'}</span>
          </div>
        </div>

        {isAdmin ? (
          <p className="profile-admin-note">
            У администратора нет дополнительных данных в профиле.
          </p>
        ) : loading ? (
          <div className="profile-loading">Загрузка…</div>
        ) : (
          <form onSubmit={handleSave} className="profile-form">
            {error && <div className="form-error">{error}</div>}

            <section className="profile-section">
              <h3>Навыки</h3>
              <p className="profile-section-hint">Оцените себя от 0 до 10</p>
              <div className="skill-grid">
                {SKILLS.map(({ key, label }) => (
                  <div key={key} className="skill-row">
                    <div className="skill-row-top">
                      <span className="skill-label">{label}</span>
                      <span className="skill-value">{skills[key]}</span>
                    </div>
                    <input
                      type="range"
                      min={0}
                      max={10}
                      value={skills[key]}
                      onChange={(e) => setSkill(key, Number(e.target.value))}
                      className="skill-range"
                    />
                  </div>
                ))}
              </div>
            </section>

            <section className="profile-section">
              <h3>Профессиональные цели</h3>
              <div className="chip-row">
                {PROFESSIONAL_GOALS.map(({ value, label }) => (
                  <button
                    type="button"
                    key={value}
                    className={`chip chip--toggle${goals.has(value) ? ' chip--active' : ''}`}
                    onClick={() => toggleGoal(value)}
                  >
                    {label}
                  </button>
                ))}
              </div>
            </section>

            <section className="profile-section">
              <h3>Стиль обучения</h3>
              <div className="chip-row">
                {LEARNING_STYLES.map(({ value, label }) => (
                  <button
                    type="button"
                    key={value}
                    className={`chip chip--toggle${learningStyle === value ? ' chip--active' : ''}`}
                    onClick={() => setLearningStyle(value)}
                  >
                    {label}
                  </button>
                ))}
              </div>
            </section>

            <section className="profile-section">
              <h3>Сертификаты</h3>
              <div className="chip-row">
                <button
                  type="button"
                  className={`chip chip--toggle${certificates === 1 ? ' chip--active' : ''}`}
                  onClick={() => setCertificates(1)}
                >
                  Есть
                </button>
                <button
                  type="button"
                  className={`chip chip--toggle${certificates === 0 ? ' chip--active' : ''}`}
                  onClick={() => setCertificates(0)}
                >
                  Нет
                </button>
              </div>
            </section>

            <div className="form-actions">
              <button type="button" className="btn btn-ghost" onClick={handleClose} disabled={saving}>
                Отмена
              </button>
              <button type="submit" className="btn btn-primary" disabled={saving}>
                {saving ? 'Сохранение…' : 'Сохранить'}
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  )
}
