package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"os"
	"strings"
	"time"
)

var client = &http.Client{Timeout: 30 * time.Second}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: htw-calendar <feeds.txt> <out.ics>")
		os.Exit(2)
	}
	if err := run(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(feedsPath, outPath string) error {
	urls, err := readFeeds(feedsPath)
	if err != nil {
		return err
	}
	// An empty list would publish an empty calendar over the last good one.
	if len(urls) == 0 {
		return errors.New("no feeds")
	}

	var timezone []string
	var events []string
	// Errors name the feed by its position: the run log of a public repo is public.
	for i, u := range urls {
		body, err := fetch(u)
		if err != nil {
			return fmt.Errorf("feed %d: %w", i+1, err)
		}
		tz, ev, err := split(body)
		if err != nil {
			return fmt.Errorf("feed %d: %w", i+1, err)
		}
		if timezone == nil {
			timezone = tz
		}
		events = append(events, ev...)
	}

	lines := []string{
		"BEGIN:VCALENDAR",
		"VERSION:2.0",
		"PRODID:-//Reconnact//htw-calendar//EN",
		"CALSCALE:GREGORIAN",
		"METHOD:PUBLISH",
		"X-WR-CALNAME:HTW",
		"X-WR-TIMEZONE:Europe/Berlin",
	}
	lines = append(lines, timezone...)
	lines = append(lines, events...)
	lines = append(lines, "END:VCALENDAR")

	// RFC 5545 wants CRLF line endings.
	return os.WriteFile(outPath, []byte(strings.Join(lines, "\r\n")+"\r\n"), 0o644)
}

func readFeeds(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var urls []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		urls = append(urls, line)
	}
	return urls, s.Err()
}

func fetch(url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		// url.Error carries the URL in its message.
		if ue, ok := errors.AsType[*neturl.Error](err); ok {
			return "", ue.Err
		}
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}
	b, err := io.ReadAll(resp.Body)
	return string(b), err
}

// split returns the VTIMEZONE block and the VEVENT lines of one feed.
// LSF answers a wrong or expired id with status 200 and an empty calendar,
// so a feed without events fails the run instead of dropping a course.
func split(body string) (timezone, events []string, err error) {
	lines := strings.Split(strings.ReplaceAll(body, "\r\n", "\n"), "\n")
	if len(lines) == 0 || lines[0] != "BEGIN:VCALENDAR" {
		return nil, nil, fmt.Errorf("not an iCalendar file")
	}

	var block *[]string
	count := 0
	for _, line := range lines {
		switch line {
		case "BEGIN:VTIMEZONE":
			block = &timezone
		case "BEGIN:VEVENT":
			block = &events
			count++
		}
		if block != nil {
			*block = append(*block, line)
		}
		if line == "END:VTIMEZONE" || line == "END:VEVENT" {
			block = nil
		}
	}
	if count == 0 {
		return nil, nil, fmt.Errorf("no events")
	}
	return timezone, events, nil
}
