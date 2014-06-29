package debug

import "math/rand"
import "strconv"
import "strings"
import "regexp"
import "bufio"
import "sync"
import "time"
import "fmt"
import "net"
import "log"
import "os"
import "io"

// Pattern regexp.
var reg *regexp.Regexp

// Whether or not any are enabled.
var enabled = false

// Writer used when outputting debug information.
var Writer io.Writer = os.Stderr

// Mutex for enabling/disabling.
var m sync.Mutex

// Debugger function.
type DebugFunction func(string, ...interface{})

// Terminal colors used at random.
var colors []string = []string{
	"31",
	"32",
	"33",
	"34",
	"35",
	"36",
}

// Initialize with DEBUG environment variable.
func init() {
	env := os.Getenv("DEBUG")

	if "" != env {
		Enable(env)
	}

	go serve()
}

// Serve debugging information of a temporary unix domain socket.
func serve() {
	path := fmt.Sprintf("/tmp/debug-%d.sock", os.Getpid())
	lsock, err := net.Listen("unix", path)

	if err != nil {
		log.Printf("debug: failed to create %s", path)
		return
	}

	for {
		sock, err := lsock.Accept()
		if err != nil {
			log.Printf("debug: failed to accept connection")
			continue
		}

		go talk(sock)
	}
}

// Chat with the socket about receving debug lines.
// TODO: support multiple listeners
func talk(sock net.Conn) {
	prev := Writer
	Writer = sock

	for {
		r := bufio.NewReader(sock)
		b, _, err := r.ReadLine()
		if err != nil {
			log.Printf("debug: failed read command, disabling")
			Disable()
			Writer = prev
			return
		}

		cmd := string(b)

		switch strings.Trim(cmd, " ") {
		case "quit", "q":
			log.Printf("debug: quit")
			Disable()
			sock.Close()
			Writer = prev
			return
		case "disable", "d":
			log.Printf("debug: disabling")
			Disable()
		default:
			log.Printf("debug: enabling %q", cmd)
			Enable(cmd)
		}
	}
}

// Enable the given debug `pattern`. Patterns take a glob-like form,
// for example if you wanted to enable everything, just use "*", or
// if you had a library named mongodb you could use "mongodb:connection",
// or "mongodb:*". Multiple matches can be made with a comma, for
// example "mongo*,redis*".
//
// This function is thread-safe.
func Enable(pattern string) {
	m.Lock()
	defer m.Unlock()
	pattern = regexp.QuoteMeta(pattern)
	pattern = strings.Replace(pattern, "\\*", ".*?", -1)
	pattern = strings.Replace(pattern, ",", "|", -1)
	pattern = "^(" + pattern + ")$"
	reg = regexp.MustCompile(pattern)
	enabled = true
}

// Disable all pattern matching. This function is thread-safe.
func Disable() {
	m.Lock()
	defer m.Unlock()
	enabled = false
}

// Debug creates a debug function for `name` which you call
// with printf-style arguments in your application or library.
func Debug(name string) DebugFunction {
	prevGlobal := time.Now()
	color := colors[rand.Intn(len(colors))]
	prev := time.Now()

	return func(format string, args ...interface{}) {
		if !enabled {
			return
		}

		if !reg.MatchString(name) {
			return
		}

		d := deltas(prevGlobal, prev, color)
		fmt.Fprintf(Writer, d+" \033["+color+"m"+name+"\033[0m - "+format+"\n", args...)
		prevGlobal = time.Now()
		prev = time.Now()
	}
}

// Return formatting for deltas.
func deltas(prevGlobal, prev time.Time, color string) string {
	now := time.Now()
	global := now.Sub(prevGlobal).Nanoseconds()
	delta := now.Sub(prev).Nanoseconds()
	ts := now.UTC().Format("15:04:05.000")
	deltas := fmt.Sprintf("%s %-6s \033["+color+"m%-6s", ts, humanizeNano(global), humanizeNano(delta))
	return deltas
}

// Humanize nanoseconds to a string.
func humanizeNano(n int64) string {
	var suffix string

	switch {
	case n > 1e9:
		n /= 1e9
		suffix = "s"
	case n > 1e6:
		n /= 1e6
		suffix = "ms"
	case n > 1e3:
		n /= 1e3
		suffix = "us"
	default:
		suffix = "ns"
	}

	return strconv.Itoa(int(n)) + suffix
}
