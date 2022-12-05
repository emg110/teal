package lsp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"strconv"

	"github.com/dragmz/teal"
	"github.com/pkg/errors"
)

type lsp struct {
	docs     map[string]string
	shutdown bool
	exit     bool

	tp *textproto.Reader
	w  *bufio.Writer

	debug *bufio.Writer
}

type LspOption func(l *lsp) error

func WithDebug(w io.Writer) LspOption {
	return func(l *lsp) error {
		l.debug = bufio.NewWriter(w)
		return nil
	}
}

func New(r io.Reader, w io.Writer, opts ...LspOption) (*lsp, error) {
	l := &lsp{
		tp:   textproto.NewReader(bufio.NewReader(r)),
		w:    bufio.NewWriter(w),
		docs: map[string]string{},
	}

	for _, opt := range opts {
		err := opt(l)
		if err != nil {
			return nil, errors.Wrap(err, "failed to set lsp option")
		}
	}

	return l, nil
}

type jsonRpcHeader struct {
	JsonRpc string      `json:"jsonrpc"`
	Id      interface{} `json:"id"`
	Method  string      `json:"method"`
}

type jsonRpcResponse struct {
	JsonRpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   interface{} `json:"error,omitempty"`
	Id      interface{} `json:"id"`
}

type lspRequest[T any] struct {
	Params T `json:"params"`
}

type lspDiagnosticProvider struct {
	InterFileDependencies bool `json:"interFileDependencies"`
	WorkspaceDiagnostics  bool `json:"workspaceDiagnostics"`
}

type lspServerCapabilities struct {
	TextDocumentSync   int                    `json:"textDocumentSync,omitempty"`
	DiagnosticProvider *lspDiagnosticProvider `json:"diagnosticProvider,omitempty"`
}

type lspInitializeResult struct {
	Capabilities *lspServerCapabilities `json:"capabilities"`
}

type lspInitializeClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type lspInitializeRequestParams struct {
	ProcessId  int                      `json:"id"`
	ClientInfo *lspInitializeClientInfo `json:"clientInfo"`
}

type lspDidOpenTextDocument struct {
	Uri     string `json:"uri"`
	Version int    `json:"version"`
	Text    string `json:"text"`
}

type lspDidOpenParams struct {
	TextDocument *lspDidOpenTextDocument `json:"textDocument"`
}

type lspDidChangeTextDocument struct {
	Uri     string `json:"uri"`
	Version int    `json:"version"`
}

type lspContentChange struct {
	Text string `json:"text"`
}

type lspDidChangeParams struct {
	TextDocument   *lspDidChangeTextDocument `json:"textDocument"`
	ContentChanges []*lspContentChange       `json:"contentChanges"`
}

type lspDidSaveTextDocument struct {
	Uri string `json:"uri"`
}

type lspDidSaveParams struct {
	TextDocument *lspDidSaveTextDocument `json:"textDocument"`
}

type lspDidChange lspRequest[*lspDidChangeParams]
type lspDidOpen lspRequest[*lspDidOpenParams]
type lspDidSave lspRequest[*lspDidSaveParams]

type lspDiagnosticRequestTextDocument struct {
	Uri string `json:"uri"`
}

type lspDiagnosticRequestParams struct {
	TextDocument *lspDiagnosticRequestTextDocument `json:"textDocument"`
}

type lspDiagnosticRequest lspRequest[*lspDiagnosticRequestParams]

type lspDidCloseTextDocument struct {
	Uri string `json:"uri"`
}

type lspDidCloseRequestParams struct {
	TextDocument *lspDidCloseTextDocument `json:"textDocument"`
}

type lspDidCloseRequest lspRequest[*lspDidCloseRequestParams]

type lspPosition struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type lspRange struct {
	Start lspPosition `json:"start"`
	End   lspPosition `json:"end"`
}

type lspDiagnostic struct {
	Range    lspRange `json:"range"`
	Severity *int     `json:"severity,omitempty"`
	Message  string   `json:"message"`
}

type lspPublishDiagnostic struct {
	Uri         string          `json:"uri"`
	Diagnostics []lspDiagnostic `json:"diagnostics"`
}

type lspNotification struct {
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

func read[T any](b []byte) (T, error) {
	var v T

	err := json.Unmarshal(b, &v)
	if err != nil {
		return v, err
	}

	return v, nil
}

func (l *lsp) doDiagnostic(uri string) error {
	text := l.docs[uri]

	diags := []lspDiagnostic{}

	p, errs := teal.Parse(text)
	if len(errs) > 0 {
		for _, pe := range errs {
			sev := new(int)
			*sev = 1

			diags = append(diags, lspDiagnostic{
				Range: lspRange{
					Start: lspPosition{
						Line: pe.Line() - 1,
					},
					End: lspPosition{
						Line: pe.Line() - 1,
					},
				},
				Severity: sev,
				Message:  fmt.Sprintf("%s", pe),
			})
		}
	} else {
		ls := teal.Compile(p)
		for _, le := range ls.Lint() {
			sev := new(int)
			*sev = 2

			diags = append(diags, lspDiagnostic{
				Range: lspRange{
					Start: lspPosition{
						Line: le.Line() - 1,
					},
					End: lspPosition{
						Line: le.Line() - 1,
					},
				},
				Severity: sev,
				Message:  fmt.Sprintf("%s", le),
			})
		}
	}

	return l.write(lspNotification{
		JsonRpc: "2.0",
		Method:  "textDocument/publishDiagnostics",
		Params: &lspPublishDiagnostic{
			Uri:         uri,
			Diagnostics: diags,
		},
	})
}

func (l *lsp) handle(h jsonRpcHeader, b []byte) error {
	if l.shutdown {
		return errors.New("server is shut down")
	}

	switch h.Method {
	case "initialized":
	case "$/cancelRequest":

	case "exit":
		l.exit = true
	case "shutdown":
		l.shutdown = true

	case "textDocument/didSave":
		req, err := read[lspDidSave](b)
		if err != nil {
			return err
		}

		return l.doDiagnostic(req.Params.TextDocument.Uri)

	case "textDocument/didClose":
		req, err := read[lspDidCloseRequest](b)
		if err != nil {
			return err
		}

		delete(l.docs, req.Params.TextDocument.Uri)

	case "textDocument/diagnostic":
		req, err := read[lspDiagnosticRequest](b)
		if err != nil {
			return err
		}

		return l.doDiagnostic(req.Params.TextDocument.Uri)

	case "textDocument/didOpen":
		req, err := read[lspDidOpen](b)
		if err != nil {
			return err
		}

		l.docs[req.Params.TextDocument.Uri] = req.Params.TextDocument.Text
		return l.doDiagnostic(req.Params.TextDocument.Uri)

	case "textDocument/didChange":
		req, err := read[lspDidChange](b)
		if err != nil {
			return err
		}

		for _, ch := range req.Params.ContentChanges {
			l.docs[req.Params.TextDocument.Uri] = ch.Text
		}

		return l.doDiagnostic(req.Params.TextDocument.Uri)

	case "initialize":
		return l.write(&jsonRpcResponse{
			JsonRpc: "2.0",
			Id:      h.Id,
			Result: &lspInitializeResult{
				Capabilities: &lspServerCapabilities{
					TextDocumentSync: 1,
				},
			},
		})
	default:
		return errors.New("unknown method")
	}

	return nil
}

func (l *lsp) write(v interface{}) error {
	rb, err := json.Marshal(v)
	if err != nil {
		return errors.Wrap(err, "failed to marshal response")
	}

	l.trace(fmt.Sprintf("OUT: %s", string(rb)))

	h := http.Header{}
	h.Set("Content-Length", strconv.Itoa(len(rb)))

	err = h.Write(l.w)
	if err != nil {
		return errors.Wrap(err, "failed to write response headers")
	}

	_, err = l.w.Write([]byte("\r\n"))
	if err != nil {
		return errors.Wrap(err, "failed to write")
	}

	_, err = l.w.Write(rb)
	if err != nil {
		return errors.Wrap(err, "failed to write response body")
	}

	err = l.w.Flush()
	if err != nil {
		return errors.Wrap(err, "failed to flush")
	}

	return nil
}

func (l *lsp) trace(s string) {
	if l.debug == nil {
		return
	}

	l.debug.WriteString(s)
	l.debug.WriteString("\n")

	l.debug.Flush()
}

func (l *lsp) Run() error {
	l.trace("TEAL LSP running..")
	defer func() {
		l.trace("TEAL LSP exited.")
	}()

	for !l.exit {
		err := func() error {
			mh, err := l.tp.ReadMIMEHeader()
			if err != nil {
				return errors.Wrap(err, "failed to read request headers")
			}

			h := http.Header(mh)

			length, err := strconv.Atoi(h.Get("Content-Length"))
			if err != nil {
				return errors.Wrap(err, "failed to parse content length")
			}

			data := make([]byte, length)
			_, err = io.ReadFull(l.tp.R, data)
			if err != nil {
				return errors.Wrap(err, "failed to read content body")
			}

			l.trace(fmt.Sprintf("IN: %s", string(data)))

			var jh jsonRpcHeader
			err = json.Unmarshal(data, &jh)
			if err != nil {
				return errors.Wrap(err, "failed to unmarshal json rpc header")
			}

			err = l.handle(jh, data)
			if err != nil {
				return errors.Wrap(err, "failed to handle request")
			}

			return nil
		}()

		if err != nil {
			l.trace(fmt.Sprintf("ERR: %s", err))
		}

	}

	return nil
}
