package main

import (
	"fmt"
)

type FeatureToggle struct{}

func (FeatureToggle) IsEnabled(feature, eligible string) (bool, string) {
	return true, "variant_a"
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
			// How to introduce an AB test here?
			SubscribeFMTemplate{},
			SubscribeDefaultTemplate{},
		},
	}

	abTestChain := ABTestChain{
		[]ABTest{
			SubscribeInternationalCheckIn{},
		},
	}

	featureToggle := FeatureToggle{}

	//////////

	batch := InvitationBatch{} // Dummy batch

	for _, id := range batch.ids {
		// Create invitation
		dto := HydratorDTO{eligibleID: id, isMember: false}

		hydratorChain.hydrate(&dto)

		var templateData TemplateData
		var experimentData *ExperimentData
		if !dto.isMember {
			templateData, experimentData, _ = signupTemplateChain.choose(featureToggle, dto)
		} else {
			templateData, experimentData, _ = subscribeTemplateChain.choose(featureToggle, dto)
		}

		if experimentData == nil {
			templateData, experimentData, _ = abTestChain.choose(featureToggle, dto, templateData)
		}

		// Convert to KNS event
		// Persist and publish events
		fmt.Printf("%v %v", templateData, experimentData)
	}

}
