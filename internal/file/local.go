package file

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

// getCachePath returns the cache file path in XDG data directory
func getCachePath() (string, error) {
	cacheFilePath, err := xdg.DataFile("chaosplayr" + "/" + CacheFile)
	if err != nil {
		log.Fatalf("error obtaining data file through XDG: %v\n", err)
		return "", err
	}
	return cacheFilePath, nil
}

// updateCache saves the feed content to the cache file
func updateCache(content string) error {
	cachePath, err := getCachePath()
	if err != nil {
		return fmt.Errorf("error getting cache path: %v", err)
	}

	err = os.MkdirAll(filepath.Dir(cachePath), 0o755)
	if err != nil {
		return fmt.Errorf("error creating cache directory: %v", err)
	}

	err = os.WriteFile(cachePath, []byte(content), 0o644)
	if err != nil {
		return fmt.Errorf("error writing to cache file: %v", err)
	}

	return nil
}

// getFavoritesFilePath returns the path to the favorites file in the XDG data directory
func getFavoritesFilePath() (string, error) {
	favoritesFilePath, err := xdg.DataFile("chaosplayr/favorites.txt")
	if err != nil {
		return "", err
	}
	return favoritesFilePath, nil
}

// createFavorites creates an empty favorites file if it doesn't exist
func CreateFavorites() error {
	favoritesFilePath, err := getFavoritesFilePath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(favoritesFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

// IsInFavorites checks if a URL is in the favorites list
func IsInFavorites(url string) (bool, error) {
	favoritesFilePath, err := getFavoritesFilePath()
	if err != nil {
		return false, err
	}

	file, err := os.Open(favoritesFilePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == url {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

// addToFavorites adds a URL to the favorites list
func addToFavorites(url string) error {
	favoritesFilePath, err := getFavoritesFilePath()
	if err != nil {
		return err
	}

	if exists, err := IsInFavorites(url); err != nil {
		return err
	} else if exists {
		return nil // URL already in favorites
	}

	file, err := os.OpenFile(favoritesFilePath, os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(url + "\n")
	return err
}

// removeFromFavorites removes a URL from the favorites list
func removeFromFavorites(url string) error {
	favoritesFilePath, err := getFavoritesFilePath()
	if err != nil {
		return err
	}

	file, err := os.Open(favoritesFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != url {
			lines = append(lines, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	file, err = os.Create(favoritesFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// toggleFavorites toggles the presence of a URL in the favorites list
func ToggleFavorites(url string) error {
	exists, err := IsInFavorites(url)
	if err != nil {
		return err
	}

	if exists {
		return removeFromFavorites(url)
	} else {
		return addToFavorites(url)
	}
}
