package australia

func New() *AustraliaEditionData {
	return &AustraliaEditionData{
		cards: buildCards(),
		rules: buildRules(),
	}
}
