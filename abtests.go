package main

type ABTest interface {
	accept(featureToggle FeatureToggle, dto HydratorDTO, template TemplateData) (TemplateData, *ExperimentData, error)
}

type ABTestChain struct {
	abTests []ABTest
}

func (tc ABTestChain) choose(featureToggle FeatureToggle, dto HydratorDTO, template TemplateData) (TemplateData, *ExperimentData, error) {
	for _, abTest := range tc.abTests {
		updatedTemplate, experiment, err := abTest.accept(featureToggle, dto, template)
		if err != nil {
			return TemplateData{}, nil, err
		}
		if experiment != nil {
			return updatedTemplate, experiment, nil
		}
	}

	return template, nil, nil
}

///
/// Running AB Tests
///

type SubscribeInternationalCheckIn struct{}

func (SubscribeInternationalCheckIn) accept(featureToggle FeatureToggle, dto HydratorDTO, template TemplateData) (TemplateData, *ExperimentData, error) {
	if !dto.isMember {
		return template, nil, nil
	}

	if !dto.clientOrder.hasInternationCheckin {
		return template, nil, nil
	}

	enabled, variant := featureToggle.IsEnabled("InternationalCheckIn", dto.eligibleID)
	if !enabled {
		return template, nil, nil
	}

	experimentData := &ExperimentData{
		experimentEnabled: true,
		name:              "InternationalCheckIn",
		scenario:          "B",
	}

	if variant == "variant_a" {
		template = TemplateData{
			template: "subscribe_international_checkin",
			subject:  template.subject,
		}
	}

	return template, experimentData, nil
}
