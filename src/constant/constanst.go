package constant

const (
	// User
	AdminRoleName      string = "admin"
	DefaultRoleName    string = "default"
	DefaultUserName    string = "admin"
	RedisOtpDefaultKey string = "otp"

	// Claims
	AuthorizationHeaderKey string = "Authorization"
	UserIdKey              string = "UserId"
	UID           		   string = "UID"
	Level            	   string = "Level"
	Otp                    string = "OTP"
	State        	       string = "State"
	EmailKey               string = "Email"
	RolesKey               string = "Roles"
	ExpireTimeKey          string = "Exp"
	IssuedAtKey            string = "Iat"
)

// Keys for Labels
const (
	LabelKeyPhone    string = "phone"
	LabelKeyProfile  string = "profile"
	LabelKeyEmail    string = "email"
	LabelKeyDocument string = "document"
	LabelKeyLiveness string = "liveness"
)


// Values for Labels
const (
	LabelValueSubmitted string = "submitted"
	LabelValueVerified  string = "verified"
	LabelValueRejected  string = "rejected"
)

// Scopes for Labels
const (
	LabelScopePrivate string = "private"
	LabelScopePublic  string = "public"
)

// Descriptions for Labels
const (
	LabelDescriptionProfileSubmitted   string = "User profile information has been submitted for verification."
	LabelDescriptionProfileVerified     string = "User profile has been successfully verified."
	LabelDescriptionProfileRejected      string = "User profile submission has been rejected due to missing or incorrect information."

	LabelDescriptionPhoneSubmitted    string = "User phone number has been submitted for verification"
	LabelDescriptionPhoneVerified     string = "User phone number has been verified"
	LabelDescriptionPhoneRejected     string = "User phone number verification has been rejected"
	
	LabelDescriptionEmailSubmitted    string = "User email address has been submitted for verification"
	LabelDescriptionEmailVerified     string = "User email address has been verified"
	LabelDescriptionEmailRejected     string = "User email address verification has been rejected"
	
	LabelDescriptionDocumentSubmitted string = "User document has been submitted for verification"
	LabelDescriptionDocumentVerified  string = "User document has been verified"
	LabelDescriptionDocumentRejected  string = "User document verification has been rejected"
	
	LabelDescriptionLivenessSubmitted string = "User liveness has been submitted for verification"
	LabelDescriptionLivenessVerified  string = "User liveness has been verified"
	LabelDescriptionLivenessRejected  string = "User liveness verification has been rejected"
)