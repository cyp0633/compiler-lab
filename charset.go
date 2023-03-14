package main

// 字符集
type Charset struct {
	// 字符集 ID
	IndexID int
	// 段 ID。对每个字符集，段 ID 从 0 递增
	SegmentID int
	// 起始字符
	FromChar rune
	// 终止字符
	ToChar rune
}

// 字符集表
var CharsetTable []*Charset

// 字符的范围运算
//
// 生成一个字符集，包含从 fromChar 到 toChar 的所有字符
func rangeChars(fromChar rune, toChar rune) (indexID int) {
	var c = Charset{FromChar: fromChar, ToChar: toChar}
	c.IndexID = maxIndexID() + 1
	c.SegmentID = 0
	CharsetTable = append(CharsetTable, &c)
	return c.IndexID
}

// 字符集的并运算
//
// 生成一个字符集，包含 c1 和 c2
func unionChars(c1 rune, c2 rune) (indexID int) {
	var cs1 = Charset{FromChar: c1, ToChar: c1}
	cs1.IndexID = maxIndexID() + 1
	cs1.SegmentID = 0
	var cs2 = Charset{FromChar: c2, ToChar: c2}
	cs2.IndexID = maxIndexID() + 1 // 属于同一个 charset 的不同段，index 一样
	cs2.SegmentID = 1
	CharsetTable = append(CharsetTable, &cs1, &cs2)
	return cs1.IndexID
}

// 字符集与字符之间的并运算
//
// 将原有的字符集与新字符合并
func unionCharsetAndChar(indexID int, c rune) (newIndexID int) {
	// 原来的字符集可能不只有一段
	var oldCharset []*Charset
	for _, csTemp := range CharsetTable {
		if csTemp.IndexID == indexID {
			oldCharset = append(oldCharset, csTemp)
		}
	}
	maxID := maxIndexID()
	// 将老的字符集拷贝一份（不懂为什么非要创建新的）
	newCharset := copyCharset(oldCharset, maxID+1)
	if len(oldCharset) == 0 {
		return -1
	}
	// 遍历老的字符集各段，看看能不能合并进去
	for _, csTemp := range oldCharset {
		// 在中间，无需其他操作即可合并
		if c >= csTemp.FromChar && c <= csTemp.ToChar {
			CharsetTable = append(CharsetTable, newCharset...)
			return maxID + 1
		}
	}
	// 看看是不是在边缘，合进来
	// 不能与上个放在同一个循环里，防止两个重合
	for index, csTemp := range oldCharset {
		if c == csTemp.FromChar-1 || c == csTemp.ToChar+1 {
			// 将新字符集对应段的边缘合并
			if c == csTemp.FromChar-1 {
				newCharset[index].FromChar = c
			} else {
				newCharset[index].ToChar = c
			}
			CharsetTable = append(CharsetTable, newCharset...)
			return maxID + 1
		}
	}
	// 完全没法合并了
	// 按顺序插入新的字符集
	for index, csTemp := range newCharset {
		if c < csTemp.FromChar {
			var csNew = Charset{IndexID: maxID + 1, SegmentID: index, FromChar: c, ToChar: c}
			newCharset = append(newCharset[:index], append([]*Charset{&csNew}, newCharset[index:]...)...)
			// 重新设置段 ID
			for i := index + 1; i < len(newCharset); i++ {
				newCharset[i].SegmentID = i
			}
			CharsetTable = append(CharsetTable, newCharset...)
			return maxID + 1
		}
	}
	// 也许是最后一个
	var csNew = Charset{IndexID: maxID + 1, SegmentID: len(newCharset), FromChar: c, ToChar: c}
	newCharset = append(newCharset, &csNew)
	CharsetTable = append(CharsetTable, newCharset...)
	return maxID + 1
}

// 字符集与字符集之间的并运算
//
// 将两个字符集合并
func unionTwoCharsets(charsetID1, charsetID2 int) (newIndexID int) {
	// 找两个旧字符集
	var charset1, charset2, newCharset []*Charset
	newIndexID = maxIndexID() + 1
	for _, csTemp := range CharsetTable {
		if csTemp.IndexID == charsetID1 {
			charset1 = append(charset1, csTemp)
		}
		if csTemp.IndexID == charsetID2 {
			charset2 = append(charset2, csTemp)
		}
	}
	if len(charset1) == 0 || len(charset2) == 0 {
		return -1
	}
	// 将两个字符集依次合并到新字符集
	segmentID := 0
	for len(charset1) > 0 || len(charset2) > 0 {
		var csNew Charset
		switch {
		// charset1 没了
		case len(charset1) == 0:
			csNew = Charset{IndexID: newIndexID, SegmentID: segmentID, FromChar: charset2[0].FromChar, ToChar: charset2[0].ToChar}
			// 从 charset2 中移除
			charset2 = charset2[1:]
		// charset2 没了
		case len(charset2) == 0:
			csNew = Charset{IndexID: newIndexID, SegmentID: segmentID, FromChar: charset1[0].FromChar, ToChar: charset1[0].ToChar}
			// 从 charset1 中移除
			charset1 = charset1[1:]
		// 不重合，charset1[0] 在前，而且接不上
		case charset1[0].ToChar < charset2[0].FromChar-1:
			csNew = Charset{IndexID: newIndexID, SegmentID: segmentID, FromChar: charset1[0].FromChar, ToChar: charset1[0].ToChar}
			// 从 charset1 中移除
			charset1 = charset1[1:]
		// 不重合，charset2[0] 在前，而且接不上
		case charset2[0].ToChar < charset1[0].FromChar-1:
			csNew = Charset{IndexID: newIndexID, SegmentID: segmentID, FromChar: charset2[0].FromChar, ToChar: charset2[0].ToChar}
			// 从 charset2 中移除
			charset2 = charset2[1:]
		// 重合
		default:
			// 取两端
			var fromChar, toChar rune
			if charset1[0].FromChar < charset2[0].FromChar {
				fromChar = charset1[0].FromChar
			} else {
				fromChar = charset2[0].FromChar
			}
			if charset1[0].ToChar > charset2[0].ToChar {
				toChar = charset1[0].ToChar
			} else {
				toChar = charset2[0].ToChar
			}
			csNew = Charset{IndexID: newIndexID, SegmentID: segmentID, FromChar: fromChar, ToChar: toChar}
			// 从 charset1 和 charset2 中移除
			charset1 = charset1[1:]
			charset2 = charset2[1:]
		}
		// 看看能否与上一 segment 合并
		if len(newCharset) > 0 && newCharset[len(newCharset)-1].ToChar >= csNew.FromChar-1 {
			newCharset[len(newCharset)-1].ToChar = csNew.ToChar // 这里就不需要多加一个段了
		} else {
			// 合并不了，再加一个段
			newCharset = append(newCharset, &csNew)
			segmentID++
		}
	}
	CharsetTable = append(CharsetTable, newCharset...)
	return
}

// 字符集与字符之间的差运算
//
// 将一个字符从一个字符集中移除
func difference(charsetID int, c rune) (newIndexID int) {
	newIndexID = maxIndexID() + 1
	var newCharset []*Charset
	for _, csTemp := range CharsetTable {
		if csTemp.IndexID == charsetID {
			// 先在这里把 c 移除，后期重新分配段 ID
			switch {
			// 如果这一段就一个 c，就不加进去了
			case csTemp.FromChar == c && csTemp.ToChar == c:
				continue
			// 起始是 c
			case csTemp.FromChar == c:
				var csNew = Charset{IndexID: newIndexID, SegmentID: csTemp.SegmentID, FromChar: c + 1, ToChar: csTemp.ToChar}
				newCharset = append(newCharset, &csNew)
			// 结束是 c
			case csTemp.ToChar == c:
				var csNew = Charset{IndexID: newIndexID, SegmentID: csTemp.SegmentID, FromChar: csTemp.FromChar, ToChar: c - 1}
				newCharset = append(newCharset, &csNew)
			// 中间是 c
			case csTemp.FromChar < c && csTemp.ToChar > c:
				var csNew1 = Charset{IndexID: newIndexID, SegmentID: csTemp.SegmentID, FromChar: csTemp.FromChar, ToChar: c - 1}
				var csNew2 = Charset{IndexID: newIndexID, SegmentID: csTemp.SegmentID, FromChar: c + 1, ToChar: csTemp.ToChar}
				newCharset = append(newCharset, &csNew1, &csNew2)
			// 其他情况，就是不包含 c
			default:
				newCharset = append(newCharset, csTemp)
			}
		}
	}
	// 重新拷贝
	newCharset = copyCharset(newCharset, newIndexID)
	// 重新分配段 ID
	// Go range 的第二个值是拷贝，之前可把我坑坏了
	for i := range newCharset {
		newCharset[i].SegmentID = i
	}
	CharsetTable = append(CharsetTable, newCharset...)
	return
}

// 将一个字符集复制一份
func copyCharset(oldCharset []*Charset, newIndex int) (newCharset []*Charset) {
	for _, csTemp := range oldCharset {
		var csNew = Charset{IndexID: newIndex, SegmentID: csTemp.SegmentID, FromChar: csTemp.FromChar, ToChar: csTemp.ToChar}
		newCharset = append(newCharset, &csNew)
	}
	return
}

// 获得最大的字符集 ID
func maxIndexID() (maxID int) {
	if len(CharsetTable) == 0 {
		return -1
	}
	maxID = CharsetTable[len(CharsetTable)-1].IndexID
	return
}
