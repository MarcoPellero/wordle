package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

/*
 * /spawn: creates a new game session, sends back the session uuid as a set-cookie token
 * /kill: deletes a game session, identified by the request's uuid cookie
 * note that sessions automatically die some time after they've been created

 * /remove: removes a word from the server's wordlist
 * /add: adds a word to the server's wordlist
 * these endpoints are basically only useful for game bots, or in general systems where the wordlist shifts

 * /guess: decides the best guess for a specific session, identified by the request's uuid cookie
 */

type Session struct {
	candidates []string
	created_at time.Time
}

type Server struct {
	wordlist []string
	sessions map[uuid.UUID]Session
	mutex    sync.Mutex
}

var cache func([]byte) string

func dump_wordlist(wordlist []string, path string) {
	file, err := os.Create(path)
	if err != nil {
		panic(fmt.Sprintf("Couldn't create wordlist file at %s", path))
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for i, word := range wordlist {
		writer.WriteString(word)
		if i+1 != len(wordlist) {
			writer.WriteByte('\n')
		}
	}
}

func (server *Server) spawn(res http.ResponseWriter, req *http.Request) {
	new_uuid := uuid.New()
	new_session := Session{make([]string, len(server.wordlist)), time.Now()}
	copy(new_session.candidates, server.wordlist)

	server.mutex.Lock()
	server.sessions[new_uuid] = new_session
	server.mutex.Unlock()

	http.SetCookie(res, &http.Cookie{Name: "session-id", Value: new_uuid.String()})
	res.WriteHeader(200)
	fmt.Printf("/spawn %s\n", new_uuid.String())
}

func (server *Server) kill(res http.ResponseWriter, req *http.Request) {
	uuid_cookie, _ := req.Cookie("session-id")
	uuid, _ := uuid.Parse(uuid_cookie.Value)

	server.mutex.Lock()
	if _, ok := server.sessions[uuid]; !ok {
		res.WriteHeader(401)
		return
	}
	delete(server.sessions, uuid)
	server.mutex.Unlock()

	res.WriteHeader(200)
	fmt.Printf("/kill %s\n", uuid.String())
}

func (server *Server) remove(res http.ResponseWriter, req *http.Request) {
	bad_word := req.URL.Query().Get("word")
	if bad_word == "" {
		res.WriteHeader(400)
		return
	}

	was_found := false
	for i, x := range server.wordlist {
		if bad_word == x {
			server.wordlist[i] = server.wordlist[len(server.wordlist)-1]
			server.wordlist = server.wordlist[:len(server.wordlist)-1]
			was_found = true
			break
		}
	}

	if was_found {
		fmt.Printf("/remove %s\n", bad_word)
	}

	uuid_str := req.Header.Get("session-id")
	uuid, _ := uuid.Parse(uuid_str)

	server.mutex.Lock()
	session, ok := server.sessions[uuid]
	server.mutex.Unlock()
	if !ok {
		return
	}

	candidates := session.candidates
	for i, x := range candidates {
		if bad_word == x {
			candidates[i] = candidates[len(candidates)-1]
			session.candidates = candidates[:len(candidates)-1]
			break
		}
	}

	res.WriteHeader(200)
}

func (server *Server) add(res http.ResponseWriter, req *http.Request) {
	new_word := req.URL.Query().Get("word")
	if new_word == "" || len(new_word) != len(server.wordlist[0]) {
		res.WriteHeader(400)
		return
	}

	if !slices.Contains(server.wordlist, new_word) {
		server.wordlist = append(server.wordlist, new_word)
	}

	res.WriteHeader(200)
	fmt.Printf("/add %s\n", new_word)
}

func (server *Server) guess(res http.ResponseWriter, req *http.Request) {
	uuid_cookie, _ := req.Cookie("session-id")
	uuid, _ := uuid.Parse(uuid_cookie.Value)

	server.mutex.Lock()
	session, ok := server.sessions[uuid]
	server.mutex.Unlock()
	if !ok {
		res.WriteHeader(401)
		return
	}

	last_guess := req.URL.Query().Get("guess")
	pattern := req.URL.Query().Get("pattern")
	if len(last_guess) != len(server.wordlist[0]) || len(pattern) != len(server.wordlist[0]) {
		res.WriteHeader(400)
		return
	}

	candidates := session.candidates
	candidates = get_candidates(candidates, last_guess, []byte(pattern))
	server.mutex.Lock()
	server.sessions[uuid] = Session{candidates, session.created_at}
	server.mutex.Unlock()

	var next_guess Guess
	if last_guess == "sarti" {
		next_guess.word = cache([]byte(pattern))
	} else {
		var err error
		next_guess, err = get_optimal_guess(candidates, server.wordlist)
		if err != nil {
			res.WriteHeader(500)
			return
		}
	}

	res.Write([]byte(next_guess.word))
	fmt.Printf("/guess %s [%f] [%d solutions]\n", next_guess.word, next_guess.entropy, len(candidates))
}

func bot_server(wordlist_path string, cache_path string) {
	server := Server{read_wordlist(wordlist_path), make(map[uuid.UUID]Session), sync.Mutex{}}
	cache = build_cache(cache_path)

	http.HandleFunc("/spawn", server.spawn)
	http.HandleFunc("/kill", server.kill)
	http.HandleFunc("/remove", server.remove)
	http.HandleFunc("/add", server.add)
	http.HandleFunc("/guess", server.guess)
	go http.ListenAndServe(":8081", nil)

	for {
		server.mutex.Lock()
		for uuid, session := range server.sessions {
			if time.Since(session.created_at) >= 30*time.Second {
				delete(server.sessions, uuid)
			}
		}
		server.mutex.Unlock()

		fmt.Printf("[%d] [%d]\r", len(server.sessions), len(server.wordlist))
		dump_wordlist(server.wordlist, wordlist_path)
		time.Sleep(100 * time.Millisecond)
	}
}

func filter_wordlist_server(wordlist_path string) {
	server := Server{read_wordlist(wordlist_path), make(map[uuid.UUID]Session), sync.Mutex{}}
	filter_idx := 0

	http.HandleFunc("/spawn", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "session-id", Value: "filter-wordlist"})
		w.WriteHeader(200)
	})

	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/remove", server.remove)
	http.HandleFunc("/add", server.add)
	http.HandleFunc("/guess", func(w http.ResponseWriter, r *http.Request) {
		server.mutex.Lock()
		defer server.mutex.Unlock()

		w.Write([]byte(server.wordlist[filter_idx]))
		filter_idx++
	})
	go http.ListenAndServe(":8081", nil)

	for {
		server.mutex.Lock()
		for uuid, session := range server.sessions {
			if time.Since(session.created_at) >= 2*time.Minute {
				delete(server.sessions, uuid)
			}
		}
		server.mutex.Unlock()

		fmt.Printf("[%d] [%d]         \r", len(server.sessions), len(server.wordlist))
		dump_wordlist(server.wordlist, wordlist_path)
		time.Sleep(100 * time.Millisecond)
	}
}
