package pdk

import "testing"

func Test_Card(t *testing.T) {
	c := Card(0x18)
	var cc Card
	cc = 0x19
	t.Log(c.getCardIndex())
	t.Log(cc.getCardIndex())

	cList := Cards{0x03,0x04,0x05,0x06,0x07}
	//ret := cList.straight()
	//t.Log(ret)
	cList.getType()
}



