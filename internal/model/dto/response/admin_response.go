package response;

// AdminResponse represents an admin user and their associated establishment.
type AdminResponse struct {
    User         *UserResponse         `json:"user"`
    Establishment *EstablishmentResponse `json:"establishment"`
}