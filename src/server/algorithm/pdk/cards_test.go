package pdk

import "testing"

func Test_Card(t *testing.T) {
	c := Card(0x18)
	var cc Card
	cc = 0x19
	t.Log(c.getCardIndex())
	t.Log(cc.getCardIndex())
}

