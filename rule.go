package picrosssolver

type Rule interface {
	Deduce(HintedCells) []Cell
}

type ExtractMatchRule struct{}

func (r ExtractMatchRule) Deduce(hc HintedCells) []Cell {
	if len(hc.Hints) == 1 && hc.Hints[0] == len(hc.Cells) {
		return filledCells(len(hc.Cells), CellBlack)
	}
	return nil
}

type ZeroHintRule struct{}

func (r ZeroHintRule) Deduce(hc HintedCells) []Cell {
	if len(hc.Hints) == 1 && hc.Hints[0] == 0 {
		return filledCells(len(hc.Cells), CellWhite)
	}
	return nil
}

// 黒と白の配置が一意に決まる
type MinimumSpacingRule struct{}

func (r MinimumSpacingRule) Deduce(hc HintedCells) []Cell {

	// 判定部分
	var sum int
	for _, h := range hc.Hints {
		sum += h
	}
	if sum+(len(hc.Hints)-1) != len(hc.Cells) {
		return nil
	}

	// 生成部分
	deduced := make([]Cell, len(hc.Cells))
	var last int
	for i, hint := range hc.Hints {
		for range hint {
			deduced[last] = CellBlack
			last++
		}
		if i != len(hc.Hints)-1 {
			deduced[last] = CellWhite
			last++
		}
	}
	return deduced
}

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
