package strutil

// 敏感词过滤
type WordFilter struct {
	tree *TrieTree
}

func NewWordFilter() *WordFilter {
	return &WordFilter{
		tree: NewTrieTree(),
	}
}

func (wf *WordFilter) InsertWord(word string) {
	wf.tree.InsertWord(word)
}

func (wf *WordFilter) InsertWords(words []string) {
	wf.tree.InsertWords(words)
}

func (wf *WordFilter) Clear() {
	wf.tree.Clear()
}

// 是否包含敏感词
func (wf *WordFilter) Contains(str string) bool {
	return wf.tree.ContainsWord(str)
}

// 找到str中包含的所有敏感词，并替换为replace
func (wf *WordFilter) Filter(str string, replace rune) string {
	return wf.tree.ReplaceWords(str, replace)
}
