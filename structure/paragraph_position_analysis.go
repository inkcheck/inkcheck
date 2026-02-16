package structure

import (
	"math"

	"github.com/inkcheck/config"
	"github.com/inkcheck/shared"
)

type ParagraphPositionResult struct {
	OpeningLength    int
	ClosingLength    int
	BodyMeanLength   float64
	OpeningDeviation float64
	ClosingDeviation float64
	Uniform          bool
}

func ParagraphPositionAnalysis(cfg config.Config, text string) ParagraphPositionResult {
	paragraphs := shared.SplitParagraphs(text)
	n := len(paragraphs)

	if n < 3 {
		return ParagraphPositionResult{}
	}

	openLen := shared.CountWords(paragraphs[0])
	closeLen := shared.CountWords(paragraphs[n-1])

	bodyLengths := make([]float64, n-2)
	for i := 1; i < n-1; i++ {
		bodyLengths[i-1] = float64(shared.CountWords(paragraphs[i]))
	}
	bodyMean := shared.Mean(bodyLengths)

	openDev := 0.0
	closeDev := 0.0
	if bodyMean > 0 {
		openDev = math.Abs(float64(openLen)-bodyMean) / bodyMean
		closeDev = math.Abs(float64(closeLen)-bodyMean) / bodyMean
	}

	return ParagraphPositionResult{
		OpeningLength:    openLen,
		ClosingLength:    closeLen,
		BodyMeanLength:   bodyMean,
		OpeningDeviation: openDev,
		ClosingDeviation: closeDev,
		Uniform:          openDev < cfg.UniformityDeviation && closeDev < cfg.UniformityDeviation,
	}
}
