# htw-calendar

Merges my HTW course groups' LSF iCal feeds into one calendar, published hourly to GitHub Pages
under an unguessable file name.

- `FEEDS` secret: the LSF iCal links, one per line. New semester -> replace them.
- `FEED_NAME` secret: the file name, without `.ics`. Leaked -> set a new one; the old URL dies on the next run.
