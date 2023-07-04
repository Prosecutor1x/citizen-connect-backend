package model

type ProblemData struct {
	// IssueId          string        `bson:"_id"`
	IssueTitle       string        `json:"issue_name"`
	IssueType        string        `json:"issue_type"`
	IssueDescription string        `json:"issue_description"`
	IssueLevel       string        `json:"issue_level"`
	IssueLocation    IssueLocation `json:"issue_location"`
	IssueProgress    string        `json:"issue_progress"`
	IssueRaiserId    string        `json:"issue_raiser_id"`
	IssueDate        string        `json:"issue_date"`
	IssueImages      []string      `json:"issue_images"`
	IssueVideos      []string      `json:"issue_videos"`
	IssueComments    []Comment     `json:"issue_comments"`
}

type IssueRaiserDetails struct {
	IssueRaiserName         string `json:"issue_raiser_name"`
	IssueRaiserId           string `json:"issue_raiser_id"`
	IssueRaiserMail         string `json:"issue_raiser_mail"`
	IssueRaiserPhone        string `json:"issue_raiser_phone"`
	IssueRaiserProfilePhoto string `json:"issue_raiser_profile_photo"`
}

type Comment struct {
	Body        string             `json:"body"`
	IssueRaiser IssueRaiserDetails `json:"issue_raiser"`
	CommentType   string           `json:"issue_type"`
}

type UserData struct {
	// UserId           string `bson:"_id"`
	UserName         string `json:"user_name"`
	UserEmail        string `json:"user_email"`
	Gender           string `json:"gender"`
	UserPhone        string `json:"user_phone"`
	UserProfilePhoto string `json:"user_profile_photo"`
	UserLocation     string `json:"user_location"`
	UserAge          string `json:"user_age"`
	UserVerified     bool   `json:"user_verified"`
	UserIdProof      string `json:"user_id_proof"`
}

type Phone struct {
	Phone string `json:"user_phone"`
}

type IssueLocation struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}
