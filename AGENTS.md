> **Dieses Repo wird im TOWER gebaut.** Melde dich an, bevor du liest, planst oder eine Datei anfasst:
> `tower anmelden <rufzeichen> "<deine Rolle>"` — die Antwort ist die Lage, und sie ist aktueller als jede Datei.
> Das Protokoll (Vorgang, Arbeitsbaum, Punkte, Sichtung, Bericht, dev) liefert der TOWER: `http://tower:8000/tower/doku/PROTOKOLL.md`.
> Ohne `tower` im PATH: `sh ~/.tower/kit/<pin>/tower anmelden …` (Pin: `kit` in `.tower/airfield.json`) — `tower kit verteiler` legt `tower` in den PATH, `tower hilfe` nennt jeden Befehl.
> Ist der TOWER stumm, warnt der Hook und merkt vor; gesperrt wird nur, wenn er antwortet und ablehnt.
> Was DIESES Airfield ist, steht darunter. Was der TOWER ist, steht dort.

# htw-calendar

Merges the LSF iCal feeds of the HTW course groups into one calendar. `.github/workflows/publish.yml`
runs hourly and on every push to `main`: it fetches the feeds, merges them (`main.go`) and publishes
one `.ics` on GitHub Pages, which a Google calendar subscribes to. Nothing runs anywhere else.

- **`main` is live.** A push to `main` republishes the calendar, so `main` must build:
  `sh ops/test.sh` (go vet and go build) passes before the Freigabe.
- **This repo is public.** The feed URLs and the file name are repository secrets (`FEEDS`,
  `FEED_NAME`, see `README.md`) and never appear in a file, a commit or a log line; errors name a
  feed by its position. No personal data, no paths of a machine.
