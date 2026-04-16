import { apiFetch } from './auth'

export function fetchProfileCompletion() {
  return apiFetch('/me/profile-completion')
}

export function fetchUserPreferences() {
  return apiFetch('/me/info')
}

export function updateUserPreferences(payload) {
  return apiFetch('/me/edit-info', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}
