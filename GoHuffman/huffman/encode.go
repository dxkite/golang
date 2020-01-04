// Huffman PackStream
// 小文件压缩没问题
// 大文件当心内存爆炸 （所有数据全部载入内存的）
package huffman

import (
	"errors"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"math"
	"os"
)

var (
	ErrCrc32       = errors.New("error crc32")
	ErrByteLength  = errors.New("error byte length")
	ErrByteTable   = errors.New("error byte table")
	ErrMagicNumber = errors.New("error magic number")
)

const PackStreamHeadLength = 8 + 8 + 8 + 8 + 4 + 2
const MagicNumber = uint64(0x31584448)

type PackStream struct {
	Crc32        uint32
	TreeLength   uint8
	Tree         *HuffmanTree
	StreamLength uint64
	Stream       *ByteStream
}

// 压缩文件
func EncodeFile(inputFile, outputFile string) (err error) {
	if data, err := ioutil.ReadFile(inputFile); err != nil {
		return err
	} else {
		encode, err := EncodeByte(data)
		if err != nil {
			return err
		}
		fmt.Printf("[+] compression ratio %.2f%%, size from %d -> %d ", (float64(len(encode))/float64(len(data)))*100, len(data), len(encode))
		return ioutil.WriteFile(outputFile, encode, os.ModePerm)
	}
}

// 解压文件
func DecodeFile(inputFile, outputFile string) (err error) {
	if data, err := ioutil.ReadFile(inputFile); err != nil {
		return err
	} else {
		decode, err := DecodeByte(data)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(outputFile, decode, os.ModePerm)
	}
}

// 根据字典树压缩
func TreeEncode(tree *HuffmanTree, data []byte) (stream *ByteStream, err error) {
	stream = NewEmptyByteStream()
	for _, item := range data {
		if _, code, ok := tree.FindCode(item); ok {
			stream.WriteBitString(code)
		} else {
			return nil, ErrByteTable
		}
	}
	return
}

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

// 编码字节流
// 字典:CRC32:DATA
func EncodeByte(data []byte) (output []byte, err error) {
	table := CreateHuffmanByteTable([]byte(data))
	tree := BuildHuffmanTree(table)
	if byteStream, err := TreeEncode(tree, data); err != nil {
		return nil, err
	} else {
		stream := &PackStream{
			Tree:   tree,
			Stream: byteStream,
		}
		stream.StreamLength = uint64(len(data))
		return EncodePackStream(stream), nil
	}
}

func EncodePackStream(pack *PackStream) []byte {
	stream := NewEmptyByteStream()
	encodeTree := EncodeHuffmanTree(pack.Tree)
	stream.WriteUint64(MagicNumber)                                                                 // 文件幻数
	stream.WriteUint64(uint64(PackStreamHeadLength + len(encodeTree.data) + len(pack.Stream.data))) // 总长度
	stream.WriteUint64(pack.StreamLength)                                                           // 压缩前大小
	stream.WriteUint64(uint64(len(pack.Stream.data)))                                               // 压缩后大小
	stream.WriteUint32(crc32.ChecksumIEEE(append(encodeTree.data, pack.Stream.data...)))            // 总数据校验
	stream.WriteUint16(uint16(len(encodeTree.data)))                                                // 表占位长度
	stream.AppendByte(encodeTree.data)                                                              // 编码表
	stream.AppendByte(pack.Stream.data)                                                             // 编码数据
	return stream.data
}

// 解码字节流
func DecodeByte(input []byte) (data []byte, err error) {
	stream := NewByteStream(input, uint64(len(input)*4))
	magic := stream.ReadUint64()
	if magic != MagicNumber {
		err = ErrMagicNumber
		return
	}
	total := stream.ReadUint64()
	// 校验数据大小
	if uint64(len(input)) != total {
		err = ErrByteLength
		return
	}
	byteRawLength := stream.ReadUint64()
	byteEncodeLength := stream.ReadUint64()
	// 压缩比
	//fmt.Printf("compress %.2f%%\n", float64(byteRawLength)/float64(byteEncodeLength))
	crc32Sum := stream.ReadUint32()
	tableLen := stream.ReadUint16()
	// 校验CRC32
	if crc32Sum != crc32.ChecksumIEEE(input[PackStreamHeadLength:]) {
		err = ErrCrc32
		return
	}
	// 生成字典树
	tree := DecodeHuffmanTree(NewByteStream(input[PackStreamHeadLength:PackStreamHeadLength+tableLen], uint64(PackStreamHeadLength+tableLen*4)))
	// 根据树解码数据
	data = TreeDecode(tree, NewByteStream(input[PackStreamHeadLength+tableLen:], byteEncodeLength*4), byteRawLength)
	return
}

// 编码频度树
func EncodeHuffmanTree(tree *HuffmanTree) (stream *ByteStream) {
	stream = NewEmptyByteStream()
	storeLen := int(math.Ceil(math.Log2(float64(getCodeMaxLen(tree.table)))))
	stream.WriteByte(uint8(storeLen))            // max = 8
	stream.WriteByte(uint8(len(tree.table) - 1)) // max = 256
	for _, pair := range tree.table {
		stream.WriteByte(pair.Data)
		stream.WriteByteLen(uint8(len(pair.Code)), storeLen)
		stream.WriteBitString(pair.Code)
	}
	return
}

// 解码频度树
func DecodeHuffmanTree(stream *ByteStream) (tree *HuffmanTree) {
	readLen := stream.ReadByte()
	var treeLen = int(stream.ReadByte())
	var huffmanTable = HuffmanTable{}
	for i := 0; i <= treeLen; i++ {
		data := stream.ReadByte()
		size := stream.ReadByteLen(int(readLen))
		bits := stream.ReadBitString(int(size))
		huffmanTable = append(huffmanTable, HuffmanPair{
			Pair: Pair{
				Data:  data,
				Count: uint(0xff - i), // 伪顺序
			},
			Code: bits,
		})
	}
	return RebuildHuffmanTree(huffmanTable)
}

func getCodeMaxLen(table HuffmanTable) int {
	maxLength := 0
	for _, pair := range table {
		if len(pair.Code) > maxLength {
			maxLength = len(pair.Code)
		}
	}
	return maxLength
}
