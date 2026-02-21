package main

import "errors"

type Template interface {
	accept(dto HydratorDTO) (*TemplateData, error)
}

type TemplateChain struct {
	templates []Template
}

func (tc TemplateChain) choose(dto HydratorDTO) (TemplateData, error) {
	for _, template := range tc.templates {
		templateData, error := template.accept(dto)
		if error != nil {
			return TemplateData{}, error
		}

		if templateData != nil {
			return *templateData, nil
		}
	}

	return TemplateData{}, errors.New("No default template configured. Should never happen")
}

///
/// Signup templates
///

type SignupWHPlusTemplate struct{}

func (SignupWHPlusTemplate) accept(dto HydratorDTO) (*TemplateData, error) {
	if !dto.clientOrder.hasWHPlus {
		return nil, nil
	}

	bestPlan, err := dto.GetBestPlan()
	if err != nil {
		return nil, err
	}

	if bestPlan.discountedPrice == 0 {
		template := TemplateData{
			template: "signup_wh_plus",
			subject:  "default",
		}
		return &template, nil
	}

	return nil, nil
}

type SignupDigitalTemplate struct{}

func (SignupDigitalTemplate) accept(dto HydratorDTO) (*TemplateData, error) {
	if dto.clientOrder.hasDigitalPlan {
		template := TemplateData{
			template: "signup_wh_plus",
			subject:  "default",
		}
		return &template, nil
	}

	return nil, nil
}

type SignupDefaultTemplate struct{}

func (SignupDefaultTemplate) accept(dto HydratorDTO) (*TemplateData, error) {
	template := TemplateData{
		template: "default",
		subject:  "default",
	}

	return &template, nil
}

///
/// Subscribe templates
///

type SubscribeWHPlusFMTemplate struct{}

func (SubscribeWHPlusFMTemplate) accept(dto HydratorDTO) (*TemplateData, error) {
	if !dto.clientOrder.hasWHPlus && !dto.clientOrder.hasFM {
		return nil, nil
	}

	bestPlan, err := dto.GetBestPlan()
	if err != nil {
		return nil, err
	}

	if bestPlan.discountedPrice == 0 {
		template := TemplateData{
			template: "subscribe_wh_plus_fm",
			subject:  "default",
		}
		return &template, nil
	}

	return nil, nil
}

type SubscribeWHPlusTemplate struct{}

func (SubscribeWHPlusTemplate) accept(dto HydratorDTO) (*TemplateData, error) {
	if !dto.clientOrder.hasWHPlus {
		return nil, nil
	}

	bestPlan, err := dto.GetBestPlan()
	if err != nil {
		return nil, err
	}

	if bestPlan.discountedPrice == 0 {
		template := TemplateData{
			template: "subscribe_wh_plus",
			subject:  "default",
		}
		return &template, nil
	}

	return nil, nil
}

type SubscribeFMTemplate struct{}

func (SubscribeFMTemplate) accept(dto HydratorDTO) (*TemplateData, error) {
	if dto.clientOrder.hasFM {
		template := TemplateData{
			template: "subscribe_fm",
			subject:  "default",
		}
		return &template, nil
	}

	return nil, nil
}

type SubscribeInternationalCheckingTemplate struct{}

func (SubscribeInternationalCheckingTemplate) accept(dto HydratorDTO) (*TemplateData, error) {
	if dto.clientOrder.hasInternationCheckin {
		template := TemplateData{
			template: "subscribe_international_checkin",
			subject:  "default",
		}
		return &template, nil
	}

	return nil, nil
}

type SubscribeDefaultTemplate struct{}

func (SubscribeDefaultTemplate) accept(dto HydratorDTO) (*TemplateData, error) {
	template := TemplateData{
		template: "subscribe_default",
		subject:  "default",
	}
	return &template, nil
}
