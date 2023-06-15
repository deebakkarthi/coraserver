package db

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Using the same conn as prod
// TODO: Maybe consider using a testing account in the db
const testDSN = "cora:@/cora?parseTime=true"

// setupTestDB creates tables and inserts test data
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", testDSN)
	if err != nil {
		t.Skip("MySQL test database not available:", err)
	}

	// Create test tables with exact schema
	createTables := []string{
		`CREATE TABLE IF NOT EXISTS slot (
			id INT,
			stime TIME NOT NULL, 
			etime TIME NOT NULL,
			PRIMARY KEY (id)
		)`,
		`CREATE TABLE IF NOT EXISTS subject (
			id CHAR(8),
			name VARCHAR(64) NOT NULL, 
			PRIMARY KEY (id)
		)`,
		`CREATE TABLE IF NOT EXISTS faculty (
			id CHAR(254),
			name VARCHAR(64) NOT NULL,
			PRIMARY KEY (id)
		)`,
		`CREATE TABLE IF NOT EXISTS static (
			class_id CHAR(4),
			day ENUM ("MON", "TUE", "WED", "THU", "FRI"), 
			slot_id INT, 
			faculty_id CHAR(254),
			subject_id CHAR(8),
			FOREIGN KEY (slot_id) REFERENCES slot (id), 
			FOREIGN KEY (faculty_id) REFERENCES faculty (id), 
			FOREIGN KEY (subject_id) REFERENCES subject (id), 
			PRIMARY KEY (class_id, day, slot_id)
		)`,
		`CREATE TABLE IF NOT EXISTS dynamic (
			class_id CHAR(4),
			date DATE, 
			slot_id INT, 
			faculty_id CHAR(254) NOT NULL, 
			subject_id CHAR(8) NOT NULL, 
			FOREIGN KEY (faculty_id) REFERENCES faculty (id), 
			FOREIGN KEY (slot_id) REFERENCES slot (id), 
			FOREIGN KEY (subject_id) REFERENCES subject (id), 
			PRIMARY KEY (class_id, date, slot_id)
		)`,
	}

	// Drop and recreate tables to ensure clean state
	dropTables := []string{
		"SET FOREIGN_KEY_CHECKS = 0",
		"DROP TABLE IF EXISTS dynamic",
		"DROP TABLE IF EXISTS static",
		"DROP TABLE IF EXISTS faculty",
		"DROP TABLE IF EXISTS subject",
		"DROP TABLE IF EXISTS slot",
		"SET FOREIGN_KEY_CHECKS = 1",
	}

	for _, query := range dropTables {
		if _, err := db.Exec(query); err != nil {
			t.Logf("Warning: Could not drop table: %v", err)
		}
	}

	for _, query := range createTables {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("Failed to create table: %v", err)
		}
	}

	// Insert test data
	insertTestData(t, db)
	return db
}

func insertTestData(t *testing.T, db *sql.DB) {
	// Insert slots with time information
	slotData := []struct {
		id    int
		stime string
		etime string
	}{
		{1, "08:50", "09:40"},
		{2, "09:40", "10:30"},
		{3, "10:40", "11:30"},
		{4, "11:30", "12:20"},
		{5, "13:40", "14:30"},
		{6, "14:30", "15:20"},
		{7, "15:20", "16:10"},
		{8, "16:10", "17:00"},
	}

	for _, slot := range slotData {
		_, err := db.Exec("INSERT INTO slot (id, stime, etime) VALUES (?, ?, ?)",
			slot.id, slot.stime, slot.etime)
		if err != nil {
			t.Fatalf("Failed to insert slot: %v", err)
		}
	}

	// Insert subjects (subset for testing)
	subjectData := []struct {
		id   string
		name string
	}{
		{"19CSE311", "Computer Security"},
		{"19CSE312", "Distributed Systems"},
		{"19CSE313", "Principles of Programming Languages"},
		{"19CSE314", "Software Engineering"},
		{"19CSE332", "Information Security"},
		{"19CSE434", "Image and Video Analysis"},
		{"19CSE352", "Business Analytics"},
		{"19CSE446", "Internet of Things"},
		{"19CSE435", "Computer Vision"},
		{"19CSE356", "Social Network Analytics"},
		{"19CSE456", "Neural Networks and Deep Learning"},
		{"FREE", "Free Period"},
	}

	for _, subject := range subjectData {
		_, err := db.Exec("INSERT INTO subject (id, name) VALUES (?, ?)",
			subject.id, subject.name)
		if err != nil {
			t.Fatalf("Failed to insert subject: %v", err)
		}
	}

	// Insert faculty (subset for testing)
	facultyData := []struct {
		id   string
		name string
	}{
		{"FREE", "No Faculty"},
		{"s_padmavathi@cb.amrita.edu", "Padmavathi.S"},
		{"n_harini@cb.amrita.edu", "Dr.Harini.N"},
		{"g_jeyakumar@cb.amrita.edu", "Dr.Jeyakumar.G"},
		{"r_aarthi@cb.amrita.edu", "Aarthi.R"},
		{"tr_swapna@cb.amrita.edu", "Dr.Swapna.T.R"},
		{"c_arunkumar@cb.amrita.edu", "Dr.Arunkumar.C"},
		{"t_gireeshkumar@cb.amrita.edu", "Dr.Gireesh Kumar T"},
		{"d_bharathi@cb.amrita.edu", "Ms.Bharathi.D"},
		{"test.faculty@test.com", "Test Faculty"},
	}

	for _, faculty := range facultyData {
		_, err := db.Exec("INSERT INTO faculty (id, name) VALUES (?, ?)",
			faculty.id, faculty.name)
		if err != nil {
			t.Fatalf("Failed to insert faculty: %v", err)
		}
	}

	// Insert static timetable data (subset focusing on A104 and C203 for testing)
	staticData := []struct {
		classID   string
		day       string
		slotID    int
		facultyID string
		subjectID string
	}{
		// A104 Monday schedule
		{"A104", "MON", 1, "s_padmavathi@cb.amrita.edu", "19CSE356"},
		{"A104", "MON", 2, "d_bharathi@cb.amrita.edu", "19CSE312"},
		{"A104", "MON", 3, "d_bharathi@cb.amrita.edu", "19CSE313"},
		{"A104", "MON", 4, "FREE", "FREE"},
		{"A104", "MON", 5, "FREE", "FREE"},
		{"A104", "MON", 6, "d_bharathi@cb.amrita.edu", "19CSE312"},
		{"A104", "MON", 7, "d_bharathi@cb.amrita.edu", "19CSE312"},
		{"A104", "MON", 8, "FREE", "FREE"},

		// A104 Tuesday schedule
		{"A104", "TUE", 1, "FREE", "FREE"},
		{"A104", "TUE", 2, "t_gireeshkumar@cb.amrita.edu", "19CSE311"},
		{"A104", "TUE", 3, "FREE", "FREE"},
		{"A104", "TUE", 4, "FREE", "FREE"},
		{"A104", "TUE", 5, "FREE", "FREE"},
		{"A104", "TUE", 6, "FREE", "FREE"},
		{"A104", "TUE", 7, "FREE", "FREE"},
		{"A104", "TUE", 8, "FREE", "FREE"},

		// C203 Monday schedule
		{"C203", "MON", 1, "s_padmavathi@cb.amrita.edu", "19CSE435"},
		{"C203", "MON", 2, "n_harini@cb.amrita.edu", "19CSE311"},
		{"C203", "MON", 3, "g_jeyakumar@cb.amrita.edu", "19CSE312"},
		{"C203", "MON", 4, "r_aarthi@cb.amrita.edu", "19CSE434"},
		{"C203", "MON", 5, "FREE", "FREE"},
		{"C203", "MON", 6, "g_jeyakumar@cb.amrita.edu", "19CSE312"},
		{"C203", "MON", 7, "g_jeyakumar@cb.amrita.edu", "19CSE312"},
		{"C203", "MON", 8, "FREE", "FREE"},

		// C203 Tuesday schedule
		{"C203", "TUE", 1, "tr_swapna@cb.amrita.edu", "19CSE313"},
		{"C203", "TUE", 2, "g_jeyakumar@cb.amrita.edu", "19CSE312"},
		{"C203", "TUE", 3, "c_arunkumar@cb.amrita.edu", "19CSE314"},
		{"C203", "TUE", 4, "c_arunkumar@cb.amrita.edu", "19CSE314"},
		{"C203", "TUE", 5, "FREE", "FREE"},
		{"C203", "TUE", 6, "FREE", "FREE"},
		{"C203", "TUE", 7, "FREE", "FREE"},
		{"C203", "TUE", 8, "FREE", "FREE"},
	}

	for _, data := range staticData {
		_, err := db.Exec(
			"INSERT INTO static (class_id, day, slot_id, faculty_id, subject_id) VALUES (?, ?, ?, ?, ?)",
			data.classID, data.day, data.slotID, data.facultyID, data.subjectID,
		)
		if err != nil {
			t.Fatalf("Failed to insert static data: %v", err)
		}
	}
}

func teardownTestDB(t *testing.T, db *sql.DB) {
	// Clean up test data
	dropTables := []string{
		"SET FOREIGN_KEY_CHECKS = 0",
		"DELETE FROM dynamic",
		"DELETE FROM static",
		"DELETE FROM faculty",
		"DELETE FROM subject",
		"DELETE FROM slot",
		"SET FOREIGN_KEY_CHECKS = 1",
	}
	for _, query := range dropTables {
		if _, err := db.Exec(query); err != nil {
			t.Logf("Warning: Could not clean table: %v", err)
		}
	}
	db.Close()
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestGetFreeClass(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	// Test date: Monday
	testDate := time.Date(2023, 6, 12, 0, 0, 0, 0, time.UTC) // Monday

	tests := []struct {
		name     string
		slot     int
		date     time.Time
		expected []string
	}{
		{
			name:     "Get free classes for slot 4 on Monday",
			slot:     4,
			date:     testDate,
			expected: []string{"A104"}, // A104 has FREE in slot 4, C203 has a class
		},
		{
			name:     "Get free classes for slot 5 on Monday",
			slot:     5,
			date:     testDate,
			expected: []string{"A104", "C203"}, // Both have FREE in slot 5
		},
		{
			name:     "Get free classes for slot 8 on Monday",
			slot:     8,
			date:     testDate,
			expected: []string{"A104", "C203"}, // Both have FREE in slot 8
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFreeClass(tt.slot, tt.date)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d classes, got %d. Expected: %v, Got: %v",
					len(tt.expected), len(result), tt.expected, result)
			}

			// Check if expected classes are present (order might vary)
			for _, expected := range tt.expected {
				found := false
				for _, actual := range result {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected class %s not found in result %v", expected, result)
				}
			}
		})
	}
}

func TestGetFreeSlot(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	testDate := time.Date(2023, 6, 12, 0, 0, 0, 0, time.UTC) // Monday

	tests := []struct {
		name     string
		class    string
		date     time.Time
		expected []int
	}{
		{
			name:     "Get free slots for A104 on Monday",
			class:    "A104",
			date:     testDate,
			expected: []int{4, 5, 8}, // slots 4, 5, 8 are FREE
		},
		{
			name:     "Get free slots for C203 on Monday",
			class:    "C203",
			date:     testDate,
			expected: []int{5, 8}, // slots 5, 8 are FREE
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFreeSlot(tt.class, tt.date)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d slots, got %d: %v", len(tt.expected), len(result), result)
			}

			for _, expected := range tt.expected {
				found := false
				for _, actual := range result {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected slot %d not found in result %v", expected, result)
				}
			}
		})
	}
}

func TestMultiFreeSlot(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	testDate := time.Date(2023, 6, 13, 0, 0, 0, 0, time.UTC) // Tuesday

	tests := []struct {
		name      string
		startSlot int
		endSlot   int
		date      time.Time
		expected  []string
	}{
		{
			name:      "Get classes free for slots 5-8 on Tuesday",
			startSlot: 5,
			endSlot:   8,
			date:      testDate,
			expected:  []string{"A104", "C203"}, // Both have slots 5-8 free on Tuesday
		},
		{
			name:      "Get classes free for slot 1 on Tuesday",
			startSlot: 1,
			endSlot:   1,
			date:      testDate,
			expected:  []string{"A104"}, // Only A104 has slot 1 free on Tuesday
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MultiFreeSlot(tt.startSlot, tt.endSlot, tt.date)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d classes, got %d: %v", len(tt.expected), len(result), result)
			}

			for _, expected := range tt.expected {
				found := false
				for _, actual := range result {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected class %s not found in result %v", expected, result)
				}
			}
		})
	}
}

func TestGetTimetableByDay(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	testDate := time.Date(2023, 6, 12, 0, 0, 0, 0, time.UTC) // Monday

	tests := []struct {
		name          string
		class         string
		date          time.Time
		expectedCount int
	}{
		{
			name:          "Get timetable for A104 on Monday",
			class:         "A104",
			date:          testDate,
			expectedCount: 8, // Should return 8 subjects (one for each slot)
		},
		{
			name:          "Get timetable for C203 on Monday",
			class:         "C203",
			date:          testDate,
			expectedCount: 8, // Should return 8 subjects (one for each slot)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTimetableByDay(tt.class, tt.date)

			if len(result) != tt.expectedCount {
				t.Errorf("Expected %d subjects, got %d: %v", tt.expectedCount, len(result), result)
			}
		})
	}
}

func TestGetAllSlot(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	result := GetAllSlot()
	expected := 8 // We inserted 8 slots

	if len(result) != expected {
		t.Errorf("Expected %d slots, got %d", expected, len(result))
	}

	// Check if all slots 1-8 are present
	for i := 1; i <= 8; i++ {
		found := false
		for _, slot := range result {
			if slot == i {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected slot %d not found in result", i)
		}
	}
}

func TestGetAllClass(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	result := GetAllClass()
	expectedClasses := []string{"A104", "C203"}

	if len(result) != len(expectedClasses) {
		t.Errorf("Expected %d classes, got %d", len(expectedClasses), len(result))
	}

	for _, expected := range expectedClasses {
		found := false
		for _, actual := range result {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected class %s not found in result %v", expected, result)
		}
	}
}

func TestGetAllSubject(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	result := GetAllSubject()
	// Should return all subjects except "FREE"
	expectedSubjects := []string{
		"19CSE311", "19CSE312", "19CSE313", "19CSE314", "19CSE332",
		"19CSE434", "19CSE352", "19CSE446", "19CSE435", "19CSE356", "19CSE456",
	}

	if len(result) != len(expectedSubjects) {
		t.Errorf("Expected %d subjects, got %d", len(expectedSubjects), len(result))
	}

	// Ensure "FREE" is not included
	for _, subject := range result {
		if subject == "FREE" {
			t.Error("FREE subject should not be included in GetAllSubject result")
		}
	}
}

func TestBooking(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	testDate := time.Date(2023, 6, 12, 0, 0, 0, 0, time.UTC) // Monday

	tests := []struct {
		name          string
		class         string
		date          time.Time
		slot          int
		faculty       string
		subject       string
		expectedRows  int64
		shouldSucceed bool
	}{
		{
			name:          "Book free slot successfully",
			class:         "A104",
			date:          testDate,
			slot:          4, // This is FREE in our test data for A104 Monday
			faculty:       "test.faculty@test.com",
			subject:       "19CSE311",
			expectedRows:  1,
			shouldSucceed: true,
		},
		{
			name:          "Try to book occupied slot",
			class:         "A104",
			date:          testDate,
			slot:          2, // This is occupied in our test data
			faculty:       "test.faculty@test.com",
			subject:       "19CSE312",
			expectedRows:  0,
			shouldSucceed: true, // Function succeeds but no rows affected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rowsAffected, err := Booking(tt.class, tt.date, tt.slot, tt.faculty, tt.subject)

			if tt.shouldSucceed && err != nil {
				t.Errorf("Expected booking to succeed, got error: %v", err)
			}

			if rowsAffected != tt.expectedRows {
				t.Errorf("Expected %d rows affected, got %d", tt.expectedRows, rowsAffected)
			}
		})
	}
}

func TestCancelBooking(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	testDate := time.Date(2023, 6, 12, 0, 0, 0, 0, time.UTC)

	// First create a booking
	_, err := Booking("A104", testDate, 4, "test.faculty@test.com", "19CSE311")
	if err != nil {
		t.Fatalf("Failed to create test booking: %v", err)
	}

	// Now test canceling it
	err = CancelBooking("A104", testDate, 4)
	if err != nil {
		t.Errorf("Failed to cancel booking: %v", err)
	}

	// Test canceling non-existent booking (should not error)
	err = CancelBooking("A104", testDate, 4)
	if err != nil {
		t.Errorf("Canceling non-existent booking should not error: %v", err)
	}
}

func TestGetBooking(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	testDate := time.Date(2023, 6, 12, 0, 0, 0, 0, time.UTC)
	faculty := "test.faculty@test.com"

	// Create some test bookings
	_, err := Booking("A104", testDate, 4, faculty, "19CSE311")
	if err != nil {
		t.Fatalf("Failed to create test booking: %v", err)
	}

	result := GetBooking(faculty)

	if len(result) != 1 {
		t.Errorf("Expected 1 booking, got %d", len(result))
	}

	if len(result) > 0 {
		booking := result[0]
		if booking.Class != "A104" {
			t.Errorf("Expected class A104, got %s", booking.Class)
		}
		if booking.Faculty != faculty {
			t.Errorf("Expected faculty %s, got %s", faculty, booking.Faculty)
		}
		if booking.Subject != "19CSE311" {
			t.Errorf("Expected subject 19CSE311, got %s", booking.Subject)
		}
		if booking.Slot != 4 {
			t.Errorf("Expected slot 4, got %d", booking.Slot)
		}
	}
}

func TestMultiBooking(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	testDate := time.Date(2023, 6, 13, 0, 0, 0, 0, time.UTC) // Tuesday

	tests := []struct {
		name         string
		class        string
		date         time.Time
		startSlot    int
		endSlot      int
		faculty      string
		subject      string
		expectedRows int64
	}{
		{
			name:         "Book multiple free slots",
			class:        "A104",
			date:         testDate,
			startSlot:    5,
			endSlot:      8, // Slots 5-8 are FREE for A104 on Tuesday
			faculty:      "test.faculty@test.com",
			subject:      "19CSE312",
			expectedRows: 4,
		},
		{
			name:         "Book range with some occupied slots",
			class:        "C203",
			date:         testDate,
			startSlot:    4,
			endSlot:      6, // Slot 4 is occupied, slots 5-6 are FREE
			faculty:      "test.faculty@test.com",
			subject:      "19CSE313",
			expectedRows: 2, // Only slots 5-6 should be booked
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rowsAffected, err := MultiBooking(tt.class, tt.date, tt.startSlot, tt.endSlot, tt.faculty, tt.subject)

			if err != nil {
				t.Errorf("MultiBooking failed: %v", err)
			}

			if rowsAffected != tt.expectedRows {
				t.Errorf("Expected %d rows affected, got %d", tt.expectedRows, rowsAffected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkGetFreeClass(b *testing.B) {
	// Create a dummy testing.T for setup
	dummyT := &testing.T{}
	db := setupTestDB(dummyT)
	defer teardownTestDB(dummyT, db)

	testDate := time.Date(2023, 6, 12, 0, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetFreeClass(4, testDate)
	}
}

func BenchmarkBooking(b *testing.B) {
	dummyT := &testing.T{}
	db := setupTestDB(dummyT)
	defer teardownTestDB(dummyT, db)

	testDate := time.Date(2023, 6, 12, 0, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Booking("A104", testDate, 4, "test.faculty@test.com", "19CSE311")
		CancelBooking("A104", testDate, 4) // Clean up for next iteration
	}
}
