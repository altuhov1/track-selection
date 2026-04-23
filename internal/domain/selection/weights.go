package selection

type CriteriaWeights struct {
	ProfessionalGoals float64 `json:"professional_goals"`
	Employment        float64 `json:"employment"`
	AlumniReviews     float64 `json:"alumni_reviews"`
	Difficulty        float64 `json:"difficulty"`
	Certificates      float64 `json:"certificates"`
	LearningStyle     float64 `json:"learning_style"`
	DesiredTechSkills float64 `json:"desired_tech_skills"`
	DesiredMathSkills float64 `json:"desired_math_skills"`
	DesiredSoftSkills float64 `json:"desired_soft_skills"`
}

func DefaultWeights() CriteriaWeights {
	return CriteriaWeights{
		ProfessionalGoals: 5.0,
		Employment:        3.0,
		AlumniReviews:     3.0,
		Difficulty:        2.0,
		Certificates:      2.0,
		LearningStyle:     2.0,
		DesiredTechSkills: 1.0,
		DesiredMathSkills: 1.0,
		DesiredSoftSkills: 1.0,
	}
}
