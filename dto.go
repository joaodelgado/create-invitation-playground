package main

type InvitationBatch struct {
	ids []string
}

type ClientOrder struct {
	name                  string
	country               string
	locale                string
	hasWHPlus             bool
	hasFM                 bool
	hasDigitalPlan        bool
	hasInternationCheckin bool
}

type Membership struct {
	locale string
	name   string
}

type Plan struct {
	planName           string
	originalPrice      float32
	discountedPrice    float32
	discountPercentage float32
}

type HydratorDependencies struct{}

func (HydratorDependencies) LoadBestPlan() (Plan, error) {
	// Would be an external call
	return Plan{
		planName:           "Silver",
		originalPrice:      10,
		discountedPrice:    8,
		discountPercentage: 20,
	}, nil
}

type HydratorDTO struct {
	eligibleID  string
	isMember    bool
	clientOrder *ClientOrder
	membership  *Membership
	_bestPlan   *Plan

	deps HydratorDependencies
}

func (dto HydratorDTO) GetBestPlan() (Plan, error) {
	if dto._bestPlan == nil {
		// External request
		bestPlan, err := dto.deps.LoadBestPlan()
		if err != nil {
			return Plan{}, err
		}

		dto._bestPlan = &bestPlan
	}
	return *dto._bestPlan, nil
}

type ExperimentData struct {
	experimentEnabled bool
	name              string
	scenario          string
}

type InvitationContext struct {
	clientName      string
	userName        string
	recommendedPlan Plan
}

type TemplateData struct {
	template string
	subject  string
}
