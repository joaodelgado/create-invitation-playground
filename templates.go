package main

import "errors"

type Template interface {
	accept(dto HydratorDTO) *TemplateData
}

type TemplateChain struct {
	templates []Template
}

func (tc TemplateChain) choose(dto HydratorDTO) (TemplateData, error) {
	for _, template := range tc.templates {
		templateData := template.accept(dto)
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

func (SignupWHPlusTemplate) accept(dto HydratorDTO) *TemplateData {
	if dto.clientOrder.hasWHPlus && dto.bestPlan.discountedPrice == 0 {
		template := TemplateData{
			template: "signup_wh_plus",
			subject:  "default",
		}
		return &template
	}

	return nil
}

type SignupDigitalTemplate struct{}

func (SignupDigitalTemplate) accept(dto HydratorDTO) *TemplateData {
	if dto.clientOrder.hasDigitalPlan {
		template := TemplateData{
			template: "signup_wh_plus",
			subject:  "default",
		}
		return &template
	}

	return nil
}

type SignupDefaultTemplate struct{}

func (SignupDefaultTemplate) accept(dto HydratorDTO) *TemplateData {
	template := TemplateData{
		template: "default",
		subject:  "default",
	}

	return &template
}

///
/// Subscribe templates
///

type SubscribeWHPlusFMTemplate struct{}

func (SubscribeWHPlusFMTemplate) accept(dto HydratorDTO) *TemplateData {
	if dto.clientOrder.hasWHPlus && dto.bestPlan.discountedPrice == 0 && dto.clientOrder.hasFM {
		template := TemplateData{
			template: "subscribe_wh_plus_fm",
			subject:  "default",
		}
		return &template
	}

	return nil
}

type SubscribeWHPlusTemplate struct{}

func (SubscribeWHPlusTemplate) accept(dto HydratorDTO) *TemplateData {
	if dto.clientOrder.hasWHPlus && dto.bestPlan.discountedPrice == 0 {
		template := TemplateData{
			template: "subscribe_wh_plus",
			subject:  "default",
		}
		return &template
	}

	return nil
}

type SubscribeFMTemplate struct{}

func (SubscribeFMTemplate) accept(dto HydratorDTO) *TemplateData {
	if dto.clientOrder.hasFM {
		template := TemplateData{
			template: "subscribe_fm",
			subject:  "default",
		}
		return &template
	}

	return nil
}

type SubscribeInternationalCheckingTemplate struct{}

func (SubscribeInternationalCheckingTemplate) accept(dto HydratorDTO) *TemplateData {
	if dto.clientOrder.hasInternationCheckin {
		template := TemplateData{
			template: "subscribe_international_checkin",
			subject:  "default",
		}
		return &template
	}

	return nil
}

type SubscribeDefaultTemplate struct{}

func (SubscribeDefaultTemplate) accept(dto HydratorDTO) *TemplateData {
	template := TemplateData{
		template: "subscribe_default",
		subject:  "default",
	}
	return &template
}
