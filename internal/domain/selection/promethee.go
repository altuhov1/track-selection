package selection

import (
	"math"
	"sort"
)

type TrackScore struct {
	TrackID        string             `json:"track_id"`
	TrackName      string             `json:"track_name"`
	Score          float64            `json:"score"`
	Rank           int                `json:"rank"`
	CriteriaScores map[string]float64 `json:"criteria_scores"`
}

type PrometheeInput struct {
	TrackID           string
	TrackName         string
	ProfessionalGoals []int
	Employment        int
	AlumniReviews     int
	Difficulty        int
	HasCertificates   int
	LearningStyle     int
	DesiredTechSkills int
	DesiredMathSkills int
	DesiredSoftSkills int
	Requirements      []Requirement
}

type Requirement struct {
	Subject  string
	MinGrade int
}

type StudentData struct {
	ProfessionalGoals []int
	Grades            Grades
	Skills            Skills
	LearningStyle     int
	Certificates      int
}

type Grades struct {
	Informatics              int
	Programming              int
	ForeignLanguage          int
	Physics                  int
	AIG                      int
	MathAnalysis             int
	AlgorithmsDataStructures int
	Databases                int
	DiscreteMath             int
}

type Skills struct {
	Databases              int
	SystemArchitecture     int
	AlgorithmicProgramming int
	PublicSpeaking         int
	Testing                int
	Analytics              int
	MachineLearning        int
	OSKnowledge            int
	ResearchProjects       int
}

type PrometheeCalculator struct {
	weights CriteriaWeights
}

func NewPrometheeCalculator(weights CriteriaWeights) *PrometheeCalculator {
	return &PrometheeCalculator{weights: weights}
}

func (p *PrometheeCalculator) CalculateScores(tracks []PrometheeInput, student StudentData) []TrackScore {
	var scores []TrackScore

	for _, track := range tracks {
		if !p.meetsRequirements(track, student) {
			continue
		}

		criteriaScores := make(map[string]float64)

		// Профессиональные цели
		criteriaScores["professional_goals"] = p.calcProfessionalGoalsMatch(track.ProfessionalGoals, student.ProfessionalGoals)

		// Перспективы трудоустройства
		criteriaScores["employment"] = float64(track.Employment) / 10.0

		// Отзывы выпускников
		criteriaScores["alumni_reviews"] = float64(track.AlumniReviews) / 10.0

		// Сложность
		criteriaScores["difficulty"] = p.calcDifficultyMatch(track.Difficulty, student.Grades)

		// Сертификаты
		criteriaScores["certificates"] = p.calcCertificatesScore(track.HasCertificates, student.Certificates)

		// Стиль обучения
		criteriaScores["learning_style"] = p.calcLearningStyleMatch(track.LearningStyle, student.LearningStyle)

		// Желаемые навыки
		criteriaScores["desired_tech_skills"] = p.calcTechSkillsMatch(track.DesiredTechSkills, student.Grades, student.Skills)
		criteriaScores["desired_math_skills"] = p.calcMathSkillsMatch(track.DesiredMathSkills, student.Grades, student.Skills)
		criteriaScores["desired_soft_skills"] = p.calcSoftSkillsMatch(track.DesiredSoftSkills, student.Grades, student.Skills)

		totalScore := p.calculateWeightedSum(criteriaScores)

		scores = append(scores, TrackScore{
			TrackID:        track.TrackID,
			TrackName:      track.TrackName,
			Score:          totalScore,
			CriteriaScores: criteriaScores,
		})
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	for i := range scores {
		scores[i].Rank = i + 1
	}

	return scores
}

func (p *PrometheeCalculator) meetsRequirements(track PrometheeInput, student StudentData) bool {
	for _, req := range track.Requirements {
		grade := p.getGradeBySubject(req.Subject, student.Grades)
		if grade < req.MinGrade {
			return false
		}
	}
	if len(track.ProfessionalGoals) > 0 && len(student.ProfessionalGoals) > 0 {
		match := false
		for _, tg := range track.ProfessionalGoals {
			for _, sg := range student.ProfessionalGoals {
				if tg == sg {
					match = true
					break
				}
			}
			if match {
				break
			}
		}
		if !match {
			return false
		}
	}

	return true
}

// calcProfessionalGoalsMatch вычисляет совпадение профессиональных целей
func (p *PrometheeCalculator) calcProfessionalGoalsMatch(trackGoals, studentGoals []int) float64 {
	if len(trackGoals) == 0 || len(studentGoals) == 0 {
		return 0
	}

	matchCount := 0
	for _, tg := range trackGoals {
		for _, sg := range studentGoals {
			if tg == sg {
				matchCount++
				break
			}
		}
	}

	return float64(matchCount) / float64(len(trackGoals))
}

func (p *PrometheeCalculator) calcDifficultyMatch(trackDifficulty int, studentGrades Grades) float64 {
	avgGrade := p.calculateAverageGrade(studentGrades)

	var recommendedDifficulty int
	if avgGrade < 3.0 {
		recommendedDifficulty = 1
	} else if avgGrade < 4.0 {
		recommendedDifficulty = 2
	} else {
		recommendedDifficulty = 4
	}

	diff := math.Abs(float64(trackDifficulty - recommendedDifficulty))
	match := 1.0 - diff/4.0
	if match < 0 {
		match = 0
	}
	return match
}

func (p *PrometheeCalculator) calcLearningStyleMatch(trackStyle, studentStyle int) float64 {
	if trackStyle == studentStyle {
		return 1.0
	}
	return 0.0
}

func (p *PrometheeCalculator) calcCertificatesScore(trackHasCerts, studentWantsCerts int) float64 {
	if studentWantsCerts == 1 && trackHasCerts == 0 {
		return 0.0
	}
	return 1.0
}

func skillsMatchScore(studentAvg float64, trackDesired int) float64 {
	if trackDesired == 0 {
		return 1.0
	}
	score := studentAvg / float64(trackDesired)
	if score > 1.0 {
		return 1.0
	}
	return score
}

func (p *PrometheeCalculator) calcMathSkillsMatch(trackDesired int, grades Grades, skills Skills) float64 {
	avg := float64(grades.AIG+grades.MathAnalysis+grades.DiscreteMath+skills.MachineLearning) / 4.0
	return skillsMatchScore(avg, trackDesired)
}

func (p *PrometheeCalculator) calcTechSkillsMatch(trackDesired int, grades Grades, skills Skills) float64 {
	avg := float64(grades.Programming+grades.Informatics+grades.AlgorithmsDataStructures+grades.Databases+
		skills.SystemArchitecture+skills.AlgorithmicProgramming+skills.Testing+skills.OSKnowledge+skills.Databases) / 9.0
	return skillsMatchScore(avg, trackDesired)
}

func (p *PrometheeCalculator) calcSoftSkillsMatch(trackDesired int, grades Grades, skills Skills) float64 {
	avg := float64(skills.PublicSpeaking+skills.Analytics+grades.ForeignLanguage) / 3.0
	return skillsMatchScore(avg, trackDesired)
}

func (p *PrometheeCalculator) calculateWeightedSum(criteriaScores map[string]float64) float64 {
	totalWeight := 0.0
	weightedSum := 0.0

	weights := map[string]float64{
		"professional_goals":  p.weights.ProfessionalGoals,
		"employment":          p.weights.Employment,
		"alumni_reviews":      p.weights.AlumniReviews,
		"difficulty":          p.weights.Difficulty,
		"certificates":        p.weights.Certificates,
		"learning_style":      p.weights.LearningStyle,
		"desired_tech_skills": p.weights.DesiredTechSkills,
		"desired_math_skills": p.weights.DesiredMathSkills,
		"desired_soft_skills": p.weights.DesiredSoftSkills,
	}

	for key, score := range criteriaScores {
		if weight, ok := weights[key]; ok {
			weightedSum += score * weight
			totalWeight += weight
		}
	}

	if totalWeight == 0 {
		return 0
	}
	return weightedSum / totalWeight
}

func (p *PrometheeCalculator) getGradeBySubject(subject string, grades Grades) int {
	switch subject {
	case "informatics":
		return grades.Informatics
	case "programming":
		return grades.Programming
	case "foreign_language":
		return grades.ForeignLanguage
	case "physics":
		return grades.Physics
	case "aig":
		return grades.AIG
	case "math_analysis":
		return grades.MathAnalysis
	case "algorithms_data_structures":
		return grades.AlgorithmsDataStructures
	case "databases":
		return grades.Databases
	case "discrete_math":
		return grades.DiscreteMath
	default:
		return 0
	}
}

func (p *PrometheeCalculator) calculateAverageGrade(grades Grades) float64 {
	sum := float64(grades.Informatics + grades.Programming + grades.ForeignLanguage +
		grades.Physics + grades.AIG + grades.MathAnalysis +
		grades.AlgorithmsDataStructures + grades.Databases + grades.DiscreteMath)
	return sum / 9.0
}
