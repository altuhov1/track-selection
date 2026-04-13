import { useState, useEffect, useCallback } from 'react'
import { login, register } from '../services/auth'

const ERROR_MESSAGES = {
  UNAUTHORIZED:  'Неверный email или пароль',
  EMAIL_EXISTS:  'Этот email уже зарегистрирован',
  WEAK_PASSWORD: 'Пароль должен быть не менее 6 символов',
  INVALID_EMAIL: 'Некорректный email адрес',
}

export default function Modal({ initialTab, onClose, onSuccess }) {
  const [tab, setTab]         = useState(initialTab)
  const [loading, setLoading] = useState(false)
  const [error, setError]     = useState('')

  // Login fields
  const [loginEmail, setLoginEmail]       = useState('')
  const [loginPassword, setLoginPassword] = useState('')

  // Register fields
  const [regFirstName, setRegFirstName] = useState('')
  const [regLastName, setRegLastName]   = useState('')
  const [regEmail, setRegEmail]         = useState('')
  const [regPassword, setRegPassword]   = useState('')

  const handleClose = useCallback(() => {
    if (!loading) onClose()
  }, [loading, onClose])

  // Close on Escape
  useEffect(() => {
    const onKey = (e) => { if (e.key === 'Escape') handleClose() }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [handleClose])

  // Lock body scroll
  useEffect(() => {
    document.body.style.overflow = 'hidden'
    return () => { document.body.style.overflow = '' }
  }, [])

  function switchTab(t) {
    setTab(t)
    setError('')
  }

  async function handleLogin(e) {
    e.preventDefault()
    if (!loginEmail || !loginPassword) { setError('Заполните все поля'); return }
    setError(''); setLoading(true)
    try {
      const user = await login(loginEmail, loginPassword)
      onSuccess(user)
    } catch (err) {
      setError(ERROR_MESSAGES[err.code] || err.message || 'Ошибка входа')
    } finally {
      setLoading(false)
    }
  }

  async function handleRegister(e) {
    e.preventDefault()
    if (!regFirstName || !regLastName || !regEmail || !regPassword) {
      setError('Заполните все поля'); return
    }
    if (regPassword.length < 6) {
      setError('Пароль должен быть не менее 6 символов'); return
    }
    setError(''); setLoading(true)
    try {
      const user = await register(regFirstName, regLastName, regEmail, regPassword)
      onSuccess(user)
    } catch (err) {
      setError(ERROR_MESSAGES[err.code] || err.message || 'Ошибка регистрации')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="modal-overlay" onClick={(e) => e.target === e.currentTarget && handleClose()}>
      <div className="modal" role="dialog" aria-modal="true">

        <div className="modal-header">
          <h2>{tab === 'login' ? 'Войти' : 'Регистрация'}</h2>
          <button className="modal-close" onClick={handleClose} aria-label="Закрыть">&times;</button>
        </div>

        <div className="modal-tabs">
          <button className={`modal-tab${tab === 'login' ? ' active' : ''}`} onClick={() => switchTab('login')}>
            Вход
          </button>
          <button className={`modal-tab${tab === 'register' ? ' active' : ''}`} onClick={() => switchTab('register')}>
            Регистрация
          </button>
        </div>

        {error && <div className="form-error">{error}</div>}

        {tab === 'login' ? (
          <form onSubmit={handleLogin} noValidate>
            <div className="form-group">
              <label htmlFor="login-email">Email</label>
              <input
                id="login-email" type="email" placeholder="ivan@example.com"
                value={loginEmail} onChange={(e) => setLoginEmail(e.target.value)}
                autoComplete="email" required
              />
            </div>
            <div className="form-group">
              <label htmlFor="login-password">Пароль</label>
              <input
                id="login-password" type="password" placeholder="••••••••"
                value={loginPassword} onChange={(e) => setLoginPassword(e.target.value)}
                autoComplete="current-password" required
              />
            </div>
            <button type="submit" className="btn btn-primary form-submit" disabled={loading}>
              {loading ? 'Загрузка...' : 'Войти'}
            </button>
          </form>
        ) : (
          <form onSubmit={handleRegister} noValidate>
            <div className="form-row">
              <div className="form-group">
                <label htmlFor="reg-first-name">Имя</label>
                <input
                  id="reg-first-name" type="text" placeholder="Иван"
                  value={regFirstName} onChange={(e) => setRegFirstName(e.target.value)}
                  autoComplete="given-name" required
                />
              </div>
              <div className="form-group">
                <label htmlFor="reg-last-name">Фамилия</label>
                <input
                  id="reg-last-name" type="text" placeholder="Петров"
                  value={regLastName} onChange={(e) => setRegLastName(e.target.value)}
                  autoComplete="family-name" required
                />
              </div>
            </div>
            <div className="form-group">
              <label htmlFor="reg-email">Email</label>
              <input
                id="reg-email" type="email" placeholder="ivan@example.com"
                value={regEmail} onChange={(e) => setRegEmail(e.target.value)}
                autoComplete="email" required
              />
            </div>
            <div className="form-group">
              <label htmlFor="reg-password">
                Пароль <small style={{ color: 'var(--text-muted)' }}>минимум 6 символов</small>
              </label>
              <input
                id="reg-password" type="password" placeholder="••••••••"
                value={regPassword} onChange={(e) => setRegPassword(e.target.value)}
                autoComplete="new-password" required minLength={6}
              />
            </div>
            <button type="submit" className="btn btn-primary form-submit" disabled={loading}>
              {loading ? 'Загрузка...' : 'Создать аккаунт'}
            </button>
          </form>
        )}

      </div>
    </div>
  )
}
