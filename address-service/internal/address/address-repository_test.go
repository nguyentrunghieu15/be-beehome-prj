package address

import (
	"fmt"
	"testing"
)

func Test_fixToTextSearchQuery(t *testing.T) {
	fmt.Println(fixToTextSearchQuery("Cẩm Điền, Cẩm            Giàng, Hải Dương"))
}
