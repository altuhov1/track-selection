const TOKEN_KEY = 'jwt_token'
const USER_KEY  = 'user_data'

export const getToken = () => localStorage.getItem(TOKEN_KEY)
export const getUser  = () => { const s = localStorage.getItem(USER_KEY); return s ? JSON.parse(s) : null }
export const saveToken = (t) => localStorage.setItem(TOKEN_KEY, t)
export const saveUser  = (u) => localStorage.setItem(USER_KEY, JSON.stringify(u))
export const clearAuth = () => { localStorage.removeItem(TOKEN_KEY); localStorage.removeItem(USER_KEY) }

export async function apiFetch(path, options = {}) {
  const token = getToken()
  const headers = { 'Content-Type': 'application/json', ...(options.headers || {}) }
  if (token) headers['Authorization'] = `Bearer ${token}`

  const res = await fetch('/api' + path, { ...options, headers })

  if (res.status === 401) {
    clearAuth()
    const err = Object.assign(new Error('UNAUTHORIZED'), { status: 401, code: 'UNAUTHORIZED' })
    throw err
  }

  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    const msg = data?.error?.message || `HTTP ${res.status}`
    throw Object.assign(new Error(msg), { code: data?.error?.code, status: res.status })
  }
  return data
}

export async function fetchMe() {
  return apiFetch('/me')
}

export async function login(email, password) {
  const { token } = await apiFetch('/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
  saveToken(token)
  const user = await fetchMe()
  saveUser(user)
  return user
}

export async function register(first_name, last_name, email, password, role = 'student') {
  await apiFetch('/register', {
    method: 'POST',
    body: JSON.stringify({ first_name, last_name, email, password, role }),
  })
  return login(email, password)
}
