package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestNewGame tests the game creation functionality
func TestNewGame(t *testing.T) {
	t.Parallel()
	store := &GameStore{games: make(map[string]*Game)}
	game := store.NewGame()

	tests := []struct {
		name      string
		checkFunc func() bool
		errorMsg  string
	}{
		{"game is not nil", func() bool { return game != nil }, "NewGame returned nil"},
		{"game ID is not empty", func() bool { return game.ID != "" }, "Game ID should not be empty"},
		{"game word is not empty", func() bool { return game.Word != "" }, "Game word should not be empty"},
		{"word length is 5", func() bool { return len(game.Word) == 5 }, "Word length should be 5"},
		{"max attempts is 6", func() bool { return game.MaxAttempts == 6 }, "MaxAttempts should be 6"},
		{
			"new game not won or lost",
			func() bool { return !game.IsWon && !game.IsLost },
			"New game should not be won or lost",
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if !tt.checkFunc() {
				t.Error(tt.errorMsg)
			}
		})
	}
}

// TestGameStoreGetGame tests retrieving games from the store
func TestGameStoreGetGame(t *testing.T) {
	t.Parallel()
	store := &GameStore{games: make(map[string]*Game)}
	game := store.NewGame()

	tests := []struct {
		name    string
		gameID  string
		wantErr bool
		errMsg  string
	}{
		{"existing game", game.ID, false, "Game should exist in store"},
		{"nonexistent game", "nonexistent-id", true, "Nonexistent game should not be found"},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, exists := store.GetGame(tt.gameID)
			if tt.wantErr && exists {
				t.Error(tt.errMsg)
			} else if !tt.wantErr && !exists {
				t.Error(tt.errMsg)
			}
		})
	}
}

// TestCheckGuess tests the guess checking logic
func TestCheckGuess(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		word          string
		guess         string
		expectedType  string // "correct", "present", "absent"
		checkPosition int    // -1 means check all positions
	}{
		{"all correct", "HOUSE", "HOUSE", "correct", -1},
		{"all absent", "HOUSE", "BRAIN", "absent", -1},
		{"has letter in wrong position", "HOUSE", "EATEN", "present", -1},
		{"first letter correct", "HOUSE", "HATER", "correct", 0},
		{"last letter correct", "HOUSE", "ANIME", "correct", 4},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			game := &Game{
				Word:        tt.word,
				MaxAttempts: 6,
				Attempts:    []Attempt{},
			}

			attempt := game.CheckGuess(tt.guess)

			if tt.checkPosition == -1 {
				// Check all positions for specific type
				switch tt.expectedType {
				case "correct":
					for i, status := range attempt.Status {
						if status.Status != "correct" {
							t.Errorf("Position %d should be correct, got %s", i, status.Status)
						}
					}
				case "absent":
					for i, status := range attempt.Status {
						if status.Status != "absent" {
							t.Errorf("Position %d should be absent, got %s", i, status.Status)
						}
					}
				case "present":
					// For "present" type, just verify at least one letter has this status
					found := false
					for _, status := range attempt.Status {
						if status.Status == "present" {
							found = true
							break
						}
					}
					if !found {
						t.Error("Expected at least one letter with 'present' status (wrong position)")
					}
				}
			} else {
				// Check specific position
				if attempt.Status[tt.checkPosition].Status != tt.expectedType {
					t.Errorf("Position %d: expected %s, got %s",
						tt.checkPosition, tt.expectedType, attempt.Status[tt.checkPosition].Status)
				}
			}
		})
	}
}

// TestMakeGuess tests different game scenarios
func TestMakeGuess(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		word        string
		guesses     []string
		expectWon   bool
		expectLost  bool
		expectCount int
	}{
		{
			name:        "win on first guess",
			word:        "HOUSE",
			guesses:     []string{"HOUSE"},
			expectWon:   true,
			expectLost:  false,
			expectCount: 1,
		},
		{
			name:        "lose after 6 wrong guesses",
			word:        "HOUSE",
			guesses:     []string{"BRAIN", "LIGHT", "CRANE", "FRUIT", "MEDAL", "PIZZA"},
			expectWon:   false,
			expectLost:  true,
			expectCount: 6,
		},
		{
			name:        "win on third attempt",
			word:        "HOUSE",
			guesses:     []string{"BRAIN", "LIGHT", "HOUSE"},
			expectWon:   true,
			expectLost:  false,
			expectCount: 3,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			game := &Game{
				Word:        tt.word,
				MaxAttempts: 6,
				Attempts:    []Attempt{},
				IsWon:       false,
				IsLost:      false,
			}

			for _, guess := range tt.guesses {
				game.MakeGuess(guess)
			}

			if game.IsWon != tt.expectWon {
				t.Errorf("Expected IsWon=%v, got %v", tt.expectWon, game.IsWon)
			}
			if game.IsLost != tt.expectLost {
				t.Errorf("Expected IsLost=%v, got %v", tt.expectLost, game.IsLost)
			}
			if len(game.Attempts) != tt.expectCount {
				t.Errorf("Expected %d attempts, got %d", tt.expectCount, len(game.Attempts))
			}
		})
	}
}

// TestMakeGuessAfterGameOver tests that guesses after game over are ignored
func TestMakeGuessAfterGameOver(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		isWon  bool
		isLost bool
	}{
		{"after win", true, false},
		{"after loss", false, true},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			game := &Game{
				Word:        "HOUSE",
				MaxAttempts: 6,
				Attempts:    []Attempt{},
				IsWon:       tt.isWon,
				IsLost:      tt.isLost,
			}

			game.MakeGuess("BRAIN")

			if len(game.Attempts) != 0 {
				t.Error("No attempts should be added after game is over")
			}
		})
	}
}

// TestHTTPHandlers tests HTTP handler responses
func TestHTTPHandlers(t *testing.T) {
	// Note: Not using t.Parallel() because this test modifies global gameStore
	tests := []struct {
		name           string
		handler        http.HandlerFunc
		method         string
		path           string
		expectStatus   int
		expectContains string
	}{
		{
			name:         "index handler redirects",
			handler:      IndexHandler,
			method:       http.MethodGet,
			path:         "/",
			expectStatus: http.StatusSeeOther,
		},
		{
			name:           "health handler returns OK",
			handler:        HealthHandler,
			method:         http.MethodGet,
			path:           "/health",
			expectStatus:   http.StatusOK,
			expectContains: "healthy",
		},
		{
			name:         "guess handler rejects GET",
			handler:      GuessHandler,
			method:       http.MethodGet,
			path:         "/guess",
			expectStatus: http.StatusMethodNotAllowed,
		},
	}

	gameStore = &GameStore{games: make(map[string]*Game)}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Not using t.Parallel() due to shared gameStore
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			tt.handler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectStatus {
				t.Errorf("Expected status %d, got %d", tt.expectStatus, resp.StatusCode)
			}

			if tt.expectContains != "" {
				body := w.Body.String()
				if !strings.Contains(body, tt.expectContains) {
					t.Errorf("Response should contain '%s'", tt.expectContains)
				}
			}
		})
	}
}

// TestHealthHandler tests the health check endpoint details
func TestHealthHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		checkField string
		checkValue string
	}{
		{"content type is JSON", "content-type", "application/json"},
		{"contains status field", "body", "status"},
		{"contains service field", "body", "wordle-game"},
		{"status is healthy", "body", "healthy"},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			w := httptest.NewRecorder()

			HealthHandler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.checkField == "content-type" {
				contentType := resp.Header.Get("Content-Type")
				if !strings.Contains(contentType, tt.checkValue) {
					t.Errorf("Expected Content-Type to contain %s, got %s", tt.checkValue, contentType)
				}
			} else if tt.checkField == "body" {
				body := w.Body.String()
				if !strings.Contains(body, tt.checkValue) {
					t.Errorf("Body should contain '%s'", tt.checkValue)
				}
			}
		})
	}
}

// TestGuessHandlerValidation tests guess validation logic
func TestGuessHandlerValidation(t *testing.T) {
	// Note: Not using t.Parallel() because this test modifies global gameStore
	gameStore = &GameStore{games: make(map[string]*Game)}
	game := gameStore.NewGame()

	tests := []struct {
		name           string
		guess          string
		expectAttempts int
		expectRedirect bool
	}{
		{"too short", "ABC", 0, true},
		{"too long", "ABCDEF", 0, true},
		{"valid length", "HOUSE", 1, true},
		{"empty guess", "", 0, true},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Not using t.Parallel() due to shared game instance
			// Reset game attempts
			game.Attempts = []Attempt{}

			form := url.Values{}
			form.Add("game_id", game.ID)
			form.Add("guess", tt.guess)

			req := httptest.NewRequest(http.MethodPost, "/guess", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			GuessHandler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.expectRedirect && resp.StatusCode != http.StatusSeeOther {
				t.Errorf("Expected redirect, got status %d", resp.StatusCode)
			}

			retrievedGame, _ := gameStore.GetGame(game.ID)
			if len(retrievedGame.Attempts) != tt.expectAttempts {
				t.Errorf("Expected %d attempts, got %d", tt.expectAttempts, len(retrievedGame.Attempts))
			}
		})
	}
}

// TestConcurrentGameCreation tests thread-safe game creation
func TestConcurrentGameCreation(t *testing.T) {
	store := &GameStore{games: make(map[string]*Game)}

	const numGames = 100
	done := make(chan bool, numGames)

	for i := 0; i < numGames; i++ {
		go func() {
			game := store.NewGame()
			if game == nil {
				t.Error("NewGame returned nil")
			}
			done <- true
		}()
	}

	for i := 0; i < numGames; i++ {
		<-done
	}

	if len(store.games) != numGames {
		t.Errorf("Expected %d games, got %d", numGames, len(store.games))
	}
}

// TestWordListValidity tests that all words in the word list are valid
func TestWordListValidity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		checkFunc func(word string, idx int) error
		errorFmt  string
	}{
		{
			name: "word length is 5",
			checkFunc: func(word string, idx int) error {
				if len(word) != 5 {
					t.Errorf("Word at index %d has length %d, expected 5: %s", idx, len(word), word)
				}
				return nil
			},
		},
		{
			name: "word is uppercase",
			checkFunc: func(word string, idx int) error {
				if strings.ToUpper(word) != word {
					t.Errorf("Word at index %d is not uppercase: %s", idx, word)
				}
				return nil
			},
		},
	}

	if len(wordList) == 0 {
		t.Fatal("Word list should not be empty")
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			for i, word := range wordList {
				tt.checkFunc(word, i)
			}
		})
	}
}
