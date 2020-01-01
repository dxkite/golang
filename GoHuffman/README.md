# Huffman 压缩算法

- 核心： 使用二叉树编码，出现频度高的字符编码越短，用短编码代替长编码来实现数据无损压缩
- 问题：
  - 构建Huffman树
  - 编码压缩数据
  - 存储频度树
  - 读写编码数据

## 0x00 Huffman编码实现

### 树结构

```go
type HuffmanNode struct {
	Value       byte // 当前节点值
	Left, Right *HuffmanNode // 左右子树
	count       uint // 节点权值
}
```

#### 树生成算法

1. 创建N个基础节点（每个节点权值为字符出现频率，1<= N <= 256）
2. 从下而上选择权值最小的两个节点构成新节点

#### 算法实现

```go
list := createHuffmanList(table)
for len(list) >= 2 {
    sort.Sort(list)
    left, right := list[0], list[1]
    newTree := mergeNodeToNew(left.count+right.count, left, right)
    list = append(list[2:], newTree)
}
```

现有如下数据:

```go
package huffman

import "testing"

func TestEncodeHuffmanTree(t *testing.T) {
	
}
```

计算Huffman编码:

```
num	hex	count	huffman
1	74	7	1111
2	65	7	1110
...
28	09	1	000101
29	7D	1	000100
30	64	1	001010
```

## 0x01 构建二进制Huffman链

关键操作：

使用位或操作，可以实现将特殊位置1，使用取反+位与操作，可以实现将特殊位置0，如 ` ^1000 = 0111 --> 0111 & 1010 = 0010` 

> Go中使用^代替其他语言中的~按位取反

```go
if bit == 1 {
    b.data[offset] = b.data[offset] | (1 << (offsetBit))  // 置 1
} else {
    b.data[offset] = b.data[offset] & (^(1 << (offsetBit))) // 置 0 
}
```

## 0x02 存储频度树

存储Huffman树，采用的是链式存储，由于Huffman编码的长度不会很长，为了节省长度，采用了如下的结构

```
|-----------|--------|---------|---------|-------|
|	uint8   | uint8  | bit * n | bit * k |   ... |
|-----------|--------|---------|---------|-------|
   ^            ^         ^            ^
   |			└- 字典    └- 编码长度   └- 编码 
   └- 树的整块长度, 0~255 分别表示 1-256个     
```

其中， n 表示为 `ceil(log2(max(code)))`, 即能够存储表示最大Huffman编码长度的最小二进制长度。k为单个个编码长度

```go
// 编码频度树
func EncodeHuffmanTree(tree *HuffmanTree) (stream *ByteStream) {
	stream = NewEmptyByteStream()
	storeLen := int(math.Ceil(math.Log2(float64(getCodeMaxLen(tree.table)))))
	stream.WriteByte(uint8(storeLen))        // max = 8
	stream.WriteByte(uint8(len(tree.table))) // max = 256, 0 = 256
	for _, pair := range tree.table {
		stream.WriteByte(pair.Data)
		stream.WriteByteLen(uint8(len(pair.Code)), storeLen)
		stream.WriteBitString(pair.Code)
	}
	return
}
```

 ## 0x03 编码压缩数据

根据频度表，读取字节，获取Huffman编码，并写入字节流：

```go
// 根据字典树压缩
func TreeEncode(tree *HuffmanTree, data []byte) (stream *ByteStream, err error) {
	stream = NewEmptyByteStream()
	for _, item := range data {
		if _, code, ok := tree.FindCode(item); ok {
			stream.WriteBitString(code)
		} else {
			return nil, errByteTable
		}
	}
	return
}
```

## 0x04 解码压缩数据

根据频度表逐字节解码数据

```go
// 根据字典树解压
func TreeDecode(tree *HuffmanTree, stream *ByteStream, byteLen uint64) []byte {
	output := []byte{}
	for i := uint64(0); i < byteLen; i++ {
		output = append(output, SeekNextByte(tree.HuffmanNode, stream))
	}
	return output
}

func SeekNextByte(tree *HuffmanNode, stream *ByteStream) byte {
	if tree.Left == nil && tree.Right == nil {
		return tree.Value
	}
	if stream.ReadBit() == 0 {
		return SeekNextByte(tree.Left, stream)
	} else {
		return SeekNextByte(tree.Right, stream)
	}
}
```

