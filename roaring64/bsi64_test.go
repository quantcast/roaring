package roaring64

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
    "math"
    "fmt"
)

func TestSetAndGet(t *testing.T) {

	bsi := NewBSI(999, 0)
	require.NotNil(t, bsi.bA)
	assert.Equal(t, 10, len(bsi.bA))

	bsi.SetValue(1, 8)
	gv, ok := bsi.GetValue(1)
	assert.True(t, ok)
	assert.Equal(t, int64(8), gv)
}

func setup() *BSI {

	bsi := NewBSI(100, 0)
	// Setup values
	for i := 0; i < int(bsi.MaxValue); i++ {
		bsi.SetValue(uint64(i), int64(i))
	}
	return bsi
}

func TestEQ(t *testing.T) {

	bsi := setup()
	eq := bsi.CompareValue(0, EQ, 50, 0, nil)
	assert.Equal(t, uint64(1), eq.GetCardinality())

	assert.True(t, eq.ContainsInt(50))
}

func TestLT(t *testing.T) {

	bsi := setup()
	lt := bsi.CompareValue(0, LT, 50, 0, nil)
	assert.Equal(t, uint64(50), lt.GetCardinality())

	i := lt.Iterator()
	for i.HasNext() {
		v := i.Next()
		assert.Less(t, uint64(v), uint64(50))
	}
}

func TestGT(t *testing.T) {

	bsi := setup()
	gt := bsi.CompareValue(0, GT, 50, 0, nil)
	assert.Equal(t, uint64(49), gt.GetCardinality())

	i := gt.Iterator()
	for i.HasNext() {
		v := i.Next()
		assert.Greater(t, uint64(v), uint64(50))
	}
}

func TestGE(t *testing.T) {

	bsi := setup()
	ge := bsi.CompareValue(0, GE, 50, 0, nil)
	assert.Equal(t, uint64(50), ge.GetCardinality())

	i := ge.Iterator()
	for i.HasNext() {
		v := i.Next()
		assert.GreaterOrEqual(t, uint64(v), uint64(50))
	}
}

func TestLE(t *testing.T) {

	bsi := setup()
	le := bsi.CompareValue(0, LE, 50, 0, nil)
	assert.Equal(t, uint64(51), le.GetCardinality())

	i := le.Iterator()
	for i.HasNext() {
		v := i.Next()
		assert.LessOrEqual(t, uint64(v), uint64(50))
	}
}

func TestRange(t *testing.T) {

	bsi := setup()
	set := bsi.CompareValue(0, RANGE, 45, 55, nil)
	assert.Equal(t, uint64(11), set.GetCardinality())

	i := set.Iterator()
	for i.HasNext() {
		v := i.Next()
		assert.GreaterOrEqual(t, uint64(v), uint64(45))
		assert.LessOrEqual(t, uint64(v), uint64(55))
	}
}

func TestExists(t *testing.T) {

	bsi := NewBSI(10, 0)
	// Setup values
	for i := 1; i < int(bsi.MaxValue); i++ {
		bsi.SetValue(uint64(i), int64(i))
	}

	assert.Equal(t, uint64(9), bsi.GetCardinality())
	assert.False(t, bsi.ValueExists(uint64(0)))
	bsi.SetValue(uint64(0), int64(0))
	assert.Equal(t, uint64(10), bsi.GetCardinality())
	assert.True(t, bsi.ValueExists(uint64(0)))
}

func TestSum(t *testing.T) {

	bsi := setup()
	set := bsi.CompareValue(0, RANGE, 45, 55, nil)

	sum, count := bsi.Sum(set)
	assert.Equal(t, uint64(11), count)
	assert.Equal(t, int64(550), sum)
}

func TestTranspose(t *testing.T) {

	bsi := NewBSI(100, 0)
	// Setup values
	for i := 0; i < int(bsi.MaxValue); i++ {
		bsi.SetValue(uint64(i+100), int64(i))
	}

	set := bsi.Transpose()
	assert.Equal(t, uint64(100), set.GetCardinality())

	i := set.Iterator()
	j := 0
	for i.HasNext() {
		v := i.Next()
		assert.Equal(t, uint64(v), uint64(j))
		j++
	}
}

func TestAutoSize(t *testing.T) {

	bsi := NewDefaultBSI()
	for i := 0; i < 100; i++ {
		bsi.SetValue(uint64(i), int64(i))
	}

	require.NotNil(t, bsi.bA)
	assert.Equal(t, 7, bsi.BitCount())

	for i := 0; i < 100; i++ {
		gv, ok := bsi.GetValue(uint64(i))
		assert.True(t, ok)
		assert.Equal(t, int64(i), gv)
	}
}

func TestParOr(t *testing.T) {

	bsi1 := NewDefaultBSI()
	for i := 0; i < 100; i++ {
		bsi1.SetValue(uint64(i), int64(i))
	}
	bsi2 := NewDefaultBSI()
	for i := 0; i < 100; i++ {
		bsi2.SetValue(uint64(i+100), int64(i+100))
	}
	bsi1.ParOr(0, bsi2)
	for i := 0; i < 200; i++ {
		gv, ok := bsi1.GetValue(uint64(i))
		assert.True(t, ok)
		assert.Equal(t, int64(i), gv)
	}
	assert.Equal(t, uint64(200), bsi1.eBM.GetCardinality())
}

func TestNewBSIRetainSet(t *testing.T) {

	bsi := setup()
	foundSet := BitmapOf(50)
	newBSI := bsi.NewBSIRetainSet(foundSet)
	assert.Equal(t, uint64(1), newBSI.GetCardinality())
	val, ok := newBSI.GetValue(50)
	assert.True(t, ok)
	assert.Equal(t, val, int64(50))
}

func TestLargeFile(t *testing.T) {

	datEBM, err := ioutil.ReadFile("./testdata/age/EBM")
	require.Nil(t, err)
	dat1, err := ioutil.ReadFile("./testdata/age/1")
	require.Nil(t, err)
	dat2, err := ioutil.ReadFile("./testdata/age/2")
	require.Nil(t, err)
	dat3, err := ioutil.ReadFile("./testdata/age/3")
	require.Nil(t, err)
	dat4, err := ioutil.ReadFile("./testdata/age/4")
	require.Nil(t, err)
	dat5, err := ioutil.ReadFile("./testdata/age/5")
	require.Nil(t, err)
	dat6, err := ioutil.ReadFile("./testdata/age/6")
	require.Nil(t, err)
	dat7, err := ioutil.ReadFile("./testdata/age/7")
	require.Nil(t, err)
	dat8, err := ioutil.ReadFile("./testdata/age/8")
	require.Nil(t, err)

	b := [][]byte{datEBM, dat1, dat2, dat3, dat4, dat5, dat6, dat7, dat8}

	bsi := NewDefaultBSI()
	//bsi.RunOptimize()
	err = bsi.UnmarshalBinary(b)
	for i := 0; i < bsi.BitCount(); i++ {
		//assert.True(t, bsi.bA[i].HasRunCompression())
		//bsi.bA[i].RunOptimize()
	}
	//assert.True(t, bsi.eBM.HasRunCompression())
	require.Nil(t, err)

	resultA := bsi.CompareValue(0, EQ, 55, 0, nil)
	assert.Equal(t, uint64(574600), resultA.GetCardinality())

	resultB := bsi.BatchEqual(0, []int64{55, 57})
	assert.Equal(t, uint64(574600+515233), resultB.GetCardinality())

	bsi.ClearValues(resultA)
	resultC := bsi.BatchEqual(0, []int64{55, 57})
	assert.Equal(t, uint64(515233), resultC.GetCardinality())

}

func TestClone(t *testing.T) {
	bsi := setup()
	clone := bsi.Clone()
	for i := 0; i < int(bsi.MaxValue); i++ {
		a, _ := bsi.GetValue(uint64(i))
		b, _ := clone.GetValue(uint64(i))
		assert.Equal(t, a, b)
	}
}

func TestAdd(t *testing.T) {
	bsi := NewDefaultBSI()
	// Setup values
	for i := 1; i <= 10; i++ {
		bsi.SetValue(uint64(i), int64(i))
	}
	clone := bsi.Clone()
	bsi.Add(clone)
	assert.Equal(t, uint64(10), bsi.GetCardinality())
	for i := 1; i <= 10; i++ {
		a, _ := bsi.GetValue(uint64(i))
		b, _ := clone.GetValue(uint64(i))
		assert.Equal(t, b*2, a)
	}

}

func TestIncrement(t *testing.T) {
	bsi := setup()
	bsi.IncrementAll()
	for i := 0; i < int(bsi.MaxValue); i++ {
		a, _ := bsi.GetValue(uint64(i))
		assert.Equal(t, int64(i+1), a)
	}
	bsi.Increment(BitmapOf(0))
	x, _ := bsi.GetValue(uint64(0))
	assert.Equal(t, int64(2), x)
	for i := 1; i < int(bsi.MaxValue); i++ {
		a, _ := bsi.GetValue(uint64(i))
		assert.Equal(t, int64(i+1), a)
	}
}

func TestTransposeWithCounts(t *testing.T) {
	bsi := setup()
	bsi.SetValue(101, 50)
	transposed := bsi.TransposeWithCounts2(0, bsi.GetExistenceBitmap())
	a, ok := transposed.GetValue(uint64(50))
	assert.True(t, ok)
	assert.Equal(t, int64(2), a)
}


func testTranspose2(t *testing.T) {
logo := []uint64 {
0b0000000000000000000000000000000000000000000100000000000000000000,
0b0000000000000000000000000000000000000000011100000000000000000000,
0b0000000000000000000000000000000000000000111110000000000000000000,
0b0000000000000000000000000000000000000001111111000000000000000000,
0b0000000000000000000000000000000000000000111111100000000000000000,
0b0000000000000000000000000000000000000000111111100000000000000000,
0b0000000000000000000000000000000000000000011111110000000000000000,
0b0000000000000000000000000000000000000000001111111000000000000000,
0b0000000000000000000000000000000000000000001111111100000000000000,
0b0000000000000000000000000000000010000000000111111100000000000000,
0b0000000000000000000000000000000011100000000011111110000000000000,
0b0000000000000000000000000000000111110000000001111111000000000000,
0b0000000000000000000000000000001111111000000001111111100000000000,
0b0000000000000000000000000000011111111100000000111111100000000000,
0b0000000000000000000000000000001111111110000000011111110000000000,
0b0000000000000000000000000000000011111111100000001111111000000000,
0b0000000000000000000000000000000001111111110000001111111100000000,
0b0000000000000000000000000000000000111111111000000111111100000000,
0b0000000000000000000000000000000000011111111100000011111110000000,
0b0000000000000000000000000000000000001111111110000001111111000000,
0b0000000000000000000000000000000000000011111111100001111111100000,
0b0000000000000000000000001100000000000001111111110000111111100000,
0b0000000000000000000000001111000000000000111111111000011111110000,
0b0000000000000000000000011111110000000000011111111100001111100000,
0b0000000000000000000000011111111100000000001111111110001111000000,
0b0000000000000000000000111111111111000000000011111111100110000000,
0b0000000000000000000000011111111111110000000001111111110000000000,
0b0000000000000000000000000111111111111100000000111111111000000000,
0b0000000000000000000000000001111111111111100000011111110000000000,
0b0000000000000000000000000000011111111111111000001111100000000000,
0b0000000000000000000000000000000111111111111110000011000000000000,
0b0000000000000000000000000000000001111111111111100000000000000000,
0b0000000000000000000000000000000000001111111111111000000000000000,
0b0000000000000000000000000000000000000011111111111100000000000000,
0b0000000000000000000111000000000000000000111111111100000000000000,
0b0000000000000000000111111110000000000000001111111000000000000000,
0b0000000000000000000111111111111100000000000011111000000000000000,
0b0000000000000000000111111111111111110000000000110000000000000000,
0b0000000000000000001111111111111111111111100000000000000000000000,
0b0000000000000000001111111111111111111111111111000000000000000000,
0b0000000000000000000000011111111111111111111111100000000000000000,
0b0000001111110000000000000001111111111111111111100000111111000000,
0b0000001111110000000000000000000011111111111111100000111111000000,
0b0000001111110000000000000000000000000111111111100000111111000000,
0b0000001111110000000000000000000000000000001111000000111111000000,
0b0000001111110000000000000000000000000000000000000000111111000000,
0b0000001111110000000000000000000000000000000000000000111111000000,
0b0000001111110000001111111111111111111111111111000000111111000000,
0b0000001111110000001111111111111111111111111111000000111111000000,
0b0000001111110000001111111111111111111111111111000000111111000000,
0b0000001111110000001111111111111111111111111111000000111111000000,
0b0000001111110000001111111111111111111111111111000000111111000000,
0b0000001111110000001111111111111111111111111111000000111111000000,
0b0000001111110000000000000000000000000000000000000000111111000000,
0b0000001111110000000000000000000000000000000000000000111111000000,
0b0000001111110000000000000000000000000000000000000000111111000000,
0b0000001111110000000000000000000000000000000000000000111111000000,
0b0000001111110000000000000000000000000000000000000000111111000000,
0b0000001111111111111111111111111111111111111111111111111111000000,
0b0000001111111111111111111111111111111111111111111111111111000000,
0b0000001111111111111111111111111111111111111111111111111111000000,
0b0000001111111111111111111111111111111111111111111111111111000000,
0b0000001111111111111111111111111111111111111111111111111111000000,
0b0000001111111111111111111111111111111111111111111111111111000000,
};

    bsi := NewBSI(math.MaxInt64, 0)
   
    for colID := 0; colID < 64; colID++ {
        var value int64 = 0
        for i := 0; i < bsi.BitCount(); i++ {
            if logo[colID] & (1 << uint64(i)) != 0 {
                value |= (1 << uint64(i))
            }
        }
        bsi.SetValue(uint64(colID), value)
    }

    //printBits(logo)
    //trans := bsi.Transpose()
    //printBSIBits(trans)
}

func printBSIBits(bsi *BSI) {
    for i := 0; i < 64; i++ {
        for j := bsi.BitCount(); j >= 0; j-- {
            v, _ := bsi.GetValue(uint64(i))
            dig := "0"
            if (v >> j) & 1 == 1 {
                dig = "1"
            }
            fmt.Printf("%s", dig)
        }
        fmt.Printf("\n")
    }
}

func printBits(bits []uint64) {
    for i := 0; i < len(bits); i++ {
        for j := 63; j >= 0; j-- {
            dig := "0"
            if (bits[i] >> j) & 1 == 1 {
                dig = "1"
            }
            fmt.Printf("%s", dig)
        }
        fmt.Printf("\n")
    }
}

