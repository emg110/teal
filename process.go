package teal

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/algorand/go-algorand-sdk/types"
	"github.com/pkg/errors"
)

type NewOpArgType int

const (
	OpArgTypeNone = iota
	OpArgTypeConstInt
	OpArgTypeUint64
	OpArgTypeUint8
	OpArgTypeInt8
	OpArgTypeBytes
	OpArgTypeTxnField
	OpArgTypeItxnField
	OpArgTypeTxnaField
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
	OpArgTypeBase64EncodingField
	OpArgTypeBlockField
	OpArgTypeEcGroupField
)

var OpArgTypes = []NewOpArgType{
	OpArgTypeTxnField,
	OpArgTypeItxnField,
	OpArgTypeTxnaField,
	OpArgTypeJSONRefField,
	OpArgTypeEcdsaCurve,
	OpArgTypeGlobalField,
	OpArgTypeAssetHoldingField,
	OpArgTypeAssetParamsField,
	OpArgTypeAppParamsField,
	OpArgTypeAcctParamsField,
	OpArgTypeVrfStandard,
	OpArgTypeBase64EncodingField,
	OpArgTypeBlockField,
	OpArgTypeEcGroupField,
}

var OpArgVals map[NewOpArgType][]opItemArgVal
var OpValFieldNames map[NewOpArgType]map[int]string

func init() {
	res := map[NewOpArgType][]opItemArgVal{}
	res2 := map[NewOpArgType]map[int]string{}

	for _, t := range OpArgTypes {
		vals := []opItemArgVal{}
		vals2 := map[int]string{}

		switch t {
		case OpArgTypeTxnaField:
			for name, spec := range txnFieldSpecByName {
				if spec.array {
					vals = append(vals, opItemArgVal{
						Value:   uint64(spec.field),
						Name:    name,
						Docs:    spec.Note(),
						Version: spec.version,
					})
					vals2[int(spec.field)] = spec.field.String()
				}
			}
		case OpArgTypeItxnField:
			for name, spec := range txnFieldSpecByName {
				if spec.itxVersion > 0 {
					vals = append(vals, opItemArgVal{
						Value:   uint64(spec.field),
						Name:    name,
						Docs:    spec.Note(),
						Version: spec.itxVersion,
					})
					vals2[int(spec.field)] = spec.field.String()
				}
			}
		case OpArgTypeTxnField:
			for name, spec := range txnFieldSpecByName {
				vals = append(vals, opItemArgVal{
					Value:   uint64(spec.field),
					Name:    name,
					Docs:    spec.Note(),
					Version: spec.version,
				})
				vals2[int(spec.field)] = spec.field.String()
			}
		case OpArgTypeAcctParamsField:
			for name, spec := range acctParamsFieldSpecByName {
				vals = append(vals, opItemArgVal{
					Value:   uint64(spec.field),
					Name:    name,
					Docs:    spec.Note(),
					Version: spec.version,
				})
				vals2[int(spec.field)] = spec.field.String()
			}

		case OpArgTypeAppParamsField:
			for name, spec := range appParamsFieldSpecByName {
				vals = append(vals, opItemArgVal{
					Value:   uint64(spec.field),
					Name:    name,
					Docs:    spec.Note(),
					Version: spec.version,
				})
				vals2[int(spec.field)] = spec.field.String()
			}

		case OpArgTypeAssetHoldingField:
			for name, spec := range assetHoldingFieldSpecByName {
				vals = append(vals, opItemArgVal{
					Value:   uint64(spec.field),
					Name:    name,
					Docs:    spec.Note(),
					Version: spec.version,
				})
				vals2[int(spec.field)] = spec.field.String()
			}

		case OpArgTypeAssetParamsField:
			for name, spec := range assetParamsFieldSpecByName {
				vals = append(vals, opItemArgVal{
					Value:   uint64(spec.field),
					Name:    name,
					Docs:    spec.Note(),
					Version: spec.version,
				})
				vals2[int(spec.field)] = spec.field.String()
			}

		case OpArgTypeEcdsaCurve:
			for name, spec := range ecdsaCurveSpecByName {
				vals = append(vals, opItemArgVal{
					Value:   uint64(spec.field),
					Name:    name,
					Docs:    spec.Note(),
					Version: spec.version,
				})
				vals2[int(spec.field)] = spec.field.String()
			}

		case OpArgTypeGlobalField:
			for name, spec := range globalFieldSpecByName {
				vals = append(vals, opItemArgVal{
					Value:   uint64(spec.field),
					Name:    name,
					Docs:    spec.Note(),
					Version: spec.version,
				})
				vals2[int(spec.field)] = spec.field.String()
			}

		case OpArgTypeJSONRefField:
			for name, spec := range jsonRefSpecByName {
				vals = append(vals, opItemArgVal{
					Value:   uint64(spec.field),
					Name:    name,
					Docs:    spec.Note(),
					Version: spec.version,
				})
				vals2[int(spec.field)] = spec.field.String()
			}

		case OpArgTypeVrfStandard:
			for name, spec := range vrfStandardSpecByName {
				vals = append(vals, opItemArgVal{
					Value:   uint64(spec.field),
					Name:    name,
					Docs:    spec.Note(),
					Version: spec.version,
				})
				vals2[int(spec.field)] = spec.field.String()
			}

		case OpArgTypeBase64EncodingField:
			for name, spec := range base64EncodingSpecByName {
				vals = append(vals, opItemArgVal{
					Value:   uint64(spec.field),
					Name:    name,
					Docs:    spec.Note(),
					Version: spec.version,
				})
				vals2[int(spec.field)] = spec.field.String()
			}

		case OpArgTypeEcGroupField:
			for name, spec := range ecGroupSpecByName {
				vals = append(vals, opItemArgVal{
					Value:   uint64(spec.field),
					Name:    name,
					Docs:    spec.Note(),
					Version: spec.Version(),
				})
				vals2[int(spec.field)] = spec.field.String()
			}

		case OpArgTypeBlockField:
			for name, spec := range blockFieldSpecByName {
				vals = append(vals, opItemArgVal{
					Value:   uint64(spec.field),
					Name:    name,
					Docs:    spec.Note(),
					Version: spec.version,
				})
				vals2[int(spec.field)] = spec.field.String()
			}
		}

		res[t] = vals
		res2[t] = vals2
	}

	OpArgVals = res
	OpValFieldNames = res2
}

func (t NewOpArgType) String() string {
	switch t {
	default:
		return "(none)"
	case OpArgTypeConstInt:
		return "uint64"
	case OpArgTypeUint64:
		return "uint64"
	case OpArgTypeUint8:
		return "uint8"
	case OpArgTypeInt8:
		return "int8"
	case OpArgTypeBytes:
		return "bytes"
	case OpArgTypeTxnaField:
		return "transaction array field index"
	case OpArgTypeTxnField:
		return "transaction field index"
	case OpArgTypeItxnField:
		return "internal transaction field index"
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
	case OpArgTypeBase64EncodingField:
		return "base64 encoding"
	case OpArgTypeBlockField:
		return "block field"
	case OpArgTypeEcGroupField:
		return "EC group field index"
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

	SigVersion uint64
	AppVersion uint64

	Args []opItemArg

	ArgsSig string
	FullSig string

	Parse processFunc

	Doc     string
	FullDoc string
}

type opListItem struct {
	Name  string
	Parse processFunc
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
	{"app_local_get", opAppLocalGet},
	{"app_local_get_ex", opAppLocalGetEx},
	{"app_global_get", opAppGlobalGet},
	{"app_global_get_ex", opAppGlobalGetEx},
	{"app_local_put", opAppLocalPut},
	{"app_global_put", opAppGlobalPut},
	{"app_local_del", opAppLocalDel},
	{"app_global_del", opAppGlobalDel},
	{"asset_holding_get", opAssetHoldingGet},
	{"asset_params_get", opAssetParamsGet},
	{"app_params_get", opAppParamsGet},
	{"acct_params_get", opAcctParamsGet},
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
	{"ec_add", opEcAdd},
	{"ec_scalar_mul", opEcScalarMul},
	{"ec_pairing_check", opEcPairingCheck},
	{"ec_multi_exp", opEcMultiExp},
	{"ec_subgroup_check", opEcSubgroupCheck},
	{"ec_map_to", opEcMapTo},
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

type ProcessContext interface {
	comment(text string)
	emit(op Op)

	minVersion(v uint64)
	modeMinVersion(mode ProgramMode, v uint64)

	mustReadEcGroup(name string) EcGroup
	mustReadBase64Encoding(name string) Base64Encoding
	mustReadPragma(name string) uint64
	mustReadAddr(name string) string
	mustReadSignature(name string) string
	mustReadTxnField(name string) TxnField
	mustReadItxnField(name string) TxnField
	mustReadTxnaField(name string) TxnField
	mustReadJsonRef(name string) JSONRefType
	maybeReadUint8(name string) (uint8, bool)
	mustReadEcdsaCurveIndex(name string) EcdsaCurve
	readUint64Array(name string) []uint64
	mustReadUint8(name string) uint8
	mustReadInt8(name string) int8
	readBytesArray(name string) [][]byte
	mustReadGlobalField(name string) GlobalField
	mustReadInt(name string) uint64
	mustReadConstInt(name string) uint64
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
	args       []opItemArg
	sigVersion uint64
	appVersion uint64
	optional   bool
}

func (c *docContext) arg(a opItemArg) {
	if c.optional {
		a.Optional = true
	} else if a.Optional {
		c.optional = true
	}

	c.args = append(c.args, a)
}

func (c *docContext) comment(text string) {
}

func (c *docContext) emit(op Op) {
}

func (c *docContext) minVersion(v uint64) {
	c.appVersion = v
	c.sigVersion = v
}

func (c *docContext) modeMinVersion(mode ProgramMode, v uint64) {
	switch mode {
	case ModeSig:
		c.sigVersion = v
	case ModeApp:
		c.appVersion = v
	}
}

func (c *docContext) mustReadPragma(name string) (v uint64) {
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

func (c *docContext) mustReadEcGroup(name string) (v EcGroup) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeEcGroupField,
	})

	return
}

func (c *docContext) mustReadBase64Encoding(name string) (v Base64Encoding) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeBase64EncodingField,
	})

	return
}

func (c *docContext) mustReadJsonRef(name string) (v JSONRefType) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeJSONRefField,
	})

	return
}

type fieldContext int

const (
	txnFieldContext = iota
	txnaFieldContext
	itxnFieldContext
)

func (c *docContext) mustReadItxnField(name string) (f TxnField) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeItxnField,
	})

	return
}

func (c *docContext) mustReadTxnaField(name string) (f TxnField) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeTxnaField,
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
		Type: OpArgTypeEcdsaCurve,
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

func (c *docContext) mustReadConstInt(name string) (v uint64) {
	c.arg(opItemArg{
		Name: name,
		Type: OpArgTypeConstInt,
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
		Type: OpArgTypeBlockField,
	})

	return
}

type ProgramMode int

const (
	ModeNone = iota
	ModeApp
	ModeSig
)

func (m ProgramMode) String() string {
	switch m {
	case ModeApp:
		return "application"
	case ModeSig:
		return "logicsig"
	default:
		return "(unknown)"
	}
}

type Call struct {
	Op   *CallSubExpr
	Line int
}

type parserContext struct {
	mode    ProgramMode
	version uint64

	ops  []Op
	args *arguments
	diag []Diagnostic

	nums []Token
	strs []Token
	keys []Token
	mcrs []Token
	refs []Token

	vtok   *Token
	protos map[string]*ProtoExpr
	refc   map[string]int

	// current state
	line     int
	label    *LabelExpr
	comments []string
}

func (c *parserContext) comment(text string) {
	c.comments = append(c.comments, text)
}

func (c *parserContext) emit(op Op) {
	c.ops = append(c.ops, op)

	switch op := op.(type) {
	case usesLabels:
		for _, lbl := range op.Labels() {
			c.refc[lbl.Name]++
		}
	}

	switch op := op.(type) {
	case *ProtoExpr:
		if c.label != nil {
			c.protos[c.label.Name] = op
		}
	case *LabelExpr:
		c.label = op
	}
}

func (c *parserContext) minVersion(v uint64) {

}

func (c *parserContext) modeMinVersion(mode ProgramMode, v uint64) {
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

func (c *parserContext) mustReadPragma(argName string) uint64 {
	var version uint64
	c.mcrs = append(c.mcrs, c.args.Curr())

	name := c.mustRead("name")
	switch name {
	case "version":
		c.mcrs = append(c.mcrs, c.args.Curr())
		v := c.mustReadInt("version value")

		tok := c.args.Curr()
		c.vtok = &tok

		if v < 1 {
			c.failCurr(errors.New("version must be at least 1"))
		}

		version = v
	default:
		c.failCurr(errors.Errorf("unexpected #pragma: %s", c.args.Text()))
	}

	c.version = version

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

	f, isconst, err := readAcctParams(c.version, c.args.Text())

	if err != nil {
		c.failCurr(err)
	}

	if isconst {
		c.strs = append(c.strs, c.args.Curr())
	} else {
		c.nums = append(c.nums, c.args.Curr())
	}

	return f
}

func (c *parserContext) mustReadAssetHoldingField(name string) AssetHoldingField {
	c.mustReadArg(name)

	f, isconst, err := readAssetHoldingField(c.version, c.args.Text())
	if err != nil {
		c.failCurr(err)
	}

	if isconst {
		c.strs = append(c.strs, c.args.Curr())
	} else {
		c.nums = append(c.nums, c.args.Curr())
	}

	return f
}

func (c *parserContext) mustReadAssetParamsField(name string) AssetParamsField {
	c.mustReadArg(name)

	f, isconst, err := readAssetParamsField(c.version, c.args.Text())
	if err != nil {
		c.failCurr(err)
	}

	if isconst {
		c.strs = append(c.strs, c.args.Curr())
	} else {
		c.nums = append(c.nums, c.args.Curr())
	}

	return f
}

func (c *parserContext) mustReadBlockField(name string) BlockField {
	c.mustReadArg(name)

	f, isconst, err := readBlockField(c.version, c.args.Text())
	if err != nil {
		c.failCurr(err)
	}

	if isconst {
		c.strs = append(c.strs, c.args.Curr())
	} else {
		c.nums = append(c.nums, c.args.Curr())
	}
	return f
}

func (c *parserContext) mustReadGlobalField(name string) GlobalField {
	c.mustReadArg(name)

	f, isconst, err := readGlobalField(c.version, c.args.Text())
	if err != nil {
		c.failCurr(err)
	}

	if isconst {
		c.strs = append(c.strs, c.args.Curr())
	} else {
		c.nums = append(c.nums, c.args.Curr())
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

	f, isconst, err := readVrfVerifyField(c.version, c.args.Text())
	if err != nil {
		c.failCurr(err)
	}

	if isconst {
		c.strs = append(c.strs, c.args.Curr())
	} else {
		c.nums = append(c.nums, c.args.Curr())
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

	f, isconst, err := readAppParamsField(c.version, c.args.Text())
	if err != nil {
		c.failCurr(err)
	}

	if isconst {
		c.strs = append(c.strs, c.args.Curr())
	} else {
		c.nums = append(c.nums, c.args.Curr())
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

func (c *parserContext) parseConstUint64(name string) uint64 {
	v, err := readConstInt(c.args)
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

func (c *parserContext) parseEcGroup(name string) EcGroup {
	v, isconst, err := readEcGroupField(c.version, c.args.Text())
	if err != nil {
		c.failCurr(errors.Wrapf(err, "failed to parse EC group field: %s", name))
	}

	if isconst {
		c.strs = append(c.strs, c.args.Curr())
	} else {
		c.nums = append(c.nums, c.args.Curr())
	}

	return v
}

func (c *parserContext) parseBase64Encoding(name string) Base64Encoding {
	v, isconst, err := readBase64EncodingField(c.version, c.args.Text())
	if err != nil {
		c.failCurr(errors.Wrapf(err, "failed to parse base64 encoding field: %s", name))
	}

	if isconst {
		c.strs = append(c.strs, c.args.Curr())
	} else {
		c.nums = append(c.nums, c.args.Curr())
	}

	return v
}

func (c *parserContext) parseJsonRef(name string) JSONRefType {
	v, isconst, err := readJsonRefField(c.version, c.args.Text())
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

func (c *parserContext) parseTxnField(tc fieldContext, name string) TxnField {
	v, isconst, err := readTxnField(tc, c.version, c.args.Text())

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
	v, isconst, err := readEcdsaCurveIndex(c.version, c.args.Text())

	if err != nil {
		c.failCurr(errors.Wrapf(err, "failed to parse ESCDS curve index: %s", name))
	}

	if isconst {
		c.strs = append(c.strs, c.args.Curr())
	} else {
		c.nums = append(c.nums, c.args.Curr())
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

func (c *parserContext) mustReadConstInt(name string) uint64 {
	c.mustReadArg(name)
	return c.parseConstUint64(name)
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
	return c.parseTxnField(txnFieldContext, name)
}

func (c *parserContext) mustReadItxnField(name string) TxnField {
	c.mustReadArg(name)
	return c.parseTxnField(itxnFieldContext, name)
}

func (c *parserContext) mustReadTxnaField(name string) TxnField {
	c.mustReadArg(name)
	return c.parseTxnField(txnaFieldContext, name)
}

func (c *parserContext) mustReadEcGroup(name string) EcGroup {
	c.mustReadArg(name)
	return c.parseEcGroup(name)
}

func (c *parserContext) mustReadBase64Encoding(name string) Base64Encoding {
	c.mustReadArg(name)
	return c.parseBase64Encoding(name)
}

func (c *parserContext) mustReadJsonRef(name string) JSONRefType {
	c.mustReadArg(name)
	return c.parseJsonRef(name)
}

func (c *parserContext) mustReadEcdsaCurveIndex(name string) EcdsaCurve {
	c.mustReadArg(name)
	return c.parseEcdsaCurveIndex(name)
}

type opItems struct {
	Items map[string]opItem
}

type OpContext struct {
	Name    string
	Version uint64
}

func (d opItems) Get(c OpContext) (opItem, bool) {
	item, ok := d.Items[c.Name]
	return item, ok
}

var Ops = func() *opItems {
	d := &opItems{
		Items: map[string]opItem{},
	}

	for _, info := range opsList {
		c := &docContext{
			appVersion: 1,
			sigVersion: 1,
			args:       []opItemArg{},
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

			if optional && len(fullnames) > 0 {
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

			if optional && len(shortnames) > 0 {
				shortnames[len(shortnames)-1] += "]"
			}
		}

		sigargs := strings.Join(shortnames, " ")
		sigfull := fmt.Sprintf("%s %s", info.Name, strings.Join(fullnames, " "))

		d.Items[info.Name] = opItem{
			Name: info.Name,

			SigVersion: c.sigVersion,
			AppVersion: c.appVersion,

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

type processFunc func(c ProcessContext)

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

func opPragma(c ProcessContext) {
	version := c.mustReadPragma("version")
	c.emit(&PragmaExpr{Version: uint8(version)})
}

func opAddr(c ProcessContext) {
	value := c.mustReadAddr("address")
	c.emit(&AddrExpr{Address: value})
}

func opByte(c ProcessContext) {
	value := c.mustReadBytes("value")
	c.emit(&ByteExpr{Value: value})
}

func opInt(c ProcessContext) {
	value := c.mustReadConstInt("value")
	c.emit(&IntExpr{Value: value})
}

func opMethod(c ProcessContext) {
	value := c.mustReadSignature("signature")
	c.emit(&MethodExpr{Signature: value})
}

func opErr(c ProcessContext) {
	c.emit(Err)
}

func opSHA256(c ProcessContext) {
	c.minVersion(1)
	c.emit(Sha256)
}

func opKeccak256(c ProcessContext) {
	c.minVersion(1)
	c.emit(Keccak256)
}

func opSHA512_256(c ProcessContext) {
	c.minVersion(1)
	c.emit(Sha512256)
}
func opEd25519Verify(c ProcessContext) {
	c.modeMinVersion(ModeSig, 1)
	c.modeMinVersion(ModeApp, 5)
	c.emit(ED25519Verify)
}

func opEcdsaVerify(c ProcessContext) {
	c.minVersion(5)
	curve := c.mustReadEcdsaCurveIndex("curve index")
	c.emit(&EcdsaVerifyExpr{Index: curve})
}
func opEcdsaPkDecompress(c ProcessContext) {
	c.minVersion(5)
	curve := c.mustReadEcdsaCurveIndex("curve index")
	c.emit(&EcdsaPkDecompressExpr{Index: curve})
}

func opEcdsaPkRecover(c ProcessContext) {
	c.minVersion(5)
	curve := c.mustReadEcdsaCurveIndex("curve index")
	c.emit(&EcdsaPkRecoverExpr{Index: curve})
}

func opPlus(c ProcessContext) {
	c.emit(PlusOp)
}
func opMinus(c ProcessContext) {
	c.emit(MinusOp)
}
func opDiv(c ProcessContext) {
	c.emit(Div)
}
func opMul(c ProcessContext) {
	c.emit(Mul)
}
func opLt(c ProcessContext) {
	c.emit(Lt)
}
func opGt(c ProcessContext) {
	c.emit(Gt)
}
func opLe(c ProcessContext) {
	c.emit(Le)
}
func opGe(c ProcessContext) {
	c.emit(Ge)
}
func opAnd(c ProcessContext) {
	c.emit(And)
}
func opOr(c ProcessContext) {
	c.emit(Or)
}
func opEq(c ProcessContext) {
	c.emit(Eq)
}
func opNeq(c ProcessContext) {
	c.emit(Neq)
}
func opNot(c ProcessContext) {
	c.emit(Not)
}
func opLen(c ProcessContext) {
	c.emit(Len)
}
func opItob(c ProcessContext) {
	c.emit(Itob)
}
func opBtoi(c ProcessContext) {
	c.emit(Btoi)
}
func opModulo(c ProcessContext) {
	c.emit(Modulo)
}
func opBitOr(c ProcessContext) {
	c.emit(BitOr)
}
func opBitAnd(c ProcessContext) {
	c.emit(BitAnd)
}
func opBitXor(c ProcessContext) {
	c.emit(BitXor)
}
func opBitNot(c ProcessContext) {
	c.emit(BitNot)
}
func opMulw(c ProcessContext) {
	c.emit(Mulw)
}
func opAddw(c ProcessContext) {
	c.minVersion(2)
	c.emit(Addw)
}
func opDivModw(c ProcessContext) {
	c.minVersion(4)
	c.emit(DivModw)
}
func opIntConstBlock(c ProcessContext) {
	values := c.readUint64Array("value")

	c.emit(&IntcBlockExpr{Values: values})
}

func opIntConstLoad(c ProcessContext) {
	value := c.mustReadUint8("value")
	c.emit(&IntcExpr{Index: uint8(value)})
}

func opIntConst0(c ProcessContext) {
	c.emit(Intc0)
}
func opIntConst1(c ProcessContext) {
	c.emit(Intc1)
}
func opIntConst2(c ProcessContext) {
	c.emit(Intc2)
}
func opIntConst3(c ProcessContext) {
	c.emit(Intc3)
}
func opByteConstBlock(c ProcessContext) {
	values := c.readBytesArray("bytes")

	c.emit(&BytecBlockExpr{Values: values})
}

func opByteConstLoad(c ProcessContext) {
	value := c.mustReadUint8("index")
	c.emit(&BytecExpr{Index: uint8(value)})
}

func opByteConst0(c ProcessContext) {
	c.emit(Bytec0)
}
func opByteConst1(c ProcessContext) {
	c.emit(Bytec1)
}
func opByteConst2(c ProcessContext) {
	c.emit(Bytec2)
}
func opByteConst3(c ProcessContext) {
	c.emit(Bytec3)
}
func opArg(c ProcessContext) {
	c.modeMinVersion(ModeSig, 1)
	c.modeMinVersion(ModeApp, 0)
	value := c.mustReadUint8("index")
	c.emit(&ArgExpr{Index: uint8(value)})
}
func opArg0(c ProcessContext) {
	c.modeMinVersion(ModeSig, 1)
	c.modeMinVersion(ModeApp, 0)
	c.emit(Arg0)
}
func opArg1(c ProcessContext) {
	c.modeMinVersion(ModeSig, 1)
	c.modeMinVersion(ModeApp, 0)
	c.emit(Arg1)
}
func opArg2(c ProcessContext) {
	c.modeMinVersion(ModeSig, 1)
	c.modeMinVersion(ModeApp, 0)
	c.emit(Arg2)
}
func opArg3(c ProcessContext) {
	c.modeMinVersion(ModeSig, 1)
	c.modeMinVersion(ModeApp, 0)
	c.emit(Arg3)
}
func opTxn(c ProcessContext) {
	f := c.mustReadTxnField("f")

	i, ok := c.maybeReadUint8("i")
	if ok {
		c.emit(&TxnaExpr{Field: f, Index: i})
	} else {
		c.emit(&TxnExpr{Field: f})
	}
}
func opGlobal(c ProcessContext) {
	field := c.mustReadGlobalField("field")

	c.emit(&GlobalExpr{Field: field})
}

func opGtxn(c ProcessContext) {
	t := c.mustReadUint8("t")
	f := c.mustReadTxnField("f")
	i, ok := c.maybeReadUint8("i")
	if ok {
		c.emit(&GtxnaExpr{Group: uint8(t), Field: f, Index: i})
	} else {
		c.emit(&GtxnExpr{Group: uint8(t), Field: f})
	}
}

func opLoad(c ProcessContext) {
	value := c.mustReadUint8("i")
	c.emit(&LoadExpr{Index: uint8(value)})
}
func opStore(c ProcessContext) {
	value := c.mustReadUint8("i")
	c.emit(&StoreExpr{Index: uint8(value)})
}

func opTxna(c ProcessContext) {
	c.minVersion(2)

	f := c.mustReadTxnaField("f")
	i := c.mustReadUint8("i")

	c.emit(&TxnaExpr{Field: f, Index: i})
}

func opGtxna(c ProcessContext) {
	c.minVersion(2)

	t := c.mustReadUint8("t")
	f := c.mustReadTxnaField("f")
	i := c.mustReadUint8("i")

	c.emit(&GtxnaExpr{Group: uint8(t), Field: f, Index: uint8(i)})
}
func opGtxns(c ProcessContext) {
	c.minVersion(3)

	f := c.mustReadTxnField("f")
	i, ok := c.maybeReadUint8("i")

	if ok {
		c.emit(&GtxnsaExpr{Field: f, Index: i})
	} else {
		c.emit(&GtxnsExpr{Field: f})
	}
}

func opGtxnsa(c ProcessContext) {
	c.minVersion(3)

	f := c.mustReadTxnaField("f")
	i := c.mustReadUint8("i")

	c.emit(&GtxnsaExpr{Field: f, Index: uint8(i)})
}

func opGload(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 4)

	t := c.mustReadUint8("t")
	value := c.mustReadUint8("i")

	c.emit(&GloadExpr{Group: uint8(t), Index: uint8(value)})
}

func opGloads(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 4)
	value := c.mustReadUint8("i")
	c.emit(&GloadsExpr{Index: uint8(value)})
}

func opGaid(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 4)

	t := c.mustReadUint8("t")
	c.emit(&GaidExpr{Group: uint8(t)})
}
func opGaids(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 4)

	c.emit(Gaids)
}
func opLoads(c ProcessContext) {
	c.minVersion(5)

	c.emit(Loads)
}
func opStores(c ProcessContext) {
	c.minVersion(5)
	c.emit(Stores)
}
func opBnz(c ProcessContext) {
	name := c.mustReadLabel("label")
	c.emit(&BnzExpr{Label: &LabelExpr{Name: name}})
}
func opBz(c ProcessContext) {
	c.minVersion(2)
	name := c.mustReadLabel("label")
	c.emit(&BzExpr{Label: &LabelExpr{Name: name}})
}
func opB(c ProcessContext) {
	c.minVersion(2)
	name := c.mustReadLabel("label")
	c.emit(&BExpr{Label: &LabelExpr{Name: name}})
}
func opReturn(c ProcessContext) {
	c.minVersion(2)
	c.emit(Return)
}
func opAssert(c ProcessContext) {
	c.minVersion(3)
	c.emit(Assert)
}

func opBury(c ProcessContext) {
	c.minVersion(8)

	n := c.mustReadUint8("n")
	c.emit(&BuryExpr{Depth: n})
}

func opPopN(c ProcessContext) {
	c.minVersion(8)

	n := c.mustReadUint8("n")
	c.emit(&PopNExpr{Depth: n})
}
func opDupN(c ProcessContext) {
	c.minVersion(8)

	n := c.mustReadUint8("n")
	c.emit(&DupNExpr{Count: n})
}

func opPop(c ProcessContext) {
	c.emit(Pop)
}
func opDup(c ProcessContext) {
	c.emit(Dup)
}
func opDup2(c ProcessContext) {
	c.minVersion(2)
	c.emit(Dup2)
}
func opDig(c ProcessContext) {
	c.minVersion(3)
	value := c.mustReadUint8("n")
	c.emit(&DigExpr{Index: uint8(value)})
}
func opSwap(c ProcessContext) {
	c.minVersion(3)
	c.emit(Swap)
}
func opSelect(c ProcessContext) {
	c.minVersion(3)
	c.emit(Select)
}
func opCover(c ProcessContext) {
	c.minVersion(5)
	value := c.mustReadUint8("n")
	c.emit(&CoverExpr{Depth: uint8(value)})
}
func opUncover(c ProcessContext) {
	c.minVersion(5)
	value := c.mustReadUint8("index")
	c.emit(&UncoverExpr{Depth: uint8(value)})
}
func opConcat(c ProcessContext) {
	c.minVersion(2)
	c.emit(Concat)
}
func opSubstring(c ProcessContext) {
	c.minVersion(2)
	start := c.mustReadUint8("s")
	end := c.mustReadUint8("e")
	c.emit(&SubstringExpr{Start: uint8(start), End: uint8(end)})
}
func opSubstring3(c ProcessContext) {
	c.minVersion(2)
	c.emit(Substring3)
}
func opGetBit(c ProcessContext) {
	c.minVersion(3)
	c.emit(GetBit)
}
func opSetBit(c ProcessContext) {
	c.minVersion(3)
	c.emit(SetBit)
}
func opGetByte(c ProcessContext) {
	c.minVersion(3)
	c.emit(GetByte)
}
func opSetByte(c ProcessContext) {
	c.minVersion(3)
	c.emit(SetByte)
}

func opReplace(c ProcessContext) {
	c.minVersion(7)
	s, ok := c.maybeReadUint8("s")
	if !ok {
		c.emit(&Replace3Expr{})
	} else {
		c.emit(&Replace2Expr{Start: s})
	}
}

func opExtract(c ProcessContext) {
	c.minVersion(5)
	s, ok := c.maybeReadUint8("s")
	if !ok {
		c.emit(&Extract3Expr{})
	} else {
		l := c.mustReadUint8("l")
		c.emit(&ExtractExpr{Start: uint8(s), Length: uint8(l)})
	}
}

func opExtract3(c ProcessContext) {
	c.minVersion(5)
	c.emit(Extract3)
}
func opExtract16Bits(c ProcessContext) {
	c.minVersion(5)
	c.emit(Extract16Bits)
}
func opExtract32Bits(c ProcessContext) {
	c.minVersion(5)
	c.emit(Extract32Bits)
}
func opExtract64Bits(c ProcessContext) {
	c.minVersion(5)
	c.emit(Extract64Bits)
}
func opReplace2(c ProcessContext) {
	c.minVersion(7)
	value := c.mustReadUint8("s")
	c.emit(&Replace2Expr{Start: uint8(value)})
}
func opReplace3(c ProcessContext) {
	c.minVersion(7)
	c.emit(Replace3)
}
func opBase64Decode(c ProcessContext) {
	c.minVersion(7)
	value := c.mustReadBase64Encoding("e")
	c.emit(&Base64DecodeExpr{Index: uint8(value)})
}
func opJSONRef(c ProcessContext) {
	c.minVersion(7)
	value := c.mustReadJsonRef("r")
	c.emit(&JsonRefExpr{Index: uint8(value)})
}
func opBalance(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 2)

	c.emit(Balance)
}
func opAppOptedIn(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 2)
	c.emit(AppOptedIn)
}
func opAppLocalGet(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 2)
	c.emit(AppLocalGet)
}
func opAppLocalGetEx(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 2)
	c.emit(AppLocalGetEx)
}
func opAppGlobalGet(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 2)
	c.emit(AppGlobalGet)
}
func opAppGlobalGetEx(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 2)
	c.emit(AppGlobalGetEx)
}
func opAppLocalPut(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 2)
	c.emit(AppLocalPut)
}
func opAppGlobalPut(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 2)
	c.emit(AppGlobalPut)
}
func opAppLocalDel(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 2)
	c.emit(AppLocalDel)
}
func opAppGlobalDel(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 2)
	c.emit(AppGlobalDel)
}
func opAssetHoldingGet(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 2)
	f := c.mustReadAssetHoldingField("f")
	// TODO report semantics

	c.emit(&AssetHoldingGetExpr{Field: f})
}
func opAssetParamsGet(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 2)
	field := c.mustReadAssetParamsField("f")
	c.emit(&AssetParamsGetExpr{Field: field})
}
func opAppParamsGet(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 5)
	f := c.mustReadAppParamsField("f")
	c.emit(&AppParamsGetExpr{Field: f})
}
func opAcctParamsGet(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 6)
	f := c.mustReadAcctParamsField("f")
	c.emit(&AcctParamsGetExpr{Field: f})
}
func opMinBalance(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 3)
	c.emit(MinBalanceOp)
}
func opPushBytes(c ProcessContext) {
	c.minVersion(3)
	value := c.mustReadBytes("value")
	c.emit(&PushBytesExpr{Value: value})
}
func opPushInt(c ProcessContext) {
	c.minVersion(3)
	value := c.mustReadInt("value")
	c.emit(&PushIntExpr{Value: value})
}
func opPushBytess(c ProcessContext) {
	c.minVersion(8)
	bss := c.readBytesArray("value")
	c.emit(&PushBytessExpr{
		Bytess: bss,
	})
}
func opPushInts(c ProcessContext) {
	c.minVersion(8)
	iss := c.readUint64Array("value")
	c.emit(&PushIntsExpr{
		Ints: iss,
	})
}
func opEd25519VerifyBare(c ProcessContext) {
	c.minVersion(7)
	c.emit(Ed25519VerifyBare)
}
func opCallSub(c ProcessContext) {
	c.minVersion(4)
	name := c.mustReadLabel("label")
	c.emit(&CallSubExpr{Label: &LabelExpr{Name: name}})
}
func opRetSub(c ProcessContext) {
	c.minVersion(4)
	c.emit(RetSub)
}
func opProto(c ProcessContext) {
	c.minVersion(8)
	a := c.mustReadUint8("a")
	r := c.mustReadUint8("r")

	c.emit(&ProtoExpr{Args: uint8(a), Results: uint8(r)})
}
func opFrameDig(c ProcessContext) {
	c.minVersion(8)
	value := c.mustReadInt8("index")
	c.emit(&FrameDigExpr{Index: value})
}
func opFrameBury(c ProcessContext) {
	c.minVersion(8)
	value := c.mustReadInt8("index")
	c.emit(&FrameBuryExpr{Index: value})
}
func opSwitch(c ProcessContext) {
	c.minVersion(8)
	names := c.readLabelsArray("label")

	var labels []*LabelExpr
	for _, name := range names {
		labels = append(labels, &LabelExpr{Name: name})
	}

	c.emit(&SwitchExpr{Targets: labels})
}
func opMatch(c ProcessContext) {
	c.minVersion(8)
	names := c.readLabelsArray("label")

	var labels []*LabelExpr
	for _, name := range names {
		labels = append(labels, &LabelExpr{Name: name})
	}

	c.emit(&MatchExpr{Targets: labels})
}
func opShiftLeft(c ProcessContext) {
	c.minVersion(4)
	c.emit(ShiftLeft)
}
func opShiftRight(c ProcessContext) {
	c.minVersion(4)
	c.emit(ShiftRight)
}
func opSqrt(c ProcessContext) {
	c.minVersion(4)
	c.emit(Sqrt)
}
func opBitLen(c ProcessContext) {
	c.minVersion(4)
	c.emit(BitLen)
}
func opExp(c ProcessContext) {
	c.minVersion(4)
	c.emit(Exp)
}
func opExpw(c ProcessContext) {
	c.minVersion(4)
	c.emit(Expw)
}
func opBytesSqrt(c ProcessContext) {
	c.minVersion(6)
	c.emit(Bsqrt)
}
func opDivw(c ProcessContext) {
	c.minVersion(6)
	c.emit(Divw)
}
func opSHA3_256(c ProcessContext) {
	c.minVersion(7)
	c.emit(Sha3256)
}
func opEcAdd(c ProcessContext) {
	c.minVersion(9)
	v := c.mustReadEcGroup("curve")
	c.emit(&EcAddExpr{Group: v})
}
func opEcScalarMul(c ProcessContext) {
	c.minVersion(9)
	v := c.mustReadEcGroup("curve")
	c.emit(&EcScalarMul{Group: v})
}
func opEcPairingCheck(c ProcessContext) {
	c.minVersion(9)
	v := c.mustReadEcGroup("curve")
	c.emit(&EcPairingCheckExpr{Group: v})
}

func opEcMultiExp(c ProcessContext) {
	c.minVersion(9)
	v := c.mustReadEcGroup("curve")
	c.emit(&EcMultiExpExpr{Group: v})
}

func opEcSubgroupCheck(c ProcessContext) {
	c.minVersion(9)
	v := c.mustReadEcGroup("curve")
	c.emit(&EcSubgroupCheckExpr{Group: v})
}

func opEcMapTo(c ProcessContext) {
	c.minVersion(9)
	v := c.mustReadEcGroup("curve")
	c.emit(&EcMapToExpr{Group: v})
}

func opBytesPlus(c ProcessContext) {
	c.minVersion(4)
	c.emit(BytesPlus)
}
func opBytesMinus(c ProcessContext) {
	c.minVersion(4)
	c.emit(BytesMinus)
}
func opBytesDiv(c ProcessContext) {
	c.minVersion(4)
	c.emit(BytesDiv)
}
func opBytesMul(c ProcessContext) {
	c.minVersion(4)
	c.emit(BytesMul)
}
func opBytesLt(c ProcessContext) {
	c.minVersion(4)
	c.emit(BytesLt)
}
func opBytesGt(c ProcessContext) {
	c.minVersion(4)
	c.emit(BytesGt)
}
func opBytesLe(c ProcessContext) {
	c.minVersion(4)
	c.emit(BytesLe)
}
func opBytesGe(c ProcessContext) {
	c.minVersion(4)
	c.emit(BytesGe)
}
func opBytesEq(c ProcessContext) {
	c.minVersion(4)
	c.emit(BytesEq)
}
func opBytesNeq(c ProcessContext) {
	c.minVersion(4)
	c.emit(BytesNeq)
}
func opBytesModulo(c ProcessContext) {
	c.minVersion(4)
	c.emit(BytesModulo)
}
func opBytesBitOr(c ProcessContext) {
	c.minVersion(4)
	c.emit(BytesBitOr)
}
func opBytesBitAnd(c ProcessContext) {
	c.minVersion(4)
	c.emit(BytesBitAnd)
}
func opBytesBitXor(c ProcessContext) {
	c.minVersion(4)
	c.emit(BytesBitXor)
}
func opBytesBitNot(c ProcessContext) {
	c.minVersion(4)
	c.emit(BytesBitNot)
}
func opBytesZero(c ProcessContext) {
	c.minVersion(4)
	c.emit(BytesZero)
}
func opLog(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 5)
	c.emit(Log)
}
func opTxBegin(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 5)
	c.emit(ItxnBegin)
}
func opItxnField(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 5)
	f := c.mustReadItxnField("f")
	c.emit(&ItxnFieldExpr{Field: f})
}
func opItxnSubmit(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 5)
	c.emit(ItxnSubmit)
}
func opItxn(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 5)
	f := c.mustReadTxnField("f")
	c.emit(&ItxnExpr{Field: f})
}
func opItxna(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 5)
	f := c.mustReadTxnaField("f")
	i := c.mustReadUint8("i")
	c.emit(&ItxnaExpr{Field: f, Index: i})
}
func opItxnNext(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 6)
	c.emit(ItxnNext)
}
func opGitxn(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 6)
	t := c.mustReadUint8("t")
	f := c.mustReadTxnField("f")
	c.emit(&GitxnExpr{Index: uint8(t), Field: f})
}
func opGitxna(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 6)
	t := c.mustReadInt("t")
	f := c.mustReadTxnaField("f")
	i := c.mustReadUint8("i")

	c.emit(&GitxnaExpr{Group: uint8(t), Field: f, Index: uint8(i)})

}
func opBoxCreate(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 8)
	c.emit(BoxCreate)
}
func opBoxExtract(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 8)
	c.emit(BoxExtract)
}
func opBoxReplace(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 8)
	c.emit(BoxReplace)
}
func opBoxDel(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 8)
	c.emit(BoxDel)
}
func opBoxLen(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 8)
	c.emit(BoxLen)
}
func opBoxGet(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 8)
	c.emit(BoxGet)
}
func opBoxPut(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 8)
	c.emit(BoxPut)
}
func opTxnas(c ProcessContext) {
	c.minVersion(5)
	f := c.mustReadTxnaField("f")
	c.emit(&TxnasExpr{Field: f})
}
func opGtxnas(c ProcessContext) {
	c.minVersion(5)
	t := c.mustReadUint8("t")
	f := c.mustReadTxnaField("f")
	c.emit(&GtxnasExpr{Index: t, Field: f})
}
func opGtxnsas(c ProcessContext) {
	c.minVersion(5)
	f := c.mustReadTxnaField("f")
	c.emit(&GtxnsasExpr{Field: f})
}
func opArgs(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 5)

	c.emit(Args)
}
func opGloadss(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 6)
	c.emit(Gloadss)
}
func opItxnas(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 6)
	f := c.mustReadTxnaField("f")
	c.emit(&ItxnasExpr{Field: f})
}
func opGitxnas(c ProcessContext) {
	c.modeMinVersion(ModeSig, 0)
	c.modeMinVersion(ModeApp, 6)
	t := c.mustReadUint8("t")
	f := c.mustReadTxnaField("f")
	c.emit(&GitxnasExpr{Index: t, Field: f})
}
func opVrfVerify(c ProcessContext) {
	c.minVersion(7)
	f := c.mustReadVrfVerifyField("f")
	c.emit(&VrfVerifyExpr{Field: f})
}
func opBlock(c ProcessContext) {
	c.minVersion(7)
	f := c.mustReadBlockField("f")
	c.emit(&BlockExpr{Field: f})
}

type Line []Token

func (ln Line) Begin() int {
	switch len(ln) {
	case 0:
		return 0
	default:
		return ln[0].b
	}
}

func (ln Line) End() int {
	switch len(ln) {
	case 0:
		return 0
	default:
		return ln[len(ln)-1].e
	}
}

func (ln Line) ImmAt(pos int) (Token, int, bool) {
	for idx, tok := range ln {
		if idx > 0 {
			if pos >= tok.Begin() && pos <= tok.End() {
				return tok, idx - 1, true
			}
		}
	}

	return Token{}, 0, false
}

type RequiredVersion struct {
	Line  int
	Begin int
	End   int

	Version uint64
}

func (v RequiredVersion) StartLine() int {
	return v.Line
}

func (v RequiredVersion) EndLine() int {
	return v.Line
}

func (v RequiredVersion) StartCharacter() int {
	return v.Begin
}

func (v RequiredVersion) EndCharacter() int {
	return v.End
}

type ProcessResult struct {
	Mode ProgramMode

	Version      uint64
	VersionToken *Token
	Versions     []RequiredVersion

	Diagnostics []Diagnostic

	MissRefs   []Token
	Symbols    []Symbol
	SymbolRefs []Token

	Tokens  []Token
	Listing Listing
	Lines   []Line

	Ops []Token

	Numbers  []Token
	Strings  []Token
	Keywords []Token
	Macros   []Token

	Redundants []RedundantLine

	RefCounts map[string]int
}

func (r ProcessResult) SymbolsForRefWithin(rg Range) []Symbol {
	var res []Symbol

	refs := r.SymbolRefsWithin(rg)

	if len(refs) > 0 {
		ref := refs[0]

		for _, sym := range r.Symbols {
			if sym.Name() == ref.String() {
				res = append(res, sym)
			}
		}
	}

	return res
}

func (r ProcessResult) SymbolsWithin(rg Range) []Symbol {
	var res []Symbol

	for _, sym := range r.Symbols {
		if Overlaps(rg, sym) {
			res = append(res, sym)
		}
	}

	return res
}

func (r ProcessResult) SymbolRefsWithin(rg Range) []Token {
	var res []Token

	for _, ref := range r.SymbolRefs {
		if Overlaps(rg, ref) {
			res = append(res, ref)
		}
	}

	return res
}

func (r ProcessResult) getOp(name string) (opItem, bool) {
	return Ops.Get(OpContext{
		Name:    name,
		Version: r.Version,
	})
}

type opItemArgVal struct {
	NoValue bool
	Value   uint64

	Name      string
	Docs      string
	Signature string
	Version   uint64
}

type NamedInlayHint struct {
	T    Token
	Name string
}

type DecodedInlayHint struct {
	T     Token
	Value string
}

type SignatureInlayHint struct {
	Line  int
	Start int
	End   int

	Signature string
}

type InlayHints struct {
	Named   []NamedInlayHint
	Decoded []DecodedInlayHint
}

type InlayHint struct {
	Line      int
	Character int
	Label     string
}

func (r ProcessResult) ArgFieldName(t NewOpArgType, v int) string {
	ns, ok := OpValFieldNames[t]
	if !ok {
		return ""
	}

	n, ok := ns[v]
	if !ok {
		return ""
	}

	return n
}

func (r ProcessResult) InlayHints(rg Range) InlayHints {
	var ihs InlayHints

	for li := rg.StartLine(); li <= rg.EndLine(); li++ {
		if li >= len(r.Lines) {
			continue
		}

		ln := r.Lines[li]

		var ok bool
		var spec opItem

		if len(ln) > 0 {
			spec, ok = r.getOp(ln[0].String())
		}

		for i, tok := range ln {
			if !Overlaps(tok, rg) {
				continue
			}

			func() {
				if !ok {
					return
				}

				if i == 0 {
					return
				}

				idx := i - 1

				if idx >= len(spec.Args) {
					return
				}

				iv, err := strconv.Atoi(tok.v)
				if err != nil {
					return
				}

				arg := spec.Args[idx]

				name := r.ArgFieldName(arg.Type, iv)

				if name == "" {
					return
				}

				ihs.Named = append(ihs.Named, NamedInlayHint{
					T:    tok,
					Name: name,
				})
			}()

			func() {
				if tok.Type() != TokenValue {
					return
				}

				s := tok.String()
				if !strings.HasPrefix(s, "0x") {
					return
				}

				bs, err := hex.DecodeString(s[2:])
				if err != nil {
					return
				}

				ds := string(bs)

				if func() bool {
					for _, r := range ds {
						if r > unicode.MaxASCII || !unicode.IsPrint(r) {
							return false
						}
					}
					return true
				}() {
					ihs.Decoded = append(ihs.Decoded, DecodedInlayHint{
						T:     tok,
						Value: ds,
					})
				}
			}()
		}
	}

	return ihs
}

func (r ProcessResult) ArgVals(arg opItemArg) []opItemArgVal {
	var res []opItemArgVal

	switch arg.Type {
	case OpArgTypeLabel:
		for _, sym := range r.Symbols {
			res = append(res, opItemArgVal{
				NoValue:   true,
				Name:      sym.Name(),
				Docs:      sym.Docs(),
				Signature: sym.Signature(),
			})
		}
	case OpArgTypeConstInt:
		for name, value := range txnTypeMap {
			if value != 0 {
				res = append(res, opItemArgVal{
					Name:  name,
					Value: value,
				})
			}
		}
		for name, value := range onCompletionMap {
			res = append(res, opItemArgVal{
				Name:  name,
				Value: value,
			})
		}
	default:
		for _, v := range OpArgVals[arg.Type] {
			if r.Version >= v.Version {
				res = append(res, v)
			}
		}
	}

	return res
}

func (r ProcessResult) SymByName(name string) []Symbol {
	var res []Symbol
	for _, sym := range r.Symbols {
		if sym.Name() == name {
			res = append(res, sym)
		}
	}
	return res
}

func (r ProcessResult) SymRefByName(name string) []Token {
	var res []Token
	for _, sym := range r.SymbolRefs {
		if sym.String() == name {
			res = append(res, sym)
		}
	}
	return res
}

func (r ProcessResult) SymOrRefAt(rg Range) string {
	for _, sym := range r.Symbols {
		if Overlaps(rg, sym) {
			return sym.Name()
		}
	}

	for _, ref := range r.SymbolRefs {
		if Overlaps(rg, ref) {
			return ref.String()
		}
	}

	return ""
}

func (r ProcessResult) ArgAt(l int, ch int) (opItemArg, int, bool) {
	var res opItemArg

	if l >= len(r.Lines) {
		return res, -1, false
	}

	ln := r.Lines[l]

	curr := len(ln) - 1

	_, idx, ok := ln.ImmAt(ch)
	if ok {
		curr = idx
	}

	op := ln[0]
	info, ok := Ops.Get(OpContext{
		Name:    op.String(),
		Version: r.Version,
	})

	if !ok {
		return res, -1, false
	}

	if len(info.Args) > 0 {
		if curr >= len(info.Args) {
			if info.Args[len(info.Args)-1].Array {
				curr = len(info.Args) - 1
			}
		}
	}

	if curr >= len(info.Args) {
		return res, curr, false
	}

	arg := info.Args[curr]

	return arg, curr, true
}

func (r ProcessResult) ArgValsAt(l int, ch int) []opItemArgVal {
	var res []opItemArgVal

	arg, _, ok := r.ArgAt(l, ch)
	if !ok {
		return res
	}

	res = r.ArgVals(arg)

	return res
}

func (r ProcessResult) DocAt(l int, ch int) string {
	if l >= len(r.Lines) {
		return ""
	}

	ln := r.Lines[l]

	for i, t := range ln {
		if t.b > ch || t.End() < ch {
			continue
		}

		if i == 0 {
			info, ok := r.getOp(t.String())
			if ok {
				return info.FullDoc
			}
		} else {
			tok, idx, ok := ln.ImmAt(ch)
			if ok {
				info, ok := r.getOp(ln[0].String())
				if ok {
					if len(info.Args) > 0 && idx >= len(info.Args) && info.Args[len(info.Args)-1].Array {
						idx = len(info.Args) - 1
					}

					if idx < len(info.Args) {
						arg := info.Args[idx]
						switch arg.Type {
						case OpArgTypeTxnaField:
							spec, ok := txnFieldSpecByName[tok.String()]
							if ok {
								return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
							}
							v, err := strconv.Atoi(tok.String())
							if err == nil {
								spec, ok = txnFieldSpecByField(TxnField(v))
								if ok {
									return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
								}
							}

						case OpArgTypeItxnField:
							spec, ok := txnFieldSpecByName[tok.String()]
							if ok {
								return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
							}
							v, err := strconv.Atoi(tok.String())
							if err == nil {
								spec, ok = txnFieldSpecByField(TxnField(v))
								if ok {
									return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
								}
							}

						case OpArgTypeTxnField:
							spec, ok := txnFieldSpecByName[tok.String()]
							if ok {
								return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
							}
							v, err := strconv.Atoi(tok.String())
							if err == nil {
								spec, ok = txnFieldSpecByField(TxnField(v))
								if ok {
									return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
								}
							}
						case OpArgTypeAcctParamsField:
							spec, ok := acctParamsFieldSpecByName[tok.String()]
							if ok {
								return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
							}
							v, err := strconv.Atoi(tok.String())
							if err == nil {
								spec, ok = acctParamsFieldSpecByField(AcctParamsField(v))
								if ok {
									return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
								}
							}
						case OpArgTypeAppParamsField:
							spec, ok := appParamsFieldSpecByName[tok.String()]
							if ok {
								return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
							}
							v, err := strconv.Atoi(tok.String())
							if err == nil {
								spec, ok = appParamsFieldSpecByField(AppParamsField(v))
								if ok {
									return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
								}
							}

						case OpArgTypeAssetHoldingField:
							spec, ok := assetHoldingFieldSpecByName[tok.String()]
							if ok {
								return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
							}
							v, err := strconv.Atoi(tok.String())
							if err == nil {
								spec, ok = assetHoldingFieldSpecByField(AssetHoldingField(v))
								if ok {
									return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
								}
							}

						case OpArgTypeAssetParamsField:
							spec, ok := assetParamsFieldSpecByName[tok.String()]
							if ok {
								return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
							}
							v, err := strconv.Atoi(tok.String())
							if err == nil {
								spec, ok = assetParamsFieldSpecByField(AssetParamsField(v))
								if ok {
									return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
								}
							}

						case OpArgTypeBase64EncodingField:
							spec, ok := base64EncodingSpecByName[tok.String()]
							if ok {
								return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
							}
							v, err := strconv.Atoi(tok.String())
							if err == nil {
								spec, ok = base64EncodingSpecByField(Base64Encoding(v))
								if ok {
									return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
								}
							}

						case OpArgTypeBlockField:
							spec, ok := blockFieldSpecByName[tok.String()]
							if ok {
								return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
							}
							v, err := strconv.Atoi(tok.String())
							if err == nil {
								spec, ok = blockFieldSpecByField(BlockField(v))
								if ok {
									return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
								}
							}

						case OpArgTypeEcdsaCurve:
							spec, ok := ecdsaCurveSpecByName[tok.String()]
							if ok {
								return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
							}
							v, err := strconv.Atoi(tok.String())
							if err == nil {
								spec, ok = ecdsaCurveSpecByField(EcdsaCurve(v))
								if ok {
									return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
								}
							}

						case OpArgTypeGlobalField:
							spec, ok := globalFieldSpecByName[tok.String()]
							if ok {
								return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
							}
							v, err := strconv.Atoi(tok.String())
							if err == nil {
								spec, ok = globalFieldSpecByField(GlobalField(v))
								if ok {
									return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
								}
							}

						case OpArgTypeJSONRefField:
							spec, ok := jsonRefSpecByName[tok.String()]
							if ok {
								return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
							}
							v, err := strconv.Atoi(tok.String())
							if err == nil {
								spec, ok = jsonRefSpecByField(JSONRefType(v))
								if ok {
									return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
								}
							}

						case OpArgTypeVrfStandard:
							spec, ok := vrfStandardSpecByName[tok.String()]
							if ok {
								return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
							}
							v, err := strconv.Atoi(tok.String())
							if err == nil {
								spec, ok = vrfStandardSpecByField(VrfStandard(v))
								if ok {
									return fmt.Sprintf("%s = %d\r\n%s", spec.field.String(), spec.field, spec.Note())
								}
							}
						}
					}
				}
			}
		}

		break
	}

	return ""
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
		version: 1,
		ops:     []Op{},
		mode:    ModeApp,
		protos:  map[string]*ProtoExpr{},
		refc:    map[string]int{},
	}

	var ts []Token
	ts, c.diag = readTokens(source)

	lines := []Line{}

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
		for i := 1; i < len(l); i++ {
			t := l[i]
			if t.Type() == TokenComment {
				lines[li] = l[:i]
			}
		}
	}

	var lts []Line
	var ops []Token
	var lsyms []*labelSymbol
	var vers []RequiredVersion

	for line, l := range lines {
		c.line = line
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

			if c.args.Curr().Type() == TokenComment {
				if strings.TrimSpace(c.args.Curr().String()) == "#pragma mode logicsig" {
					c.mode = ModeSig
				} else {
					c.comment(c.args.Curr().String())
				}
				c.emit(Empty)
				return
			} else if strings.HasSuffix(c.args.Text(), ":") {
				name := c.args.Text()
				name = name[:len(name)-1]
				if len(name) == 0 {
					c.failCurr(errors.New("missing label name"))
					return
				}

				t := c.args.Curr()
				lsyms = append(lsyms, &labelSymbol{
					n:    name,
					l:    t.l,
					b:    t.b, // TODO: what about whitespaces before label name?
					e:    t.e,
					docs: strings.Join(c.comments, "\n"),
				})

				c.emit(&LabelExpr{Name: name})
				c.comments = nil

				return
			}

			name := c.args.Text()
			switch c.args.Text() {
			case "":
				c.emit(Empty)
			case "#pragma":
				opPragma(c)
			default:
				info, ok := Ops.Get(OpContext{
					Name:    name,
					Version: c.version,
				})
				if ok {
					curr := c.args.Curr()
					ops = append(ops, curr)

					var min uint64
					switch c.mode {
					case ModeApp:
						min = info.AppVersion
					case ModeSig:
						min = info.SigVersion
					}

					if min == 0 {
						c.diag = append(c.diag, lintError{
							error: errors.Errorf("opcode not available in the current mode: %s", c.mode),
							l:     curr.l,
							b:     curr.b,
							e:     curr.e,
							s:     DiagErr,
						})
					}

					if min > c.version {
						c.diag = append(c.diag, lintError{
							error: errors.Errorf("opcode requires version >= %d (current: %d)", min, c.version),
							l:     curr.l,
							b:     curr.b,
							e:     curr.e,
							s:     DiagErr,
						})

						var ln Line = c.args.ts

						vers = append(vers, RequiredVersion{
							Line:    curr.l,
							Begin:   ln.Begin(),
							End:     ln.End(),
							Version: min,
						})
					}

					info.Parse(c)
					if c.args.i < len(c.args.ts) {
						if c.args.Scan() {
							c.failCurr(errors.Errorf("too many values"))
						}
					}
				} else {
					c.failCurr(errors.Errorf("unknown opcode: %s", c.args.Text()))
				}
				return
			}
		}()

		lts = append(lts, c.args.ts)
	}

	l := &Linter{l: c.ops}
	l.Lint()

	for _, le := range l.errs {
		ln := lts[le.Line()]
		c.diag = append(c.diag, lintError{
			error: le,
			l:     le.Line(),
			b:     ln.Begin(),
			e:     ln.End(),
			s:     le.Severity(),
		})
	}

	symm := map[string]bool{}
	for _, sym := range lsyms {
		symm[sym.Name()] = true
	}

	var mrefs []Token
	for _, ref := range c.refs {
		if _, ok := symm[ref.String()]; !ok {
			mrefs = append(mrefs, ref)
		}
	}

	for i := 0; i < len(lsyms); i++ {
		sym := lsyms[i]
		proto := c.protos[sym.Name()]
		if proto != nil {
			sym.sig = fmt.Sprintf("in: %d, out: %d", proto.Args, proto.Results)
			lsyms[i] = sym
		}
	}

	syms := make([]Symbol, len(lsyms))

	for i := 0; i < len(lsyms); i++ {
		syms[i] = lsyms[i]
	}

	result := &ProcessResult{
		Mode:         c.mode,
		Version:      c.version,
		VersionToken: c.vtok,
		Diagnostics:  c.diag,
		MissRefs:     mrefs,
		Symbols:      syms,
		SymbolRefs:   c.refs,
		Tokens:       ts,
		Lines:        lts,
		Listing:      c.ops,
		Ops:          ops,
		Numbers:      c.nums,
		Strings:      c.strs,
		Keywords:     c.keys,
		Macros:       c.mcrs,
		Redundants:   l.reds,
		Versions:     vers,
		RefCounts:    c.refc,
	}

	return result
}
