package track

import (
	"time"

	"github.com/google/uuid"
)

// Course — один курс
type Course struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	IsElective  bool     `json:"is_elective"`
	Options     []string `json:"options,omitempty"`
}

// Semester — семестр с курсами
type Semester struct {
	Number  int      `json:"number"`
	Courses []Course `json:"courses"`
}

// YearTrack — информация о треке на год (для типа "single")
type YearTrack struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Semesters   []Semester `json:"semesters"`
}

// YearBranch — ветка выбора (для типа "branching")
type YearBranch struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Semesters   []Semester `json:"semesters"`
}

// YearPlan — план на один год
type YearPlan struct {
	Year     int          `json:"year"`
	Type     string       `json:"type"` // "single" или "branching"
	Track    *YearTrack   `json:"track,omitempty"`
	Branches []YearBranch `json:"branches,omitempty"`
}

// Curriculum — учебный план (по годам)
type Curriculum struct {
	Years []YearPlan `json:"years"`
}

// Requirement — требования к поступающему
type Requirement struct {
	Subject  string `json:"subject"`
	MinGrade int    `json:"min_grade"`
}

// Track — образовательный трек
type Track struct {
	ID                  string        `json:"id"`
	Name                string        `json:"name"`
	Description         string        `json:"description"`
	Curriculum          Curriculum    `json:"curriculum"`
	Requirements        []Requirement `json:"requirements"`
	Teachers            []string      `json:"teachers"`
	Difficulty          int           `json:"difficulty"`
	Type                int           `json:"type"`
	EmploymentProspects int           `json:"employment_prospects"`
	AlumniReviews       int           `json:"alumni_reviews"`
	WebLink             string        `json:"web_link"`
	HasCertificates     int           `json:"has_certificates"`
	LearningStyle       int           `json:"learning_style"`
	DesiredTechSkills   int           `json:"desired_tech_skills"`
	DesiredMathSkills   int           `json:"desired_math_skills"`
	DesiredSoftSkills   int           `json:"desired_soft_skills"`
	ProfessionalGoals   []int         `json:"professional_goals"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
}

// NewTrack создает новый трек
func NewTrack(
	name, description string,
	curriculum Curriculum,
	requirements []Requirement,
	teachers []string,
	difficulty, trackType, employmentProspects, alumniReviews int,
	webLink string,
	hasCertificates, learningStyle int,
	desiredTechSkills, desiredMathSkills, desiredSoftSkills int,
	professionalGoals []int,
) (*Track, error) {
	now := time.Now()

	return &Track{
		ID:                  uuid.New().String(),
		Name:                name,
		Description:         description,
		Curriculum:          curriculum,
		Requirements:        requirements,
		Teachers:            teachers,
		Difficulty:          difficulty,
		Type:                trackType,
		EmploymentProspects: employmentProspects,
		AlumniReviews:       alumniReviews,
		WebLink:             webLink,
		HasCertificates:     hasCertificates,
		LearningStyle:       learningStyle,
		DesiredTechSkills:   desiredTechSkills,
		DesiredMathSkills:   desiredMathSkills,
		DesiredSoftSkills:   desiredSoftSkills,
		ProfessionalGoals:   professionalGoals,
		CreatedAt:           now,
		UpdatedAt:           now,
	}, nil
}

// Update обновляет поля трека
func (t *Track) Update(updates map[string]interface{}) {
	if val, ok := updates["name"]; ok {
		t.Name = val.(string)
	}
	if val, ok := updates["description"]; ok {
		t.Description = val.(string)
	}
	if val, ok := updates["curriculum"]; ok {
		if curriculum, ok := val.(map[string]interface{}); ok {
			t.Curriculum = parseCurriculum(curriculum)
		}
	}
	if val, ok := updates["requirements"]; ok {
		t.Requirements = val.([]Requirement)
	}
	if val, ok := updates["teachers"]; ok {
		t.Teachers = val.([]string)
	}
	if val, ok := updates["difficulty"]; ok {
		t.Difficulty = int(val.(float64))
	}
	if val, ok := updates["type"]; ok {
		t.Type = int(val.(float64))
	}
	if val, ok := updates["employment_prospects"]; ok {
		t.EmploymentProspects = int(val.(float64))
	}
	if val, ok := updates["alumni_reviews"]; ok {
		t.AlumniReviews = int(val.(float64))
	}
	if val, ok := updates["web_link"]; ok {
		t.WebLink = val.(string)
	}
	if val, ok := updates["has_certificates"]; ok {
		t.HasCertificates = int(val.(float64))
	}
	if val, ok := updates["learning_style"]; ok {
		t.LearningStyle = int(val.(float64))
	}
	if val, ok := updates["desired_tech_skills"]; ok {
		t.DesiredTechSkills = int(val.(float64))
	}
	if val, ok := updates["desired_math_skills"]; ok {
		t.DesiredMathSkills = int(val.(float64))
	}
	if val, ok := updates["desired_soft_skills"]; ok {
		t.DesiredSoftSkills = int(val.(float64))
	}
	if val, ok := updates["professional_goals"]; ok {
		t.ProfessionalGoals = val.([]int)
	}

	t.UpdatedAt = time.Now()
}

// parseCurriculum преобразует map в Curriculum (с поддержкой years)
func parseCurriculum(data map[string]interface{}) Curriculum {
	var curriculum Curriculum

	// Проверяем новую структуру (с years)
	if years, ok := data["years"].([]interface{}); ok {
		for _, y := range years {
			yearMap := y.(map[string]interface{})
			yearPlan := YearPlan{}

			if yearNum, ok := yearMap["year"].(float64); ok {
				yearPlan.Year = int(yearNum)
			}
			if typeStr, ok := yearMap["type"].(string); ok {
				yearPlan.Type = typeStr
			}

			// Парсим single track
			if trackData, ok := yearMap["track"].(map[string]interface{}); ok && trackData != nil {
				track := &YearTrack{}
				if name, ok := trackData["name"].(string); ok {
					track.Name = name
				}
				if desc, ok := trackData["description"].(string); ok {
					track.Description = desc
				}
				if semesters, ok := trackData["semesters"].([]interface{}); ok {
					track.Semesters = parseSemesters(semesters)
				}
				yearPlan.Track = track
			}

			// Парсим branches
			if branches, ok := yearMap["branches"].([]interface{}); ok {
				for _, b := range branches {
					branchMap := b.(map[string]interface{})
					branch := YearBranch{}

					if name, ok := branchMap["name"].(string); ok {
						branch.Name = name
					}
					if desc, ok := branchMap["description"].(string); ok {
						branch.Description = desc
					}
					if semesters, ok := branchMap["semesters"].([]interface{}); ok {
						branch.Semesters = parseSemesters(semesters)
					}

					yearPlan.Branches = append(yearPlan.Branches, branch)
				}
			}

			curriculum.Years = append(curriculum.Years, yearPlan)
		}
		return curriculum
	}

	// Поддержка старой структуры (semesters) для обратной совместимости
	if semesters, ok := data["semesters"].([]interface{}); ok {
		curriculum.Years = []YearPlan{
			{
				Year: 0,
				Type: "single",
				Track: &YearTrack{
					Name:        "",
					Description: "",
					Semesters:   parseSemesters(semesters),
				},
			},
		}
	}

	return curriculum
}

// parseSemesters парсит семестры из JSON
func parseSemesters(semestersData []interface{}) []Semester {
	var semesters []Semester

	for _, s := range semestersData {
		semesterMap := s.(map[string]interface{})
		semester := Semester{}

		if num, ok := semesterMap["number"].(float64); ok {
			semester.Number = int(num)
		}

		if courses, ok := semesterMap["courses"].([]interface{}); ok {
			for _, c := range courses {
				courseMap := c.(map[string]interface{})
				course := Course{}

				if name, ok := courseMap["name"].(string); ok {
					course.Name = name
				}
				if desc, ok := courseMap["description"].(string); ok {
					course.Description = desc
				}
				if isElective, ok := courseMap["is_elective"].(bool); ok {
					course.IsElective = isElective
				}
				if options, ok := courseMap["options"].([]interface{}); ok {
					for _, opt := range options {
						if optStr, ok := opt.(string); ok {
							course.Options = append(course.Options, optStr)
						}
					}
				}

				semester.Courses = append(semester.Courses, course)
			}
		}

		semesters = append(semesters, semester)
	}

	return semesters
}
