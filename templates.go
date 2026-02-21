package main

import "errors"

type TemplateDefinition interface {
	accept(dto HydratorDTO) (*Template, error)
}

type TemplateChain struct {
	templates []TemplateDefinition
}

func (tc TemplateChain) choose(dto HydratorDTO) (Template, error) {
	for _, template := range tc.templates {
		template, err := template.accept(dto)
		if err != nil {
			return Template{}, err
		}
		if template != nil {
			return *template, nil
		}
	}

	return Template{}, errors.New("No default template configured. Should never happen")
}

///
/// Signup templates
///

type SignupWHPlusTemplate struct{}

func (SignupWHPlusTemplate) accept(dto HydratorDTO) (*Template, error) {
	if !dto.clientOrder.hasWHPlus {
		return nil, nil
	}

	bestPlan, err := dto.BestPlan.Get()
	if err != nil {
		return nil, err
	}

	if bestPlan.discountedPrice == 0 {
		template := Template{
			template: "signup_wh_plus",
			subject:  "default",
		}
		return &template, nil
	}

	return nil, nil
}

type SignupDigitalTemplate struct{}

func (SignupDigitalTemplate) accept(dto HydratorDTO) (*Template, error) {
	if dto.clientOrder.hasDigitalPlan {
		template := Template{
			template: "signup_wh_plus",
			subject:  "default",
		}
		return &template, nil
	}

	return nil, nil
}

type SignupDefaultTemplate struct{}

func (SignupDefaultTemplate) accept(dto HydratorDTO) (*Template, error) {
	template := Template{
		template: "signup_default",
		subject:  "default",
	}

	return &template, nil
}

///
/// Subscribe templates
///

type SubscribeWHPlusFMTemplate struct{}

func (SubscribeWHPlusFMTemplate) accept(dto HydratorDTO) (*Template, error) {
	if !dto.clientOrder.hasWHPlus || !dto.clientOrder.hasFM {
		return nil, nil
	}

	bestPlan, err := dto.BestPlan.Get()
	if err != nil {
		return nil, err
	}

	if bestPlan.discountedPrice == 0 {
		template := Template{
			template: "subscribe_wh_plus_fm",
			subject:  "default",
		}
		return &template, nil
	}

	return nil, nil
}

type SubscribeWHPlusTemplate struct{}

func (SubscribeWHPlusTemplate) accept(dto HydratorDTO) (*Template, error) {
	if !dto.clientOrder.hasWHPlus {
		return nil, nil
	}

	bestPlan, err := dto.BestPlan.Get()
	if err != nil {
		return nil, err
	}

	if bestPlan.discountedPrice == 0 {
		template := Template{
			template: "subscribe_wh_plus",
			subject:  "default",
		}
		return &template, nil
	}

	return nil, nil
}

type SubscribeFMTemplate struct{}

func (SubscribeFMTemplate) accept(dto HydratorDTO) (*Template, error) {
	if dto.clientOrder.hasFM {
		template := Template{
			template: "subscribe_fm",
			subject:  "default",
		}
		return &template, nil
	}

	return nil, nil
}

type SubscribeInternationalCheckingTemplate struct{}

func (SubscribeInternationalCheckingTemplate) accept(dto HydratorDTO) (*Template, error) {
	if dto.clientOrder.hasInternationCheckin {
		template := Template{
			template: "subscribe_international_checkin",
			subject:  "default",
		}
		return &template, nil
	}

	return nil, nil
}

type SubscribeDefaultTemplate struct{}

func (SubscribeDefaultTemplate) accept(dto HydratorDTO) (*Template, error) {
	template := Template{
		template: "subscribe_default",
		subject:  "default",
	}
	return &template, nil
}
