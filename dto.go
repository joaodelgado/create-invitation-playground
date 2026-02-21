package main

type InvitationBatch struct {
	ids []string
}

type ClientOrder struct {
	name                  string
	country               string
	locale                string
	hasWHPlus             bool
	hasFM                 bool
	hasDigitalPlan        bool
	hasInternationCheckin bool
}

type Membership struct {
	locale string
	name   string
}

type Plan struct {
	planName           string
	originalPrice      float32
	discountedPrice    float32
	discountPercentage float32
}

type HydratorDTO struct {
	eligibleID  string
	isMember    bool
	clientOrder *ClientOrder
	membership  *Membership
	bestPlan    *Plan
}

type ExperimentData struct {
	experimentEnabled bool
	name              string
	scenario          string
}

type InvitationContext struct {
	clientName      string
	userName        string
	recommendedPlan Plan
}

type TemplateData struct {
	template string
	subject  string
}
