import { useEffect, useCallback } from 'react'
import {
  getTrackStyle,
  DIFFICULTY_LABELS,
  LEARNING_STYLES,
  PROFESSIONAL_GOALS,
  SUBJECTS,
} from '../data/trackStyles'

const SUBJECT_LABEL = Object.fromEntries(SUBJECTS.map(s => [s.key, s.label]))
const LEARNING_LABEL = Object.fromEntries(LEARNING_STYLES.map(s => [s.value, s.label]))
const GOAL_LABEL = Object.fromEntries(PROFESSIONAL_GOALS.map(g => [g.value, g.label]))

export default function TrackDetailsModal({ track, onClose }) {
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

  if (!track) return null

  const style = getTrackStyle(track.type)
  const difficulty = Math.max(1, Math.min(5, track.difficulty || 1))
  const years = track.curriculum?.years || []
  const legacySemesters = years.length === 0 ? (track.curriculum?.semesters || []) : []
  const requirements = track.requirements || []
  const teachers = track.teachers || []
  const goals = track.professional_goals || []

  return (
    <div className="modal-overlay" onClick={(e) => e.target === e.currentTarget && handleClose()}>
      <div className="modal modal--wide" role="dialog" aria-modal="true">
        <div className="modal-header">
          <div className="details-title">
            <div className="details-title-icon" style={{ background: style.color }}>
              <div dangerouslySetInnerHTML={{ __html: style.icon(style.shapeColor) }} />
            </div>
            <div>
              <span className="card-category">{style.label}</span>
              <h2>{track.name}</h2>
            </div>
          </div>
          <button className="modal-close" onClick={handleClose} aria-label="Закрыть">×</button>
        </div>

        <div className="details-body">
          {track.description && (
            <section className="details-section">
              <p className="details-description">{track.description}</p>
            </section>
          )}

          <section className="details-metrics">
            <Metric label="Сложность" value={`${difficulty} / 5 · ${DIFFICULTY_LABELS[difficulty]}`} />
            {track.employment_prospects != null && (
              <Metric label="Перспективы" value={`${track.employment_prospects} / 10`} />
            )}
            {track.alumni_reviews != null && (
              <Metric label="Отзывы выпускников" value={`${track.alumni_reviews} / 10`} />
            )}
            {track.learning_style != null && LEARNING_LABEL[track.learning_style] && (
              <Metric label="Стиль обучения" value={LEARNING_LABEL[track.learning_style]} />
            )}
            {track.has_certificates === 1 && (
              <Metric label="Сертификаты" value="Да" />
            )}
            {track.desired_tech_skills != null && (
              <Metric label="Tech" value={`${track.desired_tech_skills} / 10`} />
            )}
            {track.desired_math_skills != null && (
              <Metric label="Math" value={`${track.desired_math_skills} / 10`} />
            )}
            {track.desired_soft_skills != null && (
              <Metric label="Soft" value={`${track.desired_soft_skills} / 10`} />
            )}
          </section>

          {requirements.length > 0 && (
            <section className="details-section">
              <h3>Требования к поступлению</h3>
              <ul className="details-list">
                {requirements.map((r, i) => (
                  <li key={i}>
                    <span>{SUBJECT_LABEL[r.subject] || r.subject}</span>
                    <strong>от {r.min_grade}</strong>
                  </li>
                ))}
              </ul>
            </section>
          )}

          {goals.length > 0 && (
            <section className="details-section">
              <h3>Подходит для целей</h3>
              <div className="chip-row">
                {goals.map((g) => (
                  <span key={g} className="chip">{GOAL_LABEL[g] || `Цель ${g}`}</span>
                ))}
              </div>
            </section>
          )}

          {teachers.length > 0 && (
            <section className="details-section">
              <h3>Преподаватели</h3>
              <div className="chip-row">
                {teachers.map((t, i) => <span key={i} className="chip">{t}</span>)}
              </div>
            </section>
          )}

          {(years.length > 0 || legacySemesters.length > 0) && (
            <section className="details-section">
              <h3>Учебный план</h3>
              {years.length > 0 ? (
                <div className="year-plans">
                  {years.map((yp) => (
                    <div key={yp.year} className="year-plan">
                      <div className="year-plan-header">
                        <span className="year-plan-num">{yp.year} курс</span>
                        {yp.type === 'branching' && (
                          <span className="year-plan-badge">выбор специализации</span>
                        )}
                      </div>
                      {yp.type === 'single' && yp.track && (
                        <div className="year-plan-body">
                          {yp.track.name && <p className="year-track-name">{yp.track.name}</p>}
                          {yp.track.description && <p className="year-track-desc">{yp.track.description}</p>}
                          <SemesterList semesters={yp.track.semesters || []} />
                        </div>
                      )}
                      {yp.type === 'branching' && (yp.branches || []).length > 0 && (
                        <BranchTree branches={yp.branches} />
                      )}
                    </div>
                  ))}
                </div>
              ) : (
                <div className="semesters">
                  <SemesterList semesters={legacySemesters} />
                </div>
              )}
            </section>
          )}

          {track.web_link && (
            <section className="details-section">
              <a href={track.web_link} target="_blank" rel="noreferrer" className="btn btn-outline">
                Перейти на сайт трека
              </a>
            </section>
          )}
        </div>
      </div>
    </div>
  )
}

function Metric({ label, value }) {
  return (
    <div className="metric">
      <span className="metric-label">{label}</span>
      <span className="metric-value">{value}</span>
    </div>
  )
}

const BRANCH_COLORS = ['#2563EB', '#7C3AED', '#059669', '#D97706', '#DC2626', '#0891B2']

function getNestedBranches(br) {
  const list = br?.branches ?? br?.Branches
  return Array.isArray(list) ? list : []
}

function BranchTree({ branches }) {
  if (!branches?.length) return null
  return (
    <ul className="branch-tree">
      {branches.map((br, bi) => {
        const hasSems   = Array.isArray(br.semesters) && br.semesters.length > 0
        const nested    = getNestedBranches(br)
        const hasNested = nested.length > 0
        const color = BRANCH_COLORS[bi % BRANCH_COLORS.length]
        return (
          <li key={bi} className="branch-tree-item" style={{ '--branch-color': color }}>
            <div className="branch-node">
              <div className="branch-node-head">
                <span className="branch-node-dot" />
                <span className="branch-node-title">{br.name}</span>
              </div>
              {br.description && <p className="branch-node-desc">{br.description}</p>}
              {hasSems && (
                <div className="branch-node-semesters">
                  <SemesterList semesters={br.semesters} />
                </div>
              )}
              {hasNested && (
                <div className="branch-node-choice">
                  <div className="branch-node-choice-label">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                      <polyline points="6 9 12 15 18 9"/>
                    </svg>
                    Далее выбор специализации
                  </div>
                  <BranchTree branches={nested} />
                </div>
              )}
            </div>
          </li>
        )
      })}
    </ul>
  )
}

function SemesterList({ semesters }) {
  if (!semesters.length) return null
  return (
    <div className="semesters">
      {semesters.map((sem) => (
        <div key={sem.number} className="semester">
          <h4>Семестр {sem.number}</h4>
          <ul className="course-list">
            {renderCourseItems(sem.courses || [])}
          </ul>
        </div>
      ))}
    </div>
  )
}

function renderCourseItems(courses) {
  const items = []
  let i = 0
  while (i < courses.length) {
    if (courses[i].is_elective) {
      let j = i
      while (j < courses.length && courses[j].is_elective) j++
      const runLen = j - i
      if (runLen > 1) {
        items.push(
          <li key={`note-${i}`} className="course-elective-note">
            Выберите только один из предметов по выбору
          </li>
        )
      }
      for (let k = i; k < j; k++) items.push(renderCourse(courses[k], k))
      i = j
    } else {
      items.push(renderCourse(courses[i], i))
      i++
    }
  }
  return items
}

function renderCourse(course, key) {
  return (
    <li key={key} className="course">
      <div className="course-top">
        <span className="course-name">{course.name}</span>
        {course.is_elective && <span className="course-tag">по выбору</span>}
      </div>
      {course.description && <p className="course-desc">{course.description}</p>}
      {course.is_elective && course.options?.length > 0 && (
        <div className="course-options">
          {course.options.map((o, j) => <span key={j} className="chip chip-sm">{o}</span>)}
        </div>
      )}
    </li>
  )
}
