package money

type Money int

func (m Money) ToJPY() int {
	return int(m)
}
