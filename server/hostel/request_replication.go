package hostel

// ShouldReplicate says that by default the request are not replicated.
func (r *hostelRequest) ShouldReplicate() bool{
	return false
}

// ShouldReplicate a login request must be replicated.
func (r *loginRequest) ShouldReplicate() bool{
	return true
}

// ShouldReplicate a book request must be replicated.
func (r *bookRequest) ShouldReplicate() bool{
	return true
}

// ShouldReplicate a logout request must be replicated
func (r *logoutRequest) ShouldReplicate() bool{
	return true
}
