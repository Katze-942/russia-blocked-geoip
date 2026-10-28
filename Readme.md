# Katze GeoIP

Этот форк собирает только **`geoip-katze.dat`** — лёгкую базу из 16 категорий.
Генерация остальных DAT, MMDB, SRS, MRS, Clash, SURGE, nginx и текстовых правил отключена.

Категории перечислены в `wantedList` файла [config.json](config.json).
Сборка проверяет их наличие, корректность CIDR и создаёт SHA-256.
Обновление запускается каждые 6 часов и вручную через GitHub Actions.

- [geoip-katze.dat](https://raw.githubusercontent.com/Katze-942/russia-blocked-geoip/release/geoip-katze.dat)
- [SHA-256](https://raw.githubusercontent.com/Katze-942/russia-blocked-geoip/release/geoip-katze.dat.sha256sum)

Для замены стандартной базы сохраняйте файл как `geoip.dat`
и используйте ссылку на лёгкую базу в настройках автообновления.

Форк основан на [runetfreedom/russia-blocked-geoip](https://github.com/runetfreedom/russia-blocked-geoip)
и использует генератор [runetfreedom/geoip](https://github.com/runetfreedom/geoip).
Источники данных указаны в [config.json](config.json).
