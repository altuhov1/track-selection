import { apiFetch } from './auth'

export function fetchAllTracks() {
  return apiFetch('/all-tracks')
}

export function createTrack(payload) {
  return apiFetch('/new-track', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function updateTrack(id, payload) {
  return apiFetch(`/edit-track/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  })
}

export function deleteTrack(id) {
  return apiFetch(`/delete-track/${id}`, { method: 'DELETE' })
}
