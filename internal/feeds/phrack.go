package feeds

import (
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/feeds"
	"go.uber.org/zap"
)

const (
	phrackIndexURL = "https://phrack.org"
	tenMegaInBytes = 1024 * 1024
	issueRegex     = `\/issues\/([0-9]*)\/1\.html`
	articlesRegex  = `<tr><td align=\"left\"><a href=\"([^\"]+)\">([^<]+)</a></td><td align=\"right\">([^<]+)</td></tr>`
	dateRegex      = `"datePublished":\s*"([0-9]{4}-[0-9]{2}-[0-9]{2})"`
)

type PhrackFeed struct{}

func (pf *PhrackFeed) GetID() string {
	return "phrack"
}

func (pf *PhrackFeed) GetFeed() feeds.Feed {
	now := time.Now()

	entries := pf.getEntries()

	return feeds.Feed{
		Title:       "Phrack Magazine",
		Link:        &feeds.Link{Href: phrackIndexURL},
		Description: "Phrack is written by hackers, for hackers, and offers a glimpse into the world just beyond what most people see.",
		Author:      &feeds.Author{Name: "Phrack Staff", Email: "staff[at]phrack{dot}org"},
		Created:     now,
		Items:       entries,
	}
}

func (pf *PhrackFeed) getEntries() []*feeds.Item {
	issuesURLs, err := pf.getIssuesURLs()
	if err != nil {
		zap.L().Error("Failed to get Phrack issues URLs", zap.Error(err))
		return nil
	}

	var items []*feeds.Item
	for _, issueURL := range issuesURLs {
		zap.L().Debug("Processing Phrack issue", zap.String("url", issueURL))

		articles, err := pf.fetchIssueArticlesAsRssItems(issueURL)
		if err != nil {
			zap.L().Error("Failed to fetch issue content", zap.String("url", issueURL), zap.Error(err))
			continue
		}

		items = append(items, articles...)
	}

	return items
}

func (pf *PhrackFeed) getIssuesURLs() ([]string, error) {
	resp, err := http.Get(phrackIndexURL)
	if err != nil {
		zap.L().Error("Failed to fetch Phrack index", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("Failed to read Phrack index content", zap.Error(err))
		return nil, err
	}

	issuesRe := regexp.MustCompile(issueRegex)
	matches := issuesRe.FindAllString(string(content), -1)

	return matches, nil
}

func (pf *PhrackFeed) fetchIssueArticlesAsRssItems(issueURL string) ([]*feeds.Item, error) {
	resp, err := http.Get(phrackIndexURL + issueURL)
	if err != nil {
		zap.L().Error("Failed to fetch Phrack issue", zap.String("url", issueURL), zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(io.LimitReader(resp.Body, tenMegaInBytes))
	if err != nil {
		zap.L().Error("Failed to read Phrack issue content", zap.String("url", issueURL), zap.Error(err))
		return nil, err
	}

	dateRe := regexp.MustCompile(dateRegex)
	dateMatch := dateRe.FindStringSubmatch(string(content))
	if len(dateMatch) != 2 {
		zap.L().Warn("Failed to parse Phrack issue date", zap.String("url", issueURL))
	}
	var date time.Time
	if len(dateMatch) == 2 {
		date, err = time.Parse("2006-01-02", dateMatch[1])
		if err != nil {
			zap.L().Warn("Failed to parse Phrack issue date as time", zap.String("url", issueURL), zap.Error(err))
			date = time.Now()
		}
	} else {
		date = time.Now()
	}

	issueID := regexp.MustCompile(issueRegex).FindStringSubmatch(issueURL)[1]

	var items []*feeds.Item
	articlesRe := regexp.MustCompile(articlesRegex)
	matches := articlesRe.FindAllStringSubmatch(string(content), -1)
	for _, m := range matches {
		if len(m) != 4 {
			return nil, nil
		}

		link := phrackIndexURL + m[1]
		title := m[2]
		author := m[3]

		zap.L().Debug("Parsed Phrack article", zap.String("title", title), zap.String("link", link), zap.String("author", author), zap.Time("date", date))

		item := feeds.Item{
			Title:   pf.getIssueTitle(issueID, title),
			Link:    &feeds.Link{Href: link},
			Author:  &feeds.Author{Name: author},
			Created: date,
		}
		items = append(items, &item)
	}

	return items, nil
}

func (pf *PhrackFeed) getIssueTitle(issueID, articleTitle string) string {
	return "Issue #" + issueID + ": " + articleTitle
}
