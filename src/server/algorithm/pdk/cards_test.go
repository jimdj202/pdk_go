package pdk

import "testing"

type IntList []int

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

func Test_Card1(t *testing.T){
	cards := &Cards{0x13,0x23,0x33,0x43}
	ret := cards.getType()
	t.Log(ret)


}

func Test_List(t *testing.T){
	ttt := []IntList{}
	t.Log(ttt)
	ttt2 := []IntList{{}}
	t.Log(ttt2)
}

func Test_Or(t *testing.T){
	t1 := 0x12
	tf := 0xf0
	a:= t1 & tf
	b :=  0x12 >> 2
	t.Log(a)
	t.Log(b)
}


