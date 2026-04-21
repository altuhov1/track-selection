import { useState, useRef, useEffect } from 'react'

export default function Header({ user, onOpenLogin, onOpenRegister, onLogout, onOpenSettings, onOpenProfile }) {
  const [menuOpen, setMenuOpen] = useState(false)
  const menuRef = useRef(null)

  const initials = user
    ? (user.first_name[0] + user.last_name[0]).toUpperCase()
    : ''

  useEffect(() => {
    if (!menuOpen) return
    function handleClickOutside(e) {
      if (menuRef.current && !menuRef.current.contains(e.target)) {
        setMenuOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [menuOpen])

  return (
    <header className="header">
      <div className="container">
        <div className="header-inner">
          <a href="/" className="logo" aria-label="На главную">
            <span className="logo-text">ТрекВыбор</span>
          </a>

          <nav className="header-nav">
            {user ? (
              <div className="user-menu-wrap" ref={menuRef}>
                <button
                  className="user-section"
                  onClick={() => setMenuOpen(o => !o)}
                  aria-expanded={menuOpen}
                >
                  <div className="user-avatar">{initials}</div>
                  <span className="user-full-name">{user.first_name} {user.last_name}</span>
                  <span className="user-first-name">{user.first_name}</span>
                </button>

                {menuOpen && (
                  <div className="user-dropdown">
                    <button className="user-dropdown-item" onClick={() => { setMenuOpen(false); onOpenProfile?.() }}>
                      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="8" r="4"/><path d="M4 20c0-4 3.6-7 8-7s8 3 8 7"/></svg>
                      Профиль
                    </button>
                    <button className="user-dropdown-item" onClick={() => { setMenuOpen(false); onOpenSettings() }}>
                      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
                      Настройки
                    </button>
                    <div className="user-dropdown-divider" />
                    <button
                      className="user-dropdown-item user-dropdown-item--danger"
                      onClick={() => { setMenuOpen(false); onLogout() }}
                    >
                      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
                      Выйти
                    </button>
                  </div>
                )}
              </div>
            ) : (
              <>
                <button className="btn btn-outline" onClick={onOpenLogin}>Войти</button>
                <button className="btn btn-primary" onClick={onOpenRegister}>Регистрация</button>
              </>
            )}
          </nav>
        </div>
      </div>
    </header>
  )
}
