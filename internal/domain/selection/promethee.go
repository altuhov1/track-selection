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

// PreferenceFunctionType — 6 generalized criteria from Brans & Vincke,
// see Papathanasiou & Ploskas, Fig. 3.1 (p. 60).
type PreferenceFunctionType int

const (
	UsualCriterion    PreferenceFunctionType = 1
	UShapeCriterion   PreferenceFunctionType = 2
	VShapeCriterion   PreferenceFunctionType = 3
	LevelCriterion    PreferenceFunctionType = 4
	LinearCriterion   PreferenceFunctionType = 5
	GaussianCriterion PreferenceFunctionType = 6
)

type CriterionParams struct {
	Type PreferenceFunctionType
	Q    float64
	P    float64
	S    float64
}

// preference returns P_j(d) for deviation d = g_j(a) − g_j(b).
func (cp CriterionParams) preference(d float64) float64 {
	if d <= 0 {
		return 0
	}
	switch cp.Type {
	case UsualCriterion:
		return 1
	case UShapeCriterion:
		if d > cp.Q {
			return 1
		}
		return 0
	case VShapeCriterion:
		if cp.P == 0 {
			return 1
		}
		if d > cp.P {
			return 1
		}
		return d / cp.P
	case LevelCriterion:
		if d > cp.P {
			return 1
		}
		if d > cp.Q {
			return 0.5
		}
		return 0
	case LinearCriterion:
		if d > cp.P {
			return 1
		}
		if d > cp.Q {
			return (d - cp.Q) / (cp.P - cp.Q)
		}
		return 0
	case GaussianCriterion:
		if cp.S == 0 {
			return 1
		}
		return 1 - math.Exp(-(d*d)/(2*cp.S*cp.S))
	}
	return 0
}

// criterionOrder fixes the iteration order of criteria — Go map iteration is
// non-deterministic, and PROMETHEE pairwise comparisons must be reproducible.
var criterionOrder = []string{
	"professional_goals",
	"employment",
	"alumni_reviews",
	"difficulty",
	"certificates",
	"learning_style",
	"desired_tech_skills",
	"desired_math_skills",
	"desired_soft_skills",
}

// DefaultCriterionParams: all decision-matrix values are normalised to [0, 1]
// with "higher is better". Continuous criteria use V-shape with p=1 (preference
// grows linearly with the gap, saturates at 1). Binary 0/1 criteria use Usual.
func DefaultCriterionParams() map[string]CriterionParams {
	vshape := CriterionParams{Type: VShapeCriterion, P: 1.0}
	usual := CriterionParams{Type: UsualCriterion}
	return map[string]CriterionParams{
		"professional_goals":  vshape,
		"employment":          vshape,
		"alumni_reviews":      vshape,
		"difficulty":          vshape,
		"certificates":        usual,
		"learning_style":      usual,
		"desired_tech_skills": vshape,
		"desired_math_skills": vshape,
		"desired_soft_skills": vshape,
	}
}

type PrometheeCalculator struct {
	weights         CriteriaWeights
	criterionParams map[string]CriterionParams
}

func NewPrometheeCalculator(weights CriteriaWeights) *PrometheeCalculator {
	return &PrometheeCalculator{
		weights:         weights,
		criterionParams: DefaultCriterionParams(),
	}
}

func NewPrometheeCalculatorWithParams(weights CriteriaWeights, params map[string]CriterionParams) *PrometheeCalculator {
	return &PrometheeCalculator{weights: weights, criterionParams: params}
}

// CalculateScores implements PROMETHEE II per Papathanasiou & Ploskas, Ch. 3:
//  1. Filter alternatives that fail hard requirements.
//  2. Build decision matrix g_j(a) ∈ [0, 1] (eq. 3.1, table 3.2).
//  3. For each pair (a, b), aggregate per-criterion preferences:
//       π(a, b) = Σ_j w_j · P_j(g_j(a) − g_j(b)),  Σ w_j = 1   (eq. 3.6).
//  4. Positive flow Φ⁺(a) = (1/(m−1)) Σ_x π(a, x)            (eq. 3.8).
//     Negative flow Φ⁻(a) = (1/(m−1)) Σ_x π(x, a)            (eq. 3.9).
//     Net flow      Φ(a)  = Φ⁺(a) − Φ⁻(a) ∈ [−1, 1]          (eq. 3.11).
//  5. Rank by Φ descending (PROMETHEE II rule, eq. 3.12).
//
// Score is reported as (Φ + 1) / 2 ∈ [0, 1] — a monotone transform of net
// flow that preserves the PROMETHEE II ranking exactly while keeping the
// number positive (a "match percentage", consistent with the prior API).
func (p *PrometheeCalculator) CalculateScores(tracks []PrometheeInput, student StudentData) []TrackScore {
	type candidate struct {
		input    PrometheeInput
		criteria map[string]float64
	}

	var candidates []candidate
	for _, track := range tracks {
		if !p.meetsRequirements(track, student) {
			continue
		}
		candidates = append(candidates, candidate{
			input:    track,
			criteria: p.evaluateCriteria(track, student),
		})
	}

	n := len(candidates)
	if n == 0 {
		return nil
	}

	weights := p.normalizedWeights()

	if n == 1 {
		var ws float64
		for _, name := range criterionOrder {
			ws += weights[name] * candidates[0].criteria[name]
		}
		return []TrackScore{{
			TrackID:        candidates[0].input.TrackID,
			TrackName:      candidates[0].input.TrackName,
			Score:          ws,
			Rank:           1,
			CriteriaScores: candidates[0].criteria,
		}}
	}

	pi := make([][]float64, n)
	for i := range pi {
		pi[i] = make([]float64, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			var aggregated float64
			for _, name := range criterionOrder {
				w := weights[name]
				if w == 0 {
					continue
				}
				d := candidates[i].criteria[name] - candidates[j].criteria[name]
				aggregated += w * p.criterionParams[name].preference(d)
			}
			pi[i][j] = aggregated
		}
	}

	scores := make([]TrackScore, n)
	for i := 0; i < n; i++ {
		var phiPlus, phiMinus float64
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			phiPlus += pi[i][j]
			phiMinus += pi[j][i]
		}
		phiPlus /= float64(n - 1)
		phiMinus /= float64(n - 1)
		net := phiPlus - phiMinus
		scores[i] = TrackScore{
			TrackID:        candidates[i].input.TrackID,
			TrackName:      candidates[i].input.TrackName,
			Score:          (net + 1) / 2,
			CriteriaScores: candidates[i].criteria,
		}
	}

	sort.SliceStable(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})
	for i := range scores {
		scores[i].Rank = i + 1
	}
	return scores
}

// evaluateCriteria builds one row of the decision matrix. Every value is
// normalised to [0, 1], maximised — the orientation PROMETHEE expects when
// computing d = g(a) − g(b).
func (p *PrometheeCalculator) evaluateCriteria(track PrometheeInput, student StudentData) map[string]float64 {
	return map[string]float64{
		"professional_goals":  p.calcProfessionalGoalsMatch(track.ProfessionalGoals, student.ProfessionalGoals),
		"employment":          float64(track.Employment) / 10.0,
		"alumni_reviews":      float64(track.AlumniReviews) / 10.0,
		"difficulty":          p.calcDifficultyMatch(track.Difficulty, student.Grades),
		"certificates":        p.calcCertificatesScore(track.HasCertificates, student.Certificates),
		"learning_style":      p.calcLearningStyleMatch(track.LearningStyle, student.LearningStyle),
		"desired_tech_skills": p.calcTechSkillsMatch(track.DesiredTechSkills, student.Grades, student.Skills),
		"desired_math_skills": p.calcMathSkillsMatch(track.DesiredMathSkills, student.Grades, student.Skills),
		"desired_soft_skills": p.calcSoftSkillsMatch(track.DesiredSoftSkills, student.Grades, student.Skills),
	}
}

// normalizedWeights returns w_j with Σ w_j = 1, as required by eq. (3.2).
func (p *PrometheeCalculator) normalizedWeights() map[string]float64 {
	raw := map[string]float64{
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
	var total float64
	for _, w := range raw {
		total += w
	}
	if total == 0 {
		return raw
	}
	for k, w := range raw {
		raw[k] = w / total
	}
	return raw
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
