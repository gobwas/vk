package vk

//go:generate easyjson -all

type Friend struct {
	ID                     int          `json:"id"`
	FirstName              string       `json:"first_name"`
	LastName               string       `json:"last_name"`
	Nickname               string       `json:"nickname"`
	Domain                 string       `json:"domain"`
	Sex                    int          `json:"sex"`
	Bdate                  string       `json:"bdate"`
	City                   City         `json:"city"`
	Country                Country      `json:"country"`
	Timezone               int          `json:"timezone"`
	Photo50                string       `json:"photo_50"`
	Photo100               string       `json:"photo_100"`
	Photo200Orig           string       `json:"photo_200_orig"`
	HasMobile              int          `json:"has_mobile"`
	Contacts               Contact      `json:"contacts"`
	Education              Education    `json:"education"`
	Online                 int          `json:"online"`
	Relation               int          `json:"relation"`
	LastSeen               LastSeen     `json:"last_seen"`
	Status                 string       `json:"status"`
	CanWritePrivateMessage int          `json:"can_write_private_message"`
	CanSeeAllPosts         int          `json:"can_see_all_posts"`
	CanPost                int          `json:"can_post"`
	Universities           []University `json:"universities"`
}

type LastSeen struct {
	Time     int `json:"time"`
	Platform int `json:"platform"`
}

type University struct {
	ID              int    `json:"id"`
	Country         int    `json:"country"`
	City            int    `json:"city"`
	Name            string `json:"name"`
	Faculty         int    `json:"faculty"`
	FacultyName     string `json:"faculty_name"`
	Chair           int    `json:"chair"`
	ChairName       string `json:"chair_name"`
	Graduation      int    `json:"graduation"`
	EducationForm   string `json:"education_form"`
	EducationStatus string `json:"education_status"`
}

type Education struct {
	University     int    `json:"university"`
	UniversityName string `json:"university_name"`
	Faculty        int    `json:"faculty"`
	FacultyName    string `json:"faculty_name"`
	Graduation     int    `json:"graduation"`
}

type Country struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type Contact struct {
	Mobile string `json:"mobile_phone"`
	Home   string `json:"home_phone"`
}

type Friends struct {
	Count int      `json:"count"`
	Items []Friend `json:"items"`
}

type City struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}
