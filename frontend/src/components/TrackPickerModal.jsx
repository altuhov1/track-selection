import { useEffect, useCallback } from 'react'
import { getTrackStyle } from '../data/trackStyles'

export default function TrackPickerModal({ tracks, selectedTrack, selectingId, error, onPick, onClose }) {
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

  return (
    <div className="modal-overlay" onClick={(e) => e.target === e.currentTarget && handleClose()}>
      <div className="modal modal--picker" role="dialog" aria-modal="true" aria-label="Выбор трека">
        <div className="modal-header">
          <h2>Выбор трека</h2>
          <button className="modal-close" onClick={handleClose} aria-label="Закрыть">×</button>
        </div>

        <p className="picker-hint">Выберите один трек для обучения</p>

        {error && <div className="picker-error">{error}</div>}

        <ul className="picker-list">
          {tracks.length === 0 ? (
            <li className="picker-empty">Треки не найдены</li>
          ) : tracks.map(track => {
            const style = getTrackStyle(track.type)
            const isSelected = selectedTrack?.id === track.id
            const isLoading  = selectingId === track.id

            return (
              <li key={track.id}>
                <button
                  className={`picker-item${isSelected ? ' picker-item--selected' : ''}`}
                  onClick={() => onPick(track)}
                  disabled={!!selectingId}
                >
                  <span className="picker-item-dot" style={{ background: style.color }} />
                  <span className="picker-item-label">
                    <span className="picker-item-category">{style.label}</span>
                    <span className="picker-item-name">{track.name}</span>
                  </span>
                  <span className="picker-item-end">
                    {isLoading ? (
                      <span className="picker-spinner" />
                    ) : isSelected ? (
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                        <polyline points="20 6 9 17 4 12"/>
                      </svg>
                    ) : null}
                  </span>
                </button>
              </li>
            )
          })}
        </ul>
      </div>
    </div>
  )
}
