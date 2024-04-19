package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"html"
	"regexp"
	"strings"
	"unicode"

	//"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username,omitempty" bson:"username,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	UserType string             `json:"userType,omitempty" bson:"userType,omitempty"`
	Code     string             `json:"code,omitempty" bson:"code,omitempty"`
}

type Patients struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
}

type Doctors struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Date        string             `json:"date,omitempty" bson:"date,omitempty"`
	Time        string             `json:"time,omitempty" bson:"time,omitempty"`
	IsAvailable bool               `json:"isAvailable,omitempty" bson:"isAvailable,omitempty"`
}
type Appointments struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PatientName  string             `json:"patientName,omitempty" bson:"patientName,omitempty"`
	DoctorName   string             `json:"doctorName,omitempty" bson:"doctorName,omitempty"`
	SelectedDate string             `json:"date,omitempty" bson:"date,omitempty"`
	SelectedTime string             `json:"time,omitempty" bson:"time,omitempty"`
}
type UpdateDoctors struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	NewName string             `json:"newName,omitempty" bson:"newName,omitempty"`
	Name    string             `json:"name,omitempty" bson:"name,omitempty"`
	Date    string             `json:"date,omitempty" bson:"date,omitempty"`
	Time    string             `json:"time,omitempty" bson:"time,omitempty"`
	OldDate string             `json:"oldDate,omitempty" bson:"oldDate,omitempty"`
	OldTime string             `json:"oldTime,omitempty" bson:"oldTime,omitempty"`
}
type UpdateSlot struct {
	Name    string `json:"doctorName,omitempty" bson:"doctorName,omitempty"`
	Date    string `json:"date,omitempty" bson:"date,omitempty"`
	Time    string `json:"time,omitempty" bson:"time,omitempty"`
	OldDate string `json:"oldDate,omitempty" bson:"oldDate,omitempty"`
	OldTime string `json:"oldTime,omitempty" bson:"oldTime,omitempty"`
}
type ReservationEvent struct {
	DoctorID  string `json:"doctorName"`
	PatientID string `json:"patientName"`
	Operation string `json:"operation"` // ReservationCreated, ReservationUpdated, ReservationCancelled
}

type CustomClaims struct {
	Username string `json:"username,omitempty" bson:"username,omitempty"`
	jwt.StandardClaims
}

var jwtKey = []byte("jasmn")

var client *mongo.Client

/////////////////////////////////////////////////////////////////////////////////////////////////////

func containsUpperCase(s string) bool {
	for _, c := range s {
		if unicode.IsUpper(c) {
			return true
		}
	}
	return false
}

func containsDigit(s string) bool {
	for _, c := range s {
		if unicode.IsDigit(c) {
			return true
		}
	}
	return false
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

// sign up
func SignUPEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	//globalCodeValidation := "Doctor23165790"

	//ENCODE SENSITIVE DATA(PASSWORD)
	globalCodeValidation := "$2a$10$gmQfqxsQYv045tyt4X1qPuXftXsScCbVDz5Nd9BdoMrJMQMw6B1UK"

	//globalCodeValidationhashed, errr := bcrypt.GenerateFromPassword([]byte(globalCodeValidation), bcrypt.DefaultCost)
	//if errr != nil {
	//	fmt.Println("Error hashing password:", errr)
	//	return
	//}
	//fmt.Println("globalCodeValidationhashedconverted:", string(globalCodeValidationhashed))

	var user User
	json.NewDecoder(request.Body).Decode(&user)

	collection := client.Database("Clinic").Collection("Users")

	email := user.Username
	password := user.Password
	vCode := user.Code

	// Encoding the user-supplied data before sending the response
	encodedEmail := html.EscapeString(email)
	encodedPassword := html.EscapeString(password)

	fmt.Println("encooodedd Email:", encodedEmail)
	fmt.Println("encoodeddd Password:", encodedPassword)

	filter := bson.M{"username": email}

	fmt.Println("Email:", email)
	fmt.Println("Password:", password)
	fmt.Println("codee:", vCode)
	fmt.Println("userrr typeee:", user.UserType)

	if email == "" || password == "" {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	if len(email) > 50 || len(password) > 20 {
		response.WriteHeader(http.StatusLengthRequired)
		return
	}
	// Format validation for email using regular expression
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(encodedEmail) {
		response.WriteHeader(http.StatusNotAcceptable)
		return
	}
	passwordRegex := regexp.MustCompile(`^[A-Za-z\d@%+/-]{8,20}$`)
	if !passwordRegex.MatchString(encodedPassword) || !containsUpperCase(encodedPassword) || !containsDigit(encodedPassword) {
		response.WriteHeader(http.StatusNotAcceptable)
		return
	}

	var results bson.M
	err := collection.FindOne(context.Background(), filter).Decode(&results)
	if err != nil {
		if user.UserType == "patient" {
			user.Code = ""

			//ENCODE SENSITIVE DATA(PASSWORD)
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				fmt.Println("Error hashing password:", err)
				return
			}
			fmt.Println("hashed password:", hashedPassword)

			user.Password = string(hashedPassword)

			fmt.Println("user.Password :", user.Password)

			collection := client.Database("Clinic").Collection("Users")
			ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
			result, _ := collection.InsertOne(ctx, user)
			json.NewEncoder(response).Encode(result)

			collectionPatient := client.Database("Clinic").Collection("Patients")
			ctx2, _ := context.WithTimeout(context.Background(), 100*time.Second)
			result2, _ := collectionPatient.InsertOne(ctx2, bson.M{"name": email})
			json.NewEncoder(response).Encode(result2)
			response.WriteHeader(http.StatusOK)
			return
		} else if user.UserType == "doctor" {
			codeRegex := regexp.MustCompile(`^[A-Za-z\d]+$`)
			if !codeRegex.MatchString(vCode) {
				response.WriteHeader(http.StatusPreconditionFailed)
				return
			}
			if err := bcrypt.CompareHashAndPassword([]byte(globalCodeValidation), []byte(vCode)); err != nil {
				// Hashes don't match
				response.WriteHeader(http.StatusLocked)
				return
			} else {
				user.Code = ""
				collection := client.Database("Clinic").Collection("Users")
				ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)

				//ENCODE SENSITIVE DATA(PASSWORD)
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				if err != nil {
					fmt.Println("Error hashing password:", err)
					return
				}

				user.Password = string(hashedPassword)

				//inserting the account
				result, _ := collection.InsertOne(ctx, user)
				json.NewEncoder(response).Encode(result)
				response.WriteHeader(http.StatusOK)
				return
			}
		} else {
			response.WriteHeader(http.StatusMultipleChoices)
			return
		}

	} else if err == nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

func test(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	response.Write([]byte("hello credentials"))

}

//////////////////////////////////////////////////////////////////////////////////////////////////////

// sign in
func SignIN(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	// Decode JSON body
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{"error": "Invalid JSON"}`))
		return
	}
	email := input.Username
	password := input.Password
	collection := client.Database("Clinic").Collection("Users")

	if email == "" || password == "" {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	if len(email) > 50 || len(password) > 20 {
		response.WriteHeader(http.StatusLengthRequired)
		return
	}

	// Encoding the user-supplied data before sending the response
	encodedEmail := html.EscapeString(email)
	encodedPassword := html.EscapeString(password)
	// Format validation for email using regular expression
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(encodedEmail) {
		response.WriteHeader(http.StatusNotAcceptable)
		return
	}

	passwordRegex := regexp.MustCompile(`^[A-Za-z\d@%+/-]{8,20}$`)
	if !passwordRegex.MatchString(encodedPassword) {
		response.WriteHeader(http.StatusNotAcceptable)
		return
	}

	////ENCODE SENSITIVE DATA(PASSWORD)

	filter := bson.M{"username": email}
	var result bson.M
	err := collection.FindOne(context.Background(), filter).Decode(&result)

	storedPassword := result["password"].(string)

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte("Invalid credentials"))
		} else {
			log.Printf("error:%v", err)
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		}

	} else if result["userType"] == "doctor" {

		response.WriteHeader(http.StatusOK)

		// Generate JWT token
		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &CustomClaims{
			Username: input.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
		// Send JWT token in response
		response.Header().Set("Content-Type", "application/json")
		json.NewEncoder(response).Encode(map[string]string{"token": signedToken, "typeee": "doctor"})
		return

	} else if result["userType"] == "patient" {

		response.WriteHeader(http.StatusOK)

		//response.Write([]byte("patient"))
		// Generate JWT token
		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &CustomClaims{
			Username: input.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
		// Send JWT token in response
		response.Header().Set("Content-Type", "application/json")
		json.NewEncoder(response).Encode(map[string]string{"token": signedToken, "typeee": "patient"})
		return
	} else {
		response.WriteHeader(http.StatusLocked)
		return
	}

}

//////////////////////////////////////////////////////////////////////////////////////////////////////

// set schedule for the doctor
func SetDoctorSchudule(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var Doc Doctors
	json.NewDecoder(request.Body).Decode(&Doc)
	collection := client.Database("Clinic").Collection("Doctors")
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
	Doc.IsAvailable = true
	result, _ := collection.InsertOne(ctx, Doc)
	json.NewEncoder(response).Encode(result)
}

// Convert SetDoctorSchudule function into an http.Handler
type DoctorHandlerFunc func(http.ResponseWriter, *http.Request)

func (f DoctorHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

// cancel the patient's appointment
func CancelReservation(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var Appointment Appointments
	json.NewDecoder(request.Body).Decode(&Appointment)
	collection := client.Database("Clinic").Collection("Appointments")
	DocCollection := client.Database("Clinic").Collection("Doctors")
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
	filter := bson.M{
		"name":        Appointment.DoctorName,
		"date":        Appointment.SelectedDate,
		"time":        Appointment.SelectedTime,
		"isAvailable": false,
	}
	update := bson.M{
		"$set": bson.M{"isAvailable": true},
	}
	updateResult, updateErr := DocCollection.UpdateOne(ctx, filter, update)
	//updateResult, updateErr := DocCollection.UpdateOne(ctx,
	//bson.M{"name": Appointment.DoctorName, "date": Appointment.SelectedDate, "time": Appointment.SelectedTime, "isavailable": false},
	//bson.D{{"$set", bson.M{"isavailable": true}}})
	//fmt.Println("Filter: %+v", bson.M{"name": Appointment.DoctorName, "date": Appointment.SelectedDate, "time": Appointment.SelectedTime})
	//err := produceKafkaMessage(Appointment.DoctorName, Appointment.PatientName, "ReservationCanceled")
	//if err != nil {
	//	// Handle error
	//	fmt.Println("Error producing Kafka message:", err)
	//}

	if updateErr != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "error": "Failed to update doctor status" }`))
		return
	}
	if updateResult.ModifiedCount == 0 {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{ "error": "Doctor not found or status not updated" }`))
		return
	}
	result, err := collection.DeleteOne(ctx, Appointment)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "error": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(result)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////
func GetAllReservation(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var matchingAppointments []Appointments
	//var Appointment Appointments
	//json.NewDecoder(request.Body).Decode(&Appointment)
	collection := client.Database("Clinic").Collection("Appointments")
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
	cursor, err := collection.Find(ctx, bson.M{ /*"patientName": Appointment.PatientName*/ })
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var appointment Appointments
		cursor.Decode(&appointment)
		matchingAppointments = append(matchingAppointments, appointment)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(matchingAppointments)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

// get all patient's appointments
func GetAllDrSlots(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var matchingSlots []Doctors
	//var Appointment Appointments
	//json.NewDecoder(request.Body).Decode(&Appointment)
	collection := client.Database("Clinic").Collection("Doctors")
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var doc Doctors
		cursor.Decode(&doc)
		matchingSlots = append(matchingSlots, doc)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(matchingSlots)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

// get available slots
func GetAllSlots(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	//var req ReserveAppointmentRequest
	var doc []Doctors

	collection := client.Database("Clinic").Collection("Doctors")
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
	cursor, err := collection.Find(ctx, bson.M{"isAvailable": true})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var doctor Doctors
		cursor.Decode(&doctor)
		doc = append(doc, doctor)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(doc)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

// create appointment
func ReserveAppointment(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var appointment Appointments
	err := json.NewDecoder(request.Body).Decode(&appointment)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{ "error": "Invalid JSON" }`))
		return
	}
	// Validate the selected date and time
	if appointment.SelectedDate == "" || appointment.SelectedTime == "" {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{ "error": "Selected date and time are required" }`))
		return
	}
	// Store the appointment in the "Appointments" collection
	collection := client.Database("Clinic").Collection("Appointments")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, appointment)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "error": "` + err.Error() + `" }`))
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Write([]byte(fmt.Sprintf(`{ "message": "Appointment created with ID %s" }`, result.InsertedID)))

	///////////////////////////////////
	filter := bson.M{
		"name":        appointment.DoctorName,
		"date":        appointment.SelectedDate,
		"time":        appointment.SelectedTime,
		"isAvailable": true,
	}
	update := bson.M{
		"$set": bson.M{"isAvailable": false},
	}
	DocCollection := client.Database("Clinic").Collection("Doctors")
	ctx1, _ := context.WithTimeout(context.Background(), 100*time.Second)

	//err = produceKafkaMessage(appointment.DoctorName, appointment.PatientName, "ReservationCreated")
	//if err != nil {
	//	// Handle error
	//	fmt.Println("Error producing Kafka message:", err)
	//}
	updateResult, err := DocCollection.UpdateOne(ctx1, filter, update)
	if updateResult.ModifiedCount == 0 {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{ "error": "appointments not updated" }`))
		return
	}
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "error": "Failed to update slot availability" }`))
		return
	}

}

///////////////////////////////////////////////////////////////////////////////////////////////////////

// update Appointment
func UpdateReservationDoctor(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var Appointment UpdateDoctors
	json.NewDecoder(request.Body).Decode(&Appointment)
	collection := client.Database("Clinic").Collection("Appointments")
	DocCollection := client.Database("Clinic").Collection("Doctors")
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
	filter := bson.M{
		"name":        Appointment.NewName,
		"date":        Appointment.Date,
		"time":        Appointment.Time,
		"isAvailable": true,
	}
	update := bson.M{
		"$set": bson.M{"isAvailable": false},
	}
	updateResult, err := DocCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "error": "Failed to update slot availability" }`))
		return
	}
	// return old slot to true
	filterOldApp := bson.M{
		"name":        Appointment.Name,
		"date":        Appointment.OldDate,
		"time":        Appointment.OldTime,
		"isAvailable": false,
	}
	update = bson.M{
		"$set": bson.M{"isAvailable": true},
	}
	updateResult, err = DocCollection.UpdateOne(ctx, filterOldApp, update)
	fmt.Println("Update:", update)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "error": "Failed to update slot availability" }`))
		return
	}
	// update appointments in appointments collection
	update = bson.M{
		"$set": bson.M{"doctorName": Appointment.NewName, "date": Appointment.Date, "time": Appointment.Time},
	}
	updateResult, err = collection.UpdateOne(ctx, bson.M{"doctorName": Appointment.Name, "date": Appointment.OldDate, "time": Appointment.OldTime}, update)
	fmt.Println(updateResult)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "error": "Failed Appointments" }`))
		return
	}
	//err = produceKafkaMessage(Appointment.Name, "Patient", "ReservationUpdated")
	//if err != nil {
	//	// Handle error
	//	fmt.Println("Error producing Kafka message:", err)
	//}
	if updateResult.ModifiedCount == 0 {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{ "error": "appointments not updated" }`))
		return
	} else {
		response.WriteHeader(http.StatusOK)
		response.Write([]byte(fmt.Sprintf(`{ "message": "Appointment Updated Succesfully" }`)))
		return

	}

}

///////////////////////////////////////////////////////////////////////////////////////////////////////

// update the solt for the same doctor
func UpdateReservationSlot(response http.ResponseWriter, request *http.Request) {

	response.Header().Add("content-type", "application/json")
	var Appointment UpdateSlot
	json.NewDecoder(request.Body).Decode(&Appointment)
	collection := client.Database("Clinic").Collection("Appointments")
	DocCollection := client.Database("Clinic").Collection("Doctors")
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)

	// return old slot of old to true. make is available to true again
	filterNew := bson.M{
		"name":        Appointment.Name,
		"date":        Appointment.OldDate,
		"time":        Appointment.OldTime,
		"isAvailable": false,
	}
	fmt.Println(filterNew)
	update := bson.M{
		"$set": bson.M{"isAvailable": true},
	}
	fmt.Println(update)
	updateResult, err := DocCollection.UpdateOne(ctx, filterNew, update)
	fmt.Println(updateResult.ModifiedCount)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "error": "Failed to update slot availability" }`))
		return
	}
	// update state of new doc from is available =true to false
	filter := bson.M{
		"name":        Appointment.Name,
		"date":        Appointment.Date,
		"time":        Appointment.Time,
		"isAvailable": true,
	}
	fmt.Println(filter)
	update = bson.M{
		"$set": bson.M{"isAvailable": false},
	}
	fmt.Println(update)
	updateResult, err = DocCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "error": "Failed to update slot availability" }`))
		return
	}
	fmt.Println(updateResult.ModifiedCount)

	// update new slot in appointment
	update = bson.M{
		"$set": bson.M{"date": Appointment.Date, "time": Appointment.Time},
	}
	fmt.Println(update)
	updateResult, err = collection.UpdateOne(ctx, bson.M{"doctorName": Appointment.Name, "date": Appointment.OldDate, "time": Appointment.OldTime}, update)
	fmt.Println(Appointment.Name, Appointment.OldDate, Appointment.OldTime)

	//err = produceKafkaMessage(Appointment.Name, "Patient", "ReservationUpdated")
	//if err != nil {
	//	// Handle error
	//	fmt.Println("Error producing Kafka message:", err)
	//}
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "error": "Failed to update appointment" }`))
		return
	}
	if updateResult.ModifiedCount == 0 {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{ "error": "slot not updated" }`))
		return
	} else {
		response.WriteHeader(http.StatusOK)
		response.Write([]byte(fmt.Sprintf(`{ "message": "Appointment Updated Succesfully" }`)))
	}
}

///////////////////////////////////////////////////////////////////////////////////////////

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the Authorization header from the request
		authHeader := r.Header.Get("Authorization")

		// Check if the Authorization header exists
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Check if the Authorization header has the format "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		// Retrieve the JWT token from the Authorization header
		tokenString := tokenParts[1]

		// Parse and validate the JWT token
		claims := &CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

// New handler function for doctor schedule with auth middleware
func SetDoctorSchuduleWithAuth(response http.ResponseWriter, request *http.Request) {
	authMiddleware(http.HandlerFunc(SetDoctorSchudule)).ServeHTTP(response, request)
}

// ////////////////////////////////////////////////
func main() {
	fmt.Println("starting the app")
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)

	router := mux.NewRouter()
	router.HandleFunc("/test", test).Methods("GET")
	router.HandleFunc("/SignUP", SignUPEndPoint).Methods("POST")
	router.HandleFunc("/SignIN", SignIN).Methods("POST")
	router.HandleFunc("/doctor/SetSchudule", SetDoctorSchuduleWithAuth).Methods("POST")

	router.HandleFunc("/doctor/AllSlots", GetAllDrSlots).Methods("GET")

	router.HandleFunc("/patient/CancelReservation", CancelReservation).Methods("POST")
	router.HandleFunc("/patient/AllReservation", GetAllReservation).Methods("GET")
	router.HandleFunc("/patient/Getslot", GetAllSlots).Methods("GET")
	router.HandleFunc("/patient/ReserveAppointment", ReserveAppointment).Methods("POST")
	router.HandleFunc("/patient/UpdateReservation/Doctor", UpdateReservationDoctor).Methods("POST")
	router.HandleFunc("/patient/UpdateReservation/Slot", UpdateReservationSlot).Methods("POST")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"}) // Replace with the actual origin of your frontend
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	// Use the CORS middleware
	corsHandler := handlers.CORS(originsOk, headersOk, methodsOk)(router)
	// Start the server with the CORS middleware enabled
	err := http.ListenAndServe(":12345", corsHandler)
	if err != nil {
		log.Fatal(err)
	}
}
