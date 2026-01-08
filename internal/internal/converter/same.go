package converter

type Same struct{}

func (s *Same) Convert(from string) string {
	return from
}
