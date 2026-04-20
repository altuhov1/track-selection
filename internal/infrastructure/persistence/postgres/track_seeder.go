package postgres

import (
	"context"
	"log/slog"
	"time"
	"track-selection/internal/domain/track"

	"github.com/google/uuid"
)

func SeedTracks(ctx context.Context, repo *TrackRepository) {
	existing, _ := repo.FindAll(ctx)
	if len(existing) > 0 {
		slog.Info("Tracks already seeded, skipping")
		return
	}

	now := time.Now()

	// ==================== МАТЕМАТИЧЕСКИЙ ТРЕК ====================
	mathCurriculum := track.Curriculum{
		Years: []track.YearPlan{
			{
				Year: 3,
				Type: "single",
				Track: &track.YearTrack{
					Name:        "Математический трек",
					Description: "Углубленное изучение математики",
					Semesters: []track.Semester{
						{
							Number: 3,
							Courses: []track.Course{
								{Name: "Специальные главы математического анализа", Description: "Углубленное изучение матанализа", IsElective: false},
							},
						},
						{
							Number: 4,
							Courses: []track.Course{
								{Name: "Алгебраические структуры", Description: "Группы, кольца, поля", IsElective: false},
							},
						},
						{
							Number: 5,
							Courses: []track.Course{
								{Name: "Статистический анализ и основы биостатистики", Description: "Статистические методы", IsElective: false},
								{Name: "Введение в методы машинного обучения, Ч.1-2", Description: "Основы ML", IsElective: false},
							},
						},
						{
							Number: 6,
							Courses: []track.Course{
								{Name: "Численное моделирование", Description: "Численные методы", IsElective: false},
								{Name: "Выбор направления ML", Description: "Машинное или Глубокое обучение", IsElective: true, Options: []string{"Машинное обучение", "Глубокое обучение"}},
							},
						},
						{
							Number: 7,
							Courses: []track.Course{
								{Name: "Глубокое обучение в NLP", Description: "Обработка естественного языка", IsElective: false},
								{Name: "Введение в квантовые вычисления", Description: "Квантовые алгоритмы", IsElective: false},
							},
						},
						{
							Number: 8,
							Courses: []track.Course{
								{Name: "Генеративные нейронные сети и LLM модели", Description: "GPT, GAN", IsElective: false},
								{Name: "Нечеткая логика", Description: "Fuzzy logic", IsElective: false},
							},
						},
					},
				},
			},
		},
	}

	// ==================== ИНЖЕНЕРНЫЙ ТРЕК ====================
	engCurriculum := track.Curriculum{
		Years: []track.YearPlan{
			{
				Year: 3,
				Type: "single",
				Track: &track.YearTrack{
					Name:        "Инженерный трек - базовый",
					Description: "Базовые дисциплины",
					Semesters: []track.Semester{
						{
							Number: 3,
							Courses: []track.Course{
								{Name: "Специальные главы математического анализа", Description: "Углубленное изучение матанализа", IsElective: false},
							},
						},
						{
							Number: 4,
							Courses: []track.Course{
								{Name: "Алгебраические структуры", Description: "Группы, кольца, поля", IsElective: false},
							},
						},
					},
				},
			},
			{
				Year: 5,
				Type: "branching",
				Branches: []track.YearBranch{
					{
						Name:        "Искусственный интеллект",
						Description: "Специализация по ИИ",
						Semesters: []track.Semester{
							{
								Number: 5,
								Courses: []track.Course{
									{Name: "Машинное обучение", Description: "Основы ML", IsElective: false},
								},
							},
							{
								Number: 6,
								Courses: []track.Course{
									{Name: "Глубокое обучение", Description: "Нейронные сети", IsElective: false},
									{Name: "Выбор подподтрека", Description: "Выберите направление", IsElective: true, Options: []string{"Обработка сигналов", "NLP", "Компьютерное зрение", "Промпт-инжиниринг"}},
								},
							},
							{
								Number: 7,
								Courses: []track.Course{
									{Name: "Специализация по выбору", Description: "Продолжение выбранного подтрека", IsElective: true, Options: []string{"Обработка аудиосигналов", "Глубокое обучение в NLP", "Основы компьютерного зрения", "Коммуникация с ИИ"}},
								},
							},
							{
								Number: 8,
								Courses: []track.Course{
									{Name: "Завершающий курс", Description: "По выбранной специализации", IsElective: true, Options: []string{"Визуализация данных", "Генеративные нейронные сети и LLMs", "Федеративное обучение", "Разработка цифровых сервисов на базе ИИ"}},
								},
							},
						},
					},
					{
						Name:        "Аналитик",
						Description: "Специализация по аналитике",
						Semesters: []track.Semester{
							{
								Number: 5,
								Courses: []track.Course{
									{Name: "Анализ требований, ч.1-2", Description: "Требования к ПО", IsElective: false},
									{Name: "Выбор направления", Description: "Бизнес-анализ или Системный анализ", IsElective: true, Options: []string{"Бизнес-анализ", "Системный анализ"}},
								},
							},
							{
								Number: 7,
								Courses: []track.Course{
									{Name: "Проектирование графического пользовательского интерфейса", Description: "UI/UX дизайн", IsElective: false},
								},
							},
							{
								Number: 8,
								Courses: []track.Course{
									{Name: "Анализ процессов / Process Mining", Description: "Process Mining", IsElective: false},
									{Name: "Визуализация данных", Description: "Data Visualization", IsElective: false},
								},
							},
						},
					},
					{
						Name:        "Разработка ПО",
						Description: "Специализация по разработке",
						Semesters: []track.Semester{
							{
								Number: 5,
								Courses: []track.Course{
									{Name: "Выбор направления разработки", Description: "Выберите направление", IsElective: true, Options: []string{"Основы бэкенд-разработки на Python", "Основы фронтенд-разработки", "Моб. разработка на Android"}},
								},
							},
							{
								Number: 6,
								Courses: []track.Course{
									{Name: "Выбор подподтрека", Description: "Выберите специализацию", IsElective: true, Options: []string{"Высоконагруженные системы", "Блокчейн", "Большие Данные"}},
								},
							},
							{
								Number: 7,
								Courses: []track.Course{
									{Name: "Специализация по выбору", Description: "Продолжение выбранного подтрека", IsElective: true, Options: []string{"Параллельные системы", "Архитектура вычислительных сетей", "Большие данные"}},
								},
							},
							{
								Number: 8,
								Courses: []track.Course{
									{Name: "Завершающий курс", Description: "По выбранной специализации", IsElective: true, Options: []string{"Распределенные вычислительные системы", "Сети блокчейн", "Визуализация данных", "Федеративное обучение"}},
								},
							},
						},
					},
				},
			},
		},
	}

	// ==================== АППАРАТНО-ПРОГРАММНОЕ ОБЕСПЕЧЕНИЕ ====================
	hardwareCurriculum := track.Curriculum{
		Years: []track.YearPlan{
			{
				Year: 3,
				Type: "single",
				Track: &track.YearTrack{
					Name:        "АПО - базовый",
					Description: "Базовые дисциплины",
					Semesters: []track.Semester{
						{
							Number: 3,
							Courses: []track.Course{
								{Name: "Основы программирования на языке Ассемблера", Description: "Низкоуровневое программирование", IsElective: false},
							},
						},
						{
							Number: 4,
							Courses: []track.Course{
								{Name: "Низкоуровневое программирование", Description: "Работа с памятью, указатели", IsElective: false},
							},
						},
					},
				},
			},
			{
				Year: 5,
				Type: "branching",
				Branches: []track.YearBranch{
					{
						Name:        "Сети",
						Description: "Сетевая специализация",
						Semesters: []track.Semester{
							{
								Number: 5,
								Courses: []track.Course{
									{Name: "Моб. разработка на Android, ч.1-2", Description: "Android разработка", IsElective: false},
									{Name: "Основы сетевых технологий", Description: "TCP/IP, маршрутизация", IsElective: false},
								},
							},
							{
								Number: 6,
								Courses: []track.Course{
									{Name: "Инженерная инфраструктура корпоративных сетей", Description: "Корпоративные сети", IsElective: false},
								},
							},
							{
								Number: 7,
								Courses: []track.Course{
									{Name: "Сетевое проектирование и администрирование", Description: "Проектирование сетей", IsElective: false},
								},
							},
							{
								Number: 8,
								Courses: []track.Course{
									{Name: "Мониторинг и управление оборудованием в сетях", Description: "Управление сетевым оборудованием", IsElective: false},
								},
							},
						},
					},
					{
						Name:        "Средства АСУ",
						Description: "Автоматизация",
						Semesters: []track.Semester{
							{
								Number: 5,
								Courses: []track.Course{
									{Name: "Современные средства автоматизации", Description: "Автоматизация процессов", IsElective: false},
								},
							},
							{
								Number: 6,
								Courses: []track.Course{
									{Name: "ПЛК в задачах автоматизации", Description: "Программируемые логические контроллеры", IsElective: false},
								},
							},
							{
								Number: 7,
								Courses: []track.Course{
									{Name: "Интеллектуальные устройства", Description: "Умные устройства", IsElective: false},
								},
							},
							{
								Number: 8,
								Courses: []track.Course{
									{Name: "Микропроцессорная техника", Description: "Микроконтроллеры", IsElective: false},
								},
							},
						},
					},
				},
			},
		},
	}

	// ==================== ТЕСТИРОВАНИЕ ====================
	testingCurriculum := track.Curriculum{
		Years: []track.YearPlan{
			{
				Year: 3,
				Type: "single",
				Track: &track.YearTrack{
					Name:        "Тестирование - базовый",
					Description: "Базовые дисциплины тестирования",
					Semesters: []track.Semester{
						{
							Number: 3,
							Courses: []track.Course{
								{Name: "Основы тестирования программного обеспечения", Description: "Введение в тестирование", IsElective: false},
							},
						},
						{
							Number: 4,
							Courses: []track.Course{
								{Name: "Тестирование программного обеспечения", Description: "Методы тестирования", IsElective: false},
							},
						},
					},
				},
			},
			{
				Year: 5,
				Type: "branching",
				Branches: []track.YearBranch{
					{
						Name:        "Бизнес-аналитик",
						Description: "Аналитическая специализация",
						Semesters: []track.Semester{
							{
								Number: 5,
								Courses: []track.Course{
									{Name: "Анализ требований, ч.1-2", Description: "Требования к ПО", IsElective: false},
									{Name: "Выбор направления", Description: "Бизнес-анализ или Системный анализ", IsElective: true, Options: []string{"Бизнес-анализ", "Системный анализ"}},
								},
							},
							{
								Number: 7,
								Courses: []track.Course{
									{Name: "Проектирование графического пользовательского интерфейса", Description: "UI/UX дизайн", IsElective: false},
								},
							},
							{
								Number: 8,
								Courses: []track.Course{
									{Name: "Качество и метрология программного обеспечения", Description: "Качество ПО", IsElective: false},
									{Name: "Инженерный документооборот", Description: "Документооборот", IsElective: false},
								},
							},
						},
					},
					{
						Name:        "Web-разработка",
						Description: "Веб-специализация",
						Semesters: []track.Semester{
							{
								Number: 5,
								Courses: []track.Course{
									{Name: "Основы фронтенд-разработки, ч.1-2", Description: "Frontend", IsElective: false},
								},
							},
							{
								Number: 7,
								Courses: []track.Course{
									{Name: "Проектирование графического пользовательского интерфейса", Description: "UI/UX дизайн", IsElective: false},
								},
							},
							{
								Number: 8,
								Courses: []track.Course{
									{Name: "Качество и метрология программного обеспечения", Description: "Качество ПО", IsElective: false},
									{Name: "Инженерный документооборот", Description: "Документооборот", IsElective: false},
								},
							},
						},
					},
					{
						Name:        "Тестировщик",
						Description: "Техническая специализация",
						Semesters: []track.Semester{
							{
								Number: 5,
								Courses: []track.Course{
									{Name: "Автоматизация тестирования", Description: "Автотесты", IsElective: false},
								},
							},
							{
								Number: 6,
								Courses: []track.Course{
									{Name: "Промышленное тестирование ПО", Description: "Тестирование в индустрии", IsElective: false},
								},
							},
							{
								Number: 7,
								Courses: []track.Course{
									{Name: "Проектирование графического пользовательского интерфейса", Description: "UI/UX дизайн", IsElective: false},
								},
							},
							{
								Number: 8,
								Courses: []track.Course{
									{Name: "Качество и метрология программного обеспечения", Description: "Качество ПО", IsElective: false},
									{Name: "Инженерный документооборот", Description: "Документооборот", IsElective: false},
								},
							},
						},
					},
				},
			},
		},
	}

	tracks := []*track.Track{
		{
			ID:                  uuid.New().String(),
			Name:                "Математический трек",
			Description:         "Углубленное изучение математики, статистики и машинного обучения. Подготовка к работе в Data Science и AI.",
			Curriculum:          mathCurriculum,
			Requirements:        []track.Requirement{{Subject: "math_analysis", MinGrade: 4}, {Subject: "aig", MinGrade: 4}},
			Teachers:            []string{"Проф. Иванов", "Доц. Петрова"},
			Difficulty:          4,
			Type:                1,
			EmploymentProspects: 9,
			AlumniReviews:       9,
			WebLink:             "https://example.com/math-track",
			HasCertificates:     1,
			LearningStyle:       3,
			DesiredTechSkills:   6,
			DesiredMathSkills:   9,
			DesiredSoftSkills:   5,
			ProfessionalGoals:   []int{2, 3},
			CreatedAt:           now,
			UpdatedAt:           now,
		},
		{
			ID:                  uuid.New().String(),
			Name:                "Инженерный трек",
			Description:         "Разработка ПО, искусственный интеллект, аналитика. Широкий выбор специализаций.",
			Curriculum:          engCurriculum,
			Requirements:        []track.Requirement{{Subject: "programming", MinGrade: 4}},
			Teachers:            []string{"Проф. Смирнов", "Доц. Козлова"},
			Difficulty:          3,
			Type:                2,
			EmploymentProspects: 9,
			AlumniReviews:       8,
			WebLink:             "https://example.com/eng-track",
			HasCertificates:     1,
			LearningStyle:       2,
			DesiredTechSkills:   8,
			DesiredMathSkills:   6,
			DesiredSoftSkills:   7,
			ProfessionalGoals:   []int{1, 2},
			CreatedAt:           now,
			UpdatedAt:           now,
		},
		{
			ID:                  uuid.New().String(),
			Name:                "Аппаратно-программное обеспечение",
			Description:         "Низкоуровневое программирование, сети, автоматизация, микропроцессорная техника.",
			Curriculum:          hardwareCurriculum,
			Requirements:        []track.Requirement{{Subject: "programming", MinGrade: 3}, {Subject: "os_knowledge", MinGrade: 3}},
			Teachers:            []string{"Проф. Орлов", "Доц. Соколова"},
			Difficulty:          4,
			Type:                2,
			EmploymentProspects: 8,
			AlumniReviews:       8,
			WebLink:             "https://example.com/hardware-track",
			HasCertificates:     1,
			LearningStyle:       1,
			DesiredTechSkills:   8,
			DesiredMathSkills:   4,
			DesiredSoftSkills:   5,
			ProfessionalGoals:   []int{1, 4},
			CreatedAt:           now,
			UpdatedAt:           now,
		},
		{
			ID:                  uuid.New().String(),
			Name:                "Тестирование ПО",
			Description:         "Качество ПО, автоматизация тестирования, аналитика требований.",
			Curriculum:          testingCurriculum,
			Requirements:        []track.Requirement{{Subject: "programming", MinGrade: 3}},
			Teachers:            []string{"Проф. Михайлов", "Доц. Андреева"},
			Difficulty:          2,
			Type:                3,
			EmploymentProspects: 8,
			AlumniReviews:       7,
			WebLink:             "https://example.com/testing-track",
			HasCertificates:     1,
			LearningStyle:       2,
			DesiredTechSkills:   6,
			DesiredMathSkills:   3,
			DesiredSoftSkills:   8,
			ProfessionalGoals:   []int{1, 4},
			CreatedAt:           now,
			UpdatedAt:           now,
		},
	}

	for _, t := range tracks {
		if err := repo.Save(ctx, t); err != nil {
			slog.Error("Failed to seed track", "name", t.Name, "error", err)
		} else {
			slog.Info("Seeded track", "name", t.Name, "id", t.ID)
		}
	}
}
