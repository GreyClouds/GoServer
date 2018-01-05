package wordstock

var stock *DFA

func init() {
	stock = new(DFA)
}

func Configurate(filename string) {
	stock.BuildTree(filename)
}

// 敏感词校验
func ValidWord(plain string) bool {
	return !stock.IsContain(plain)
}
