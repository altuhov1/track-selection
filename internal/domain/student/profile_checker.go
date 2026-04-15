package student

type ProfileChecker struct {
}

func NewProfileChecker() *ProfileChecker {
	return &ProfileChecker{}
}

func (c *ProfileChecker) IsProfileComplete(prefs *Preferences) bool {
	// Проверяем профессиональные цели
	if len(prefs.ProfessionalGoals) == 0 {
		return false
	}

	// Проверяем академические оценки (все должны быть от 2 до 5)
	if !c.areGradesComplete(prefs) {
		return false
	}

	// Проверяем навыки (хотя бы один не 0)
	if c.areSkillsEmpty(prefs) {
		return false
	}

	// Проверяем стиль обучения
	if prefs.LearningStyle == 0 {
		return false
	}

	// Проверяем сертификаты (должны быть 0 или 1, но не проверяем на заполнение)
	// Сертификаты не обязательны для полноты профиля

	return true
}

func (c *ProfileChecker) areGradesComplete(prefs *Preferences) bool {
	g := &prefs.Grades

	// Проверяем, что все оценки заполнены (не 0) и в диапазоне 2-5
	if g.Informatics < 2 || g.Informatics > 5 {
		return false
	}
	if g.Programming < 2 || g.Programming > 5 {
		return false
	}
	if g.ForeignLanguage < 2 || g.ForeignLanguage > 5 {
		return false
	}
	if g.Physics < 2 || g.Physics > 5 {
		return false
	}
	if g.AIG < 2 || g.AIG > 5 {
		return false
	}
	if g.MathAnalysis < 2 || g.MathAnalysis > 5 {
		return false
	}
	if g.AlgorithmsDataStructures < 2 || g.AlgorithmsDataStructures > 5 {
		return false
	}
	if g.Databases < 2 || g.Databases > 5 {
		return false
	}
	if g.DiscreteMath < 2 || g.DiscreteMath > 5 {
		return false
	}

	return true
}

func (c *ProfileChecker) areSkillsEmpty(prefs *Preferences) bool {
	s := &prefs.Skills
	return s.Databases == 0 && s.SystemArchitecture == 0 &&
		s.AlgorithmicProgramming == 0 && s.PublicSpeaking == 0 &&
		s.Testing == 0 && s.Analytics == 0 &&
		s.MachineLearning == 0 && s.OSKnowledge == 0 &&
		s.ResearchProjects == 0
}
