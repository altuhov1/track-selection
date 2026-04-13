export const TRACKS = [
  {
    id: 1, category: 'dev', categoryLabel: 'Разработка',
    title: 'Frontend разработчик',
    level: 'beginner', levelLabel: 'С нуля',
    isFree: true, duration: '8 месяцев',
    color: '#EEF2FF', shapeColor: '#818CF8', shape: 'bracket',
  },
  {
    id: 2, category: 'dev', categoryLabel: 'Разработка',
    title: 'Backend разработчик',
    level: 'beginner', levelLabel: 'С нуля',
    isFree: false, duration: '10 месяцев',
    color: '#F0FDF4', shapeColor: '#4ADE80', shape: 'server',
  },
  {
    id: 3, category: 'data', categoryLabel: 'Данные',
    title: 'Аналитик данных',
    level: 'beginner', levelLabel: 'С нуля',
    isFree: true, duration: '7 месяцев',
    color: '#FFF7ED', shapeColor: '#FB923C', shape: 'chart',
  },
  {
    id: 4, category: 'data', categoryLabel: 'Данные',
    title: 'Инженер машинного обучения',
    level: 'middle', levelLabel: 'Базовый',
    isFree: false, duration: '12 месяцев',
    color: '#FDF4FF', shapeColor: '#E879F9', shape: 'brain',
  },
  {
    id: 5, category: 'design', categoryLabel: 'Дизайн',
    title: 'UX/UI дизайнер',
    level: 'beginner', levelLabel: 'С нуля',
    isFree: true, duration: '6 месяцев',
    color: '#FFF1F2', shapeColor: '#FB7185', shape: 'pen',
  },
  {
    id: 6, category: 'devops', categoryLabel: 'DevOps',
    title: 'DevOps инженер',
    level: 'middle', levelLabel: 'Базовый',
    isFree: false, duration: '9 месяцев',
    color: '#F0F9FF', shapeColor: '#38BDF8', shape: 'cloud',
  },
  {
    id: 7, category: 'dev', categoryLabel: 'Разработка',
    title: 'Мобильный разработчик',
    level: 'beginner', levelLabel: 'С нуля',
    isFree: false, duration: '10 месяцев',
    color: '#F0FDF4', shapeColor: '#34D399', shape: 'mobile',
  },
  {
    id: 8, category: 'data', categoryLabel: 'Данные',
    title: 'Python разработчик',
    level: 'beginner', levelLabel: 'С нуля',
    isFree: false, duration: '10 месяцев',
    color: '#FEF9C3', shapeColor: '#FACC15', shape: 'python',
  },
]

export const SHAPES = {
  bracket: (c) => `<svg viewBox="0 0 72 72" fill="none" xmlns="http://www.w3.org/2000/svg">
    <rect x="10" y="10" width="52" height="52" rx="14" fill="${c}" opacity=".25"/>
    <path d="M28 20 L18 36 L28 52" stroke="${c}" stroke-width="5" stroke-linecap="round" stroke-linejoin="round"/>
    <path d="M44 20 L54 36 L44 52" stroke="${c}" stroke-width="5" stroke-linecap="round" stroke-linejoin="round"/>
  </svg>`,
  server: (c) => `<svg viewBox="0 0 72 72" fill="none" xmlns="http://www.w3.org/2000/svg">
    <rect x="10" y="10" width="52" height="52" rx="14" fill="${c}" opacity=".25"/>
    <rect x="16" y="20" width="40" height="12" rx="4" fill="${c}" opacity=".7"/>
    <rect x="16" y="36" width="40" height="12" rx="4" fill="${c}" opacity=".4"/>
    <circle cx="48" cy="26" r="2.5" fill="${c}"/>
    <circle cx="48" cy="42" r="2.5" fill="${c}" opacity=".5"/>
  </svg>`,
  chart: (c) => `<svg viewBox="0 0 72 72" fill="none" xmlns="http://www.w3.org/2000/svg">
    <rect x="10" y="10" width="52" height="52" rx="14" fill="${c}" opacity=".25"/>
    <rect x="16" y="40" width="8" height="16" rx="2" fill="${c}"/>
    <rect x="28" y="30" width="8" height="26" rx="2" fill="${c}" opacity=".7"/>
    <rect x="40" y="20" width="8" height="36" rx="2" fill="${c}" opacity=".5"/>
    <path d="M14 52 L58 52" stroke="${c}" stroke-width="2" stroke-linecap="round"/>
  </svg>`,
  brain: (c) => `<svg viewBox="0 0 72 72" fill="none" xmlns="http://www.w3.org/2000/svg">
    <rect x="10" y="10" width="52" height="52" rx="14" fill="${c}" opacity=".25"/>
    <ellipse cx="32" cy="34" rx="12" ry="14" fill="${c}" opacity=".5"/>
    <ellipse cx="44" cy="34" rx="10" ry="14" fill="${c}" opacity=".7"/>
    <line x1="36" y1="20" x2="36" y2="48" stroke="${c}" stroke-width="2"/>
  </svg>`,
  pen: (c) => `<svg viewBox="0 0 72 72" fill="none" xmlns="http://www.w3.org/2000/svg">
    <rect x="10" y="10" width="52" height="52" rx="14" fill="${c}" opacity=".25"/>
    <rect x="32" y="16" width="10" height="22" rx="5" fill="${c}" opacity=".7" transform="rotate(15 32 16)"/>
    <path d="M20 52 L28 44 L36 52 Z" fill="${c}" opacity=".5"/>
    <circle cx="38" cy="38" r="8" stroke="${c}" stroke-width="3" fill="none"/>
  </svg>`,
  cloud: (c) => `<svg viewBox="0 0 72 72" fill="none" xmlns="http://www.w3.org/2000/svg">
    <rect x="10" y="10" width="52" height="52" rx="14" fill="${c}" opacity=".25"/>
    <path d="M20 44 Q16 44 16 38 Q16 30 24 29 Q24 22 32 22 Q40 22 42 29 Q50 28 52 35 Q54 42 48 44 Z" fill="${c}" opacity=".7"/>
  </svg>`,
  mobile: (c) => `<svg viewBox="0 0 72 72" fill="none" xmlns="http://www.w3.org/2000/svg">
    <rect x="10" y="10" width="52" height="52" rx="14" fill="${c}" opacity=".25"/>
    <rect x="24" y="16" width="24" height="40" rx="6" fill="${c}" opacity=".5"/>
    <rect x="28" y="21" width="16" height="24" rx="2" fill="${c}" opacity=".5"/>
    <circle cx="36" cy="51" r="2.5" fill="${c}"/>
  </svg>`,
  python: (c) => `<svg viewBox="0 0 72 72" fill="none" xmlns="http://www.w3.org/2000/svg">
    <rect x="10" y="10" width="52" height="52" rx="14" fill="${c}" opacity=".25"/>
    <path d="M26 20 Q26 14 36 14 Q46 14 46 20 L46 32 Q46 38 36 38 L26 38 Q18 38 18 44 L18 52 Q18 58 28 58 L44 58 Q52 58 52 52 L52 44 Q52 38 42 38 L36 38" stroke="${c}" stroke-width="3.5" stroke-linecap="round" fill="none" opacity=".8"/>
    <circle cx="30" cy="26" r="3" fill="${c}"/>
    <circle cx="42" cy="46" r="3" fill="${c}" opacity=".6"/>
  </svg>`,
}
