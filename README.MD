# chaosplayr

chaosplayr is a small CLI tool that displays a searchable list of the last 2 years of media.ccc.de content. Can open MPV on keypress

## Features

- Search list of videos & authors, and keywords
- view summary from video
- Bookmark certain videos in a local list to store
- caching of the rss feed


## To be implemented

- if mpv has a callback mechanism, we could store watched time in that list and "continue" videos 
- parameters for either sd or hd and mpv command line switches

## Installation

* Prerequsites
    * go >= 1.23.2
    * mpv (tested with version 0.39.0)

```bash
  go install github.com/burgr033/chaosplayr/cmd/chaosplayr
```

### Build locally and run

```bash
  go build ./cmd/chaosplayr && ./chaosplayr
```

## Acknowledgements

 - [Bubbletea TUI framework](https://github.com/charmbracelet/bubbletea)
 - [Media.ccc.de](https://media.ccc.de)


## Contributors

- [@burgr033](https://www.github.com/burgr033)
- [@katrinleinweber](https://www.github.com/katrinleinweber)


## Contributing

Code & Resource contributions are always welcome!
