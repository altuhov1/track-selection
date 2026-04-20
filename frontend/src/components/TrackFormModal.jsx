import { useState, useEffect, useCallback } from 'react'
import { createTrack, updateTrack } from '../services/tracks'
import {
  TRACK_STYLES,
  SUBJECTS,
  LEARNING_STYLES,
  PROFESSIONAL_GOALS,
} from '../data/trackStyles'

const MAX_SEMESTERS = 12

function emptyCourse() {
  return { name: '', description: '', is_elective: false, options: [] }
}

function emptySemester(number = 1) {
  return { number, courses: [emptyCourse()] }
}

function emptyTrack() {
  return { name: '', description: '', semesters: [emptySemester(1)] }
}

function emptyBranch() {
  return { name: '', description: '', semesters: [emptySemester(1)] }
}

function emptyYear(year = 3, type = 'single') {
  return {
    year,
    type,
    track: type === 'single' ? emptyTrack() : emptyTrack(),
    branches: type === 'branching' ? [emptyBranch(), emptyBranch()] : [],
  }
}

function normalizeSemesters(list) {
  if (!Array.isArray(list) || list.length === 0) return [emptySemester(1)]
  return list.map((s) => ({
    number: s.number ?? 1,
    courses: (s.courses || []).map((c) => ({
      name: c.name ?? '',
      description: c.description ?? '',
      is_elective: !!c.is_elective,
      options: Array.isArray(c.options) ? c.options : [],
    })),
  }))
}

function initialYears(curriculum) {
  const years = curriculum?.years
  if (Array.isArray(years) && years.length > 0) {
    return years.map((y) => ({
      year: y.year ?? 3,
      type: y.type === 'branching' ? 'branching' : 'single',
      track: {
        name: y.track?.name ?? '',
        description: y.track?.description ?? '',
        semesters: normalizeSemesters(y.track?.semesters),
      },
      branches: Array.isArray(y.branches) && y.branches.length > 0
        ? y.branches.map((b) => ({
            name: b.name ?? '',
            description: b.description ?? '',
            semesters: normalizeSemesters(b.semesters),
          }))
        : [emptyBranch(), emptyBranch()],
    }))
  }

  // Legacy fallback: plain semesters → single year
  const legacy = curriculum?.semesters
  if (Array.isArray(legacy) && legacy.length > 0) {
    return [{
      year: 3,
      type: 'single',
      track: { name: '', description: '', semesters: normalizeSemesters(legacy) },
      branches: [emptyBranch(), emptyBranch()],
    }]
  }

  return [emptyYear(3, 'single')]
}

export default function TrackFormModal({ track, onClose, onSaved }) {
  const isEdit = !!track

  const [name, setName]               = useState(track?.name || '')
  const [description, setDescription] = useState(track?.description || '')
  const [webLink, setWebLink]         = useState(track?.web_link || '')
  const [type, setType]               = useState(track?.type || 1)
  const [difficulty, setDifficulty]   = useState(track?.difficulty || 3)
  const [employment, setEmployment]   = useState(track?.employment_prospects || 5)
  const [alumni, setAlumni]           = useState(track?.alumni_reviews || 5)
  const [learningStyle, setLearningStyle] = useState(track?.learning_style || 3)
  const [certificates, setCertificates]   = useState(track?.has_certificates ?? 0)
  const [techSkills, setTechSkills]   = useState(track?.desired_tech_skills || 5)
  const [mathSkills, setMathSkills]   = useState(track?.desired_math_skills || 5)
  const [softSkills, setSoftSkills]   = useState(track?.desired_soft_skills || 5)
  const [goals, setGoals]             = useState(new Set(track?.professional_goals || []))
  const [teachers, setTeachers]       = useState((track?.teachers || []).join(', '))
  const [requirements, setRequirements] = useState(track?.requirements || [])
  const [years, setYears]             = useState(() => initialYears(track?.curriculum))

  const [saving, setSaving] = useState(false)
  const [error, setError]   = useState('')

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

  function toggleGoal(v) {
    setGoals(prev => {
      const next = new Set(prev)
      next.has(v) ? next.delete(v) : next.add(v)
      return next
    })
  }

  function addRequirement() {
    setRequirements(r => [...r, { subject: SUBJECTS[0].key, min_grade: 3 }])
  }
  function updateRequirement(i, patch) {
    setRequirements(r => r.map((req, idx) => idx === i ? { ...req, ...patch } : req))
  }
  function removeRequirement(i) {
    setRequirements(r => r.filter((_, idx) => idx !== i))
  }

  function updateYear(yi, patch) {
    setYears(ys => ys.map((y, i) => i === yi ? { ...y, ...patch } : y))
  }
  function addYear() {
    setYears(ys => {
      const used = new Set(ys.map(y => y.year))
      let next = 3
      while (used.has(next)) next++
      return [...ys, emptyYear(next, 'single')]
    })
  }
  function removeYear(yi) {
    setYears(ys => ys.filter((_, i) => i !== yi))
  }
  function setYearType(yi, newType) {
    setYears(ys => ys.map((y, i) => {
      if (i !== yi) return y
      if (y.type === newType) return y
      if (newType === 'branching') {
        const branches = y.branches && y.branches.length > 0 ? y.branches : [emptyBranch(), emptyBranch()]
        return { ...y, type: 'branching', branches }
      }
      const track = y.track && y.track.semesters?.length ? y.track : emptyTrack()
      return { ...y, type: 'single', track }
    }))
  }
  function updateYearTrack(yi, patch) {
    setYears(ys => ys.map((y, i) => i === yi ? { ...y, track: { ...y.track, ...patch } } : y))
  }
  function addBranch(yi) {
    setYears(ys => ys.map((y, i) => i === yi ? { ...y, branches: [...(y.branches || []), emptyBranch()] } : y))
  }
  function updateBranch(yi, bi, patch) {
    setYears(ys => ys.map((y, i) => {
      if (i !== yi) return y
      return { ...y, branches: y.branches.map((b, bj) => bj === bi ? { ...b, ...patch } : b) }
    }))
  }
  function removeBranch(yi, bi) {
    setYears(ys => ys.map((y, i) => {
      if (i !== yi) return y
      return { ...y, branches: y.branches.filter((_, bj) => bj !== bi) }
    }))
  }

  async function handleSubmit(e) {
    e.preventDefault()
    setError('')

    const curriculum = {
      years: years.map((y) => {
        const base = { year: Number(y.year), type: y.type }
        if (y.type === 'single') {
          return {
            ...base,
            track: {
              name: y.track.name.trim(),
              description: y.track.description.trim(),
              semesters: serializeSemesters(y.track.semesters),
            },
          }
        }
        return {
          ...base,
          branches: y.branches.map((b) => ({
            name: b.name.trim(),
            description: b.description.trim(),
            semesters: serializeSemesters(b.semesters),
          })),
        }
      }),
    }

    const payload = {
      name: name.trim(),
      description: description.trim(),
      curriculum,
      requirements,
      teachers: teachers.split(',').map(t => t.trim()).filter(Boolean),
      difficulty: Number(difficulty),
      type: Number(type),
      employment_prospects: Number(employment),
      alumni_reviews: Number(alumni),
      learning_style: Number(learningStyle),
      has_certificates: Number(certificates),
      desired_tech_skills: Number(techSkills),
      desired_math_skills: Number(mathSkills),
      desired_soft_skills: Number(softSkills),
      professional_goals: Array.from(goals),
      web_link: webLink.trim(),
    }

    setSaving(true)
    try {
      if (isEdit) await updateTrack(track.id, payload)
      else        await createTrack(payload)
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
          <h2>{isEdit ? 'Редактировать трек' : 'Новый трек'}</h2>
          <button className="modal-close" onClick={handleClose} aria-label="Закрыть">×</button>
        </div>

        <form onSubmit={handleSubmit} className="profile-form">
          {error && <div className="form-error">{error}</div>}

          <section className="profile-section">
            <h3>Основное</h3>
            <div className="form-group">
              <label>Название</label>
              <input type="text" value={name} onChange={(e) => setName(e.target.value)} required />
            </div>
            <div className="form-group">
              <label>Описание</label>
              <textarea
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={3}
                required
              />
            </div>
            <div className="form-group">
              <label>Ссылка</label>
              <input type="text" value={webLink} onChange={(e) => setWebLink(e.target.value)} />
            </div>
          </section>

          <section className="profile-section">
            <h3>Тип и параметры</h3>

            <p className="profile-section-hint">Тип (определяет цвет/иконку карточки)</p>
            <div className="chip-row">
              {Object.values(TRACK_STYLES).map((s) => (
                <button
                  type="button"
                  key={s.key}
                  className={`chip chip--toggle${type === s.key ? ' chip--active' : ''}`}
                  onClick={() => setType(s.key)}
                >
                  {s.label}
                </button>
              ))}
            </div>

            <div className="form-row form-row-3" style={{ marginTop: 16 }}>
              <NumField label="Сложность (1–5)" min={1} max={5} value={difficulty} onChange={setDifficulty} />
              <NumField label="Перспективы (1–10)" min={1} max={10} value={employment} onChange={setEmployment} />
              <NumField label="Отзывы (1–10)" min={1} max={10} value={alumni} onChange={setAlumni} />
            </div>

            <div className="form-row form-row-3" style={{ marginTop: 8 }}>
              <NumField label="Tech (1–10)" min={1} max={10} value={techSkills} onChange={setTechSkills} />
              <NumField label="Math (1–10)" min={1} max={10} value={mathSkills} onChange={setMathSkills} />
              <NumField label="Soft (1–10)" min={1} max={10} value={softSkills} onChange={setSoftSkills} />
            </div>

            <p className="profile-section-hint" style={{ marginTop: 16 }}>Стиль обучения</p>
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

            <p className="profile-section-hint" style={{ marginTop: 16 }}>Сертификаты</p>
            <div className="chip-row">
              <button
                type="button"
                className={`chip chip--toggle${certificates === 1 ? ' chip--active' : ''}`}
                onClick={() => setCertificates(1)}
              >Есть</button>
              <button
                type="button"
                className={`chip chip--toggle${certificates === 0 ? ' chip--active' : ''}`}
                onClick={() => setCertificates(0)}
              >Нет</button>
            </div>
          </section>

          <section className="profile-section">
            <h3>Цели</h3>
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
            <h3>Преподаватели</h3>
            <div className="form-group">
              <label>Через запятую</label>
              <input type="text" value={teachers} onChange={(e) => setTeachers(e.target.value)} placeholder="Иванов И.И., Петров П.П." />
            </div>
          </section>

          <section className="profile-section">
            <h3>Требования</h3>
            <div className="req-list">
              {requirements.map((r, i) => (
                <div key={i} className="req-row">
                  <select
                    value={r.subject}
                    onChange={(e) => updateRequirement(i, { subject: e.target.value })}
                  >
                    {SUBJECTS.map(s => <option key={s.key} value={s.key}>{s.label}</option>)}
                  </select>
                  <input
                    type="number"
                    min={2}
                    max={5}
                    value={r.min_grade}
                    onChange={(e) => updateRequirement(i, { min_grade: Number(e.target.value) })}
                  />
                  <button type="button" className="btn-ghost-sm btn-ghost-sm--danger" onClick={() => removeRequirement(i)}>
                    Удалить
                  </button>
                </div>
              ))}
              <button type="button" className="btn btn-outline" onClick={addRequirement}>
                + Требование
              </button>
            </div>
          </section>

          <section className="profile-section">
            <h3>Учебный план</h3>
            <p className="profile-section-hint">
              Каждый год обучения — либо один общий трек, либо выбор специализации (подтреков).
            </p>

            <div className="year-editor">
              {years.map((y, yi) => (
                <div key={yi} className="year-edit-block">
                  <div className="year-edit-head">
                    <label className="sem-number">
                      <span>Год</span>
                      <input
                        type="number"
                        min={1}
                        max={12}
                        value={y.year}
                        onChange={(e) => updateYear(yi, { year: Number(e.target.value) })}
                      />
                    </label>

                    <div className="chip-row">
                      <button
                        type="button"
                        className={`chip chip--toggle${y.type === 'single' ? ' chip--active' : ''}`}
                        onClick={() => setYearType(yi, 'single')}
                      >Один трек</button>
                      <button
                        type="button"
                        className={`chip chip--toggle${y.type === 'branching' ? ' chip--active' : ''}`}
                        onClick={() => setYearType(yi, 'branching')}
                      >Выбор подтрека</button>
                    </div>

                    <button
                      type="button"
                      className="btn-ghost-sm btn-ghost-sm--danger"
                      onClick={() => removeYear(yi)}
                    >
                      Удалить год
                    </button>
                  </div>

                  {y.type === 'single' && (
                    <div className="year-edit-body">
                      <div className="form-group">
                        <label>Название трека на год</label>
                        <input
                          type="text"
                          value={y.track.name}
                          onChange={(e) => updateYearTrack(yi, { name: e.target.value })}
                          placeholder="Например, Базовый год"
                        />
                      </div>
                      <div className="form-group">
                        <label>Описание</label>
                        <textarea
                          rows={2}
                          value={y.track.description}
                          onChange={(e) => updateYearTrack(yi, { description: e.target.value })}
                        />
                      </div>
                      <SemesterEditor
                        semesters={y.track.semesters}
                        onChange={(semesters) => updateYearTrack(yi, { semesters })}
                      />
                    </div>
                  )}

                  {y.type === 'branching' && (
                    <div className="year-edit-body">
                      <p className="profile-section-hint">Подтреки (студент выберет один)</p>
                      {y.branches.map((b, bi) => (
                        <div key={bi} className="branch-edit-block">
                          <div className="branch-edit-head">
                            <strong>Подтрек {bi + 1}</strong>
                            <button
                              type="button"
                              className="btn-ghost-sm btn-ghost-sm--danger"
                              onClick={() => removeBranch(yi, bi)}
                            >
                              Удалить подтрек
                            </button>
                          </div>
                          <div className="form-group">
                            <label>Название подтрека</label>
                            <input
                              type="text"
                              value={b.name}
                              onChange={(e) => updateBranch(yi, bi, { name: e.target.value })}
                              placeholder="Например, Искусственный интеллект"
                            />
                          </div>
                          <div className="form-group">
                            <label>Описание</label>
                            <textarea
                              rows={2}
                              value={b.description}
                              onChange={(e) => updateBranch(yi, bi, { description: e.target.value })}
                            />
                          </div>
                          <SemesterEditor
                            semesters={b.semesters}
                            onChange={(semesters) => updateBranch(yi, bi, { semesters })}
                          />
                        </div>
                      ))}
                      <button
                        type="button"
                        className="btn btn-outline btn-block"
                        onClick={() => addBranch(yi)}
                      >
                        + Подтрек
                      </button>
                    </div>
                  )}
                </div>
              ))}

              <button
                type="button"
                className="btn btn-outline btn-block"
                onClick={addYear}
              >
                + Год
              </button>
            </div>
          </section>

          <div className="form-actions">
            <button type="button" className="btn btn-ghost" onClick={handleClose} disabled={saving}>
              Отмена
            </button>
            <button type="submit" className="btn btn-primary" disabled={saving}>
              {saving ? 'Сохранение…' : (isEdit ? 'Сохранить' : 'Создать')}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

function serializeSemesters(list) {
  return (list || [])
    .slice()
    .sort((a, b) => a.number - b.number)
    .map((sem) => ({
      number: Number(sem.number),
      courses: sem.courses.map((c) => {
        const course = {
          name: c.name.trim(),
          description: c.description.trim(),
          is_elective: !!c.is_elective,
        }
        if (c.is_elective) {
          course.options = (c.options || [])
            .map((o) => (typeof o === 'string' ? o.trim() : String(o)))
            .filter(Boolean)
        }
        return course
      }),
    }))
}

function SemesterEditor({ semesters, onChange }) {
  function addSemester() {
    const used = new Set(semesters.map(s => s.number))
    let next = 1
    while (used.has(next) && next <= MAX_SEMESTERS) next++
    onChange([...semesters, emptySemester(next)])
  }
  function removeSemester(i) {
    onChange(semesters.filter((_, idx) => idx !== i))
  }
  function setSemesterNumber(i, number) {
    onChange(semesters.map((sem, idx) => idx === i ? { ...sem, number } : sem))
  }
  function addCourse(semIdx) {
    onChange(semesters.map((sem, idx) =>
      idx === semIdx ? { ...sem, courses: [...sem.courses, emptyCourse()] } : sem
    ))
  }
  function updateCourse(semIdx, courseIdx, patch) {
    onChange(semesters.map((sem, idx) =>
      idx !== semIdx ? sem : {
        ...sem,
        courses: sem.courses.map((c, ci) => ci === courseIdx ? { ...c, ...patch } : c),
      }
    ))
  }
  function removeCourse(semIdx, courseIdx) {
    onChange(semesters.map((sem, idx) =>
      idx !== semIdx ? sem : { ...sem, courses: sem.courses.filter((_, ci) => ci !== courseIdx) }
    ))
  }

  return (
    <div className="sem-editor">
      {semesters.map((sem, semIdx) => (
        <div key={semIdx} className="sem-block">
          <div className="sem-block-head">
            <label className="sem-number">
              <span>Семестр</span>
              <select
                value={sem.number}
                onChange={(e) => setSemesterNumber(semIdx, Number(e.target.value))}
              >
                {Array.from({ length: MAX_SEMESTERS }, (_, i) => i + 1).map(n => (
                  <option key={n} value={n}>{n}</option>
                ))}
              </select>
            </label>
            <button
              type="button"
              className="btn-ghost-sm btn-ghost-sm--danger"
              onClick={() => removeSemester(semIdx)}
            >
              Удалить семестр
            </button>
          </div>

          <div className="course-editor">
            {sem.courses.map((course, cIdx) => (
              <div key={cIdx} className="course-block">
                <div className="course-block-grid">
                  <div className="form-group">
                    <label>Название предмета</label>
                    <input
                      type="text"
                      value={course.name}
                      onChange={(e) => updateCourse(semIdx, cIdx, { name: e.target.value })}
                      placeholder="Например, Алгоритмы"
                    />
                  </div>
                  <div className="form-group">
                    <label>Что изучается</label>
                    <textarea
                      value={course.description}
                      onChange={(e) => updateCourse(semIdx, cIdx, { description: e.target.value })}
                      rows={2}
                      placeholder="Основные темы, содержание"
                    />
                  </div>
                </div>

                <div className="course-block-foot">
                  <label className="checkbox-inline">
                    <input
                      type="checkbox"
                      checked={course.is_elective}
                      onChange={(e) => updateCourse(semIdx, cIdx, { is_elective: e.target.checked })}
                    />
                    <span>Предмет по выбору</span>
                  </label>

                  {course.is_elective && (
                    <div className="course-options-field">
                      <label>Варианты (через запятую)</label>
                      <input
                        type="text"
                        value={(course.options || []).join(', ')}
                        onChange={(e) => updateCourse(semIdx, cIdx, {
                          options: e.target.value.split(',').map(o => o.trim()).filter(Boolean),
                        })}
                        placeholder="A2, B1, B2"
                      />
                    </div>
                  )}

                  <button
                    type="button"
                    className="btn-ghost-sm btn-ghost-sm--danger course-remove"
                    onClick={() => removeCourse(semIdx, cIdx)}
                  >
                    Удалить предмет
                  </button>
                </div>
              </div>
            ))}

            <button
              type="button"
              className="btn btn-outline btn-block"
              onClick={() => addCourse(semIdx)}
            >
              + Предмет
            </button>
          </div>
        </div>
      ))}

      <button
        type="button"
        className="btn btn-outline btn-block"
        onClick={addSemester}
        disabled={semesters.length >= MAX_SEMESTERS}
      >
        + Семестр
      </button>
    </div>
  )
}

function NumField({ label, min, max, value, onChange }) {
  return (
    <div className="form-group">
      <label>{label}</label>
      <input
        type="number"
        min={min}
        max={max}
        value={value}
        onChange={(e) => onChange(Number(e.target.value))}
      />
    </div>
  )
}
