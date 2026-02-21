package main

type ABTest interface {
	accept(featureToggle FeatureToggle, dto HydratorDTO, template TemplateData) (TemplateData, ExperimentData, bool)
}

type ABTestChain struct {
	abTests []ABTest
}

func (tc ABTestChain) choose(featureToggle FeatureToggle, dto HydratorDTO, template TemplateData) (TemplateData, *ExperimentData) {
	for _, abTest := range tc.abTests {
		updatedTemplate, experiment, accepted := abTest.accept(featureToggle, dto, template)
		if accepted {
			return updatedTemplate, &experiment
		}
	}

	return template, nil
}

///
/// Running AB Tests
///

type SignupWHPlusSubject struct{}

func (SignupWHPlusSubject) accept(featureToggle FeatureToggle, dto HydratorDTO, template TemplateData) (TemplateData, ExperimentData, bool) {
	if template.template != "signup_wh_plus" {
		return template, ExperimentData{}, false
	}

	if featureToggle.IsEnabled("SignupSubject", dto.eligibleID) {
		return TemplateData{
				template: template.template,
				subject:  "signup_wh_plus_subject",
			}, ExperimentData{
				experimentEnabled: true,
				name:              "SignupSubject",
				scenario:          "B",
			}, true
	} else {
		return template, ExperimentData{
			experimentEnabled: true,
			name:              "SignupSubject",
			scenario:          "Control",
		}, true
	}
}

type InternationalCheckIn struct{}

func (InternationalCheckIn) accept(featureToggle FeatureToggle, dto HydratorDTO, template TemplateData) (TemplateData, ExperimentData, bool) {
	if !dto.isMember || !dto.clientOrder.hasInternationCheckin {
		return template, ExperimentData{}, false
	}

	if featureToggle.IsEnabled("InternationalCheckIn", dto.eligibleID) {
		return TemplateData{
				template: "subscribe_international_checkin",
				subject:  template.subject,
			}, ExperimentData{
				experimentEnabled: true,
				name:              "InternationalCheckIn",
				scenario:          "B",
			}, true
	} else {
		return template, ExperimentData{
			experimentEnabled: true,
			name:              "InternationalCheckIn",
			scenario:          "Control",
		}, true
	}
}
