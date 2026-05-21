package main

type Hydrator interface {
	hydrate(dto *HydratorDTO)
}

type HydratorChain struct {
	hidrators []Hydrator
}

func (hc HydratorChain) hydrate(dto *HydratorDTO) {
	for _, hydrator := range hc.hidrators {
		hydrator.hydrate(dto)
	}
}

///
///
///

type ClientOrderHydrator struct{}

func (ClientOrderHydrator) hydrate(dto *HydratorDTO) {
	// Would be an external call
	dto.clientOrder = &ClientOrder{
		country:               "BR",
		locale:                "pt-BR",
		hasWHPlus:             CONFIG.hasWHPlus,
		hasFM:                 CONFIG.hasFM,
		hasDigitalPlan:        CONFIG.hasDigitalPlan,
		hasInternationCheckin: CONFIG.hasInternationCheckin,
	}
}

type MembershipHydrator struct{}

func (MembershipHydrator) hydrate(dto *HydratorDTO) {
	// Would be an external call
	dto.membership = &Membership{
		locale: "pt-BR",
	}
}

type PlanHydrator struct{}

func (PlanHydrator) hydrate(dto *HydratorDTO) {
	dto.BestPlan = LazyLoaded[Plan]{
		loader: func() (Plan, error) {
			// Would be an external call
			return Plan{
				planName:           "Silver",
				originalPrice:      10,
				discountedPrice:    0,
				discountPercentage: 100,
			}, nil
		},
	}
}

type PartnerHydrator struct{}

func (PartnerHydrator) hydrate(dto *HydratorDTO) {
	dto.RecommendedPartners = LazyLoaded[[]Partner]{
		loader: func() ([]Partner, error) {
			return []Partner{}, nil
		},
	}
}
