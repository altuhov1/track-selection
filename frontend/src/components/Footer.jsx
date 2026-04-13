export default function Footer() {
  return (
    <footer className="footer">
      <div className="container">
        <div className="footer-top">
          <div className="footer-brand">
            <p className="footer-phone">8 800 555-01-01</p>
            <p className="footer-phone-note">Звонок по России бесплатный</p>
            <div className="footer-socials">
              <a href="#" className="social-link" aria-label="ВКонтакте">ВК</a>
              <a href="#" className="social-link" aria-label="Telegram">TG</a>
              <a href="#" className="social-link" aria-label="YouTube">YT</a>
            </div>
            <p className="footer-legal-note">
              Образовательный сервис для выбора и прохождения треков обучения
              в IT-профессиях. Лицензия на образовательную деятельность №12345.
            </p>
          </div>

          <div className="footer-links">
            <div className="footer-col">
              <h4>О сервисе</h4>
              <a href="#">Оферта</a>
              <a href="#">Политика конфиденциальности</a>
              <a href="#">Пользовательское соглашение</a>
              <a href="#">Отзывы</a>
              <a href="#">Помощь</a>
              <a href="#">Блог</a>
            </div>
            <div className="footer-col">
              <h4>Студентам</h4>
              <a href="#">Все треки</a>
              <a href="#">Как выбрать трек</a>
              <a href="#">Бесплатные курсы</a>
              <a href="#">Расписание</a>
              <a href="#">Карьерный центр</a>
              <a href="#">FAQ</a>
            </div>
            <div className="footer-col">
              <h4>Компаниям</h4>
              <a href="#">Корпоративное обучение</a>
              <a href="#">Партнёрская программа</a>
              <a href="#">Реферальная программа</a>
              <a href="#">Связаться с нами</a>
            </div>
          </div>
        </div>

        <div className="footer-bottom">
          <span className="footer-copyright">© 2026 ТрекВыбор</span>
          <div className="footer-tags">
            {['Frontend', 'Backend', 'Data Science', 'DevOps', 'Дизайн', 'Мобильная разработка'].map(tag => (
              <span key={tag} className="footer-tag">{tag}</span>
            ))}
          </div>
        </div>
      </div>
    </footer>
  )
}
