package lrange

type LRange struct {
	From *int `json:"from,omitempty"`
	To   *int `json:"to,omitempty"`
}

func new(from, to *int) LRange {
	return LRange{
		From: from,
		To:   to,
	}
}

func Between(from, to int) LRange {
	return new(&from, &to)
}

func From(from int) LRange {
	return new(&from, nil)
}

func To(to int) LRange {
	return new(nil, &to)
}
