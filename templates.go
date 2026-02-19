package main

import "errors"

type Template interface {
	accept(featureToggle FeatureToggle, dto HydratorDTO) (TemplateData, ExperimentData, bool)
}

type TemplateChain struct {
	templates []Template
}

func (tc TemplateChain) choose(featureToggle FeatureToggle, dto HydratorDTO) (TemplateData, ExperimentData, error) {
	for _, template := range tc.templates {
		templateData, experimentData, accepted := template.accept(featureToggle, dto)
		if accepted {
			return templateData, experimentData, nil
		}
	}

	return TemplateData{}, ExperimentData{}, errors.New("No default template configured. Should never happen")
}

///
///
///

type WHPlusFMTemplate struct{}

func (WHPlusFMTemplate) accept(featureToggle FeatureToggle, dto HydratorDTO) (TemplateData, ExperimentData, bool) {
	if !(dto.clientOrder.hasWHPlus &&
		dto.bestPlan.discountedPrice == 0 &&
		dto.clientOrder.hasFM) {
		return TemplateData{}, ExperimentData{}, false
	}

	return TemplateData{
			template: "wellhub_plus_fm",
			subject:  "default",
		}, ExperimentData{
			experimentEnabled: false,
		},
		true
}

type WHPlusTemplate struct{}

func (WHPlusTemplate) accept(featureToggle FeatureToggle, dto HydratorDTO) (TemplateData, ExperimentData, bool) {
	if !(dto.clientOrder.hasWHPlus && dto.bestPlan.discountedPrice == 0) {
		return TemplateData{}, ExperimentData{}, false
	}

	return TemplateData{
			template: "wellhub_plus_fm",
			subject:  "default",
		}, ExperimentData{
			experimentEnabled: false,
		},
		true
}

type FMTemplate struct{}

func (FMTemplate) accept(featureToggle FeatureToggle, dto HydratorDTO) (TemplateData, ExperimentData, bool) {
	// Is acceptable?
	if !dto.clientOrder.hasFM {
		return TemplateData{}, ExperimentData{}, false
	}

	var subject string
	var experimentData ExperimentData
	if featureToggle.IsEnabled("fm_subject", dto.eligibleID) {
		subject = "new_fm_subject"
		experimentData = ExperimentData{
			experimentEnabled: true,
			name:              "fm_subject",
			group:             "subject",
			scenario:          "B",
		}
	} else {
		subject = "default"
		experimentData = ExperimentData{
			experimentEnabled: true,
			name:              "fm_subject",
			group:             "subject",
			scenario:          "A",
		}
	}

	return TemplateData{
			template: "fm",
			subject:  subject,
		}, experimentData,
		true
}

type InternationalCheckin struct{}

func (InternationalCheckin) accept(featureToggle FeatureToggle, dto HydratorDTO) (TemplateData, ExperimentData, bool) {
	if !dto.clientOrder.hasInternationCheckin {
		return TemplateData{}, ExperimentData{}, false
	}

	if featureToggle.IsEnabled("InternationalCheckin", dto.eligibleID) {
		return TemplateData{
				template: "new_default",
				subject:  "new_default",
			}, ExperimentData{
				experimentEnabled: true,
				name:              "new_default",
				group:             "template",
				scenario:          "B",
			}, true
	} else {
		return TemplateData{
				template: "default",
				subject:  "default",
			}, ExperimentData{
				experimentEnabled: true,
				name:              "new_default",
				group:             "template",
				scenario:          "A",
			}, true
	}

}

type Default struct{}

func (Default) accept(featureToggle FeatureToggle, dto HydratorDTO) (TemplateData, ExperimentData, bool) {
	return TemplateData{
			template: "default",
			subject:  "default",
		},
		ExperimentData{experimentEnabled: false},
		true
}
