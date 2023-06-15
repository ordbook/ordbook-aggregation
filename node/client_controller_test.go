package node

import (
	"fmt"
	"testing"
)

func TestClientController_GetBlockHeight(t *testing.T) {
	net := "mainnet"
	txHex := "02000000000106473fe246b9ed21f57c0ef2a7c8d0f8bd4b5967ae6ad2c11272224f4ab87eb0fd0a00000000ffffffff473fe246b9ed21f57c0ef2a7c8d0f8bd4b5967ae6ad2c11272224f4ab87eb0fd0b00000000ffffffffc4e84be5714c230c33cc63ad41d6c44b53ef58738bdd459b2546384b8c9dfa2a0000000000ffffffff0011c64febf30c04d8c081655fae7cd0508e68d4dd4ff5043f6b846cfc674bf70500000000ffffffff0011c64febf30c04d8c081655fae7cd0508e68d4dd4ff5043f6b846cfc674bf70600000000ffffffff0011c64febf30c04d8c081655fae7cd0508e68d4dd4ff5043f6b846cfc674bf70700000000ffffffff07b0040000000000001600141907f4c5548517cb6a000c8b135dbef735b2414522020000000000001600140458b9a6887237c13e5233f41d7fc17139a03e5d7017000000000000225120b9f34c0390baf96a02819c8b4f9f8a1ad9bfb4416dfee9d2653fe4146d7ff058302a000000000000160014f26957622882a010e696e6f3ece6c36c0896e8a8580200000000000016001475bb30e7d33e07dcad93fdc63cd83d13c3bc6264580200000000000016001475bb30e7d33e07dcad93fdc63cd83d13c3bc6264b806000000000000160014f26957622882a010e696e6f3ece6c36c0896e8a8024730440220518bad955392db092736cde65feff0704f4c2c55b22dfb9524b0a519aac2ae3d0220145eab2e90088e0b8c33b42f7e8a94da6540123492f73e69750561996c2b16f901210373873908981153c3b8cef9c50ab3b5a512cd073321f35c30de31e7e4e50976f002483045022100cc1829186179a7c2105a4e7faedbe9557bfa76705e3aa66e355002693c968ab40220499211f5550d9525a23361d943700d7b6a9e81e572b61c4ad30b576951d8284c01210373873908981153c3b8cef9c50ab3b5a512cd073321f35c30de31e7e4e50976f00141543841809188d981e281dd378b223c6954f944b37bb04d43ace5bc3498eade4aa3db61793a2bae7a82de32a2df46affbd2ad6cdf2f4fffea574acfd4dc79cfff8302483045022100bffdd734eb4a20b7a6c078afc4bec94b27a88eafa91929a35790909531e8f09f0220597d343c9cf3b14b1f8f94d294d23f6732fa3d33ff4fecff0aed38a50297ca5b012102854dfeba276ee439c3331b729d057899fd1a560e80cb240ba2d21b8b5d7bc09d0247304402203cb15dd92b101b57021b70daa8b8bdf220f02e63715c8577dc98a512970b193d022011f826e64dd91589a879a4e49ab2bd52a5cbe5ba0ca1ca8e61f3a0c50cdcc145012102854dfeba276ee439c3331b729d057899fd1a560e80cb240ba2d21b8b5d7bc09d02483045022100ddc1acf4ac2e48a3a86bc9747c2365e0cdb4b30b406daf4f2678b47d41ce105302203b3c1ba83036b1b27049d4ae5ee871b9323e46968d90a40884a08c13051b1a83012102854dfeba276ee439c3331b729d057899fd1a560e80cb240ba2d21b8b5d7bc09d00000000"
	c := NewClientController(net)
	//fmt.Println(c.GetBlockHeight(net))
	fmt.Println(c.BroadcastTx(net, txHex))
}