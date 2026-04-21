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
			PartnerHydrator{},
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
			SubscribeDefaultTemplate{},
		},
	}

	abTestChain := ABTestChain{
		[]ABTest{
			SignupWHPlusSubject{},
			SubscribeInternationalCheckIn{},
		},
	}

	featureToggle := FeatureToggle{}

	//////////

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

		templateData, experimentData, _ := abTestChain.choose(featureToggle, dto, templateData)

		// Convert to KNS event
		// Persist and publish events
		fmt.Printf("%v %v", templateData, experimentData)
	}

}
