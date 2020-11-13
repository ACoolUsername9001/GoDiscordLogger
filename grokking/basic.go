package grokking

import "github.com/logrusorgru/grokky"


func GetConstants() grokky.Host {
	h := grokky.New()
	// Parsing the date at the beginning of the log:
	// "Jan  1"
	//h.Must("MONTH", "[a-zA-Z]{0,5}")
	h.Must("MONTH", `\bJan(?:uary|uar)?|Feb(?:ruary|ruar)?|M(?:a|ä)?r(?:ch|z)?|Apr(?:il)?|Ma(?:y|i)?|Jun(?:e|i)?|Jul(?:y)?|Aug(?:ust)?|Sep(?:tember)?|O(?:c|k)?t(?:ober)?|Nov(?:ember)?|De(?:c|z)(?:ember)?\b`)
	h.Must("SPACE", "\\s{1,}")
	h.Must("MONTHDAY", `(?:0[1-9])|(?:[12][0-9])|(?:3[01])|[1-9]`)
	h.Must("DATE", "%{MONTH:Month}%{SPACE}%{MONTHDAY:MonthDay}")

	// Parsing the time after the date of the log:
	// "12:34:56"
	h.Must("HOUR", `2[0123]|[01]?[0-9]`)
	h.Must("MINUTE", `[0-5][0-9]`)
	h.Must("SECOND", `(?:[0-5]?[0-9]|60)(?:[:.,][0-9]+)?`)
	h.Must("TIME", "%{HOUR:Hour}:%{MINUTE:Min}:%{MINUTE:Sec}")

	// Combining the entire time prefix
	h.Must("TIMEPREFIX", "%{DATE}%{SPACE}%{TIME}")

	// Preparing basic word parser
	h.Must("WORD", "\\w{1,}")
	// Parsing the process info from the log prefix:
	// "sshd[12345]"
	h.Must("REPORTER", "%{WORD:ProcessName}\\[%{WORD:ProcessID}\\]")

	// Combining everything into a general log line prefix:
	h.Must("PREFIX", "%{TIMEPREFIX:TimePrefix}%{SPACE}%{WORD:Host}")
	h.Must("PREFIX_WITH_REPORTER", "%{TIMEPREFIX:TimePrefix}%{SPACE}%{WORD:Host}%{SPACE}%{REPORTER}:")

	// Extras:
	h.Must("UNTIL_NEW_LINE", `[^\r\n]*`)
	h.Must("POSINT", `\b[1-9][0-9]*\b`)
	h.Must("USERNAME", `[a-zA-Z0-9._-]+`)
	h.Must("HOSTNAME", `\b[0-9A-Za-z][0-9A-Za-z-]{0,62}(?:\.[0-9A-Za-z][0-9A-Za-z-]{0,62})*(\.?|\b)`)
	h.Must("IPV6", `((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(%.+)?`)
	h.Must("IPV4", `(?:(?:[0-1]?[0-9]{1,2}|2[0-4][0-9]|25[0-5])[.](?:[0-1]?[0-9]{1,2}|2[0-4][0-9]|25[0-5])[.](?:[0-1]?[0-9]{1,2}|2[0-4][0-9]|25[0-5])[.](?:[0-1]?[0-9]{1,2}|2[0-4][0-9]|25[0-5]))`)
	h.Must("IP", `%{IPV6}|%{IPV4}`)
	h.Must("IPORHOST", `%{IP}|%{HOSTNAME}`)
	h.Must("HOSTPORT", `%{IPORHOST}:%{POSINT}`)

	// SSH-related
	h.Must("PUBKEY_PASSWORD", "publickey|password")
	return h
}