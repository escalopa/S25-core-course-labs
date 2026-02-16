package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Game represents a Wordle game session
type Game struct {
	ID          string
	Word        string
	Attempts    []Attempt
	MaxAttempts int
	IsWon       bool
	IsLost      bool
}

// Attempt represents a single guess in the game
type Attempt struct {
	Word   string
	Status []LetterStatus
}

// LetterStatus represents the status of a letter in a guess
type LetterStatus struct {
	Letter string
	Status string // "correct", "present", "absent"
}

// GameStore manages active game sessions
type GameStore struct {
	mu    sync.RWMutex
	games map[string]*Game
}

var (
	gameStore = &GameStore{games: make(map[string]*Game)}
	templates *template.Template
	wordList  = []string{
		"AUDIO", "ADIEU", "ARISE", "RAISE", "SLATE",
		"HOUSE", "MOUSE", "CRANE", "POUND", "CRATE",
		"TRAIN", "BRAIN", "LEMON", "PEACH", "BREAD",
		"PLANT", "SMART", "BEACH", "CLOUD", "SPORT",
		"STONE", "LIGHT", "NIGHT", "FIGHT", "TIGHT",
	}

	// Prometheus metrics
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	gamesCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "games_created_total",
			Help: "Total number of games created",
		},
	)

	guessesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "guesses_total",
			Help: "Total number of guesses made",
		},
	)

	gamesWonTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "games_won_total",
			Help: "Total number of games won",
		},
	)

	gamesLostTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "games_lost_total",
			Help: "Total number of games lost",
		},
	)
)

func init() {
	// Register Prometheus metrics
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(gamesCreatedTotal)
	prometheus.MustRegister(guessesTotal)
	prometheus.MustRegister(gamesWonTotal)
	prometheus.MustRegister(gamesLostTotal)
}

// NewGame creates a new Wordle game
func (gs *GameStore) NewGame() *Game {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	game := &Game{
		ID:          uuid.New().String(),
		Word:        wordList[rand.Intn(len(wordList))],
		Attempts:    []Attempt{},
		MaxAttempts: 6,
		IsWon:       false,
		IsLost:      false,
	}
	gs.games[game.ID] = game
	gamesCreatedTotal.Inc()
	log.Printf("New game created - ID: %q, Word: %q", game.ID, game.Word)
	return game
}

// GetGame retrieves a game by ID
func (gs *GameStore) GetGame(id string) (*Game, bool) {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	game, exists := gs.games[id]
	return game, exists
}

// CheckGuess evaluates a guess and returns the status of each letter
func (g *Game) CheckGuess(guess string) Attempt {
	guess = strings.ToUpper(guess)
	attempt := Attempt{Word: guess, Status: make([]LetterStatus, 5)}

	// Create frequency map of letters in the target word
	letterCount := make(map[rune]int)
	for _, letter := range g.Word {
		letterCount[letter]++
	}

	// First pass: mark correct positions
	for i, letter := range guess {
		if i < len(g.Word) && rune(g.Word[i]) == letter {
			attempt.Status[i] = LetterStatus{
				Letter: string(letter),
				Status: "correct",
			}
			letterCount[letter]--
		}
	}

	// Second pass: mark present and absent
	for i, letter := range guess {
		if attempt.Status[i].Status == "" {
			if letterCount[letter] > 0 && strings.ContainsRune(g.Word, letter) {
				attempt.Status[i] = LetterStatus{
					Letter: string(letter),
					Status: "present",
				}
				letterCount[letter]--
			} else {
				attempt.Status[i] = LetterStatus{
					Letter: string(letter),
					Status: "absent",
				}
			}
		}
	}

	return attempt
}

// MakeGuess processes a player's guess
func (g *Game) MakeGuess(guess string) {
	if g.IsWon || g.IsLost {
		return
	}

	attempt := g.CheckGuess(guess)
	g.Attempts = append(g.Attempts, attempt)
	guessesTotal.Inc()

	// Check if won
	if strings.ToUpper(guess) == g.Word {
		g.IsWon = true
		gamesWonTotal.Inc()
	} else if len(g.Attempts) >= g.MaxAttempts {
		g.IsLost = true
		gamesLostTotal.Inc()
	}
}

// IndexHandler serves the main game page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(r.Method, "/"))
	defer timer.ObserveDuration()

	game := gameStore.NewGame()
	httpRequestsTotal.WithLabelValues(r.Method, "/", "303").Inc()
	http.Redirect(w, r, "/game/"+game.ID, http.StatusSeeOther)
}

// GameHandler serves the game interface
func GameHandler(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(r.Method, "/game/"))
	defer timer.ObserveDuration()

	gameID := strings.TrimPrefix(r.URL.Path, "/game/")
	game, exists := gameStore.GetGame(gameID)
	if !exists {
		httpRequestsTotal.WithLabelValues(r.Method, "/game/", "303").Inc()
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	httpRequestsTotal.WithLabelValues(r.Method, "/game/", "200").Inc()
	templates.ExecuteTemplate(w, "game.html", game)
}

// GuessHandler processes a guess submission
func GuessHandler(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(r.Method, "/guess"))
	defer timer.ObserveDuration()

	if r.Method != http.MethodPost {
		httpRequestsTotal.WithLabelValues(r.Method, "/guess", "405").Inc()
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameID := r.FormValue("game_id")
	guess := r.FormValue("guess")

	game, exists := gameStore.GetGame(gameID)
	if !exists {
		httpRequestsTotal.WithLabelValues(r.Method, "/guess", "303").Inc()
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Validate guess
	guess = strings.TrimSpace(guess)
	if len(guess) != 5 {
		httpRequestsTotal.WithLabelValues(r.Method, "/guess", "303").Inc()
		http.Redirect(w, r, "/game/"+gameID, http.StatusSeeOther)
		return
	}

	game.MakeGuess(guess)
	httpRequestsTotal.WithLabelValues(r.Method, "/guess", "303").Inc()
	http.Redirect(w, r, "/game/"+gameID, http.StatusSeeOther)
}

// HealthHandler serves health check endpoint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestsTotal.WithLabelValues(r.Method, "/health", "200").Inc()
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"healthy","service":"wordle-game"}`))
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/game/", GameHandler)
	http.HandleFunc("/guess", GuessHandler)
	http.HandleFunc("/health", HealthHandler)
	http.Handle("/metrics", promhttp.Handler())

	port := ":8080"
	log.Printf("Starting Wordle game server on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
