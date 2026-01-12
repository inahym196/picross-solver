package picrosssolver

type deducer struct {
	rules []Rule
}

func newDeducer() deducer {
	return deducer{
		[]Rule{
			ZeroHintRule{},
			MinimumSpacingRule{},
			OverlapFillRule{},
			OverlapExpansionRule{},
			EdgeExpansionRule{},
			BlockSatisfiedRule{},
			PruneImpossibleSegmentRule{},
			FillRemainingWhiteRule{},
		},
	}
}
