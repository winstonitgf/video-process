package cloudflare

type Store struct {
	storeMap map[string]string
}

func (s *Store) Get(fingerprint string) (string, bool) {
	url, ok := s.storeMap[fingerprint]
	return url, ok
}
func (s *Store) Set(fingerprint, url string) {

	if s.storeMap == nil {
		s.storeMap = make(map[string]string)
	}

	s.storeMap[fingerprint] = url
}
func (s *Store) Delete(fingerprint string) {
	delete(s.storeMap, fingerprint)
}
func (s *Store) Close() {
	s.storeMap = nil
}
