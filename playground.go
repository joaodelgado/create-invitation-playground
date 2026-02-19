package main

import (
	"fmt"
)

type FeatureToggle struct{}

func (FeatureToggle) IsEnabled(feature, eligible string) bool {
	return true
}

func main() {
	hydratorChain := HydratorChain{
		[]Hydrator{
			ClientOrderHydrator{},
			MembershipHydrator{},
			PlanHydrator{},
		},
	}

	templateChain := TemplateChain{ // TODO make explicit that the order of these templates matter
		[]Template{
			WHPlusFMTemplate{},
			WHPlusTemplate{},
			FMTemplate{},
			InternationalCheckin{},
			Default{},
		},
	}

	featureToggle := FeatureToggle{}

	batch := InvitationBatch{}

	for _, id := range batch.ids {
		dto := HydratorDTO{eligibleID: id}

		hydratorChain.hydrate(&dto)
		templateData, experimentData, _ := templateChain.choose(featureToggle, dto)

		// Persistence and publishing steps
		fmt.Printf("%v %v", templateData, experimentData)
	}

}
