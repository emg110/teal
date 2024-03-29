package teal

import (
	"fmt"
)

type RedundantLine interface {
	Line() int
	String() string
}

type RedundantLabelLine struct {
	line int
	name string
}

func (l RedundantLabelLine) Line() int {
	return l.line
}

func (l RedundantLabelLine) String() string {
	return fmt.Sprintf("Remove label '%s'", l.name)
}

type RedundantBLine struct {
	line int
}

func (l RedundantBLine) Line() int {
	return l.line
}

func (l RedundantBLine) String() string {
	return "Remove b call"
}

type LineError interface {
	error
	Line() int
	Severity() DiagnosticSeverity
}

type Linter struct {
	l Listing

	errs []LineError
	reds []RedundantLine
}

type DuplicateLabelError struct {
	l    int
	name string
}

func (e DuplicateLabelError) Line() int {
	return e.l
}

func (e DuplicateLabelError) Error() string {
	return fmt.Sprintf("duplicate label: \"%s\"", e.name)
}

func (e DuplicateLabelError) Severity() DiagnosticSeverity {
	return DiagErr
}

type UnusedLabelError struct {
	l    int
	name string
}

func (e *UnusedLabelError) Line() int {
	return e.l
}

func (e UnusedLabelError) Error() string {
	return fmt.Sprintf("unused label: \"%s\"", e.name)
}

func (e UnusedLabelError) Severity() DiagnosticSeverity {
	return DiagWarn
}

type UnreachableCodeError struct {
	l int
}

func (e UnreachableCodeError) Line() int {
	return e.l
}

func (e UnreachableCodeError) Error() string {
	return "unreachable code"
}

func (e UnreachableCodeError) Severity() DiagnosticSeverity {
	return DiagWarn
}

type BJustBeforeLabelError struct {
	l int
}

func (e BJustBeforeLabelError) Line() int {
	return e.l
}

func (e BJustBeforeLabelError) Error() string {
	return "unconditional branch just before the target label"
}

func (e BJustBeforeLabelError) Severity() DiagnosticSeverity {
	return DiagWarn
}

type EmptyLoopError struct {
	l int
}

func (e EmptyLoopError) Line() int {
	return e.l
}

func (e EmptyLoopError) Error() string {
	return "empty loop"
}

func (e EmptyLoopError) Severity() DiagnosticSeverity {
	return DiagWarn
}

type MissingLabelError struct {
	l    int
	name string
}

func (e MissingLabelError) Line() int {
	return e.l
}

func (e MissingLabelError) Error() string {
	return fmt.Sprintf("missing label: \"%s\"", e.name)
}

func (e MissingLabelError) Severity() DiagnosticSeverity {
	return DiagErr
}

type InfiniteLoopError struct {
	l int
}

func (e InfiniteLoopError) Line() int {
	return e.l
}

func (e InfiniteLoopError) Error() string {
	return "infinite loop"
}

func (e InfiniteLoopError) Severity() DiagnosticSeverity {
	return DiagErr
}

type PragmaVersionAfterInstrError struct {
	l int
}

func (e PragmaVersionAfterInstrError) Line() int {
	return e.l
}

func (e PragmaVersionAfterInstrError) Error() string {
	return "#pragma version is only allowed before instructions"
}

func (e PragmaVersionAfterInstrError) Severity() DiagnosticSeverity {
	return DiagErr
}

func (l *Linter) getLabelsUsers() map[string][]int {
	used := map[string][]int{}

	for i, o := range l.l {
		switch o2 := o.(type) {
		case usesLabels:
			for _, l := range o2.Labels() {
				used[l.Name] = append(used[l.Name], i)
			}
		}
	}

	return used
}

func (l *Linter) getAllLabels() map[string][]int {
	all := map[string][]int{}

	for i, o := range l.l {
		switch o := o.(type) {
		case *LabelExpr:
			all[o.Name] = append(all[o.Name], i)
		}
	}

	return all
}

func (l *Linter) checkUnusedLabels() {
	used := l.getLabelsUsers()
	for name, lines := range l.getAllLabels() {
		if len(used[name]) == 0 {
			for _, line := range lines {
				l.errs = append(l.errs, &UnusedLabelError{l: line, name: name})
				l.reds = append(l.reds, &RedundantLabelLine{line: line, name: name})
			}
		}
	}
}

func (l *Linter) checkOpsAfterUnconditionalBranch() {
	used := l.getLabelsUsers()

	for i := 0; i < len(l.l); i++ {
		o := l.l[i]
		unused := false
		switch o.(type) {
		case *BExpr:
			unused = true
		case *ReturnExpr:
			unused = true
		case *ErrExpr:
			unused = true
		}

		if unused {
		loop:
			for i = i + 1; i < len(l.l); i++ {
				o2 := l.l[i]
				switch o2 := o2.(type) {
				case *LabelExpr:
					if len(used[o2.Name]) > 0 {
						break loop
					}
				case Nop:
				default:
					l.errs = append(l.errs, UnreachableCodeError{i})
				}
			}
		}
	}
}

func (l *Linter) checkBranchJustBeforeLabel() {
	for i, o := range l.l {
		func() {
			if i >= len(l.l)-1 {
				return
			}

			switch o := o.(type) {
			case *BExpr:
				j := i + 1

				func() {
				loop:
					for {
						if j >= len(l.l) {
							break
						}

						n := l.l[j]
						j += 1

						switch n := n.(type) {
						case *LabelExpr:
							if n.Name == o.Label.Name {
								l.errs = append(l.errs, BJustBeforeLabelError{l: i})
								l.reds = append(l.reds, RedundantBLine{line: i})
								return
							}
							break loop
						case Nop:
						default:
							break loop
						}
					}
				}()
			default:
			}
		}()
	}
}

func (l *Linter) checkLoops() {
	used := l.getLabelsUsers()
	all := l.getAllLabels()

	for name, users := range used {
		_, ok := all[name]
		if !ok {
			for _, user := range users {
				l.errs = append(l.errs, MissingLabelError{l: user, name: name})
			}
		}
	}

	for i, o := range l.l {
		if i == 0 {
			continue
		}

		func() {
			if i == 0 {
				return
			}

			switch o := o.(type) {
			case *BExpr:
				j := i - 1

				func() {
					for {
						if j < 0 {
							break
						}

						n := l.l[j]

						switch n := n.(type) {
						case *LabelExpr:
							if n.Name == o.Label.Name {
								if !l.canEscape(j, i) {
									l.errs = append(l.errs, InfiniteLoopError{l: i})
								}
								return
							}
						}
						j -= 1
					}
				}()
			default:
			}
		}()
	}
}

func (l *Linter) canEscape(from, to int) bool {
	labels := l.getAllLabels()

	for i := from; i <= to; i++ {
		switch op := l.l[i].(type) {
		case usesLabels:
			for _, lbl := range op.Labels() {
				for _, idx := range labels[lbl.Name] {
					if idx < from || idx > to {
						// TODO: check if the target label block is escapable
						return true
					}
				}
			}
		case Terminator:
			return true
		}
	}

	return false
}

func (l *Linter) checkDuplicatedLabels() {
	labels := l.getAllLabels()
	for name, lines := range labels {
		if len(lines) > 1 {
			for _, line := range lines {
				l.errs = append(l.errs, DuplicateLabelError{l: line, name: name})
			}
		}
	}
}

func (l *Linter) checkPragma() {
	var prev Op
	for i, op := range l.l {
		switch op := op.(type) {
		case *PragmaExpr:
			if prev != nil {
				l.errs = append(l.errs, PragmaVersionAfterInstrError{
					l: i,
				})
			}
			prev = op
		case Nop:
		default:
			prev = op
		}
	}
}

func (l *Linter) Lint() {
	l.checkDuplicatedLabels()
	l.checkUnusedLabels()
	l.checkOpsAfterUnconditionalBranch()
	l.checkBranchJustBeforeLabel()
	l.checkLoops()
	l.checkPragma()
}
