package student

import (
	"time"
	errCustom "track-selection/internal/domain/shared/errors"
)

type Grades struct {
	Informatics              int `json:"informatics"`                // Информатика
	Programming              int `json:"programming"`                // Программирование
	ForeignLanguage          int `json:"foreign_language"`           // Иностранный язык
	Physics                  int `json:"physics"`                    // Физика
	AIG                      int `json:"aig"`                        // АиГ (Алгебра и геометрия)
	MathAnalysis             int `json:"math_analysis"`              // Математический анализ
	AlgorithmsDataStructures int `json:"algorithms_data_structures"` // Алгоритмы и структуры данных
	Databases                int `json:"databases"`                  // Базы данных
	DiscreteMath             int `json:"discrete_math"`              // Дискретная математика
}
type Skills struct {
	Databases              int `json:"databases"`
	SystemArchitecture     int `json:"system_architecture"`
	AlgorithmicProgramming int `json:"algorithmic_programming"`
	PublicSpeaking         int `json:"public_speaking"`
	Testing                int `json:"testing"`
	Analytics              int `json:"analytics"`
	MachineLearning        int `json:"machine_learning"`
	OSKnowledge            int `json:"os_knowledge"`
	ResearchProjects       int `json:"research_projects"`
}

type Preferences struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	ProfessionalGoals []int     `json:"professional_goals"`
	Grades            Grades    `json:"grades"`
	Skills            Skills    `json:"skills"`
	LearningStyle     int       `json:"learning_style"` // 1,2,3
	Certificates      int       `json:"certificates"`   // 0 или 1
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func GenerateRandomGrades() Grades {
	return Grades{
		Informatics:              randomGrade(),
		Programming:              randomGrade(),
		ForeignLanguage:          randomGrade(),
		Physics:                  randomGrade(),
		AIG:                      randomGrade(),
		MathAnalysis:             randomGrade(),
		AlgorithmsDataStructures: randomGrade(),
		Databases:                randomGrade(),
		DiscreteMath:             randomGrade(),
	}
}

func randomGrade() int {
	return []int{4, 5}[time.Now().Nanosecond()%2]
}

// Validate validates all preferences fields
func (p *Preferences) Validate() error {
	if err := p.validateLearningStyle(); err != nil {
		return err
	}
	if err := p.validateCertificates(); err != nil {
		return err
	}
	if err := p.validateSkills(); err != nil {
		return err
	}
	if err := p.validateGrades(); err != nil {
		return err
	}
	return nil
}

// ValidatePartial validates only fields that are present in the map
func (p *Preferences) ValidatePartial(updates map[string]interface{}) error {
	if _, ok := updates["learning_style"]; ok {
		if err := p.validateLearningStyle(); err != nil {
			return err
		}
	}
	if _, ok := updates["certificates"]; ok {
		if err := p.validateCertificates(); err != nil {
			return err
		}
	}
	if _, ok := updates["skills"]; ok {
		if err := p.validateSkills(); err != nil {
			return err
		}
	}
	if _, ok := updates["grades"]; ok {
		if err := p.validateGrades(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Preferences) validateLearningStyle() error {
	if p.LearningStyle < 1 || p.LearningStyle > 3 {
		return errCustom.ErrInvalidLearningStyle
	}
	return nil
}

func (p *Preferences) validateCertificates() error {
	if p.Certificates != 0 && p.Certificates != 1 {
		return errCustom.ErrInvalidCertificate
	}
	return nil
}

func (p *Preferences) validateSkills() error {
	s := &p.Skills
	if s.Databases < 0 || s.Databases > 10 {
		return errCustom.ErrInvalidSkillValue
	}
	if s.SystemArchitecture < 0 || s.SystemArchitecture > 10 {
		return errCustom.ErrInvalidSkillValue
	}
	if s.AlgorithmicProgramming < 0 || s.AlgorithmicProgramming > 10 {
		return errCustom.ErrInvalidSkillValue
	}
	if s.PublicSpeaking < 0 || s.PublicSpeaking > 10 {
		return errCustom.ErrInvalidSkillValue
	}
	if s.Testing < 0 || s.Testing > 10 {
		return errCustom.ErrInvalidSkillValue
	}
	if s.Analytics < 0 || s.Analytics > 10 {
		return errCustom.ErrInvalidSkillValue
	}
	if s.MachineLearning < 0 || s.MachineLearning > 10 {
		return errCustom.ErrInvalidSkillValue
	}
	if s.OSKnowledge < 0 || s.OSKnowledge > 10 {
		return errCustom.ErrInvalidSkillValue
	}
	if s.ResearchProjects < 0 || s.ResearchProjects > 10 {
		return errCustom.ErrInvalidSkillValue
	}
	return nil
}

func (p *Preferences) validateGrades() error {
	g := &p.Grades
	if g.Informatics != 0 && (g.Informatics < 2 || g.Informatics > 5) {
		return errCustom.ErrInvalidGrade
	}
	if g.Programming != 0 && (g.Programming < 2 || g.Programming > 5) {
		return errCustom.ErrInvalidGrade
	}
	if g.ForeignLanguage != 0 && (g.ForeignLanguage < 2 || g.ForeignLanguage > 5) {
		return errCustom.ErrInvalidGrade
	}
	if g.Physics != 0 && (g.Physics < 2 || g.Physics > 5) {
		return errCustom.ErrInvalidGrade
	}
	if g.AIG != 0 && (g.AIG < 2 || g.AIG > 5) {
		return errCustom.ErrInvalidGrade
	}
	if g.MathAnalysis != 0 && (g.MathAnalysis < 2 || g.MathAnalysis > 5) {
		return errCustom.ErrInvalidGrade
	}
	if g.AlgorithmsDataStructures != 0 && (g.AlgorithmsDataStructures < 2 || g.AlgorithmsDataStructures > 5) {
		return errCustom.ErrInvalidGrade
	}
	if g.Databases != 0 && (g.Databases < 2 || g.Databases > 5) {
		return errCustom.ErrInvalidGrade
	}
	if g.DiscreteMath != 0 && (g.DiscreteMath < 2 || g.DiscreteMath > 5) {
		return errCustom.ErrInvalidGrade
	}
	return nil
}

// Частичное обновление
func (p *Preferences) Merge(updates map[string]interface{}) {
	if val, ok := updates["professional_goals"]; ok {
		if goals, ok := val.([]interface{}); ok {
			p.ProfessionalGoals = make([]int, len(goals))
			for i, g := range goals {
				if f, ok := g.(float64); ok {
					p.ProfessionalGoals[i] = int(f)
				}
			}
		}
	}

	if val, ok := updates["grades"]; ok {
		if grades, ok := val.(map[string]interface{}); ok {
			p.mergeGrades(grades)
		}
	}

	if val, ok := updates["skills"]; ok {
		if skills, ok := val.(map[string]interface{}); ok {
			p.mergeSkills(skills)
		}
	}

	if val, ok := updates["learning_style"]; ok {
		if style, ok := val.(float64); ok {
			p.LearningStyle = int(style)
		}
	}

	if val, ok := updates["certificates"]; ok {
		if cert, ok := val.(float64); ok {
			p.Certificates = int(cert)
		}
	}

	p.UpdatedAt = time.Now()
}

func (p *Preferences) mergeGrades(updates map[string]interface{}) {
	if val, ok := updates["informatics"]; ok {
		p.Grades.Informatics = int(val.(float64))
	}
	if val, ok := updates["programming"]; ok {
		p.Grades.Programming = int(val.(float64))
	}
	if val, ok := updates["foreign_language"]; ok {
		p.Grades.ForeignLanguage = int(val.(float64))
	}
	if val, ok := updates["physics"]; ok {
		p.Grades.Physics = int(val.(float64))
	}
	if val, ok := updates["aig"]; ok {
		p.Grades.AIG = int(val.(float64))
	}
	if val, ok := updates["math_analysis"]; ok {
		p.Grades.MathAnalysis = int(val.(float64))
	}
	if val, ok := updates["algorithms_data_structures"]; ok {
		p.Grades.AlgorithmsDataStructures = int(val.(float64))
	}
	if val, ok := updates["databases"]; ok {
		p.Grades.Databases = int(val.(float64))
	}
	if val, ok := updates["discrete_math"]; ok {
		p.Grades.DiscreteMath = int(val.(float64))
	}
}

func (p *Preferences) mergeSkills(updates map[string]interface{}) {
	if val, ok := updates["databases"]; ok {
		p.Skills.Databases = int(val.(float64))
	}
	if val, ok := updates["system_architecture"]; ok {
		p.Skills.SystemArchitecture = int(val.(float64))
	}
	if val, ok := updates["algorithmic_programming"]; ok {
		p.Skills.AlgorithmicProgramming = int(val.(float64))
	}
	if val, ok := updates["public_speaking"]; ok {
		p.Skills.PublicSpeaking = int(val.(float64))
	}
	if val, ok := updates["testing"]; ok {
		p.Skills.Testing = int(val.(float64))
	}
	if val, ok := updates["analytics"]; ok {
		p.Skills.Analytics = int(val.(float64))
	}
	if val, ok := updates["machine_learning"]; ok {
		p.Skills.MachineLearning = int(val.(float64))
	}
	if val, ok := updates["os_knowledge"]; ok {
		p.Skills.OSKnowledge = int(val.(float64))
	}
	if val, ok := updates["research_projects"]; ok {
		p.Skills.ResearchProjects = int(val.(float64))
	}
}
