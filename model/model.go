package model

type ProblemData struct {
	IssueName        string             `json:"issue_name"`
	IssueDescription string             `json:"issue_description"`
	IssueLocation    string             `json:"issue_location"`
	IssueStatus      string             `json:"issue_status"`
	IssueRaiser      IssueRaiserDetails `json:"issue_raiser"`
	IssueDate        string             `json:"issue_date"`
	IssueImages      []string           `json:"issue_images"`
	IssueVideos      []string           `json:"issue_videos"`
	IssueComments    []Comment          `json:"issue_comments"`
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
	IssueType   string             `json:"issue_type"`
}

type UserData struct {
	UserId           string `json:"user_id"`
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
