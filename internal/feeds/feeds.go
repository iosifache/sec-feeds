package feeds

import (
	"os"
	"path"

	"github.com/gorilla/feeds"
	"go.uber.org/zap"
)

const (
	OutputFeedsDir = "feeds"
)

type Feed interface {
	GetID() string
	GetFeed() feeds.Feed
}

var AllFeeds = []Feed{
	&PhrackFeed{},
}

func GenerateAllFeeds() {
	for _, feed := range AllFeeds {
		zap.L().Info("Generating feed", zap.String("id", feed.GetID()))

		rss, err := exportFeed(feed)
		if err != nil {
			zap.L().Error("Failed to export feed", zap.String("id", feed.GetID()), zap.Error(err))
			continue
		}

		err = saveExportedFeed(feed.GetID(), rss)
		if err != nil {
			zap.L().Error("Failed to save exported feed", zap.String("id", feed.GetID()), zap.Error(err))
			continue
		}
	}
}

func exportFeed(feed Feed) (string, error) {
	feedDetails := feed.GetFeed()
	rss, err := feedDetails.ToAtom()
	if err != nil {
		zap.L().Error("Failed to convert feed to RSS", zap.Error(err))
		return "", err
	}

	return rss, nil
}

func saveExportedFeed(feedID, rss string) error {
	outputPath := path.Join(OutputFeedsDir, feedID+".xml")
	err := os.WriteFile(outputPath, []byte(rss), 0644)
	if err != nil {
		zap.L().Error("Failed to write feed to file", zap.String("path", outputPath), zap.Error(err))
		return err
	}

	zap.L().Info("Feed saved successfully", zap.String("path", outputPath))

	return nil
}
