import { useState, useEffect } from 'react'
import Header from './components/Header'
import Footer from './components/Footer'
import Modal from './components/Modal'
import Home from './pages/Home'
import { getToken, getUser, saveUser, clearAuth, fetchMe } from './services/auth'

export default function App() {
  const [user, setUser]       = useState(getUser)   // lazy init from localStorage
  const [modalTab, setModalTab] = useState(null)    // null | 'login' | 'register'

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
    </>
  )
}
