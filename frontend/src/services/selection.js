import { apiFetch } from './auth'

export const selectTrack = (trackId) =>
  apiFetch('/student/select-track', {
    method: 'POST',
    body: JSON.stringify({ track_id: trackId }),
  })

export const getSelectedTracks = () =>
  apiFetch('/student/selected-tracks')

export const unselectTrack = (trackId) =>
  apiFetch(`/student/unselect-track/${trackId}`, { method: 'DELETE' })
