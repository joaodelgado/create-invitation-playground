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
		country:        "BR",
		locale:         "pt-BR",
		hasWHPlus:      true,
		hasFM:          true,
		hasDigitalPlan: false,
	}
}

type MembershipHydrator struct{}

func (MembershipHydrator) hydrate(dto *HydratorDTO) {
	// Would be an external call
	dto.membership = &Membership{
		locale: "pt-BR",
	}
}
