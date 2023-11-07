package filter

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {

	text := "文明用语你&* 妈,逼的你这个狗日的，怎么这么傻啊。我也是服了，我日,这些话我都说不出口日"
	fmt.Println(ChangeSensitiveWords(text))
}
