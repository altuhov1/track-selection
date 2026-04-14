import { useState, useEffect } from 'react'
import Header from './components/Header'
import Footer from './components/Footer'
import Modal from './components/Modal'
import SettingsModal from './components/SettingsModal'
import Home from './pages/Home'
import { getToken, getUser, saveUser, clearAuth, fetchMe } from './services/auth'

function getSystemIsDark() {
  return window.matchMedia('(prefers-color-scheme: dark)').matches
}

function applyTheme(pref) {
  const isDark = pref === 'dark' || (pref === 'system' && getSystemIsDark())
  document.documentElement.setAttribute('data-theme', isDark ? 'dark' : 'light')
}

export default function App() {
  const [user, setUser]           = useState(getUser)
  const [modalTab, setModalTab]   = useState(null)     // null | 'login' | 'register'
  const [settingsOpen, setSettingsOpen] = useState(false)
  const [theme, setTheme]         = useState(() => localStorage.getItem('theme_pref') || 'system')

  // Apply theme on mount and when pref changes
  useEffect(() => {
    applyTheme(theme)
    localStorage.setItem('theme_pref', theme)
  }, [theme])

  // Keep system theme in sync when OS preference changes
  useEffect(() => {
    if (theme !== 'system') return
    const mq = window.matchMedia('(prefers-color-scheme: dark)')
    const handler = () => applyTheme('system')
    mq.addEventListener('change', handler)
    return () => mq.removeEventListener('change', handler)
  }, [theme])

  // On mount: if we have a token but no cached user — fetch /api/me
  useEffect(() => {
    if (getToken() && !getUser()) {
      fetchMe()
        .then(u => { saveUser(u); setUser(u) })
        .catch(() => { clearAuth(); setUser(null) })
    }
  }, [])

  function handleAuthSuccess(userData) {
    setUser(userData)
    setModalTab(null)
  }

  function handleLogout() {
    clearAuth()
    setUser(null)
  }

  return (
    <>
      <Header
        user={user}
        onOpenLogin={() => setModalTab('login')}
        onOpenRegister={() => setModalTab('register')}
        onLogout={handleLogout}
        onOpenSettings={() => setSettingsOpen(true)}
      />

      <Home />

      <Footer />

      {modalTab && (
        <Modal
          initialTab={modalTab}
          onClose={() => setModalTab(null)}
          onSuccess={handleAuthSuccess}
        />
      )}

      {settingsOpen && (
        <SettingsModal
          theme={theme}
          onThemeChange={setTheme}
          onClose={() => setSettingsOpen(false)}
        />
      )}
    </>
  )
}
