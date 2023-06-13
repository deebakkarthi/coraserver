package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/deebakkarthi/coraserver/db"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

// Global OAuth Configuration variable
var oauthConfig *oauth2.Config

const (
	configFile     = "./config.json"
	port           = ":42069"
	organizationID = "00f9cda3-075e-44e5-aa0b-aba3add6539f"
)

/*
Temporary struct to unmarshall the config data from config.json
We cannot put the values of ClientSecret in the src code. It is a privileged
piece of information. So the configuration is stored in a =config.json= file.
In order to enforce the correct types this struct is needed. On a side note,
notice that the fields start with an uppercase. This means that they are to be
exported(accessible outside this package). You may think that it is not going
to be used outside this package, that is true, but since we are adding that
json tag they will be used by the =encoding/json= package to deserialize. So
whenever you want to deserialize a json file the corresponding struct members
should always be exported.
*/

type insertResponse struct {
	Inserted bool `json:"inserted"`
}
type oauthJSONRepr struct {
	ClientID     string   `json:"clientID"`
	ClientSecret string   `json:"clientSecret"`
	RedirectURL  string   `json:"redirectURL"`
	Scopes       []string `json:"scopes"`
	Tenant       string   `json:"tenant"`
}

type graphMe struct {
	OdataContext      string   `json:"@odata.context"`
	BusinessPhones    []string `json:"businessPhones"`
	DisplayName       string   `json:"displayName"`
	GivenName         string   `json:"givenName"`
	JobTitle          string   `json:"jobTitle"`
	Mail              string   `json:"mail"`
	MobilePhone       string   `json:"mobilePhone"`
	OfficeLocation    string   `json:"officeLocation"`
	PreferredLanguage string   `json:"preferredLanguage"`
	Surname           string   `json:"surname"`
	UserPrincipalName string   `json:"userPrincipalName"`
	ID                string   `json:"id"`
}

type graphOrganization struct {
	OdataContext string                   `json:"@odata.context"`
	Value        []graphOrganizationValue `json:"value"`
}

type graphOrganizationValue struct {
	ID                                        string   `json:"id"`
	DeletedDateTime                           string   `json:"deletedDateTime"`
	BusinessPhones                            []string `json:"businessPhones"`
	City                                      string   `json:"city"`
	Country                                   string   `json:"country"`
	CountryLetterCode                         string   `json:"countryLetterCode"`
	CreatedDateTime                           string   `json:"createdDateTime"`
	DefaultUsageLocation                      string   `json:"defaultUsageLocation"`
	DisplayName                               string   `json:"displayName"`
	IsMultipleDataLocationsForServicesEnabled string   `json:"isMultipleDataLocationsForServicesEnabled"`
	MarketingNotificationEmails               []string `json:"marketingNotificationEmails"`
	OnPremisesLastSyncDateTime                string   `json:"onPremisesLastSyncDateTime"`
	OnPremisesSyncEnabled                     string   `json:"onPremisesSyncEnabled"`
	PartnerTenantType                         string   `json:"partnerTenantType"`
	PostalCode                                string   `json:"postalCode"`
	PreferredLanguage                         string   `json:"preferredLanguage"`
	SecurityComplianceNotificationMails       []string `json:"securityComplianceNotificationMails"`
	SecurityComplianceNotificationPhones      []string `json:"securityComplianceNotificationPhones"`
	State                                     string   `json:"state"`
	Street                                    string   `json:"street"`
	TechnicalNotificationMails                []string `json:"technicalNotificationMails"`
	TenantType                                string   `json:"tenantType"`
	DirectorySizeQuota                        struct {
		Used  int `json:"used"`
		Total int `json:"total"`
	} `json:"directorySizeQuota"`
	OnPremisesSyncStatus []string `json:"onPremisesSyncStatus"`
	AssignedPlans        []string `json:"assignedPlans"`
	PrivacyProfile       struct {
		ContactEmail string `json:"contactEmail"`
		StatementURL string `json:"statementUrl"`
	} `json:"privacyProfile"`
	ProvisionedPlans []string `json:"provisionedPlans"`
	VerifiedDomains  []struct {
		Capabilities string `json:"capabilities"`
		IsDefault    bool   `json:"isDefault"`
		IsInitial    bool   `json:"isInitial"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	} `json:"verifiedDomains"`
}

type oauthExchangeResponse struct {
	Name         string `json:"name"`
	Mail         string `json:"mail"`
	Organization string `json:"organization"`
}

/*
=init()= is a special type of function like =main()= that is called automatically
by the go runtime. It is used to setup things that are needed before the main
function. Here we are setting up the oauthConfig variable by unmarshalling the
=config.json= file
*/
func init() {

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("Error reading JSON file:", err)
	}

	var jsonData oauthJSONRepr
	err = json.Unmarshal(file, &jsonData)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}

	oauthConfig = &oauth2.Config{
		ClientID:     jsonData.ClientID,
		ClientSecret: jsonData.ClientSecret,
		RedirectURL:  jsonData.RedirectURL,
		Scopes:       jsonData.Scopes,
		Endpoint:     microsoft.AzureADEndpoint(jsonData.Tenant),
	}

}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/oauth/login", oauthLoginHandler)
	router.HandleFunc("/oauth/exchange", oauthExchangeHandler)
	router.HandleFunc("/db/freeclass", freeClassHandler)
	router.HandleFunc("/db/freeslot", freeSlotHandler)
	router.HandleFunc("/db/daytimetable", dayTimetableHandler)
	router.HandleFunc("/db/booking", bookingHandler)
	router.HandleFunc("/db/getAllSlot", getAllSlotHandler)
	router.HandleFunc("/db/getAllClass", getAllClassHandler)
	router.HandleFunc("/db/getAllSubject", getAllSubjectHandler)
	router.HandleFunc("/db/getBooking", getBookingHandler)
	router.HandleFunc("/db/cancelBooking", cancelBookingHandler)

	server := &http.Server{Addr: port, Handler: router}

	log.Println("Server starting on port ", port)
	log.Fatal(server.ListenAndServe())
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	randomString := make([]byte, length)
	for i := 0; i < length; i++ {
		randomString[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomString)
}

func oauthLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateRandomString(16)
	authURL := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline, oauth2.SetAuthURLParam("prompt", "select_account"))
	http.Redirect(w, r, authURL, http.StatusFound)
}

func requestGraphAPI(accessToken string, endpoint string) ([]byte, error) {
	url := "https://graph.microsoft.com/v1.0/" + endpoint
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func oauthExchangeHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Println("Error while exchanging authorization code", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	graphMeResponse, err := requestGraphAPI(token.AccessToken, "me")
	if err != nil {
		log.Println("Error getting user profile", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}
	graphOrganizationResponse, err := requestGraphAPI(token.AccessToken, "organization")
	if err != nil {
		log.Println("Error getting user organization", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}
	var organization graphOrganization
	var profile graphMe
	json.Unmarshal(graphMeResponse, &profile)
	json.Unmarshal(graphOrganizationResponse, &organization)
	if organization.Value[0].ID == organizationID {
		response := oauthExchangeResponse{
			Name:         profile.GivenName,
			Mail:         profile.Mail,
			Organization: organization.Value[0].DisplayName,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Println("Error marshalling data", err)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("This app is only for members of Amrita Vishwa Vidyapeetham"))
	}
	return
}

func freeClassHandler(w http.ResponseWriter, r *http.Request) {
	slotStr := r.URL.Query().Get("slot")
	day := r.URL.Query().Get("day")
	slot, err := strconv.Atoi(slotStr)
	if err != nil {
		http.Error(w, "Invalid slot value", http.StatusBadRequest)
		return
	}
	var classroom []string = db.GetFreeClass(slot, day)
	responseJSON, err := json.Marshal(classroom)
	if err != nil {
		log.Println("Error marshalling data", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func freeSlotHandler(w http.ResponseWriter, r *http.Request) {
	class := r.URL.Query().Get("class")
	day := r.URL.Query().Get("day")
	var slot []int = db.GetFreeSlot(class, day)
	responseJSON, err := json.Marshal(slot)
	if err != nil {
		log.Println("Error marshalling data", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func dayTimetableHandler(w http.ResponseWriter, r *http.Request) {
	class := r.URL.Query().Get("class")
	day := r.URL.Query().Get("day")
	var subject []string = db.GetTimetableByDay(class, day)
	responseJSON, err := json.Marshal(subject)
	if err != nil {
		log.Println("Error marshalling data", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func getAllSlotHandler(w http.ResponseWriter, r *http.Request) {
	var slot []int = db.GetAllSlot()
	responseJSON, err := json.Marshal(slot)
	if err != nil {
		log.Println("Error marshalling data", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func getAllClassHandler(w http.ResponseWriter, r *http.Request) {
	var class []string = db.GetAllClass()
	responseJSON, err := json.Marshal(class)
	if err != nil {
		log.Println("Error marshalling data", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func getAllSubjectHandler(w http.ResponseWriter, r *http.Request) {
	var subject []string = db.GetAllSubject()
	responseJSON, err := json.Marshal(subject)
	if err != nil {
		log.Println("Error marshalling data", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func getBookingHandler(w http.ResponseWriter, r *http.Request) {
	faculty := r.URL.Query().Get("faculty")
	var subject []db.BookingRecord = db.GetBooking(faculty)
	responseJSON, err := json.Marshal(subject)
	if err != nil {
		log.Println("Error marshalling data", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func bookingHandler(w http.ResponseWriter, r *http.Request) {
	var response insertResponse
	class := r.URL.Query().Get("class")
	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slot, err := strconv.Atoi(r.URL.Query().Get("slot"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	faculty := r.URL.Query().Get("faculty")
	subject := r.URL.Query().Get("subject")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := db.Booking(class, date, slot, faculty, subject)
	if err != nil {
		log.Println(err)
		response.Inserted = false
	} else {
		if rowsAffected > 0 {
			response.Inserted = true
		} else {
			response.Inserted = false
		}
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshalling data", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
	return
}

func cancelBookingHandler(w http.ResponseWriter, r *http.Request) {
	class := r.URL.Query().Get("class")
	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slot, err := strconv.Atoi(r.URL.Query().Get("slot"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = db.CancelBooking(class, date, slot)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/profile.html", http.StatusFound)
	return
}
