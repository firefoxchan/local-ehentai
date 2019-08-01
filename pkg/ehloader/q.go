package ehloader

import (
	"fmt"
	"strings"
)

type Q struct {
	op    string
	k     TagK
	v     TagV
	subQs []Q
}

func (q Q) Dump(prefix string, indent string, sep string) string {
	switch q.op {
	case QOpAnd, QOpOr:
		lines := make([]string, len(q.subQs)+3)
		lines[0] = fmt.Sprintf("%s[", prefix)
		lines[1] = fmt.Sprintf("%s%s", prefix+indent, strings.ToUpper(q.op))
		lines[len(q.subQs)+2] = fmt.Sprintf("%s]", prefix)
		for i, subQ := range q.subQs {
			lines[i+2] = subQ.Dump(prefix+indent, indent, sep)
		}
		return strings.Join(lines, sep)
	case QOpLike, QOpEq:
		return fmt.Sprintf("%s[%s %s %s]", prefix, q.op, q.k, q.v)
	}
	return "[INVALID]"
}

const (
	QOpAnd  = "and"
	QOpOr   = "or"
	QOpLike = "like"
	QOpEq   = "="
)

func And(qs ...Q) Q {
	return Q{
		op:    QOpAnd,
		subQs: qs,
	}
}

func Or(qs ...Q) Q {
	return Q{
		op:    QOpOr,
		subQs: qs,
	}
}

func Like(k TagK, v TagV) Q {
	return Q{
		op: QOpLike,
		k:  k,
		v:  v,
	}
}

func Eq(k TagK, v TagV) Q {
	return Q{
		op: QOpEq,
		k:  k,
		v:  v,
	}
}
