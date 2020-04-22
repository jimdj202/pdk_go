package pdk

type Cards []Card

type Card byte

var StraightMask = []uint16{15872, 7936, 3968, 1984, 992, 496, 248, 124, 62, 31}

func (c Card) getCardIndex() int{
	cardValue := c & 0x0f
	//if cardValue == 0x01 {
	//	return 14
	/*}else*/ if cardValue == 0x02{
		return 16
	}else if cardValue == 0x0e{ //小王
		return 17
	}else if cardValue == 0x0f{ //大王
		return 18
	}
	return int(cardValue)
}

func (c *Cards) getType(){
	tempCards := Cards{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}  //-a-k  大小王  15张序列
	for _,v := range *c {
		vIndex := v.getCardIndex()
		tempCards[vIndex] = tempCards[vIndex] + 1 //
	}

	counts := []int{0,0,0,0,0} //--1,2,3,4组合的数目
	cards := []Cards{{},{},{},{},{}} //数目分别为1,2,3,4的牌的序列
	for i,v := range tempCards{
		if v > 0 {
			counts[i] = counts[i] +1
			cards[v] = append(cards[v], Card(i))
		}
	}

	if counts[4] > 0 {

	}else if counts[3] > 0 {

	}else if counts[2] > 0 {

	}else if counts[1] > 0{

	}

}

//func (this *Cards) straight() uint32 {
//	var handvalue uint16
//	for _, v := range (*this) {
//		value := v & 0xF
//		if value == 0xE {
//			handvalue |= 1
//		}
//		handvalue |= (1 << (value - 1 ) )
//	}
//
//	for i := uint8(0); i < 10; i++ {
//		if handvalue&StraightMask[i] == StraightMask[i] {
//			return En(5, uint32(10-i+4))
//		}
//	}
//	return 0
//}
//
//func De(v uint32) (uint8, uint32) {
//	return uint8(v >> 24), v & 0xFFFFFF
//}
//
//func En(t uint8, v uint32) uint32 {
//	v1 := v | ( uint32(t) << 24)
//	return v1
//}

