package main

type ABTest interface {
	accept(featureToggle FeatureToggle, dto HydratorDTO, template TemplateData) (TemplateData, ExperimentData, bool, error)
}

type ABTestChain struct {
	abTests []ABTest
}

func (tc ABTestChain) choose(featureToggle FeatureToggle, dto HydratorDTO, template TemplateData) (TemplateData, *ExperimentData, error) {
	for _, abTest := range tc.abTests {
		updatedTemplate, experiment, accepted, err := abTest.accept(featureToggle, dto, template)
		if err != nil {
			return TemplateData{}, nil, err
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

func (SignupWHPlusSubject) accept(featureToggle FeatureToggle, dto HydratorDTO, template TemplateData) (TemplateData, ExperimentData, bool, error) {
	if template.template != "signup_wh_plus" {
		return template, ExperimentData{}, false, nil
	}

	if featureToggle.IsEnabled("SignupSubject", dto.eligibleID) {
		return TemplateData{
				template: template.template,
				subject:  "signup_wh_plus_subject",
			}, ExperimentData{
				experimentEnabled: true,
				name:              "SignupSubject",
				scenario:          "B",
			}, true, nil
	} else {
		return template, ExperimentData{
			experimentEnabled: true,
			name:              "SignupSubject",
			scenario:          "Control",
		}, true, nil
	}
}

type SubscribeInternationalCheckIn struct{}

func (SubscribeInternationalCheckIn) accept(featureToggle FeatureToggle, dto HydratorDTO, template TemplateData) (TemplateData, ExperimentData, bool, error) {
	if !dto.isMember {
		return template, ExperimentData{}, false, nil
	}

	clientOrder, err := dto.GetClientOrder()
	if err != nil {
		return template, ExperimentData{}, false, err
	}
	if !clientOrder.hasInternationCheckin {
		return template, ExperimentData{}, false, nil
	}

	if featureToggle.IsEnabled("InternationalCheckIn", dto.eligibleID) {
		return TemplateData{
				template: "subscribe_international_checkin",
				subject:  template.subject,
			}, ExperimentData{
				experimentEnabled: true,
				name:              "InternationalCheckIn",
				scenario:          "B",
			}, true, nil
	} else {
		return template, ExperimentData{
			experimentEnabled: true,
			name:              "InternationalCheckIn",
			scenario:          "Control",
		}, true, nil
	}
}
