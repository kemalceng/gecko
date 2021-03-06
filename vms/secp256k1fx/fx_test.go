// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package secp256k1fx

import (
	"testing"
	"time"

	"github.com/ava-labs/gecko/ids"
	"github.com/ava-labs/gecko/utils/crypto"
	"github.com/ava-labs/gecko/utils/hashing"
	"github.com/ava-labs/gecko/utils/logging"
	"github.com/ava-labs/gecko/utils/timer"
	"github.com/ava-labs/gecko/utils/codec"
)

var (
	txBytes  = []byte{0, 1, 2, 3, 4, 5}
	sigBytes = [crypto.SECP256K1RSigLen]byte{
		0x0e, 0x33, 0x4e, 0xbc, 0x67, 0xa7, 0x3f, 0xe8,
		0x24, 0x33, 0xac, 0xa3, 0x47, 0x88, 0xa6, 0x3d,
		0x58, 0xe5, 0x8e, 0xf0, 0x3a, 0xd5, 0x84, 0xf1,
		0xbc, 0xa3, 0xb2, 0xd2, 0x5d, 0x51, 0xd6, 0x9b,
		0x0f, 0x28, 0x5d, 0xcd, 0x3f, 0x71, 0x17, 0x0a,
		0xf9, 0xbf, 0x2d, 0xb1, 0x10, 0x26, 0x5c, 0xe9,
		0xdc, 0xc3, 0x9d, 0x7a, 0x01, 0x50, 0x9d, 0xe8,
		0x35, 0xbd, 0xcb, 0x29, 0x3a, 0xd1, 0x49, 0x32,
		0x00,
	}
	addrBytes = [hashing.AddrLen]byte{
		0x01, 0x5c, 0xce, 0x6c, 0x55, 0xd6, 0xb5, 0x09,
		0x84, 0x5c, 0x8c, 0x4e, 0x30, 0xbe, 0xd9, 0x8d,
		0x39, 0x1a, 0xe7, 0xf0,
	}
)

type testVM struct{ clock timer.Clock }

func (vm *testVM) Codec() codec.Codec { return codec.NewDefault() }

func (vm *testVM) Clock() *timer.Clock { return &vm.clock }

func (vm *testVM) Logger() logging.Logger { return logging.NoLog{} }

type testCodec struct{}

func (c *testCodec) RegisterStruct(interface{}) {}

type testTx struct{ bytes []byte }

func (tx *testTx) UnsignedBytes() []byte { return tx.bytes }

func TestFxInitialize(t *testing.T) {
	vm := testVM{}
	fx := Fx{}
	err := fx.Initialize(&vm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFxInitializeInvalid(t *testing.T) {
	fx := Fx{}
	err := fx.Initialize(nil)
	if err == nil {
		t.Fatalf("Should have returned an error")
	}
}

func TestFxVerifyTransfer(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	if err := fx.Bootstrapping(); err != nil {
		t.Fatal(err)
	}
	if err := fx.Bootstrapped(); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	out := &TransferOutput{
		Amt:      1,
		Locktime: 0,
		OutputOwners: OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				ids.NewShortID(addrBytes),
			},
		},
	}
	in := &TransferInput{
		Amt: 1,
		Input: Input{
			SigIndices: []uint32{0},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			sigBytes,
		},
	}

	if err := fx.VerifyTransfer(tx, in, cred, out); err != nil {
		t.Fatal(err)
	}
}

func TestFxVerifyTransferNilTx(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	out := &TransferOutput{
		Amt:      1,
		Locktime: 0,
		OutputOwners: OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				ids.NewShortID(addrBytes),
			},
		},
	}
	in := &TransferInput{
		Amt: 1,
		Input: Input{
			SigIndices: []uint32{0},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			sigBytes,
		},
	}

	if err := fx.VerifyTransfer(nil, in, cred, out); err == nil {
		t.Fatalf("Should have failed verification due to a nil tx")
	}
}

func TestFxVerifyTransferNilOutput(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	in := &TransferInput{
		Amt: 1,
		Input: Input{
			SigIndices: []uint32{0},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			sigBytes,
		},
	}

	if err := fx.VerifyTransfer(tx, in, cred, nil); err == nil {
		t.Fatalf("Should have failed verification due to a nil output")
	}
}

func TestFxVerifyTransferNilInput(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	out := &TransferOutput{
		Amt:      1,
		Locktime: 0,
		OutputOwners: OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				ids.NewShortID(addrBytes),
			},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			sigBytes,
		},
	}

	if err := fx.VerifyTransfer(tx, nil, cred, out); err == nil {
		t.Fatalf("Should have failed verification due to a nil input")
	}
}

func TestFxVerifyTransferNilCredential(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	out := &TransferOutput{
		Amt:      1,
		Locktime: 0,
		OutputOwners: OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				ids.NewShortID(addrBytes),
			},
		},
	}
	in := &TransferInput{
		Amt: 1,
		Input: Input{
			SigIndices: []uint32{0},
		},
	}

	if err := fx.VerifyTransfer(tx, in, nil, out); err == nil {
		t.Fatalf("Should have failed verification due to a nil credential")
	}
}

func TestFxVerifyTransferInvalidOutput(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	out := &TransferOutput{
		Amt:      1,
		Locktime: 0,
		OutputOwners: OutputOwners{
			Threshold: 0,
			Addrs: []ids.ShortID{
				ids.NewShortID(addrBytes),
			},
		},
	}
	in := &TransferInput{
		Amt: 1,
		Input: Input{
			SigIndices: []uint32{0},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			sigBytes,
		},
	}

	if err := fx.VerifyTransfer(tx, in, cred, out); err == nil {
		t.Fatalf("Should have errored due to an invalid output")
	}
}

func TestFxVerifyTransferWrongAmounts(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	out := &TransferOutput{
		Amt:      1,
		Locktime: 0,
		OutputOwners: OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				ids.NewShortID(addrBytes),
			},
		},
	}
	in := &TransferInput{
		Amt: 2,
		Input: Input{
			SigIndices: []uint32{0},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			sigBytes,
		},
	}

	if err := fx.VerifyTransfer(tx, in, cred, out); err == nil {
		t.Fatalf("Should have errored due to different amounts")
	}
}

func TestFxVerifyTransferTimelocked(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	out := &TransferOutput{
		Amt:      1,
		Locktime: uint64(date.Add(time.Second).Unix()),
		OutputOwners: OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				ids.NewShortID(addrBytes),
			},
		},
	}
	in := &TransferInput{
		Amt: 1,
		Input: Input{
			SigIndices: []uint32{0},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			sigBytes,
		},
	}

	if err := fx.VerifyTransfer(tx, in, cred, out); err == nil {
		t.Fatalf("Should have errored due to a timelocked output")
	}
}

func TestFxVerifyTransferTooManySigners(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	out := &TransferOutput{
		Amt:      1,
		Locktime: 0,
		OutputOwners: OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				ids.NewShortID(addrBytes),
			},
		},
	}
	in := &TransferInput{
		Amt: 1,
		Input: Input{
			SigIndices: []uint32{0, 1},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			sigBytes,
			[crypto.SECP256K1RSigLen]byte{},
		},
	}

	if err := fx.VerifyTransfer(tx, in, cred, out); err == nil {
		t.Fatalf("Should have errored due to too many signers")
	}
}

func TestFxVerifyTransferTooFewSigners(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	out := &TransferOutput{
		Amt:      1,
		Locktime: 0,
		OutputOwners: OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				ids.NewShortID(addrBytes),
			},
		},
	}
	in := &TransferInput{
		Amt: 1,
		Input: Input{
			SigIndices: []uint32{},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{},
	}

	if err := fx.VerifyTransfer(tx, in, cred, out); err == nil {
		t.Fatalf("Should have errored due to too few signers")
	}
}

func TestFxVerifyTransferMismatchedSigners(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	out := &TransferOutput{
		Amt:      1,
		Locktime: 0,
		OutputOwners: OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				ids.NewShortID(addrBytes),
			},
		},
	}
	in := &TransferInput{
		Amt: 1,
		Input: Input{
			SigIndices: []uint32{0},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			sigBytes,
			[crypto.SECP256K1RSigLen]byte{},
		},
	}

	if err := fx.VerifyTransfer(tx, in, cred, out); err == nil {
		t.Fatalf("Should have errored due to too mismatched signers")
	}
}

func TestFxVerifyTransferInvalidSignature(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	if err := fx.Bootstrapping(); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	out := &TransferOutput{
		Amt:      1,
		Locktime: 0,
		OutputOwners: OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				ids.NewShortID(addrBytes),
			},
		},
	}
	in := &TransferInput{
		Amt: 1,
		Input: Input{
			SigIndices: []uint32{0},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			[crypto.SECP256K1RSigLen]byte{},
		},
	}

	if err := fx.VerifyTransfer(tx, in, cred, out); err != nil {
		t.Fatal(err)
	}

	if err := fx.Bootstrapped(); err != nil {
		t.Fatal(err)
	}

	if err := fx.VerifyTransfer(tx, in, cred, out); err == nil {
		t.Fatalf("Should have errored due to an invalid signature")
	}
}

func TestFxVerifyTransferWrongSigner(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	if err := fx.Bootstrapping(); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	out := &TransferOutput{
		Amt:      1,
		Locktime: 0,
		OutputOwners: OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				ids.ShortEmpty,
			},
		},
	}
	in := &TransferInput{
		Amt: 1,
		Input: Input{
			SigIndices: []uint32{0},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			sigBytes,
		},
	}

	if err := fx.VerifyTransfer(tx, in, cred, out); err != nil {
		t.Fatal(err)
	}

	if err := fx.Bootstrapped(); err != nil {
		t.Fatal(err)
	}

	if err := fx.VerifyTransfer(tx, in, cred, out); err == nil {
		t.Fatalf("Should have errored due to a wrong signer")
	}
}

func TestFxVerifyOperation(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	utxo := &MintOutput{
		OutputOwners: OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				ids.NewShortID(addrBytes),
			},
		},
	}
	op := &MintOperation{
		MintInput: Input{
			SigIndices: []uint32{0},
		},
		MintOutput: MintOutput{
			OutputOwners: OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					ids.NewShortID(addrBytes),
				},
			},
		},
		TransferOutput: TransferOutput{
			Amt:      1,
			Locktime: 0,
			OutputOwners: OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					ids.NewShortID(addrBytes),
				},
			},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			sigBytes,
		},
	}

	utxos := []interface{}{utxo}
	err := fx.VerifyOperation(tx, op, cred, utxos)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFxVerifyOperationUnknownTx(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	utxo := &MintOutput{
		OutputOwners: OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				ids.NewShortID(addrBytes),
			},
		},
	}
	op := &MintOperation{
		MintInput: Input{
			SigIndices: []uint32{0},
		},
		MintOutput: MintOutput{
			OutputOwners: OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					ids.NewShortID(addrBytes),
				},
			},
		},
		TransferOutput: TransferOutput{
			Amt:      1,
			Locktime: 0,
			OutputOwners: OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					ids.NewShortID(addrBytes),
				},
			},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			sigBytes,
		},
	}

	utxos := []interface{}{utxo}
	err := fx.VerifyOperation(nil, op, cred, utxos)
	if err == nil {
		t.Fatalf("Should have errored due to an invalid tx type")
	}
}

func TestFxVerifyOperationUnknownOperation(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	utxo := &MintOutput{
		OutputOwners: OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				ids.NewShortID(addrBytes),
			},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			sigBytes,
		},
	}

	utxos := []interface{}{utxo}
	err := fx.VerifyOperation(tx, nil, cred, utxos)
	if err == nil {
		t.Fatalf("Should have errored due to an invalid operation type")
	}
}

func TestFxVerifyOperationUnknownCredential(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	utxo := &MintOutput{
		OutputOwners: OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				ids.NewShortID(addrBytes),
			},
		},
	}
	op := &MintOperation{
		MintInput: Input{
			SigIndices: []uint32{0},
		},
		MintOutput: MintOutput{
			OutputOwners: OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					ids.NewShortID(addrBytes),
				},
			},
		},
		TransferOutput: TransferOutput{
			Amt:      1,
			Locktime: 0,
			OutputOwners: OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					ids.NewShortID(addrBytes),
				},
			},
		},
	}

	utxos := []interface{}{utxo}
	err := fx.VerifyOperation(tx, op, nil, utxos)
	if err == nil {
		t.Fatalf("Should have errored due to an invalid credential type")
	}
}

func TestFxVerifyOperationWrongNumberOfUTXOs(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	utxo := &MintOutput{
		OutputOwners: OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				ids.NewShortID(addrBytes),
			},
		},
	}
	op := &MintOperation{
		MintInput: Input{
			SigIndices: []uint32{0},
		},
		MintOutput: MintOutput{
			OutputOwners: OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					ids.NewShortID(addrBytes),
				},
			},
		},
		TransferOutput: TransferOutput{
			Amt:      1,
			Locktime: 0,
			OutputOwners: OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					ids.NewShortID(addrBytes),
				},
			},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			sigBytes,
		},
	}

	utxos := []interface{}{utxo, utxo}
	err := fx.VerifyOperation(tx, op, cred, utxos)
	if err == nil {
		t.Fatalf("Should have errored due to a wrong number of utxos")
	}
}

func TestFxVerifyOperationUnknownUTXOType(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	op := &MintOperation{
		MintInput: Input{
			SigIndices: []uint32{0},
		},
		MintOutput: MintOutput{
			OutputOwners: OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					ids.NewShortID(addrBytes),
				},
			},
		},
		TransferOutput: TransferOutput{
			Amt:      1,
			Locktime: 0,
			OutputOwners: OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					ids.NewShortID(addrBytes),
				},
			},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			sigBytes,
		},
	}

	utxos := []interface{}{nil}
	err := fx.VerifyOperation(tx, op, cred, utxos)
	if err == nil {
		t.Fatalf("Should have errored due to an invalid utxo type")
	}
}

func TestFxVerifyOperationInvalidOperationVerify(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	utxo := &MintOutput{
		OutputOwners: OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				ids.NewShortID(addrBytes),
			},
		},
	}
	op := &MintOperation{
		MintInput: Input{
			SigIndices: []uint32{0},
		},
		MintOutput: MintOutput{
			OutputOwners: OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					ids.NewShortID(addrBytes),
				},
			},
		},
		TransferOutput: TransferOutput{
			Amt:      1,
			Locktime: 0,
			OutputOwners: OutputOwners{
				Threshold: 1,
			},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			sigBytes,
		},
	}

	utxos := []interface{}{utxo}
	err := fx.VerifyOperation(tx, op, cred, utxos)
	if err == nil {
		t.Fatalf("Should have errored due to a failed verify")
	}
}

func TestFxVerifyOperationMismatchedMintOutputs(t *testing.T) {
	vm := testVM{}
	date := time.Date(2019, time.January, 19, 16, 25, 17, 3, time.UTC)
	vm.clock.Set(date)
	fx := Fx{}
	if err := fx.Initialize(&vm); err != nil {
		t.Fatal(err)
	}
	tx := &testTx{
		bytes: txBytes,
	}
	utxo := &MintOutput{
		OutputOwners: OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				ids.NewShortID(addrBytes),
			},
		},
	}
	op := &MintOperation{
		MintInput: Input{
			SigIndices: []uint32{0},
		},
		MintOutput: MintOutput{
			OutputOwners: OutputOwners{},
		},
		TransferOutput: TransferOutput{
			Amt:      1,
			Locktime: 0,
			OutputOwners: OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					ids.NewShortID(addrBytes),
				},
			},
		},
	}
	cred := &Credential{
		Sigs: [][crypto.SECP256K1RSigLen]byte{
			sigBytes,
		},
	}

	utxos := []interface{}{utxo}
	err := fx.VerifyOperation(tx, op, cred, utxos)
	if err == nil {
		t.Fatalf("Should have errored due to the wrong MintOutput being created")
	}
}
