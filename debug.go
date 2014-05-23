package debug

import "math/rand"
import "strconv"
import "strings"
import "regexp"
import "time"
import "fmt"
import "os"

//
// Variables
//

var prev time.Time = time.Now()
var reg *regexp.Regexp
var env string

//
// Colors.
//

var colors []string = []string{
	"31",
	"32",
	"33",
	"34",
	"35",
	"36",
}

//
// Initialize based on
// the DEBUG env var
//

func init() {
	env = os.Getenv("DEBUG")
	env = regexp.QuoteMeta(env)
	env = strings.Replace(env, "\\*", ".*?", -1)
	env = strings.Replace(env, ",", "|", -1)
	env = "^(" + env + ")$"
	reg = regexp.MustCompile(env)
}

//
// Debug function.
//

type DebugFunction func(string, ...interface{})

//
// Noop debug function.
//

func Noop(string, ...interface{}) {

}

//
// Debug function factory.
//

func Debug(name string) DebugFunction {
	if "" == env {
		return Noop
	}

	if !reg.MatchString(name) {
		return Noop
	}

	color := colors[rand.Intn(len(colors))]

	return func(format string, args ...interface{}) {
		now := time.Now()
		delta := now.Sub(prev).Nanoseconds()
		s := fmt.Sprintf("%8s", NanoToHuman(delta))
		fmt.Printf(s+" \033["+color+"m"+name+"\033[0m - "+format+"\n", args...)
		prev = now
	}
}

//
// Convert nanoseconds to formatted string.
//

func NanoToHuman(n int64) string {
	var suffix string

	switch {
	case n > 1000000000:
		n /= 1000000000
		suffix = "s"
	case n > 1000000:
		n /= 1000000
		suffix = "ms"
	case n > 1000:
		n /= 1000
		suffix = "us"
	default:
		suffix = "ns"
	}

	return strconv.Itoa(int(n)) + suffix
}
