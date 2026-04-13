import { SHAPES } from '../data/tracks'

export default function TrackCard({ track }) {
  const shapeSvg = (SHAPES[track.shape] || SHAPES.bracket)(track.shapeColor)

  return (
    <article className="track-card">
      <div className="card-visual" style={{ background: track.color }}>
        <div dangerouslySetInnerHTML={{ __html: shapeSvg }} />
        <span className={`card-badge${track.isFree ? ' free' : ''}`}>
          {track.isFree ? 'Бесплатно' : 'Платно'}
        </span>
      </div>
      <div className="card-body">
        <span className="card-category">{track.categoryLabel}</span>
        <h3 className="card-title">{track.title}</h3>
        <p className="card-meta">{track.levelLabel} · {track.duration}</p>
      </div>
    </article>
  )
}
