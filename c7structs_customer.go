package c7api

type Customers struct {
	Customers []Customer `json:"customers"`
	Total     int        `json:"total"`
}

type Customer struct {
	Id        string `json:"id"`
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Birthdate string `json:"birthDate"`
}

type CustomerBasicInfo struct {
	Id        string `json:"id"`
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Birthdate string `json:"birthDate"`
	CustomerAddress
	Emails           []Email          `json:"emails"`
	OrderInformation OrderInformation `json:"orderInformation"`
}

type CustomerAddress struct {
	City        string `json:"city"`
	Statecode   string `json:"stateCode"`
	Zipcode     string `json:"zipCode"`
	Countrycode string `json:"countryCode"`
}

type CustomerFull struct {
	Customer
	CustomerAddress
	Avatar               string           `json:"avatar"`
	Honorific            *string          `json:"honorific"`
	Notifications        []Notification   `json:"notifications"`
	CreatedAt            string           `json:"createdAt"`
	UpdatedAt            string           `json:"updatedAt"`
	OrderInformation     OrderInformation `json:"orderInformation"`
	Loyalty              Loyalty          `json:"loyalty"`
	Emails               []Email          `json:"emails"`
	Phones               []Phone          `json:"phones"`
	Tags                 []Tag            `json:"tags"`
	Flags                []Flag           `json:"flags"`
	Groups               []Tag            `json:"groups"`
	Clubs                []Club           `json:"clubs"`
	Products             []Product        `json:"products"`
	HasAccount           bool             `json:"hasAccount"`
	LoginActivity        LoginActivity    `json:"loginActivity"`
	Emailmarketingstatus string           `json:"emailMarketingStatus"`
	Lastactivitydate     string           `json:"lastActivityDate"`
	Metadata             map[string]any   `json:"metaData"`
	Appdata              *string          `json:"appData"`
	Appsync              *string          `json:"appSync"`
}

type Flag struct {
	Id      string `json:"id"`
	Content string `json:"content"`
}

type Notification struct {
	Id       string      `json:"id"`
	Data     interface{} `json:"data"`
	Type     string      `json:"type"`
	Content  string      `json:"content"`
	ObjectId string      `json:"objectId"`
}

type OrderInformation struct {
	CurrentWebCartId        *string        `json:"currentWebCartId"`
	LastOrderId             string         `json:"lastOrderId"`
	LastOrderDate           string         `json:"lastOrderDate"`
	OrderCount              int            `json:"orderCount"`
	LifetimeValue           int            `json:"lifetimeValue"`
	LifetimeValueSeedAmount int            `json:"lifetimeValueSeedAmount"`
	YearlyValue             map[string]int `json:"yearlyValue"`
	Rank                    int            `json:"rank"`
	RankTrend               *string        `json:"rankTrend"`
	GrossProfit             int            `json:"grossProfit"`
	AcquisitionChannel      *string        `json:"acquisitionChannel"`
	CurrentClubTitle        *string        `json:"currentClubTitle"`
	DaysInCurrentClub       *int           `json:"daysInCurrentClub"`
	DaysInClub              int            `json:"daysInClub"`
	IsActiveClubMember      bool           `json:"isActiveClubMember"`
}

type Loyalty struct {
	Tier          string `json:"tier"`
	LoyaltyTierId string `json:"loyaltyTierId"`
	Points        int    `json:"points"`
}

type Email struct {
	Id     string `json:"id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type Phone struct {
	Id    string `json:"id"`
	Phone string `json:"phone"`
}

// type Tag struct {
// 	Id                 string      `json:"id"`
// 	Title              string      `json:"title"`
// 	ObjectType         string      `json:"objectType"`
// 	Type               string      `json:"type"`
// 	AppliesToCondition interface{} `json:"appliesToCondition"`
// 	CreatedAt          string      `json:"createdAt"`
// 	UpdatedAt          string      `json:"updatedAt"`
// }

type Club struct {
	ClubId           string `json:"clubId"`
	ClubTitle        string `json:"clubTitle"`
	CancelDate       string `json:"cancelDate"`
	SignupDate       string `json:"signupDate"`
	ClubMembershipId string `json:"clubMembershipId"`
}

type WineProperties struct {
	Type        *string `json:"type"`
	Region      *string `json:"region"`
	Vintage     *int    `json:"vintage"`
	Varietal    *string `json:"varietal"`
	Appellation *string `json:"appellation"`
	CountryCode *string `json:"countryCode"`
}

type ProductInner struct {
	Sku            string          `json:"sku"`
	Image          string          `json:"image"`
	Price          int             `json:"price"`
	Title          string          `json:"title"`
	Quantity       int             `json:"quantity"`
	ProductId      string          `json:"productId"`
	WineProperties *WineProperties `json:"wineProperties,omitempty"`
}

type Product struct {
	Product      ProductInner `json:"product"`
	PurchaseDate string       `json:"purchaseDate"`
}

type LoginActivity struct {
	LastLoginAt  *string `json:"lastLoginAt"`
	LoginIP      *string `json:"loginIP"`
	LastLogoutAt *string `json:"lastLogoutAt"`
}
