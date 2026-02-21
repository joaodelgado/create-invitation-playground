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

func (HydratorDependencies) LoadClientOrder() (ClientOrder, error) {
	// Would be an external call
	return ClientOrder{
		country:        "BR",
		locale:         "pt-BR",
		hasWHPlus:      true,
		hasFM:          true,
		hasDigitalPlan: false,
	}, nil
}

func (HydratorDependencies) LoadMembership() (Membership, error) {
	// Would be an external call
	return Membership{
		locale: "pt-BR",
	}, nil
}

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
	eligibleID   string
	isMember     bool
	_clientOrder *ClientOrder
	_membership  *Membership
	_bestPlan    *Plan

	deps HydratorDependencies
}

func (dto HydratorDTO) GetClientOrder() (ClientOrder, error) {
	if dto._clientOrder == nil {
		// External request
		clientOrder, err := dto.deps.LoadClientOrder()
		if err != nil {
			return ClientOrder{}, err
		}

		dto._clientOrder = &clientOrder
	}

	return *dto._clientOrder, nil
}

func (dto HydratorDTO) GetMembership() (*Membership, error) {
	if !dto.isMember {
		return nil, nil
	}
	if dto._membership == nil {
		// External request
		membership, err := dto.deps.LoadMembership()
		if err != nil {
			return nil, err
		}

		dto._membership = &membership
	}

	return dto._membership, nil
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
