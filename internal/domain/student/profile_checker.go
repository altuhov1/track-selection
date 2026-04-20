package student

type ProfileChecker struct {
}

func NewProfileChecker() *ProfileChecker {
	return &ProfileChecker{}
}

func (c *ProfileChecker) IsProfileComplete(prefs *Preferences) bool {
	if len(prefs.ProfessionalGoals) == 0 {
		return false
	}

	if !c.areGradesComplete(prefs) {
		return false
	}

	if c.areSkillsEmpty(prefs) {
		return false
	}

	if prefs.LearningStyle == 0 {
		return false
	}

	return true
}

func (c *ProfileChecker) areGradesComplete(prefs *Preferences) bool {
	g := &prefs.Grades
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
