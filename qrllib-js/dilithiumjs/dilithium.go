package dilithiumjs

import (
	"encoding/hex"
	"github.com/gopherjs/gopherjs/js"
	"github.com/theQRL/go-qrllib/common"
	"github.com/theQRL/go-qrllib/dilithium"
)

type DilithiumJS struct {
	*js.Object
	d    *dilithium.Dilithium
	pk   string `js:"pk"`
	sk   string `js:"sk"`
	seed string `js:"seed"`
}

func newDilithiumJS(d *dilithium.Dilithium) *js.Object {
	binPK := d.GetPK()
	binSK := d.GetSK()
	binSeed := d.GetSeed()

	djs := DilithiumJS{Object: js.Global.Get("Object").New()}
	djs.d = d
	djs.pk = hex.EncodeToString(binPK[:])
	djs.sk = hex.EncodeToString(binSK[:])
	djs.seed = hex.EncodeToString(binSeed[:])

	djs.Object.Set("GetPK", djs.GetPK)
	djs.Object.Set("GetSK", djs.GetSK)
	djs.Object.Set("GetSeed", djs.GetSeed)
	djs.Object.Set("GetMnemonic", djs.d.GetMnemonic)
	djs.Object.Set("GetAddress", djs.GetAddress)
	djs.Object.Set("Sign", djs.Sign)

	return djs.Object
}

func NewDilithiumJS() *js.Object {
	d := dilithium.New()
	return newDilithiumJS(d)
}

func NewDilithiumJSFromSeed(seed string) *js.Object {
	binSeed, err := hex.DecodeString(seed)
	if err != nil {
		return nil
	}
	var sizedBinSeed [common.SeedSize]uint8
	copy(sizedBinSeed[:], binSeed)
	d := dilithium.NewDilithiumFromSeed(sizedBinSeed)

	return newDilithiumJS(d)
}

func (d *DilithiumJS) GetPK() string {
	return d.pk
}

func (d *DilithiumJS) GetSK() string {
	return d.sk
}

func (d *DilithiumJS) GetSeed() string {
	return d.seed
}

func (d *DilithiumJS) GetAddress() string {
	binAddr := d.d.GetAddress()
	return hex.EncodeToString(binAddr[:])
}

func (d *DilithiumJS) Sign(message string) string {
	binSignature := d.d.Sign([]uint8(message))
	return hex.EncodeToString(binSignature[:])
}

func DilithiumVerify(message string, signature string, pk string) bool {
	binMessage := []uint8(message)
	binSignature, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}
	binPK, err := hex.DecodeString(pk)
	if err != nil {
		return false
	}

	var sizedBinPK [dilithium.PKSizePacked]uint8
	copy(sizedBinPK[:], binPK)

	return dilithium.Verify(binMessage, binSignature, &sizedBinPK)
}

func GetDilithiumAddressFromPK(pk string) string {
	binPK, err := hex.DecodeString(pk)
	if err != nil {
		return ""
	}

	var sizedBinPK [dilithium.PKSizePacked]uint8
	copy(sizedBinPK[:], binPK)

	binAddress := dilithium.GetDilithiumAddressFromPK(sizedBinPK)

	return hex.EncodeToString(binAddress[:])
}

func IsValidDilithiumAddress(address string) bool {
	binAddr, err := hex.DecodeString(address)
	if err != nil {
		return false
	}

	var sizedBinAddr [common.AddressSize]uint8
	copy(sizedBinAddr[:], binAddr)

	return dilithium.IsValidDilithiumAddress(sizedBinAddr)
}