package teal

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/algorand/go-algorand-sdk/types"
	"github.com/pkg/errors"
)

/*
	TODO: add missing specs for pseudops
*/

type NewOpArgType int

const (
	OpArgTypeNone = iota
	OpArgTypeUint64
	OpArgTypeUint8
	OpArgTypeInt8
	OpArgTypeBytes
	OpArgTypeTxnField
	OpArgTypeJSONRefField
	OpArgTypeEcdsaCurve
	OpArgTypeGlobalField
	OpArgTypeLabel
	OpArgTypeAssetHoldingField
	OpArgTypeAssetParamsField
	OpArgTypeAppParamsField
	OpArgTypeAcctParamsField
	OpArgTypeVrfStandard
	OpArgTypeSignature
	OpArgTypeAddr
	OpArgTypePragmaName
)

func (t NewOpArgType) String() string {
	switch t {
	default:
		return "(none)"
	case OpArgTypeUint64:
		return "uint64"
	case OpArgTypeUint8:
		return "uint8"
	case OpArgTypeInt8:
		return "int8"
	case OpArgTypeBytes:
		return "bytes"
	case OpArgTypeTxnField:
		return "transaction field index"
	case OpArgTypeJSONRefField:
		return "json_Ref"
	case OpArgTypeEcdsaCurve:
		return "ECDSA Curve"
	case OpArgTypeGlobalField:
		return "global field index"
	case OpArgTypeLabel:
		return "label name"
	case OpArgTypeAssetHoldingField:
		return "asset holding field index"
	case OpArgTypeAssetParamsField:
		return "asset params field index"
	case OpArgTypeAppParamsField:
		return "app params field index"
	case OpArgTypeAcctParamsField:
		return "account params field index"
	case OpArgTypeVrfStandard:
		return "parameters index"
	case OpArgTypeSignature:
		return "signature"
	case OpArgTypeAddr:
		return "address"
	case OpArgTypePragmaName:
		return "pragma name"
	}
}

type opItemArg struct {
	Name     string
	Type     NewOpArgType
	Array    bool
	Optional bool
}

type opItem struct {
	Name string

	Version uint8

	Args []opItemArg

	ArgsSig string
	FullSig string

	Parse parserFunc

	Doc     string
	FullDoc string
}

type opListItem struct {
	Name  string
	Parse parserFunc
}

var opsList = []opListItem{
	{"replace", opReplace},
	{"byte", opByte},
	{"int", opInt},
	{"method", opMethod},
	{"addr", opAddr},
	{"err", opErr},
	{"sha256", opSHA256},
	{"keccak256", opKeccak256},
	{"sha512_256", opSHA512_256},
	{"sha256", opSHA256},
	{"keccak256", opKeccak256},
	{"sha512_256", opSHA512_256},
	{"ed25519verify", opEd25519Verify},
	{"ecdsa_verify", opEcdsaVerify},
	{"ecdsa_pk_decompress", opEcdsaPkDecompress},
	{"ecdsa_pk_recover", opEcdsaPkRecover},
	{"+", opPlus},
	{"-", opMinus},
	{"/", opDiv},
	{"*", opMul},
	{"<", opLt},
	{">", opGt},
	{"<=", opLe},
	{">=", opGe},
	{"&&", opAnd},
	{"||", opOr},
	{"==", opEq},
	{"!=", opNeq},
	{"!", opNot},
	{"len", opLen},
	{"itob", opItob},
	{"btoi", opBtoi},
	{"%", opModulo},
	{"|", opBitOr},
	{"&", opBitAnd},
	{"^", opBitXor},
	{"~", opBitNot},
	{"mulw", opMulw},
	{"addw", opAddw},
	{"divmodw", opDivModw},
	{"intcblock", opIntConstBlock},
	{"intc", opIntConstLoad},
	{"intc_0", opIntConst0},
	{"intc_1", opIntConst1},
	{"intc_2", opIntConst2},
	{"intc_3", opIntConst3},
	{"bytecblock", opByteConstBlock},
	{"bytec", opByteConstLoad},
	{"bytec_0", opByteConst0},
	{"bytec_1", opByteConst1},
	{"bytec_2", opByteConst2},
	{"bytec_3", opByteConst3},
	{"arg", opArg},
	{"arg_0", opArg0},
	{"arg_1", opArg1},
	{"arg_2", opArg2},
	{"arg_3", opArg3},
	{"txn", opTxn},
	{"global", opGlobal},
	{"gtxn", opGtxn},
	{"load", opLoad},
	{"store", opStore},
	{"txna", opTxna},
	{"gtxna", opGtxna},
	{"gtxns", opGtxns},
	{"gtxnsa", opGtxnsa},
	{"gload", opGload},
	{"gloads", opGloads},
	{"gaid", opGaid},
	{"gaids", opGaids},
	{"loads", opLoads},
	{"stores", opStores},
	{"bnz", opBnz},
	{"bz", opBz},
	{"b", opB},
	{"return", opReturn},
	{"assert", opAssert},
	{"bury", opBury},
	{"popn", opPopN},
	{"dupn", opDupN},
	{"pop", opPop},
	{"dup", opDup},
	{"dup2", opDup2},
	{"dig", opDig},
	{"swap", opSwap},
	{"select", opSelect},
	{"cover", opCover},
	{"uncover", opUncover},
	{"concat", opConcat},
	{"substring", opSubstring},
	{"substring3", opSubstring3},
	{"getbit", opGetBit},
	{"setbit", opSetBit},
	{"getbyte", opGetByte},
	{"setbyte", opSetByte},
	{"extract", opExtract},
	{"extract3", opExtract3},
	{"extract_uint16", opExtract16Bits},
	{"extract_uint32", opExtract32Bits},
	{"extract_uint64", opExtract64Bits},
	{"replace2", opReplace2},
	{"replace3", opReplace3},
	{"base64_decode", opBase64Decode},
	{"json_ref", opJSONRef},
	{"balance", opBalance},
	{"balance", opBalance},
	{"app_opted_in", opAppOptedIn},
	{"app_opted_in", opAppOptedIn},
	{"app_local_get", opAppLocalGet},
	{"app_local_get", opAppLocalGet},
	{"app_local_get_ex", opAppLocalGetEx},
	{"app_local_get_ex", opAppLocalGetEx},
	{"app_global_get", opAppGlobalGet},
	{"app_global_get_ex", opAppGlobalGetEx},
	{"app_local_put", opAppLocalPut},
	{"app_local_put", opAppLocalPut},
	{"app_global_put", opAppGlobalPut},
	{"app_local_del", opAppLocalDel},
	{"app_local_del", opAppLocalDel},
	{"app_global_del", opAppGlobalDel},
	{"asset_holding_get", opAssetHoldingGet},
	{"asset_holding_get", opAssetHoldingGet},
	{"asset_params_get", opAssetParamsGet},
	{"app_params_get", opAppParamsGet},
	{"acct_params_get", opAcctParamsGet},
	{"min_balance", opMinBalance},
	{"min_balance", opMinBalance},
	{"pushbytes", opPushBytes},
	{"pushint", opPushInt},
	{"pushbytess", opPushBytess},
	{"pushints", opPushInts},
	{"ed25519verify_bare", opEd25519VerifyBare},
	{"callsub", opCallSub},
	{"retsub", opRetSub},
	{"proto", opProto},
	{"frame_dig", opFrameDig},
	{"frame_bury", opFrameBury},
	{"switch", opSwitch},
	{"match", opMatch},
	{"shl", opShiftLeft},
	{"shr", opShiftRight},
	{"sqrt", opSqrt},
	{"bitlen", opBitLen},
	{"exp", opExp},
	{"expw", opExpw},
	{"bsqrt", opBytesSqrt},
	{"divw", opDivw},
	{"sha3_256", opSHA3_256},
	{"bn256_add", opBn256Add},
	{"bn256_scalar_mul", opBn256ScalarMul},
	{"bn256_pairing", opBn256Pairing},
	{"b+", opBytesPlus},
	{"b-", opBytesMinus},
	{"b/", opBytesDiv},
	{"b*", opBytesMul},
	{"b<", opBytesLt},
	{"b>", opBytesGt},
	{"b<=", opBytesLe},
	{"b>=", opBytesGe},
	{"b==", opBytesEq},
	{"b!=", opBytesNeq},
	{"b%", opBytesModulo},
	{"b|", opBytesBitOr},
	{"b&", opBytesBitAnd},
	{"b^", opBytesBitXor},
	{"b~", opBytesBitNot},
	{"bzero", opBytesZero},
	{"log", opLog},
	{"itxn_begin", opTxBegin},
	{"itxn_field", opItxnField},
	{"itxn_submit", opItxnSubmit},
	{"itxn", opItxn},
	{"itxna", opItxna},
	{"itxn_next", opItxnNext},
	{"gitxn", opGitxn},
	{"gitxna", opGitxna},
	{"box_create", opBoxCreate},
	{"box_extract", opBoxExtract},
	{"box_replace", opBoxReplace},
	{"box_del", opBoxDel},
	{"box_len", opBoxLen},
	{"box_get", opBoxGet},
	{"box_put", opBoxPut},
	{"txnas", opTxnas},
	{"gtxnas", opGtxnas},
	{"gtxnsas", opGtxnsas},
	{"args", opArgs},
	{"gloadss", opGloadss},
	{"itxnas", opItxnas},
	{"gitxnas", opGitxnas},
	{"vrf_verify", opVrfVerify},
	{"block", opBlock},
}

type recoverable struct{}

type OpContext interface {
	emit(op Op)

	minVersion(v uint8)

	mustReadPragma(name string) uint8
	mustReadAddr(name string) string
	mustReadSignature(name string) string
	mustReadTxnField(name string) TxnField
	mustReadJsonRef(name string) JSONRefType
	maybeReadUint8(name string) (uint8, bool)
	mustReadEcdsaCurveIndex(name string) EcdsaCurve
	readUint64Array(name string) []uint64
	mustReadUint8(name string) uint8
	mustReadInt8(name string) int8
	readBytesArray(name string) [][]byte
	mustReadGlobalField(name string) GlobalField
	mustReadInt(name string) uint64
	mustReadLabel(name string) string
	readLabelsArray(name string) []string
	mustReadAssetHoldingField(name string) AssetHoldingField
	mustReadAssetParamsField(name string) AssetParamsField
	mustReadAppParamsField(name string) AppParamsField
	mustReadAcctParamsField(name string) AcctParamsField
	mustReadBytes(name string) []byte
	mustReadVrfVerifyField(name string) VrfStandard
	mustReadBlockField(name string) BlockField
}

type docContext struct {
	args     []opItemArg
	version  uint8
	optional bool
}

func (c *docContext) arg(a opItemArg) {
	if c.optional {
		a.Optional = true
	} else if a.Optional {
		c.optional = true
	}

	c.args = append(c.args, a)
}

func (c *docContext) emit(op Op) {

}

func (c *docContext) minVersion(v uint8) {
	c.version = v
}

func (c *docContext) mustReadPragma(name string) (v uint8) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypePragmaName,
	})

	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeUint8,
	})

	return
}

func (c *docContext) mustReadAddr(name string) (v string) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeAddr,
	})

	return
}

func (c *docContext) mustReadSignature(name string) (v string) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeSignature,
	})

	return
}

func (c *docContext) mustReadJsonRef(name string) (v JSONRefType) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeTxnField,
	})

	return
}

func (c *docContext) mustReadTxnField(name string) (f TxnField) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeTxnField,
	})

	return
}

func (c *docContext) maybeReadUint8(name string) (v uint8, ok bool) {
	c.arg(opItemArg{
		Name:     name,
		Type:     OpArgTypeUint8,
		Optional: true,
	})

	ok = true

	return
}

func (c *docContext) mustReadEcdsaCurveIndex(name string) (v EcdsaCurve) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeTxnField,
	})

	return
}
func (c *docContext) readUint64Array(name string) (v []uint64) {
	c.arg(opItemArg{
		Name:  name,
		Type:  OpArgTypeUint64,
		Array: true,
	})

	return
}
func (c *docContext) mustReadUint8(name string) (v uint8) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeUint8,
	})

	return
}

func (c *docContext) mustReadInt8(name string) (v int8) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeInt8,
	})

	return
}
func (c *docContext) readBytesArray(name string) (v [][]byte) {
	c.arg(opItemArg{
		Name:  name,
		Type:  OpArgTypeBytes,
		Array: true,
	})

	return
}
func (c *docContext) mustReadGlobalField(name string) (v GlobalField) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeGlobalField,
	})

	return
}
func (c *docContext) mustReadInt(name string) (v uint64) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeUint64,
	})

	return
}
func (c *docContext) mustReadLabel(name string) (v string) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeLabel,
	})

	return
}
func (c *docContext) readLabelsArray(name string) (v []string) {
	c.arg(opItemArg{
		Name:  name,
		Type:  OpArgTypeLabel,
		Array: true,
	})

	return
}

func (c *docContext) mustReadAssetHoldingField(name string) (v AssetHoldingField) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeAssetHoldingField,
	})

	return
}
func (c *docContext) mustReadAssetParamsField(name string) (v AssetParamsField) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeAssetParamsField,
	})

	return
}
func (c *docContext) mustReadAppParamsField(name string) (v AppParamsField) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeAppParamsField,
	})

	return
}
func (c *docContext) mustReadAcctParamsField(name string) (v AcctParamsField) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeAcctParamsField,
	})

	return
}
func (c *docContext) mustReadBytes(name string) (v []byte) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeBytes,
	})

	return
}
func (c *docContext) mustReadVrfVerifyField(name string) (v VrfStandard) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeVrfStandard,
	})

	return
}
func (c *docContext) mustReadBlockField(name string) (v BlockField) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeTxnField,
	})

	return
}

type parserContext struct {
	ops  []Op
	args *arguments
	diag []Diagnostic

	nums []Token
	strs []Token
	keys []Token
	mcrs []Token
	refs []Token
}

func (c *parserContext) emit(op Op) {
	c.ops = append(c.ops, op)
}

func (c *parserContext) minVersion(v uint8) {

}

func (c *parserContext) failAt(l int, b int, e int, err error) {
	c.diag = append(c.diag, parseError{l: l, b: b, e: e, error: err})
	panic(recoverable{})
}

func (c *parserContext) failToken(t Token, err error) {
	c.failAt(t.l, t.b, t.e, err)
}

func (c *parserContext) failCurr(err error) {
	c.failToken(c.args.Curr(), err)
}

func (c *parserContext) failPrev(err error) {
	c.failToken(c.args.Prev(), err)
}

func (c *parserContext) mustReadArg(name string) {
	if !c.args.Scan() {
		c.failPrev(errors.Errorf("missing arg: %s", name))
	}
}

func (c *parserContext) mustReadPragma(argName string) uint8 {
	var version uint8
	c.mcrs = append(c.mcrs, c.args.Curr())

	name := c.mustRead("name")
	switch name {
	case "version":
		c.mcrs = append(c.mcrs, c.args.Curr())
		v := c.mustReadUint8("version value")
		if v < 1 {
			c.failCurr(errors.New("version must be at least 1"))
		}
		version = v
	default:
		c.failCurr(errors.Errorf("unexpected #pragma: %s", c.args.Text()))
	}

	return version
}

func (c *parserContext) mustReadAddr(name string) string {
	value := c.mustRead("address")

	_, err := types.DecodeAddress(value)
	if err != nil {
		c.failCurr(err)
	}

	c.strs = append(c.strs, c.args.Curr())

	return value
}

func (c *parserContext) mustReadSignature(name string) string {
	c.mustReadArg(name)

	value := c.args.Text()

	b := 0
	e := len(value) - 1
	if value[b] != '"' || value[e] != '"' {
		c.failCurr(errors.New("missing quotes"))
	}

	// TODO: validate method sig
	c.strs = append(c.strs, c.args.Curr())

	return c.args.Text()
}

func (c *parserContext) maybeReadArg() bool {
	return c.args.Scan()
}

func (c *parserContext) mustReadAcctParamsField(name string) AcctParamsField {
	c.mustReadArg(name)

	f, err := readAcctParams(c.args.Text())
	if err != nil {
		c.failCurr(err)
	}

	return f
}

func (c *parserContext) mustReadAssetHoldingField(name string) AssetHoldingField {
	c.mustReadArg(name)

	f, err := readAssetHoldingField(c.args.Text())
	if err != nil {
		c.failCurr(err)
	}

	return f
}

func (c *parserContext) mustReadAssetParamsField(name string) AssetParamsField {
	c.mustReadArg(name)

	f, err := readAssetParamsField(c.args.Text())
	if err != nil {
		c.failCurr(err)
	}

	return f
}

func (c *parserContext) mustReadBlockField(name string) BlockField {
	c.mustReadArg(name)

	f, err := readBlockField(c.args.Text())
	if err != nil {
		c.failCurr(err)
	}

	return f
}

func (c *parserContext) mustReadGlobalField(name string) GlobalField {
	c.mustReadArg(name)

	f, err := readGlobalField(c.args.Text())
	if err != nil {
		c.failCurr(err)
	}

	return f
}

func (c *parserContext) mustReadLabel(name string) string {
	c.mustReadArg(name)
	c.refs = append(c.refs, c.args.Curr())
	return c.args.Text()
}

func (c *parserContext) mustReadVrfVerifyField(name string) VrfStandard {
	c.mustReadArg(name)

	f, err := readVrfVerifyField(c.args.Text())
	if err != nil {
		c.failCurr(err)
	}

	return f
}

func (c *parserContext) readBytesArray(name string) [][]byte {
	res := [][]byte{}

	for c.args.Scan() {
		bs := c.parseBytes(name)
		res = append(res, bs)
	}

	return res
}

func (c *parserContext) readLabelsArray(name string) []string {
	res := []string{}

	for c.args.Scan() {
		res = append(res, c.args.Text())
		c.refs = append(c.refs, c.args.Curr())
	}

	return res
}

func (c *parserContext) readUint64Array(name string) []uint64 {
	res := []uint64{}

	for c.args.Scan() {
		i := c.parseUint64(name)
		res = append(res, i)
		c.nums = append(c.nums, c.args.Curr())
	}

	return res
}

func (c *parserContext) mustReadAppParamsField(name string) AppParamsField {
	c.mustReadArg(name)

	f, err := readAppParamsField(c.args.Text())
	if err != nil {
		c.failCurr(err)
	}

	return f
}
func (c *parserContext) parseUint64(name string) uint64 {
	v, err := readInt(c.args)
	if err != nil {
		c.failCurr(errors.Wrapf(err, "failed to parse uint64: %s", name))
	}

	c.nums = append(c.nums, c.args.Curr())

	return v
}

func (c *parserContext) parseUint8(name string) uint8 {
	v, err := readUint8(c.args.Text())
	if err != nil {
		c.failCurr(errors.Wrapf(err, "failed to parse uint8: %s", name))
	}

	c.nums = append(c.nums, c.args.Curr())

	return v
}

func (c *parserContext) parseInt8(name string) int8 {
	v, err := readInt8(c.args.Text())
	if err != nil {
		c.failCurr(errors.Wrapf(err, "failed to parse int8: %s", name))
	}

	c.nums = append(c.nums, c.args.Curr())

	return v
}

func (c *parserContext) parseJsonRef(name string) JSONRefType {
	v, isconst, err := readJsonRefField(c.args.Text())
	if err != nil {
		c.failCurr(errors.Wrapf(err, "failed to parse JSON ref field: %s", name))
	}

	if isconst {
		c.strs = append(c.strs, c.args.Curr())
	} else {
		c.nums = append(c.nums, c.args.Curr())
	}

	return v
}

func (c *parserContext) parseTxnField(name string) TxnField {
	// TODO report semantics
	v, isconst, err := readTxnField(c.args.Text())

	if err != nil {
		c.failCurr(errors.Wrapf(err, "failed to parse txn field: %s", name))
	}

	if isconst {
		c.strs = append(c.strs, c.args.Curr())
	} else {
		c.nums = append(c.nums, c.args.Curr())
	}

	return v
}

func (c *parserContext) parseBytes(name string) []byte {
	arg := c.args.Curr().String()

	if strings.HasPrefix(arg, "base32(") || strings.HasPrefix(arg, "b32(") {
		close := strings.IndexRune(arg, ')')
		if close == -1 {
			c.failCurr(errors.New("byte base32 arg lacks close paren"))
		}

		open := strings.IndexRune(arg, '(')
		val, err := base32DecodeAnyPadding(arg[open+1 : close])
		if err != nil {
			c.failCurr(err)
		}

		c.strs = append(c.strs, c.args.Curr())
		return val
	}

	if strings.HasPrefix(arg, "base64(") || strings.HasPrefix(arg, "b64(") {
		close := strings.IndexRune(arg, ')')
		if close == -1 {
			c.failCurr(errors.New("byte base64 arg lacks close paren"))
		}

		open := strings.IndexRune(arg, '(')
		val, err := base64.StdEncoding.DecodeString(arg[open+1 : close])
		if err != nil {
			c.failCurr(err)
		}
		c.strs = append(c.strs, c.args.Curr())
		return val
	}

	if strings.HasPrefix(arg, "0x") {
		val, err := hex.DecodeString(arg[2:])
		if err != nil {
			c.failCurr(err)
		}
		c.strs = append(c.strs, c.args.Curr())
		return val
	}

	if arg == "base32" || arg == "b32" {
		c.keys = append(c.keys, c.args.Curr())

		l := c.mustRead("literal")
		val, err := base32DecodeAnyPadding(l)
		if err != nil {
			c.failCurr(err)
		}

		c.strs = append(c.strs, c.args.Curr())

		return val
	}

	if arg == "base64" || arg == "b64" {
		c.keys = append(c.keys, c.args.Curr())

		l := c.mustRead("literal")
		val, err := base64.StdEncoding.DecodeString(l)
		if err != nil {
			c.failCurr(err)
		}

		c.strs = append(c.strs, c.args.Curr())

		return val
	}

	if len(arg) > 1 && arg[0] == '"' && arg[len(arg)-1] == '"' {
		val, err := parseStringLiteral(arg)
		if err != nil {
			c.failCurr(err)
		}
		c.strs = append(c.strs, c.args.Curr())
		return val
	}

	c.failCurr(fmt.Errorf("byte arg did not parse: %v", arg))

	return nil
}

func (c *parserContext) parseEcdsaCurveIndex(name string) EcdsaCurve {
	// TODO report semantics
	v, err := readEcdsaCurveIndex(c.args.Text())
	if err != nil {
		c.failCurr(errors.Wrapf(err, "failed to parse ESCDS curve index: %s", name))
	}
	return v
}

func (c *parserContext) mustReadBytes(name string) []byte {
	c.mustReadArg(name)
	return c.parseBytes(name)
}

func (c *parserContext) mustReadInt(name string) uint64 {
	c.mustReadArg(name)
	return c.parseUint64(name)
}

func (c *parserContext) mustReadUint8(name string) uint8 {
	c.mustReadArg(name)
	return c.parseUint8(name)
}

func (c *parserContext) maybeReadUint8(name string) (uint8, bool) {
	ok := c.maybeReadArg()
	if !ok {
		return 0, false
	}

	return c.parseUint8(name), true
}

func (c *parserContext) mustReadInt8(name string) int8 {
	c.mustReadArg(name)
	return c.parseInt8(name)
}

func (c *parserContext) mustRead(name string) string {
	c.mustReadArg(name)
	return c.args.Text()
}

func (c *parserContext) mustReadTxnField(name string) TxnField {
	c.mustReadArg(name)
	return c.parseTxnField(name)
}

func (c *parserContext) mustReadJsonRef(name string) JSONRefType {
	c.mustReadArg(name)
	return c.parseJsonRef(name)
}

func (c *parserContext) mustReadEcdsaCurveIndex(name string) EcdsaCurve {
	c.mustReadArg(name)
	return c.parseEcdsaCurveIndex(name)
}

type opDocs struct {
	Items map[string]opItem
}

type OpDocContext struct {
	Name    string
	Version uint8
}

func (d opDocs) GetDoc(c OpDocContext) (opItem, bool) {
	item, ok := d.Items[c.Name]
	return item, ok
}

var OpDocs = func() *opDocs {
	d := &opDocs{
		Items: map[string]opItem{},
	}

	for _, info := range opsList {
		c := &docContext{
			version: 1,
			args:    []opItemArg{},
		}

		info.Parse(c)

		doc := opDocByName[info.Name]
		extra := opDocExtras[info.Name]

		full := doc
		if extra != "" {
			if full != "" {
				full += "\n"
			}

			full += extra
		}

		var fullnames []string
		{
			var optional bool
			for _, arg := range c.args {
				name := fmt.Sprintf("%s : %s", arg.Type.String(), arg.Name)
				if arg.Array {
					name += ", ..."
				}

				name = "{" + name + "}"

				if arg.Optional {
					if !optional {
						name = "[" + name
					}
				}

				fullnames = append(fullnames, name)
			}

			if len(fullnames) > 0 {
				fullnames[len(fullnames)-1] += "]"
			}
		}

		var shortnames []string
		{
			var optional bool
			for _, arg := range c.args {
				name := arg.Name
				if arg.Array {
					name += ", ..."
				}

				if arg.Optional {
					if !optional {
						optional = true
						name = "[" + name
					}
				}

				shortnames = append(shortnames, name)
			}

			if len(shortnames) > 0 && optional {
				shortnames[len(shortnames)-1] += "]"
			}
		}

		sigargs := strings.Join(shortnames, " ")
		sigfull := fmt.Sprintf("%s %s", info.Name, strings.Join(fullnames, " "))

		d.Items[info.Name] = opItem{
			Name: info.Name,

			Version: c.version,

			ArgsSig: sigargs,
			FullSig: sigfull,

			Args: c.args,

			Parse: info.Parse,

			Doc:     doc,
			FullDoc: full,
		}
	}

	return d
}()

type parserFunc func(c OpContext)

type asmType int

const (
	asmNone = iota
	asmInt
	asmByte
	asmIntC
	asmArg
	asmAddr
	asmMethod
)

var typeLoads = []interface{}{}
var typeStores = []interface{}{}

var asmByteCBlock = []interface{}{}
var checkByteImmArgs = []interface{}{}
var immBytess = []interface{}{}
var typeTxField = []interface{}{}

func opPragma(c OpContext) uint8 {
	version := c.mustReadPragma("version")
	c.emit(&PragmaExpr{Version: uint8(version)})

	return version
}

func opAddr(c OpContext) {
	value := c.mustReadAddr("address")
	c.emit(&AddrExpr{Address: value})
}

func opByte(c OpContext) {
	value := c.mustReadBytes("value")
	c.emit(&ByteExpr{Value: value})
}

func opInt(c OpContext) {
	value := c.mustReadInt("value")
	c.emit(&IntExpr{Value: value})
}

func opMethod(c OpContext) {
	value := c.mustReadSignature("signature")
	c.emit(&MethodExpr{Signature: value})
}

func opErr(c OpContext) {
	c.emit(Err)
}

func opSHA256(c OpContext) {
	c.minVersion(2)
	c.emit(Sha256)
}

func opKeccak256(c OpContext) {
	c.minVersion(2)
	c.emit(Keccak256)
}

func opSHA512_256(c OpContext) {
	c.minVersion(2)
	c.emit(Sha512256)
}
func opEd25519Verify(c OpContext) {
	// TODO: 1 for sig, 5 for app - needs mode support
	c.emit(ED25519Verify)
}

func opEcdsaVerify(c OpContext) {
	c.minVersion(5)
	curve := c.mustReadEcdsaCurveIndex("curve index")
	c.emit(&EcdsaVerifyExpr{Index: curve})
}
func opEcdsaPkDecompress(c OpContext) {
	c.minVersion(5)
	curve := c.mustReadEcdsaCurveIndex("curve index")
	c.emit(&EcdsaPkDecompressExpr{Index: curve})
}

func opEcdsaPkRecover(c OpContext) {
	c.minVersion(5)
	curve := c.mustReadEcdsaCurveIndex("curve index")
	c.emit(&EcdsaPkRecoverExpr{Index: curve})
}

func opPlus(c OpContext) {
	c.emit(PlusOp)
}
func opMinus(c OpContext) {
	c.emit(MinusOp)
}
func opDiv(c OpContext) {
	c.emit(Div)
}
func opMul(c OpContext) {
	c.emit(Mul)
}
func opLt(c OpContext) {
	c.emit(Lt)
}
func opGt(c OpContext) {
	c.emit(Gt)
}
func opLe(c OpContext) {
	c.emit(Le)
}
func opGe(c OpContext) {
	c.emit(Ge)
}
func opAnd(c OpContext) {
	c.emit(And)
}
func opOr(c OpContext) {
	c.emit(Or)
}
func opEq(c OpContext) {
	c.emit(Eq)
}
func opNeq(c OpContext) {
	c.emit(Neq)
}
func opNot(c OpContext) {
	c.emit(Not)
}
func opLen(c OpContext) {
	c.emit(Len)
}
func opItob(c OpContext) {
	c.emit(Itob)
}
func opBtoi(c OpContext) {
	c.emit(Btoi)
}
func opModulo(c OpContext) {
	c.emit(Modulo)
}
func opBitOr(c OpContext) {
	c.emit(Bitr)
}
func opBitAnd(c OpContext) {
	c.emit(BitAnd)
}
func opBitXor(c OpContext) {
	c.emit(BitXor)
}
func opBitNot(c OpContext) {
	c.emit(BitNot)
}
func opMulw(c OpContext) {
	c.emit(Mulw)
}
func opAddw(c OpContext) {
	c.minVersion(2)
	c.emit(Addw)
}
func opDivModw(c OpContext) {
	c.minVersion(4)
	c.emit(DivModw)
}
func opIntConstBlock(c OpContext) {
	values := c.readUint64Array("value")

	c.emit(&IntcBlockExpr{Values: values})
}

func opIntConstLoad(c OpContext) {
	value := c.mustReadUint8("value")
	c.emit(&IntcExpr{Index: uint8(value)})
}

func opIntConst0(c OpContext) {
	c.emit(Intc0)
}
func opIntConst1(c OpContext) {
	c.emit(Intc1)
}
func opIntConst2(c OpContext) {
	c.emit(Intc2)
}
func opIntConst3(c OpContext) {
	c.emit(Intc3)
}
func opByteConstBlock(c OpContext) {
	values := c.readBytesArray("bytes")

	c.emit(&BytecBlockExpr{Values: values})
}

func opByteConstLoad(c OpContext) {
	value := c.mustReadUint8("index")
	c.emit(&BytecExpr{Index: uint8(value)})
}

func opByteConst0(c OpContext) {
	c.emit(Bytec0)
}
func opByteConst1(c OpContext) {
	c.emit(Bytec1)
}
func opByteConst2(c OpContext) {
	c.emit(Bytec2)
}
func opByteConst3(c OpContext) {
	c.emit(Bytec3)
}
func opArg(c OpContext) {
	value := c.mustReadUint8("index")
	c.emit(&ArgExpr{Index: uint8(value)})
}
func opArg0(c OpContext) {
	c.emit(Arg0)
}
func opArg1(c OpContext) {
	c.emit(Arg1)
}
func opArg2(c OpContext) {
	c.emit(Arg2)
}
func opArg3(c OpContext) {
	c.emit(Arg3)
}
func opTxn(c OpContext) {
	f := c.mustReadTxnField("f")

	i, ok := c.maybeReadUint8("i")
	if ok {
		c.emit(&TxnaExpr{Field: f, Index: i})
	} else {
		c.emit(&TxnExpr{Field: f})
	}
}
func opGlobal(c OpContext) {
	field := c.mustReadGlobalField("field")

	c.emit(&GlobalExpr{Index: field})
}

func opGtxn(c OpContext) {
	t := c.mustReadUint8("t")
	f := c.mustReadTxnField("f")
	i, ok := c.maybeReadUint8("i")
	if ok {
		c.emit(&GtxnaExpr{Group: uint8(t), Field: f, Index: i})
	} else {
		c.emit(&GtxnExpr{Group: uint8(t), Field: f})
	}
}

func opLoad(c OpContext) {
	value := c.mustReadUint8("i")
	c.emit(&LoadExpr{Index: uint8(value)})
}
func opStore(c OpContext) {
	value := c.mustReadUint8("i")
	c.emit(&StoreExpr{Index: uint8(value)})
}

func opTxna(c OpContext) {
	c.minVersion(2)

	f := c.mustReadTxnField("f")
	i := c.mustReadUint8("i")

	c.emit(&TxnaExpr{Field: f, Index: i})
}

func opGtxna(c OpContext) {
	c.minVersion(2)

	t := c.mustReadUint8("t")
	f := c.mustReadTxnField("f")
	i := c.mustReadUint8("i")

	c.emit(&GtxnaExpr{Group: uint8(t), Field: f, Index: uint8(i)})
}
func opGtxns(c OpContext) {
	c.minVersion(3)

	f := c.mustReadTxnField("f")
	i, ok := c.maybeReadUint8("i")

	if ok {
		c.emit(&GtxnsaExpr{Field: f, Index: i})
	} else {
		c.emit(&GtxnsExpr{Field: f})
	}
}

func opGtxnsa(c OpContext) {
	c.minVersion(3)

	f := c.mustReadTxnField("f")
	i := c.mustReadUint8("i")

	c.emit(&GtxnsaExpr{Field: f, Index: uint8(i)})
}

func opGload(c OpContext) {
	c.minVersion(4)

	t := c.mustReadUint8("t")
	value := c.mustReadUint8("i")

	c.emit(&GloadExpr{Group: uint8(t), Index: uint8(value)})
}

func opGloads(c OpContext) {
	c.minVersion(4)

	value := c.mustReadUint8("i")
	c.emit(&GloadsExpr{Index: uint8(value)})
}

func opGaid(c OpContext) {
	c.minVersion(4)

	t := c.mustReadUint8("t")
	c.emit(&GaidExpr{Group: uint8(t)})
}
func opGaids(c OpContext) {
	c.minVersion(4)

	c.emit(Gaids)
}
func opLoads(c OpContext) {
	c.minVersion(5)

	c.emit(Loads)
}
func opStores(c OpContext) {
	c.minVersion(5)
	c.emit(Stores)
}
func opBnz(c OpContext) {
	name := c.mustReadLabel("label")
	c.emit(&BnzExpr{Label: &LabelExpr{Name: name}})
}
func opBz(c OpContext) {
	c.minVersion(2)
	name := c.mustReadLabel("label")
	c.emit(&BzExpr{Label: &LabelExpr{Name: name}})
}
func opB(c OpContext) {
	c.minVersion(2)
	name := c.mustReadLabel("label")
	c.emit(&BExpr{Label: &LabelExpr{Name: name}})
}
func opReturn(c OpContext) {
	c.minVersion(2)
	c.emit(Return)
}
func opAssert(c OpContext) {
	c.minVersion(3)
	c.emit(Assert)
}

func opBury(c OpContext) {
	c.minVersion(8)

	n := c.mustReadUint8("n")
	c.emit(&BuryExpr{Index: n})
}

func opPopN(c OpContext) {
	c.minVersion(8)

	n := c.mustReadUint8("n")
	c.emit(&PopNExpr{Index: n})
}
func opDupN(c OpContext) {
	c.minVersion(8)

	n := c.mustReadUint8("n")
	c.emit(&DupNExpr{Index: n})
}

func opPop(c OpContext) {
	c.emit(Pop)
}
func opDup(c OpContext) {
	c.emit(Dup)
}
func opDup2(c OpContext) {
	c.minVersion(2)
	c.emit(Dup2)
}
func opDig(c OpContext) {
	c.minVersion(3)
	value := c.mustReadUint8("n")
	c.emit(&DigExpr{Index: uint8(value)})
}
func opSwap(c OpContext) {
	c.minVersion(3)
	c.emit(Swap)
}
func opSelect(c OpContext) {
	c.minVersion(3)
	c.emit(Select)
}
func opCover(c OpContext) {
	c.minVersion(5)
	value := c.mustReadUint8("n")
	c.emit(&CoverExpr{Index: uint8(value)})
}
func opUncover(c OpContext) {
	c.minVersion(5)
	value := c.mustReadUint8("index")
	c.emit(&UncoverExpr{Index: uint8(value)})
}
func opConcat(c OpContext) {
	c.minVersion(2)
	c.emit(Concat)
}
func opSubstring(c OpContext) {
	c.minVersion(2)
	start := c.mustReadUint8("s")
	end := c.mustReadUint8("e")
	c.emit(&SubstringExpr{Start: uint8(start), End: uint8(end)})
}
func opSubstring3(c OpContext) {
	c.minVersion(2)
	c.emit(Substring3)
}
func opGetBit(c OpContext) {
	c.minVersion(3)
	c.emit(GetBit)
}
func opSetBit(c OpContext) {
	c.minVersion(3)
	c.emit(SetBit)
}
func opGetByte(c OpContext) {
	c.minVersion(3)
	c.emit(GetByte)
}
func opSetByte(c OpContext) {
	c.minVersion(3)
	c.emit(SetByte)
}

func opReplace(c OpContext) {
	c.minVersion(7)
	s, ok := c.maybeReadUint8("s")
	if !ok {
		c.emit(&Replace3Expr{})
	} else {
		c.emit(&Replace2Expr{Start: s})
	}
}

func opExtract(c OpContext) {
	c.minVersion(5)
	s, ok := c.maybeReadUint8("s")
	if !ok {
		c.emit(&Extract3Expr{})
	} else {
		l := c.mustReadUint8("l")
		c.emit(&ExtractExpr{Start: uint8(s), Length: uint8(l)})
	}
}

func opExtract3(c OpContext) {
	c.minVersion(5)
	c.emit(Extract3)
}
func opExtract16Bits(c OpContext) {
	c.minVersion(5)
	c.emit(Extract16Bits)
}
func opExtract32Bits(c OpContext) {
	c.minVersion(5)
	c.emit(Extract32Bits)
}
func opExtract64Bits(c OpContext) {
	c.minVersion(5)
	c.emit(Extract64Bits)
}
func opReplace2(c OpContext) {
	c.minVersion(7)
	value := c.mustReadUint8("s")
	c.emit(&Replace2Expr{Start: uint8(value)})
}
func opReplace3(c OpContext) {
	c.minVersion(7)
	c.emit(Replace3)
}
func opBase64Decode(c OpContext) {
	c.minVersion(7)
	value := c.mustReadUint8("e")
	c.emit(&Base64DecodeExpr{Index: uint8(value)})
}
func opJSONRef(c OpContext) {
	c.minVersion(7)
	value := c.mustReadJsonRef("r")
	c.emit(&JsonRefExpr{Index: uint8(value)})
}
func opBalance(c OpContext) {
	c.minVersion(2)
	c.emit(Balance)
}
func opAppOptedIn(c OpContext) {
	c.minVersion(2)
	c.emit(AppOptedIn)
}
func opAppLocalGet(c OpContext) {
	c.minVersion(2)
	c.emit(AppLocalGet)
}
func opAppLocalGetEx(c OpContext) {
	c.minVersion(2)
	c.emit(AppLocalGetEx)
}
func opAppGlobalGet(c OpContext) {
	c.minVersion(2)
	c.emit(AppGlobalGet)
}
func opAppGlobalGetEx(c OpContext) {
	c.minVersion(2)
	c.emit(AppGlobalGetEx)
}
func opAppLocalPut(c OpContext) {
	c.minVersion(2)
	c.emit(AppLocalPut)
}
func opAppGlobalPut(c OpContext) {
	c.minVersion(2)
	c.emit(AppGlobalPut)
}
func opAppLocalDel(c OpContext) {
	c.minVersion(2)
	c.emit(AppLocalDel)
}
func opAppGlobalDel(c OpContext) {
	c.minVersion(2)
	c.emit(AppGlobalDel)
}
func opAssetHoldingGet(c OpContext) {
	c.minVersion(2)
	f := c.mustReadAssetHoldingField("f")
	// TODO report semantics

	c.emit(&AssetHoldingGetExpr{Field: f})
}
func opAssetParamsGet(c OpContext) {
	c.minVersion(2)
	field := c.mustReadAssetParamsField("f")
	c.emit(&AssetParamsGetExpr{Field: field})
}
func opAppParamsGet(c OpContext) {
	c.minVersion(5)
	f := c.mustReadAppParamsField("f")
	c.emit(&AppParamsGetExpr{Field: f})
}
func opAcctParamsGet(c OpContext) {
	c.minVersion(6)
	f := c.mustReadAcctParamsField("f")
	c.emit(&AcctParamsGetExpr{Field: f})
}
func opMinBalance(c OpContext) {
	c.minVersion(3)
	c.emit(MinBalanceOp)
}
func opPushBytes(c OpContext) {
	c.minVersion(3)
	value := c.mustReadBytes("value")
	c.emit(&PushBytesExpr{Value: value})
}
func opPushInt(c OpContext) {
	c.minVersion(3)
	value := c.mustReadInt("value")
	c.emit(&PushIntExpr{Value: value})
}
func opPushBytess(c OpContext) {
	c.minVersion(8)
	bss := c.readBytesArray("bytes")
	c.emit(&PushBytessExpr{
		Bytess: bss,
	})
}
func opPushInts(c OpContext) {
	c.minVersion(8)
	iss := c.readUint64Array("bytes")
	c.emit(&PushIntsExpr{
		Ints: iss,
	})
}
func opEd25519VerifyBare(c OpContext) {
	c.minVersion(7)
	c.emit(Ed25519VerifyBare)
}
func opCallSub(c OpContext) {
	c.minVersion(4)
	name := c.mustReadLabel("label")
	c.emit(&CallSubExpr{Label: &LabelExpr{Name: name}})
}
func opRetSub(c OpContext) {
	c.minVersion(4)
	c.emit(RetSub)
}
func opProto(c OpContext) {
	c.minVersion(8)
	a := c.mustReadUint8("a")
	r := c.mustReadUint8("r")

	c.emit(&ProtoExpr{Args: uint8(a), Results: uint8(r)})
}
func opFrameDig(c OpContext) {
	c.minVersion(8)
	value := c.mustReadInt8("index")
	c.emit(&FrameDigExpr{Index: value})
}
func opFrameBury(c OpContext) {
	c.minVersion(8)
	value := c.mustReadInt8("index")
	c.emit(&FrameBuryExpr{Index: value})
}
func opSwitch(c OpContext) {
	c.minVersion(8)
	names := c.readLabelsArray("label")

	var labels []*LabelExpr
	for _, name := range names {
		labels = append(labels, &LabelExpr{Name: name})
	}

	c.emit(&SwitchExpr{Targets: labels})
}
func opMatch(c OpContext) {
	c.minVersion(8)
	names := c.readLabelsArray("label")

	var labels []*LabelExpr
	for _, name := range names {
		labels = append(labels, &LabelExpr{Name: name})
	}

	c.emit(&MatchExpr{Targets: labels})
}
func opShiftLeft(c OpContext) {
	c.minVersion(4)
	c.emit(ShiftLeft)
}
func opShiftRight(c OpContext) {
	c.minVersion(4)
	c.emit(ShiftRight)
}
func opSqrt(c OpContext) {
	c.minVersion(4)
	c.emit(Sqrt)
}
func opBitLen(c OpContext) {
	c.minVersion(4)
	c.emit(BitLen)
}
func opExp(c OpContext) {
	c.minVersion(4)
	c.emit(Exp)
}
func opExpw(c OpContext) {
	c.minVersion(4)
	c.emit(Expw)
}
func opBytesSqrt(c OpContext) {
	c.minVersion(6)
	c.emit(Bsqrt)
}
func opDivw(c OpContext) {
	c.minVersion(6)
	c.emit(Divw)
}
func opSHA3_256(c OpContext) {
	c.minVersion(7)
	c.emit(Sha3256)
}
func opBn256Add(c OpContext) {
	c.minVersion(9)
	c.emit(Bn256Add)
}
func opBn256ScalarMul(c OpContext) {
	c.minVersion(9)
	c.emit(Bn256ScalarMul)
}
func opBn256Pairing(c OpContext) {
	c.minVersion(9)
	c.emit(Bn256Pairing)
}
func opBytesPlus(c OpContext) {
	c.minVersion(4)
	c.emit(BytesPlus)
}
func opBytesMinus(c OpContext) {
	c.minVersion(4)
	c.emit(BytesMinus)
}
func opBytesDiv(c OpContext) {
	c.minVersion(4)
	c.emit(BytesDiv)
}
func opBytesMul(c OpContext) {
	c.minVersion(4)
	c.emit(BytesMul)
}
func opBytesLt(c OpContext) {
	c.minVersion(4)
	c.emit(BytesLt)
}
func opBytesGt(c OpContext) {
	c.minVersion(4)
	c.emit(BytesGt)
}
func opBytesLe(c OpContext) {
	c.minVersion(4)
	c.emit(BytesLe)
}
func opBytesGe(c OpContext) {
	c.minVersion(4)
	c.emit(BytesGe)
}
func opBytesEq(c OpContext) {
	c.minVersion(4)
	c.emit(BytesEq)
}
func opBytesNeq(c OpContext) {
	c.minVersion(4)
	c.emit(BytesNeq)
}
func opBytesModulo(c OpContext) {
	c.minVersion(4)
	c.emit(BytesModulo)
}
func opBytesBitOr(c OpContext) {
	c.minVersion(4)
	c.emit(BytesBitOr)
}
func opBytesBitAnd(c OpContext) {
	c.minVersion(4)
	c.emit(BytesBitAnd)
}
func opBytesBitXor(c OpContext) {
	c.minVersion(4)
	c.emit(BytesBitXor)
}
func opBytesBitNot(c OpContext) {
	c.minVersion(4)
	c.emit(BytesBitNot)
}
func opBytesZero(c OpContext) {
	c.minVersion(4)
	c.emit(BytesZero)
}
func opLog(c OpContext) {
	c.minVersion(5)
	c.emit(Log)
}
func opTxBegin(c OpContext) {
	c.minVersion(5)
	c.emit(ItxnBegin)
}
func opItxnField(c OpContext) {
	c.minVersion(5)
	f := c.mustReadTxnField("f")
	c.emit(&ItxnFieldExpr{Field: f})
}
func opItxnSubmit(c OpContext) {
	c.minVersion(5)
	c.emit(ItxnSubmit)
}
func opItxn(c OpContext) {
	c.minVersion(5)
	f := c.mustReadTxnField("f")
	c.emit(&ItxnExpr{Field: f})
}
func opItxna(c OpContext) {
	c.minVersion(5)
	f := c.mustReadUint8("f")
	i := c.mustReadUint8("i")
	c.emit(&ItxnaExpr{Field: f, Index: i})
}
func opItxnNext(c OpContext) {
	c.minVersion(6)
	c.emit(ItxnNext)
}
func opGitxn(c OpContext) {
	c.minVersion(6)
	t := c.mustReadUint8("t")
	f := c.mustReadTxnField("f")
	c.emit(&GitxnExpr{Index: uint8(t), Field: f})
}
func opGitxna(c OpContext) {
	c.minVersion(6)
	t := c.mustReadInt("t")
	f := c.mustReadTxnField("f")
	i := c.mustReadUint8("i")

	c.emit(&GitxnaExpr{Group: uint8(t), Field: f, Index: uint8(i)})

}
func opBoxCreate(c OpContext) {
	c.minVersion(8)
	c.emit(BoxCreate)
}
func opBoxExtract(c OpContext) {
	c.minVersion(8)
	c.emit(BoxExtract)
}
func opBoxReplace(c OpContext) {
	c.minVersion(8)
	c.emit(BoxReplace)
}
func opBoxDel(c OpContext) {
	c.minVersion(8)
	c.emit(BoxDel)
}
func opBoxLen(c OpContext) {
	c.minVersion(8)
	c.emit(BoxLen)
}
func opBoxGet(c OpContext) {
	c.minVersion(8)
	c.emit(BoxGet)
}
func opBoxPut(c OpContext) {
	c.minVersion(8)
	c.emit(BoxPut)
}
func opTxnas(c OpContext) {
	c.minVersion(5)
	f := c.mustReadTxnField("f")
	c.emit(&TxnasExpr{Field: f})
}
func opGtxnas(c OpContext) {
	c.minVersion(5)
	t := c.mustReadUint8("t")
	f := c.mustReadUint8("f")
	c.emit(&GtxnasExpr{Index: t, Field: f})
}
func opGtxnsas(c OpContext) {
	c.minVersion(5)
	f := c.mustReadTxnField("f")
	c.emit(&GtxnsasExpr{Field: f})
}
func opArgs(c OpContext) {
	c.minVersion(5)
	c.emit(Args)
}
func opGloadss(c OpContext) {
	c.minVersion(6)
	c.emit(Gloadss)
}
func opItxnas(c OpContext) {
	c.minVersion(6)
	f := c.mustReadTxnField("f")
	c.emit(&ItxnasExpr{Field: f})
}
func opGitxnas(c OpContext) {
	c.minVersion(6)
	t := c.mustReadUint8("t")
	f := c.mustReadTxnField("f")
	c.emit(&GitxnasExpr{Index: t, Field: f})
}
func opVrfVerify(c OpContext) {
	c.minVersion(7)
	f := c.mustReadVrfVerifyField("f")
	c.emit(&VrfVerifyExpr{Field: f})
}
func opBlock(c OpContext) {
	c.minVersion(7)
	f := c.mustReadBlockField("f")
	c.emit(&BlockExpr{Field: f})
}

type ProcessResult struct {
	Version     uint8
	Diagnostics []Diagnostic
	Symbols     []Symbol
	SymbolRefs  []Token
	Tokens      []Token
	Listing     Listing
	Lines       [][]Token
	Ops         []Token
	Numbers     []Token
	Strings     []Token
	Keywords    []Token
	Macros      []Token
}

func readTokens(source string) ([]Token, []Diagnostic) {
	s := &Lexer{Source: []byte(source)}

	ts := []Token{}

	for s.Scan() {
		ts = append(ts, s.Curr())
	}

	diags := make([]Diagnostic, len(s.diag))

	for i, diag := range s.diag {
		diags[i] = diag
	}

	return ts, diags
}

func Process(source string) *ProcessResult {
	c := &parserContext{
		ops: []Op{},
	}

	var ts []Token
	ts, c.diag = readTokens(source)

	lines := [][]Token{}

	p := 0
	for i := 0; i < len(ts); i++ {
		t := ts[i]

		j := i + 1
		eol := t.Type() == TokenEol

		if eol || j == len(ts) {
			k := j
			if eol {
				k--
			}

			lines = append(lines, ts[p:k])
			p = j
		}
	}

	for li, l := range lines {
		for i := 0; i < len(l); i++ {
			t := l[i]
			if t.Type() == TokenComment {
				lines[li] = l[:i]
			}
		}
	}

	var lts [][]Token
	var ops []Token
	var syms []Symbol

	version := uint8(1)

	for _, l := range lines {
		c.args = &arguments{ts: l}
		func() {
			defer func() {
				switch v := recover().(type) {
				case recoverable:
					c.emit(Empty) // consider replacing with Raw string expr
				case nil:
				default:
					fmt.Printf("unrecoverable: %v", v)
					panic(v)
				}
			}()

			if !c.args.Scan() {
				c.emit(Empty)
				return
			}

			if strings.HasSuffix(c.args.Text(), ":") {
				name := c.args.Text()
				name = name[:len(name)-1]
				if len(name) == 0 {
					c.failCurr(errors.New("missing label name"))
					return
				}

				t := c.args.Curr()
				syms = append(syms, labelSymbol{
					n: name,
					l: t.l,
					b: t.b, // TODO: what about whitespaces before label name?
					e: t.e,
				})

				c.emit(&LabelExpr{Name: name})
				return
			}

			name := c.args.Text()
			switch c.args.Text() {
			case "":
				c.emit(Empty)
			case "#pragma":
				version = opPragma(c)
			default:
				info, ok := OpDocs.GetDoc(OpDocContext{
					Name:    name,
					Version: version,
				})
				if ok {
					curr := c.args.Curr()
					ops = append(ops, curr)
					if info.Version > version {
						c.diag = append(c.diag, lintError{
							error: errors.Errorf("opcode requires version >= %d (current: %d)", info.Version, version),
							l:     curr.l,
							b:     curr.b,
							e:     curr.e,
							s:     DiagErr,
						})
					}

					info.Parse(c)
				} else {
					c.failCurr(errors.Errorf("unexpected opcode: %s", c.args.Text()))
				}
				return
			}
		}()

		lts = append(lts, c.args.ts)
	}

	l := &Linter{l: c.ops}
	l.Lint()

	for _, le := range l.res {
		var b int
		var e int

		lt := lts[le.Line()]
		if len(lt) > 0 {
			b = lt[0].b
			e = lt[len(lt)-1].e
		}

		c.diag = append(c.diag, lintError{
			error: le,
			l:     le.Line(),
			b:     b,
			e:     e,
			s:     le.Severity(),
		})
	}

	result := &ProcessResult{
		Version:     version,
		Diagnostics: c.diag,
		Symbols:     syms,
		SymbolRefs:  c.refs,
		Tokens:      ts,
		Listing:     c.ops,
		Lines:       lts,
		Ops:         ops,
		Numbers:     c.nums,
		Strings:     c.strs,
		Keywords:    c.keys,
		Macros:      c.mcrs,
	}

	return result
}
