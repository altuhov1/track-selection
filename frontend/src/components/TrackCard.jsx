import { getTrackStyle, DIFFICULTY_LABELS } from '../data/trackStyles'

export default function TrackCard({ track, onOpen, isAdmin, onEdit, onDelete }) {
  const style = getTrackStyle(track.type)
  const difficulty = Math.max(1, Math.min(5, track.difficulty || 1))

  function handleEdit(e) {
    e.stopPropagation()
    onEdit?.(track)
  }

  function handleDelete(e) {
    e.stopPropagation()
    onDelete?.(track)
  }

  return (
    <article
      className="track-card"
      onClick={() => onOpen?.(track)}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => { if (e.key === 'Enter') onOpen?.(track) }}
    >
      <div className="card-visual" style={{ background: style.color }}>
        <div dangerouslySetInnerHTML={{ __html: style.icon(style.shapeColor) }} />
      </div>

      <div className="card-body">
        <span className="card-category">{style.label}</span>
        <h3 className="card-title">{track.name}</h3>
        <p className="card-meta">{DIFFICULTY_LABELS[difficulty]}</p>

        {isAdmin && (
          <div className="card-admin-actions">
            <button className="btn-ghost-sm" onClick={handleEdit}>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 1 1 3 3L7 19l-4 1 1-4 12.5-12.5z"/></svg>
              Изменить
            </button>
            <button className="btn-ghost-sm btn-ghost-sm--danger" onClick={handleDelete}>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-2 14H7L5 6"/><path d="M10 11v6"/><path d="M14 11v6"/></svg>
              Удалить
            </button>
          </div>
        )}
      </div>
    </article>
  )
}
