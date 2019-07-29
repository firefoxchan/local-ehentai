package ehloader

type Q struct {
	op string
	k  TagK
	v  TagV
	subQs []Q
}

const (
	QOpAnd   = "and"
	QOpOr    = "or"
	QOpLike  = "like"
	QOpEq    = "="
)

func And (qs ...Q) Q {
	return Q{
		op: QOpAnd,
		subQs: qs,
	}
}

func Or (qs ...Q) Q {
	return Q{
		op: QOpOr,
		subQs: qs,
	}
}

func Like (k TagK, v TagV) Q {
	return Q{
		op: QOpLike,
		k: k,
		v: v,
	}
}

func Eq (k TagK, v TagV) Q {
	return Q{
		op: QOpEq,
		k: k,
		v: v,
	}
}
