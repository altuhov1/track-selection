// Стили карточек по номеру типа трека (track.type).
// Цвет фона (color), цвет иконки (shapeColor), тематическая картинка (icon).
// Все иконки — inline SVG, возвращаются как строка для dangerouslySetInnerHTML.

export const TRACK_STYLES = {
  1: {
    key: 1,
    label: 'Математика',
    color: '#DBEAFE',       // мягкий синий
    shapeColor: '#2563EB',
    icon: (c) => `<svg viewBox="0 0 72 72" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect x="10" y="10" width="52" height="52" rx="14" fill="${c}" opacity=".2"/>
      <path d="M20 24 L28 24 L32 48 L38 22 L48 22" stroke="${c}" stroke-width="3.5" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
      <path d="M46 44 L58 44" stroke="${c}" stroke-width="3" stroke-linecap="round"/>
      <circle cx="52" cy="36" r="2.5" fill="${c}"/>
      <path d="M44 52 L58 52" stroke="${c}" stroke-width="3" stroke-linecap="round" opacity=".6"/>
    </svg>`,
  },
  2: {
    key: 2,
    label: 'Инженерия',
    color: '#EDE9FE',       // мягкий фиолетовый
    shapeColor: '#7C3AED',
    icon: (c) => `<svg viewBox="0 0 72 72" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect x="10" y="10" width="52" height="52" rx="14" fill="${c}" opacity=".2"/>
      <path d="M44 16 L56 28 L30 54 L18 54 L18 42 Z" stroke="${c}" stroke-width="3.2" stroke-linejoin="round" fill="${c}" fill-opacity=".35"/>
      <path d="M40 20 L52 32" stroke="${c}" stroke-width="3.2" stroke-linecap="round"/>
      <path d="M22 46 L28 52" stroke="${c}" stroke-width="2.5" stroke-linecap="round"/>
    </svg>`,
  },
  3: {
    key: 3,
    label: 'Тестирование',
    color: '#FCE7F3',       // мягкий розовый
    shapeColor: '#DB2777',
    icon: (c) => `<svg viewBox="0 0 72 72" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect x="10" y="10" width="52" height="52" rx="14" fill="${c}" opacity=".2"/>
      <circle cx="24" cy="40" r="10" stroke="${c}" stroke-width="3.2" fill="${c}" fill-opacity=".25"/>
      <circle cx="48" cy="40" r="10" stroke="${c}" stroke-width="3.2" fill="${c}" fill-opacity=".25"/>
      <path d="M34 40 L38 40" stroke="${c}" stroke-width="3" stroke-linecap="round"/>
      <path d="M14 34 L20 30" stroke="${c}" stroke-width="3" stroke-linecap="round"/>
      <path d="M58 34 L52 30" stroke="${c}" stroke-width="3" stroke-linecap="round"/>
    </svg>`,
  },
}

// Фолбэк-стиль на случай неизвестного типа.
const FALLBACK_STYLE = {
  key: 0,
  label: 'Трек',
  color: '#F3F4F6',
  shapeColor: '#6B7280',
  icon: (c) => `<svg viewBox="0 0 72 72" fill="none" xmlns="http://www.w3.org/2000/svg">
    <rect x="10" y="10" width="52" height="52" rx="14" fill="${c}" opacity=".2"/>
    <path d="M22 36 L34 24 L50 40 L42 48 Z" stroke="${c}" stroke-width="3" stroke-linejoin="round" fill="${c}" fill-opacity=".3"/>
  </svg>`,
}

export function getTrackStyle(type) {
  return TRACK_STYLES[type] || FALLBACK_STYLE
}

export const DIFFICULTY_LABELS = {
  1: 'Очень легко',
  2: 'Легко',
  3: 'Средне',
  4: 'Сложно',
  5: 'Очень сложно',
}

// Справочники для анкеты пользователя и карточек трека.
export const SUBJECTS = [
  { key: 'informatics',               label: 'Информатика' },
  { key: 'programming',               label: 'Программирование' },
  { key: 'foreign_language',          label: 'Ин. язык' },
  { key: 'physics',                   label: 'Физика' },
  { key: 'aig',                       label: 'АиГ' },
  { key: 'math_analysis',             label: 'Мат. анализ' },
  { key: 'algorithms_data_structures',label: 'Алгоритмы' },
  { key: 'databases',                 label: 'Базы данных' },
  { key: 'discrete_math',             label: 'Дискр. мат.' },
]

export const SKILLS = [
  { key: 'databases',               label: 'Базы данных' },
  { key: 'system_architecture',     label: 'Архитектура' },
  { key: 'algorithmic_programming', label: 'Алгоритмы' },
  { key: 'public_speaking',         label: 'Публичные выступления' },
  { key: 'testing',                 label: 'Тестирование' },
  { key: 'analytics',               label: 'Аналитика' },
  { key: 'machine_learning',        label: 'ML' },
  { key: 'os_knowledge',            label: 'Операц. системы' },
  { key: 'research_projects',       label: 'Научная работа' },
]

export const LEARNING_STYLES = [
  { value: 1, label: 'Теория' },
  { value: 2, label: 'Практика' },
  { value: 3, label: 'Смешанный' },
]

export const PROFESSIONAL_GOALS = [
  { value: 1, label: 'Работа в крупной компании' },
  { value: 2, label: 'Стартап' },
  { value: 3, label: 'Фриланс' },
  { value: 4, label: 'Научная карьера' },
  { value: 5, label: 'Преподавание' },
  { value: 6, label: 'Управленческая роль' },
]
