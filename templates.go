package main

import "errors"

//
// Template definition
//

type Precondition func(HydratorDTO) (bool, error)
type BuildTemplate func(HydratorDTO) (Template, error)
type ABTestChecker func(FeatureToggle, HydratorDTO) (enabled bool, variant string)
type ABTestPropertyBuilder func(variant string, dto HydratorDTO, template Template) (Template, Experiment, error)
type ABTestTemplateBuilder func(variant string, dto HydratorDTO) (*Template, Experiment, error)

type TemplateDefinition struct {
	// Determines if the template definition is applicable at all
	precondition Precondition

	// Builds the base template if there is no template level AB test running
	templateBuilder *BuildTemplate

	// Modifies the base template with a property level AB test, if any
	propertyABTestEnabled ABTestChecker
	propertyABTest        *ABTestPropertyBuilder

	// Builds the template when there is an template level AB test running
	templateABTestEnabled ABTestChecker
	buildTemplateABTest   *ABTestTemplateBuilder
}

// Just an idea of how to creating a TemplateDefinition could be done.
// Not to be taken too seriously
type TemplateDefinitionBuilder struct {
	_inner *TemplateDefinition
}

func NewTemplateDefinition(precondition func(HydratorDTO) (bool, error)) TemplateDefinitionBuilder {
	disabled := func(FeatureToggle, HydratorDTO) (bool, string) { return false, "" }
	return TemplateDefinitionBuilder{
		_inner: &TemplateDefinition{
			precondition:          precondition,
			propertyABTestEnabled: disabled,
			templateABTestEnabled: disabled,
		},
	}
}

func (builder TemplateDefinitionBuilder) WithBuildTemplate(buildTemplate BuildTemplate) TemplateDefinitionBuilder {
	builder._inner.templateBuilder = &buildTemplate
	return builder
}

func (builder TemplateDefinitionBuilder) WithPropertyABTest(
	enabled ABTestChecker,
	propertyABTest ABTestPropertyBuilder,
) TemplateDefinitionBuilder {
	builder._inner.propertyABTestEnabled = enabled
	builder._inner.propertyABTest = &propertyABTest
	return builder
}

func (builder TemplateDefinitionBuilder) WithTemplateABTest(
	enabled ABTestChecker,
	buildTemplateABTest ABTestTemplateBuilder,
) TemplateDefinitionBuilder {
	builder._inner.templateABTestEnabled = enabled
	builder._inner.buildTemplateABTest = &buildTemplateABTest
	return builder
}

func (builder TemplateDefinitionBuilder) Build() TemplateDefinition {
	return *builder._inner
}

type TemplateChain struct {
	templates []TemplateDefinition
}

func (tc TemplateChain) choose(featureToggle FeatureToggle, dto HydratorDTO) (Template, *Experiment, error) {
	var currentExperiment *Experiment

	for _, templateDefinition := range tc.templates {
		accepted, err := templateDefinition.precondition(dto)
		if err != nil {
			return Template{}, nil, err
		}
		if !accepted {
			continue
		}

		if templateDefinition.templateBuilder != nil {
			template, err := (*templateDefinition.templateBuilder)(dto)
			if err != nil {
				return Template{}, nil, err
			}

			if currentExperiment == nil {
				enabled, variant := templateDefinition.propertyABTestEnabled(featureToggle, dto)
				if enabled {
					template, experiment, err := (*templateDefinition.propertyABTest)(variant, dto, template)
					if err != nil {
						return Template{}, nil, err
					}
					return template, &experiment, nil
				}
			}
			return template, currentExperiment, nil
		}

		if currentExperiment == nil {
			enabled, variant := templateDefinition.templateABTestEnabled(featureToggle, dto)
			if enabled {
				maybeTemplate, experimentData, err := (*templateDefinition.buildTemplateABTest)(variant, dto)
				if err != nil {
					return Template{}, nil, err
				}

				currentExperiment = &experimentData
				if maybeTemplate != nil {
					return *maybeTemplate, currentExperiment, nil
				}
			}
		}
	}

	return Template{}, nil, errors.New("No default template configured. Should never happen.")
}

///
/// Signup templates
///

func CreateSignupWHPlusTemplateDefinition() TemplateDefinition {
	return NewTemplateDefinition(
		func(dto HydratorDTO) (bool, error) {
			if !dto.clientOrder.hasWHPlus {
				return false, nil
			}

			bestPlan, err := dto.BestPlan.Get()
			if err != nil {
				return false, err
			}

			return bestPlan.discountedPrice == 0, nil
		},
	).WithBuildTemplate(
		func(HydratorDTO) (Template, error) {
			return Template{
				template: "signup_wh_plus",
				subject:  "default",
			}, nil
		},
	).WithPropertyABTest(
		func(featureToggle FeatureToggle, dto HydratorDTO) (bool, string) {
			return featureToggle.IsEnabled("SignupSubject", dto.eligibleID)
		},
		func(variant string, dto HydratorDTO, template Template) (Template, Experiment, error) {
			experiment := Experiment{
				name:    "SignupWHPlusSubject",
				variant: variant,
			}

			if variant == "variant_a" {
				template.subject = "signup_wh_plus_subject_variant_a"
			}

			return template, experiment, nil
		},
	).Build()
}

func CreateSignupDigitalTemplateDefinition() TemplateDefinition {
	return NewTemplateDefinition(
		func(dto HydratorDTO) (bool, error) {
			return dto.clientOrder.hasDigitalPlan, nil
		},
	).WithBuildTemplate(
		func(HydratorDTO) (Template, error) {
			return Template{
				template: "signup_digital_plan",
				subject:  "default",
			}, nil
		},
	).Build()
}

func CreateSignupDefaultTemplateDefinition() TemplateDefinition {
	return NewTemplateDefinition(
		func(dto HydratorDTO) (bool, error) {
			return true, nil
		},
	).WithBuildTemplate(
		func(HydratorDTO) (Template, error) {
			return Template{
				template: "signup_default",
				subject:  "default",
			}, nil
		},
	).Build()
}

///
/// Subscribe templates
///

func CreateSubscribeWHPlusFMTemplateDefinition() TemplateDefinition {
	return NewTemplateDefinition(
		func(dto HydratorDTO) (bool, error) {
			if !dto.clientOrder.hasWHPlus || !dto.clientOrder.hasFM {
				return false, nil
			}

			bestPlan, err := dto.BestPlan.Get()
			if err != nil {
				return false, err
			}

			return bestPlan.discountedPrice == 0, nil
		},
	).WithBuildTemplate(
		func(HydratorDTO) (Template, error) {
			return Template{
				template: "subscribe_wh_plus_fm",
				subject:  "default",
			}, nil
		},
	).Build()
}

func CreateSubscribeWHPlusTemplateDefinition() TemplateDefinition {
	return NewTemplateDefinition(
		func(dto HydratorDTO) (bool, error) {
			if !dto.clientOrder.hasWHPlus {
				return false, nil
			}

			bestPlan, err := dto.BestPlan.Get()
			if err != nil {
				return false, err
			}

			return bestPlan.discountedPrice == 0, nil
		},
	).WithBuildTemplate(
		func(HydratorDTO) (Template, error) {
			return Template{
				template: "subscribe_wh_plus",
				subject:  "default",
			}, nil
		},
	).Build()
}

func CreateSubscribeInternationalCheckinTemplateDefinition() TemplateDefinition {
	return NewTemplateDefinition(
		func(dto HydratorDTO) (bool, error) {
			return dto.clientOrder.hasInternationCheckin, nil
		},
	).WithTemplateABTest(
		func(featureToggle FeatureToggle, dto HydratorDTO) (enabled bool, variant string) {
			return featureToggle.IsEnabled("InternationalCheckIn", dto.eligibleID)
		},
		func(variant string, dto HydratorDTO) (*Template, Experiment, error) {
			experiment := Experiment{
				name:    "InternationalCheckIn",
				variant: variant,
			}

			if variant == "control" {
				return nil, experiment, nil
			} else {
				return &Template{
					template: "subscribe_international_checkin",
					subject:  "default",
				}, experiment, nil
			}
		},
	).Build()
}

func CreateSubscribeFMTemplateDefinition() TemplateDefinition {
	return NewTemplateDefinition(
		func(dto HydratorDTO) (bool, error) {
			return dto.clientOrder.hasFM, nil
		},
	).WithBuildTemplate(
		func(HydratorDTO) (Template, error) {
			return Template{
				template: "subscribe_fm",
				subject:  "default",
			}, nil
		},
	).WithPropertyABTest(
		func(featureToggle FeatureToggle, dto HydratorDTO) (bool, string) {
			return featureToggle.IsEnabled("SubscribeFMSubject", dto.eligibleID)
		},
		func(variant string, dto HydratorDTO, template Template) (Template, Experiment, error) {
			experiment := Experiment{
				name:    "SubscribeFMSubject",
				variant: variant,
			}

			if variant == "variant_a" {
				template.subject = "subscribe_fm_subject_variant_a"
			}

			return template, experiment, nil
		},
	).Build()
}

func CreateSubscribeDefaultTemplateConfig() TemplateDefinition {
	return NewTemplateDefinition(
		func(dto HydratorDTO) (bool, error) {
			return true, nil
		},
	).WithBuildTemplate(
		func(HydratorDTO) (Template, error) {
			return Template{
				template: "subscribe_default",
				subject:  "default",
			}, nil
		},
	).Build()
}
