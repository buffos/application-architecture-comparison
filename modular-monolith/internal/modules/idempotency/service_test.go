package idempotency

import "testing"

type stubStore struct {
	results map[string]Result
}

func (s *stubStore) Find(key string) (Result, bool, error) {
	result, ok := s.results[key]
	return result, ok, nil
}

func (s *stubStore) Save(key string, result Result) error {
	s.results[key] = result
	return nil
}

func TestFindReturnsStoredResult(t *testing.T) {
	store := &stubStore{
		results: map[string]Result{
			"key-1": {ReturnRequestID: "return-001"},
		},
	}
	service := NewService(store)

	result, ok, err := service.Find("key-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !ok || result.ReturnRequestID != "return-001" {
		t.Fatalf("expected stored result to be returned")
	}
}
