package picrosssolver

type Rule interface {
	Apply(line Line)
}

type ExtractMatchRule struct{}

func (r *ExtractMatchRule) Apply(line Line) {
	if len(line.Hints) != 1 {
		return
	}

	hint := line.Hints[0]
	if hint == len(line.Cells) && !line.IsAllCells(CellBlack) {
		updated := filledCells(len(line.Cells), CellBlack)
		line.WriteBack(updated)
	}
}

type ZeroHintRule struct{}

func (r *ZeroHintRule) Apply(line Line) {
	if len(line.Hints) != 1 {
		return
	}

	hint := line.Hints[0]
	if hint == 0 && !line.IsAllCells(CellWhite) {
		updated := filledCells(len(line.Cells), CellWhite)
		line.WriteBack(updated)
	}
}

// 黒と白の配置が一意に決まる
type MinimumSpacingRule struct{}

// ヒントブロックを左詰め／右詰めしたときに必ず重なる部分を黒確定
type OverlapFillRule struct{}

// 端に黒が確定した場合、ヒントサイズ分伸ばせる
type EdgeExpantionRule struct{}

// 既に黒が hint 長に達しているブロックの前後を白確定
type BlockSatisfiedRule struct{}

// 白確定セルで line を分割し、それぞれにヒントを再配分
type SegmentSplitRule struct{}

// ヒントが収まらない区間を白確定
type PruneImpossibleSegmentRule struct{}

// 長さ < 最小 hint の区間はすべて白
type TooSmallSegmentRule struct{}

// すべての hint を満たした後の残りは白
type FillRemainingWhiteRule struct{}

// 仮に黒／白を置き、矛盾が出たら逆を確定
type HypothesisRule struct{}
