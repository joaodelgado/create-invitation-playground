package main

type LazyLoaded[T any] struct {
	loaded bool
	loader func() (T, error)
	val    T
}

func (l *LazyLoaded[T]) Get() (T, error) {
	var empty T
	if !l.loaded {
		val, err := l.loader()
		if err != nil {
			return empty, err
		}
		l.val = val
	}

	return l.val, nil
}

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

type HydratorDependencies struct{}

func (HydratorDependencies) LoadBestPlan() (Plan, error) {
	// Would be an external call
	return Plan{
		planName:           "Silver",
		originalPrice:      10,
		discountedPrice:    8,
		discountPercentage: 20,
	}, nil
}

type Partner struct{}

type HydratorDTO struct {
	eligibleID          string
	isMember            bool
	clientOrder         *ClientOrder
	membership          *Membership
	BestPlan            LazyLoaded[Plan]
	RecommendedPartners LazyLoaded[[]Partner]
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
