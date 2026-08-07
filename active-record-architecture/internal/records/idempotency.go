package records

func findIdempotentReturn(db *Database, command string, key string) (*ReturnRequest, bool) {
	returnID, ok := db.idempotency[command+":"+key]
	if !ok {
		return nil, false
	}

	request, err := FindReturnRequest(db, returnID)
	if err != nil {
		return nil, false
	}
	return request, true
}

func saveIdempotentReturn(db *Database, command string, key string, request *ReturnRequest) {
	if db.idempotency == nil {
		db.idempotency = make(map[string]string)
	}
	db.idempotency[command+":"+key] = request.ID
}
