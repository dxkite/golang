// Huffman Tree
package huffman

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type HuffmanNode struct {
	Value       byte // 当前节点值
	Left, Right *HuffmanNode // 左右子树
	count       uint // 节点权值
}

type Pair struct {
	Data  byte
	Count uint
}

type HuffmanPair struct {
	Pair
	Code string
}

type HuffmanByteTable []Pair
type HuffmanTable []HuffmanPair

type HuffmanTree struct {
	*HuffmanNode
	table HuffmanTable
}

func CreateHuffmanByteTable(data []byte) HuffmanByteTable {
	var table = HuffmanByteTable{}
	var index = make(map[byte]int)
	for _, value := range data {
		if val, ok := index[value]; ok {
			table[val].Count++
		} else {
			table = append(table, Pair{value, 1})
			index[value] = len(table) - 1
		}
	}
	sort.Sort(sort.Reverse(table))
	return table
}

// 排序
func (table HuffmanByteTable) Len() int {
	return len(table)
}

func (table HuffmanByteTable) Less(i, j int) bool {
	return table[i].Count < table[j].Count
}

func (table HuffmanByteTable) Swap(i, j int) {
	table[i], table[j] = table[j], table[i]
}

// 排序
func (table HuffmanTable) Len() int {
	return len(table)
}

func (table HuffmanTable) Less(i, j int) bool {
	return table[i].Count < table[j].Count
}

func (table HuffmanTable) Swap(i, j int) {
	table[i], table[j] = table[j], table[i]
}

type huffmanList []*HuffmanNode

func (table huffmanList) Len() int {
	return len(table)
}

func (table huffmanList) Less(i, j int) bool {
	return table[i].count < table[j].count
}

func (table huffmanList) Swap(i, j int) {
	table[i], table[j] = table[j], table[i]
}

func PrintHuffmanByteTable(table HuffmanByteTable) {
	for _, val := range table {
		fmt.Printf("%q -> %d\n", val.Data, val.Count)
	}
}

func newNode(val byte, count uint) *HuffmanNode {
	return &HuffmanNode{
		Value: val,
		count: count,
		Right: nil,
		Left:  nil,
	}
}

func mergeNodeToNew(count uint, left, right *HuffmanNode) *HuffmanNode {
	return &HuffmanNode{
		Value: 0,
		Left:  left,
		Right: right,
		count: count,
	}
}

func BuildHuffmanTree(table HuffmanByteTable) *HuffmanTree {
	list := createHuffmanList(table)
	for len(list) >= 2 {
		sort.Sort(list)
		left, right := list[0], list[1]
		newTree := mergeNodeToNew(left.count+right.count, left, right)
		list = append(list[2:], newTree)
	}
	treeOrderTable := createHuffmanTable([]HuffmanPair{}, list[0], "")
	sort.Sort(sort.Reverse(treeOrderTable))
	var tree = HuffmanTree{
		HuffmanNode: list[0],
		table:       treeOrderTable,
	}
	return &tree
}

func RebuildHuffmanTree(table HuffmanTable) *HuffmanTree {
	root := &HuffmanNode{
		Value: 0,
		count: 0,
		Right: nil,
		Left:  nil,
	}
	for _, pair := range table {
		root.AppendNode(pair.Code, pair.Data)
	}
	treeOrderTable := createHuffmanTable([]HuffmanPair{}, root, "")
	sort.Sort(sort.Reverse(treeOrderTable))
	var tree = HuffmanTree{
		HuffmanNode: root,
		table:       treeOrderTable,
	}
	return &tree
}

func (n *HuffmanNode) AppendNode(code string, val byte) {
	// left = 0, right = 1
	if len(code) == 1 {
		if code == "0" {
			n.Left = &HuffmanNode{
				Value: val,
				count: 0,
				Right: nil,
				Left:  nil,
			}
		}
		if code == "1" {
			n.Right = &HuffmanNode{
				Value: val,
				count: 0,
				Right: nil,
				Left:  nil,
			}
		}
	} else {
		p := code[0]
		if p == '0' {
			if n.Left == nil {
				n.Left = &HuffmanNode{
					Value: 0,
					count: 0,
					Right: nil,
					Left:  nil,
				}
			}
			n.Left.AppendNode(code[1:], val)
		}
		if p == '1' {
			if n.Right == nil {
				n.Right = &HuffmanNode{
					Value: 0,
					count: 0,
					Right: nil,
					Left:  nil,
				}
			}
			n.Right.AppendNode(code[1:], val)
		}
	}
}

func createHuffmanList(table HuffmanByteTable) huffmanList {
	list := huffmanList{}
	for _, val := range table {
		list = append(list, newNode(val.Data, val.Count))
	}
	return list
}

func (tree HuffmanNode) IsNode() bool {
	return tree.Left == nil && tree.Right == nil
}

func (tree *HuffmanNode) String() string {
	return "hex\tcount\thuffman\n" + DumpHuffmanTable(tree, "")
}

func (tree *HuffmanTree) String() string {
	return "num\thex\tcount\thuffman\n" + DumpHuffmanTree(tree)
}

func (tree *HuffmanTree) FindByte(code string) (int, byte, bool) {
	for index, pair := range tree.table {
		if pair.Code == code {
			return index, pair.Data, true
		}
	}
	return 0, 0, false
}

func (tree *HuffmanTree) FindCode(data byte) (int, string, bool) {
	for index, pair := range tree.table {
		if pair.Data == data {
			return index, pair.Code, true
		}
	}
	return 0, "", false
}



func DumpHuffmanTree(tree *HuffmanTree) string {
	sort.Sort(sort.Reverse(tree.table))
	var dumpStr = ""
	for index, pair := range tree.table {
		dumpStr += fmt.Sprintf("%d\t%02X\t%d\t%s\n", index+1, pair.Data, pair.Count, pair.Code)
	}
	return dumpStr
}

func DisplayHuffmanTree(tree *HuffmanTree) string {
	dump := DumpHuffmanTable(tree.HuffmanNode, "")
	dumpArr := strings.Split(dump, "\n")
	data := map[uint8]string{}
	for _, val := range dumpArr {
		var hex uint8
		var count int
		var huffman string
		if n, err := fmt.Sscanf(val, "%02X\t%d\t%s", &hex, &count, &huffman); n == 3 && err == nil {
			data[hex] = val;
		} else if len(val) > 0 {
			fmt.Println("error read ", val, " as format")
		}
	}
	var dumpStr = ""
	for index, pair := range tree.table {
		dumpStr += strconv.Itoa(index+1) + "\t" + data[pair.Data] + "\n"
	}
	return dumpStr
}

// 创建Huffman表
func createHuffmanTable(table HuffmanTable, tree *HuffmanNode, prefix string) HuffmanTable {
	if tree.IsNode() {
		table = append(table, HuffmanPair{
			Pair: Pair{
				Data:  tree.Value,
				Count: tree.count,
			},
			Code: prefix,
		})
	}
	if tree.Left != nil {
		table = createHuffmanTable(table, tree.Left, prefix+"0")
	}
	if tree.Right != nil {
		table = createHuffmanTable(table, tree.Right, prefix+"1")
	}
	return table
}

func DumpHuffmanTable(tree *HuffmanNode, prefix string) string {
	if tree.IsNode() {
		return fmt.Sprintf("%02X\t%d\t%s\n", tree.Value, tree.count, prefix)
	}
	var s = ""
	if tree.Left != nil {
		s += DumpHuffmanTable(tree.Left, prefix+"0")
	}
	if tree.Right != nil {
		s += DumpHuffmanTable(tree.Right, prefix+"1")
	}
	return s
}

