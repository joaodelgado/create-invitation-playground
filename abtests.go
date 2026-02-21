package main

type ABTest interface {
	accept(featureToggle FeatureToggle, dto HydratorDTO, template Template) (Template, Experiment, bool, error)
}

type ABTestChain struct {
	abTests []ABTest
}

func (tc ABTestChain) choose(featureToggle FeatureToggle, dto HydratorDTO, template Template) (Template, *Experiment, error) {
	for _, abTest := range tc.abTests {
		updatedTemplate, experiment, accepted, err := abTest.accept(featureToggle, dto, template)
		if err != nil {
			return template, nil, err
		}
		if accepted {
			return updatedTemplate, &experiment, nil
		}
	}

	return template, nil, nil
}

///
/// Running AB Tests
///

type SignupWHPlusSubject struct{}

func (SignupWHPlusSubject) accept(featureToggle FeatureToggle, dto HydratorDTO, template Template) (Template, Experiment, bool, error) {
	if template.template != "signup_wh_plus" {
		return template, Experiment{}, false, nil
	}

	enabled, variant := featureToggle.IsEnabled("SignupSubject", dto.eligibleID)
	if !enabled {
		return template, Experiment{}, false, nil
	}

	experiment := Experiment{
		name:    "SignupSubject",
		variant: variant,
	}

	if variant == "variant_a" {
		template.subject = "signup_wh_plus_subject_variant_a"
	}

	return template, experiment, true, nil
}

type InternationalCheckIn struct{}

func (InternationalCheckIn) accept(featureToggle FeatureToggle, dto HydratorDTO, template Template) (Template, Experiment, bool, error) {
	if template.template == "subscribe_wh_plus" ||
		template.template == "subscribe_wh_plus_fm" ||
		!dto.isMember ||
		!dto.clientOrder.hasInternationCheckin {

		return template, Experiment{}, false, nil
	}

	enabled, variant := featureToggle.IsEnabled("InternationalCheckIn", dto.eligibleID)
	if !enabled {
		return template, Experiment{}, false, nil
	}

	experiment := Experiment{
		name:    "InternationalCheckIn",
		variant: variant,
	}

	if variant == "variant_a" {
		template = Template{
			template: "subscribe_international_checkin",
			subject:  template.subject,
		}
	}

	return template, experiment, true, nil
}

type SubscribeFMSubject struct{}

func (SubscribeFMSubject) accept(featureToggle FeatureToggle, dto HydratorDTO, template Template) (Template, Experiment, bool, error) {
	if template.template != "subscribe_fm" {
		return template, Experiment{}, false, nil
	}

	enabled, variant := featureToggle.IsEnabled("SubscribeFMSubject", dto.eligibleID)
	if !enabled {
		return template, Experiment{}, false, nil
	}

	experiment := Experiment{
		name:    "SubscribeFMSubject",
		variant: variant,
	}

	if variant == "variant_a" {
		template.subject = "subscribe_fm_subject"
	}

	return template, experiment, true, nil
}
