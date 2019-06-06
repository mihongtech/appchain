package state

import (
	"fmt"
	"testing"

	"github.com/mihongtech/linkchain-core/common/math"
)

func TestSample(t *testing.T) {
	var emptyCodeHash = math.HashH(nil)
	fmt.Println(emptyCodeHash)

	var emptyByte []byte
	emptyByteHash := math.HashH(emptyByte)
	fmt.Println(emptyByteHash)
}
