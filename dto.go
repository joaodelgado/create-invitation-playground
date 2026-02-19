package main

type InvitationBatch struct {
	ids []string
}

type ClientOrder struct {
	country               string
	locale                string
	hasWHPlus             bool
	hasFM                 bool
	hasDigitalPlan        bool
	hasInternationCheckin bool
}

type Membership struct {
	locale string
}

type Plan struct {
	planName           string
	originalPrice      float32
	discountedPrice    float32
	discountPercentage float32
}

type HydratorDTO struct {
	eligibleID  string
	clientOrder *ClientOrder
	membership  *Membership
	bestPlan    *Plan
}

type ExperimentData struct {
	experimentEnabled bool
	name              string
	group             string
	scenario          string
}

type TemplateData struct {
	template string
	subject  string
	// Other template data...
}
