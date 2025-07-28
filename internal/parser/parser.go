package parser

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/carlosarraes/subs-cli/pkg/models"
)

type Parser struct {
	patterns []PatternMatcher
}

type PatternMatcher struct {
	Name    string
	Regex   *regexp.Regexp
	Type    string
	Example string
}

func New() *Parser {
	return &Parser{
		patterns: compilePatterns(),
	}
}

func (p *Parser) Parse(filename string) (*models.MediaInfo, error) {
	cleanName := cleanFilename(filename)

	for _, pattern := range p.patterns {
		if matches := pattern.Regex.FindStringSubmatch(cleanName); matches != nil {
			mediaInfo, err := p.extractMediaInfo(matches, pattern)
			if err != nil {
				continue
			}
			return mediaInfo, nil
		}
	}

	return nil, fmt.Errorf("unable to parse filename '%s': expected formats like:\n"+
		"  TV Show: Series.Name.S01E01.720p.x264-GROUP.mkv\n"+
		"  TV Show with Year: Series.Name.2024.S01E01.1080p.x265-GROUP.mkv\n"+
		"  Alternative TV: Series.Name.1x01.720p.WEB-DL.mkv\n"+
		"  Movie: Movie.Name.2023.1080p.BluRay.x264-GROUP.mp4", filename)
}

func (p *Parser) extractMediaInfo(matches []string, pattern PatternMatcher) (*models.MediaInfo, error) {
	submatches := pattern.Regex.SubexpNames()
	matchMap := make(map[string]string)

	for i, name := range submatches {
		if i > 0 && i < len(matches) && name != "" {
			matchMap[name] = strings.TrimSpace(matches[i])
		}
	}

	mediaInfo := &models.MediaInfo{
		Type: pattern.Type,
	}

	if title, ok := matchMap["title"]; ok {
		mediaInfo.Title = cleanTitle(title)
	}

	if year, ok := matchMap["year"]; ok && year != "" {
		mediaInfo.Year = year
	}

	if pattern.Type == "tv" {
		season, episode, err := p.extractSeasonEpisode(matchMap)
		if err != nil {
			return nil, err
		}
		mediaInfo.Season = season
		mediaInfo.Episode = episode
		mediaInfo.Type = "episode"
	}

	if quality, ok := matchMap["quality"]; ok && quality != "" {
		mediaInfo.Quality = quality
	}

	if source, ok := matchMap["source"]; ok && source != "" {
		mediaInfo.Source, mediaInfo.Codec = extractSourceAndCodec(source)
	}

	if err := p.validateMediaInfo(mediaInfo); err != nil {
		return nil, err
	}

	return mediaInfo, nil
}

func (p *Parser) extractSeasonEpisode(matchMap map[string]string) (int, int, error) {
	var season, episode int
	var err error

	if s, ok := matchMap["season"]; ok && s != "" {
		season, err = strconv.Atoi(s)
		if err != nil || season < 1 || season > 99 {
			return 0, 0, fmt.Errorf("invalid season number: %s", s)
		}
	}

	if e, ok := matchMap["episode"]; ok && e != "" {
		episode, err = strconv.Atoi(e)
		if err != nil || episode < 1 || episode > 999 {
			return 0, 0, fmt.Errorf("invalid episode number: %s", e)
		}
	}

	if alt, ok := matchMap["alt_episode"]; ok && alt != "" && season == 0 && episode == 0 {
		if len(alt) == 3 {
			season, err = strconv.Atoi(alt[:1])
			if err == nil {
				episode, err = strconv.Atoi(alt[1:])
			}
		} else if len(alt) == 4 {
			season, err = strconv.Atoi(alt[:2])
			if err == nil {
				episode, err = strconv.Atoi(alt[2:])
			}
		}
		if err != nil || season < 1 || episode < 1 {
			return 0, 0, fmt.Errorf("invalid alternative episode format: %s", alt)
		}
	}

	if season == 0 || episode == 0 {
		return 0, 0, fmt.Errorf("season and episode must be specified for TV shows")
	}

	return season, episode, nil
}

func (p *Parser) validateMediaInfo(info *models.MediaInfo) error {
	if info.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}

	if info.Type == "episode" && (!info.HasSeasonEpisode()) {
		return fmt.Errorf("TV episodes must have valid season and episode numbers")
	}

	if info.Year != "" {
		year, err := strconv.Atoi(info.Year)
		if err != nil || year < 1900 || year > 2030 {
			return fmt.Errorf("invalid year: %s", info.Year)
		}
	}

	return nil
}

func compilePatterns() []PatternMatcher {
	return []PatternMatcher{
		{
			Name:    "TV with Year (SxxExx)",
			Type:    "tv",
			Example: "Dark.Matter.2024.S01E01.1080p.x265-ELiTE.mkv",
			Regex: regexp.MustCompile(
				`^(?P<title>.*?)\.(?P<year>\d{4})\.S(?P<season>\d{1,2})E(?P<episode>\d{1,3})(?:\.(?P<quality>\d+p))?(?:\.(?P<source>.+?))?(?:\.(?P<ext>\w+))?$`,
			),
		},

		{
			Name:    "TV with Year (xXx format)",
			Type:    "tv",
			Example: "Series.Name.2024.1x01.720p.WEB-DL.mkv",
			Regex: regexp.MustCompile(
				`^(?P<title>.*?)\.(?P<year>\d{4})\.(?P<season>\d{1,2})x(?P<episode>\d{1,3})(?:\.(?P<quality>\d+p))?(?:\.(?P<source>.+?))?(?:\.(?P<ext>\w+))?$`,
			),
		},

		{
			Name:    "TV without Year (SxxExx)",
			Type:    "tv",
			Example: "The.Office.S03E07.720p.BluRay.x264.mkv",
			Regex: regexp.MustCompile(
				`^(?P<title>.*?)\.S(?P<season>\d{1,2})E(?P<episode>\d{1,3})(?:\.(?P<quality>\d+p))?(?:\.(?P<source>.+?))?\.(?P<ext>\w+)$`,
			),
		},

		{
			Name:    "TV without Year (SxxExx, no ext)",
			Type:    "tv",
			Example: "The.Office.S03E07.720p.BluRay.x264",
			Regex: regexp.MustCompile(
				`^(?P<title>.*?)\.S(?P<season>\d{1,2})E(?P<episode>\d{1,3})(?:\.(?P<quality>\d+p))?(?:\.(?P<source>.+?))?$`,
			),
		},

		{
			Name:    "TV Alternative (xXx format)",
			Type:    "tv",
			Example: "Series.Name.1x01.720p.WEB-DL.mkv",
			Regex: regexp.MustCompile(
				`^(?P<title>.*?)\.(?P<season>\d{1,2})x(?P<episode>\d{1,3})(?:\.(?P<quality>\d+p))?(?:\.(?P<source>.+?))?(?:\.(?P<ext>\w+))?$`,
			),
		},

		{
			Name:    "TV Alternative (3-digit format)",
			Type:    "tv",
			Example: "Series.Name.101.720p.x264.mkv",
			Regex: regexp.MustCompile(
				`^(?P<title>.*?)\.(?P<alt_episode>\d{3})(?:\.(?P<quality>\d+p))?(?:\.(?P<source>.+?))?(?:\.(?P<ext>\w+))?$`,
			),
		},

		{
			Name:    "Movie",
			Type:    "movie",
			Example: "Inception.2010.1080p.BluRay.x264-SPARKS.mkv",
			Regex: regexp.MustCompile(
				`^(?P<title>.*?)\.(?P<year>\d{4})(?:\.(?P<quality>\d+p))?(?:\.(?P<source>.+?))\.(?P<ext>mp4|mkv|avi|mov|wmv|flv|webm|m4v|mpg|mpeg|3gp)$`,
			),
		},

		{
			Name:    "Movie (no extension)",
			Type:    "movie",
			Example: "Movie.Name.2023.1080p.BluRay.x264",
			Regex: regexp.MustCompile(
				`^(?P<title>.*?)\.(?P<year>\d{4})(?:\.(?P<quality>\d+p))?(?:\.(?P<source>.+?))?$`,
			),
		},

		{
			Name:    "Movie (no quality)",
			Type:    "movie",
			Example: "Movie.Name.2023.BluRay.x264-GROUP.mp4",
			Regex: regexp.MustCompile(
				`^(?P<title>.*?)\.(?P<year>\d{4})\.(?P<source>.+?)\.(?P<ext>mp4|mkv|avi|mov|wmv|flv|webm|m4v|mpg|mpeg|3gp)$`,
			),
		},
	}
}

func cleanFilename(filename string) string {
	base := filepath.Base(filename)

	cleaned := strings.ReplaceAll(base, " ", ".")

	for strings.Contains(cleaned, "..") {
		cleaned = strings.ReplaceAll(cleaned, "..", ".")
	}

	return cleaned
}

func cleanTitle(title string) string {
	clean := strings.ReplaceAll(title, ".", " ")

	clean = strings.Join(strings.Fields(clean), " ")

	return strings.TrimSpace(clean)
}

func extractSourceAndCodec(combined string) (source, codec string) {
	if combined == "" {
		return "", ""
	}

	parts := strings.Split(combined, ".")
	var sourceParts []string

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if codecPart := extractCodecFromPart(part); codecPart != "" {
			if codec == "" {
				codec = codecPart
			}
			if releaseGroup := extractReleaseGroupFromPart(part, codecPart); releaseGroup != "" {
				sourceParts = append(sourceParts, releaseGroup)
			}
		} else {
			sourceParts = append(sourceParts, part)
		}
	}

	if len(sourceParts) > 0 {
		source = strings.Join(sourceParts, ".")
	}

	return source, codec
}

func extractCodecFromPart(part string) string {
	partLower := strings.ToLower(part)
	codecs := []string{
		"x264", "x265", "h264", "h265", "hevc", "avc", "xvid", "divx",
		"vp8", "vp9", "av1", "mpeg2", "mpeg4",
	}

	for _, codec := range codecs {
		if idx := strings.Index(partLower, codec); idx >= 0 {
			return part[idx : idx+len(codec)]
		}
	}

	return ""
}

func extractReleaseGroupFromPart(part, codec string) string {
	codecLower := strings.ToLower(codec)
	partLower := strings.ToLower(part)

	if idx := strings.Index(partLower, codecLower); idx >= 0 {
		afterCodec := part[idx+len(codec):]
		afterCodec = strings.TrimPrefix(afterCodec, "-")
		afterCodec = strings.TrimPrefix(afterCodec, ".")
		if afterCodec != "" {
			return afterCodec
		}
	}

	return ""
}

func isCodec(s string) bool {
	codecs := []string{
		"x264", "x265", "h264", "h265", "hevc", "avc", "xvid", "divx",
		"vp8", "vp9", "av1", "mpeg2", "mpeg4",
	}

	for _, codec := range codecs {
		if strings.Contains(s, codec) {
			return true
		}
	}

	return false
}
