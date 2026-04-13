export default function Header({ user, onOpenLogin, onOpenRegister, onLogout }) {
  const initials = user
    ? (user.first_name[0] + user.last_name[0]).toUpperCase()
    : ''

  return (
    <header className="header">
      <div className="container">
        <div className="header-inner">
          <a href="/" className="logo" aria-label="На главную">
            <div className="logo-icon">ТВ</div>
            <span className="logo-text">ТрекВыбор</span>
          </a>

          <nav className="header-nav">
            {user ? (
              <div className="user-section">
                <div className="user-avatar">{initials}</div>
                <span className="user-full-name">{user.first_name} {user.last_name}</span>
                <span className="user-first-name">{user.first_name}</span>
                <button className="user-logout" onClick={onLogout}>Выйти</button>
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
