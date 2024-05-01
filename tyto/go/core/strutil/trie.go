package strutil

// trie树
type TrieTree struct {
	root *trieNode
}

func NewTrieTree() *TrieTree {
	return &TrieTree{
		root: newTrieNode(),
	}
}

type trieNode struct {
	children map[rune]*trieNode
	isEnd    bool
}

func newTrieNode() *trieNode {
	return &trieNode{
		children: make(map[rune]*trieNode),
		isEnd:    false,
	}
}

func (t *TrieTree) InsertWord(word string) {
	if len(word) == 0 {
		return
	}

	node := t.root
	for _, r := range word {
		if _, ok := node.children[r]; !ok {
			node.children[r] = newTrieNode()
		}

		node = node.children[r]
	}
	node.isEnd = true
}

func (t *TrieTree) InsertWords(words []string) {
	for _, word := range words {
		t.InsertWord(word)
	}
}

func (t *TrieTree) Clear() {
	t.root = newTrieNode()
}

// 判断str中是否包含树中的word
func (t *TrieTree) ContainsWord(str string) bool {
	if len(str) == 0 {
		return false
	}

	var (
		currNode *trieNode
		found    bool
		offset   int
	)

	node := t.root
	runeList := []rune(str)

	for i := 0; i < len(runeList); i++ {
		r := runeList[i]
		currNode, found = node.children[r]
		if !found {
			// 未找到，回退
			i = i - offset
			offset = 0
			node = t.root
			continue
		}

		if currNode.isEnd {
			// 找到一个word
			return true
		}

		offset++
		node = currNode
	}

	return false
}

// 找到str中包含的所有word，并替换为replace
func (t *TrieTree) ReplaceWords(str string, replace rune) string {
	if len(str) == 0 {
		return ""
	}

	var (
		currNode *trieNode
		found    bool
		offset   int
	)

	node := t.root
	runeList := []rune(str)

	for i := 0; i < len(runeList); i++ {
		r := runeList[i]
		currNode, found = node.children[r]
		if !found {
			// 未找到，回退
			i = i - offset
			offset = 0
			node = t.root
			continue
		}

		if currNode.isEnd {
			// 找到一个word
			for j := i - offset; j <= i; j++ {
				runeList[j] = replace
			}

			offset = 0
			node = t.root
			continue
		}

		offset++
		node = currNode
	}

	return string(runeList)
}

// 找到str中包含的所有word
func (t *TrieTree) FindAll(str string) []string {
	if len(str) == 0 {
		return nil
	}

	var (
		currNode *trieNode
		found    bool
		offset   int
		words    []string
	)

	node := t.root
	runeList := []rune(str)

	for i := 0; i < len(runeList); i++ {
		r := runeList[i]
		currNode, found = node.children[r]
		if !found {
			// 未找到，回退
			i = i - offset
			offset = 0
			node = t.root
			continue
		}

		if currNode.isEnd {
			// 找到一个word
			words = append(words, string(runeList[i-offset:i+1]))
			offset = 0
			node = t.root
			continue
		}

		offset++
		node = currNode
	}

	return words
}

// 找到str中包含的第一个word
func (t *TrieTree) FindFirst(str string) string {
	if len(str) == 0 {
		return ""
	}

	var (
		currNode *trieNode
		found    bool
		offset   int
	)

	node := t.root
	runeList := []rune(str)

	for i := 0; i < len(runeList); i++ {
		r := runeList[i]
		currNode, found = node.children[r]
		if !found {
			// 未找到，回退
			i = i - offset
			offset = 0
			node = t.root
			continue
		}

		if currNode.isEnd {
			// 找到一个word
			return string(runeList[i-offset : i+1])
		}

		offset++
		node = currNode
	}

	return ""
}

// 找到str中包含的N个word
func (t *TrieTree) FindN(str string, n int32) []string {
	if len(str) == 0 {
		return nil
	}

	var (
		currNode *trieNode
		found    bool
		offset   int
		words    []string
	)

	node := t.root
	runeList := []rune(str)

	for i := 0; i < len(runeList); i++ {
		r := runeList[i]
		currNode, found = node.children[r]
		if !found {
			// 未找到，回退
			i = i - offset
			offset = 0
			node = t.root
			continue
		}

		if currNode.isEnd {
			// 找到一个word
			words = append(words, string(runeList[i-offset:i+1]))
			offset = 0
			node = t.root
			if int32(len(words)) >= n {
				return words
			}
			continue
		}

		offset++
		node = currNode
	}

	return words
}
