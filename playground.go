package main

import (
	"fmt"
)

// Helper structs to facilitate testing
type ft struct {
	enabled bool
	variant string
}

type Config struct {
	isMember bool

	hasFM                 bool
	hasDigitalPlan        bool
	hasInternationCheckin bool

	hasWHPlus           bool
	planDiscountedPrice float32

	abTests map[string]ft
}

var CONFIG Config = Config{
	isMember: true,

	hasFM:                 true,
	hasDigitalPlan:        false,
	hasInternationCheckin: true,
	hasWHPlus:             true,

	planDiscountedPrice: 10,

	abTests: map[string]ft{
		"SignupSubject":        {true, "variant_a"},
		"InternationalCheckIn": {true, "variant_a"},
		"SubscribeFMSubject":   {true, "control"},
	},
}

//////

type FeatureToggle struct{}

func (FeatureToggle) IsEnabled(feature, eligible string) (bool, string) {
	f := CONFIG.abTests[feature]
	return f.enabled, f.variant
}

func main() {
	//
	// Application boot
	//

	hydratorChain := HydratorChain{
		[]Hydrator{
			ClientOrderHydrator{},
			MembershipHydrator{},
			PlanHydrator{},
			PartnerHydrator{},
		},
	}

	signupTemplateChain := TemplateChain{
		[]TemplateDefinition{
			CreateSignupWHPlusTemplateConfig(),
			CreateSignupDigitalTemplateConfig(),
			CreateSignupDefaultTemplateConfig(),
		},
	}

	subscribeTemplateChain := TemplateChain{
		[]TemplateDefinition{
			CreateSubscribeWHPlusFMTemplateConfig(),
			CreateSubscribeWHPlusTemplateConfig(),
			CreateSubscribeInternationalCheckinTemplateConfig(),
			CreateSubscribeFMTemplateConfig(),
			CreateSubscribeDefaultTemplateConfig(),
		},
	}

	featureToggle := FeatureToggle{}

	//
	// CreateInvitation
	//

	dto := HydratorDTO{eligibleID: "7343a8ea-9f4b-4ebc-aca2-5e1901869e3b", isMember: CONFIG.isMember}

	hydratorChain.hydrate(&dto)

	var template Template
	var experiment *Experiment
	if !dto.isMember {
		template, experiment, _ = signupTemplateChain.choose(featureToggle, dto)
	} else {
		template, experiment, _ = subscribeTemplateChain.choose(featureToggle, dto)
	}

	// Persistence phase
	// Convert to KNS event
	// Persist and publish events
	fmt.Printf("Template:   %v\n", template.template)
	fmt.Printf("Subject:    %v\n", template.subject)
	if experiment != nil {
		fmt.Printf("Experiment: %v - %v\n", experiment.name, experiment.variant)
	} else {

		fmt.Println("Experiment: N/A")
	}

}
