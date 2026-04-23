package track

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	IsElective  bool     `json:"is_elective"`
	Options     []string `json:"options,omitempty"`
}

type Semester struct {
	Number  int      `json:"number"`
	Courses []Course `json:"courses"`
}

type YearTrack struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Semesters   []Semester `json:"semesters"`
}

type YearBranch struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Semesters   []Semester `json:"semesters"`
	Branches    []YearBranch
}

type YearPlan struct {
	Year     int          `json:"year"`
	Type     string       `json:"type"` // "single" или "branching"
	Track    *YearTrack   `json:"track,omitempty"`
	Branches []YearBranch `json:"branches,omitempty"`
}

type Curriculum struct {
	Years []YearPlan `json:"years"`
}

type Requirement struct {
	Subject  string `json:"subject"`
	MinGrade int    `json:"min_grade"`
}

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
		if raw, ok := val.([]interface{}); ok {
			t.Requirements = parseRequirements(raw)
		}
	}
	if val, ok := updates["teachers"]; ok {
		if raw, ok := val.([]interface{}); ok {
			t.Teachers = parseStringSlice(raw)
		}
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
		if raw, ok := val.([]interface{}); ok {
			t.ProfessionalGoals = parseIntSlice(raw)
		}
	}

	t.UpdatedAt = time.Now()
}

func parseCurriculum(data map[string]interface{}) Curriculum {
	var curriculum Curriculum

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

			if branches, ok := yearMap["branches"].([]interface{}); ok {
				yearPlan.Branches = parseBranches(branches)
			}

			curriculum.Years = append(curriculum.Years, yearPlan)
		}
		return curriculum
	}

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

func parseRequirements(data []interface{}) []Requirement {
	result := make([]Requirement, 0, len(data))
	for _, r := range data {
		m, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		req := Requirement{}
		if s, ok := m["subject"].(string); ok {
			req.Subject = s
		}
		if g, ok := m["min_grade"].(float64); ok {
			req.MinGrade = int(g)
		}
		result = append(result, req)
	}
	return result
}

func parseStringSlice(data []interface{}) []string {
	result := make([]string, 0, len(data))
	for _, v := range data {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

func parseIntSlice(data []interface{}) []int {
	result := make([]int, 0, len(data))
	for _, v := range data {
		if n, ok := v.(float64); ok {
			result = append(result, int(n))
		}
	}
	return result
}

func parseBranches(branchesData []interface{}) []YearBranch {
	var branches []YearBranch

	for _, b := range branchesData {
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
		// Рекурсивно парсим вложенные ветки (подподтреки)
		if subBranches, ok := branchMap["branches"].([]interface{}); ok {
			branch.Branches = parseBranches(subBranches)
		}

		branches = append(branches, branch)
	}

	return branches
}

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
