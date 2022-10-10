/*
 * @Author: webees@qq.com
 * @Date: 2021-03-29 18:10:06
 * @Last Modified by:   webees@qq.com
 * @Last Modified time: 2021-03-29 18:10:06
 */
package fil

import (
	"fmt"
	"testing"
	"time"

	"github.com/webees/hdwallet/bip32"
	"github.com/webees/hdwallet/bip39"

	"github.com/gogf/gf/container/gmap"
)

func TestFil(t *testing.T) {
	fmt.Println("\n######################################## FIL ########################################")
	mnemonic := "owner mosquito uphold squirrel utility fat warrior wheat vital chapter shoulder horn"
	seed, _ := bip39.NewSeed(mnemonic, "")
	pvk, _ := bip32.NewMaster(seed, HDPrivateKeyID)
	pbk, _ := Xpub(pvk.String(), 0)
	fmt.Println("\n助记词：  ", mnemonic)
	fmt.Println("扩展私钥: ", pvk.String())
	fmt.Println("账户公钥：", pbk)
	addr := gmap.New(true)
	start := time.Now()
	for i := 0; i < 9; i++ {
		s, _ := Addr(pbk, uint32(i))
		addr.Set(i, s)
		fmt.Println("地址", i, "=", s)
	}
	elapsed := time.Since(start)
	fmt.Println("\n耗时：", elapsed)
}

func TestTfil(t *testing.T) {
	TEST = true
	fmt.Println("\n######################################## tFIL ########################################")
	mnemonic := "owner mosquito uphold squirrel utility fat warrior wheat vital chapter shoulder horn"
	seed, _ := bip39.NewSeed(mnemonic, "")
	pvk, _ := bip32.NewMaster(seed, HDPrivateKeyID)
	pbk, _ := Xpub(pvk.String(), 0)
	fmt.Println("\n助记词：  ", mnemonic)
	fmt.Println("扩展私钥: ", pvk.String())
	fmt.Println("账户公钥：", pbk)
	addr := gmap.New(true)
	start := time.Now()
	for i := 0; i < 9; i++ {
		s, _ := Addr(pbk, uint32(i))
		addr.Set(i, s)
		fmt.Println("地址", i, "=", s)
	}
	elapsed := time.Since(start)
	fmt.Println("\n耗时：", elapsed)
}
