package stream

//type Service struct {
//	aggregates Aggregates
//	store      Events
//}
//
//func New[R Root]() *Service {
//	var s Service
//
//	return &s
//}
//
//func (s *Service) Register(a Aggregate[Root]) error {
//	if a.store == nil {
//		a.store = s.store
//	}
//
//	n, err := NewName(a.Name)
//	if err != nil {
//		return err
//	}
//
//	s.aggregates[n] = a
//
//	return nil
//}
