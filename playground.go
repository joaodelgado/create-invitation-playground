package main

import (
	"fmt"
)

type FeatureToggle struct{}

func (FeatureToggle) IsEnabled(feature, eligible string) bool {
	return true
}

func main() {
	// TODO make explicit that the order of these chain items matter
	hydratorChain := HydratorChain{
		[]Hydrator{
			ClientOrderHydrator{},
			MembershipHydrator{},
			PlanHydrator{},
		},
	}

	signupTemplateChain := TemplateChain{
		[]Template{
			SignupWHPlusTemplate{},
			SignupDigitalTemplate{},
			SignupDefaultTemplate{},
		},
	}

	subscribeTemplateChain := TemplateChain{
		[]Template{
			SubscribeWHPlusFMTemplate{},
			SubscribeWHPlusTemplate{},
			SubscribeFMTemplate{},
			SubscribeInternationalCheckingTemplate{},
			SubscribeDefaultTemplate{},
		},
	}

	abTestChain := ABTestChain{
		[]ABTest{
			SignupWHPlusSubject{},
		},
	}

	featureToggle := FeatureToggle{}

	batch := InvitationBatch{} // Dummy batch

	for _, id := range batch.ids {
		dto := HydratorDTO{eligibleID: id, isMember: false}

		hydratorChain.hydrate(&dto)

		var templateData TemplateData
		if !dto.isMember {
			templateData, _ = signupTemplateChain.choose(dto)
		} else {
			templateData, _ = subscribeTemplateChain.choose(dto)
		}

		templateData, experimentData := abTestChain.choose(featureToggle, dto, templateData)

		// Convert to KNS event
		// Persist and publish events
		fmt.Printf("%v %v", templateData, experimentData)
	}

}
