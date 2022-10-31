# Discord Scraper
Download all media from all chats in a channel using queue and multiple worker threads.

To use:
1. Copy `config.example.yaml` to `config.yaml`.  
2. Set `use_json` in `config.yaml` to `false`.
3. Insert Discord Channel ID and Discord API token in `config.yaml`.
4. Run `discord-scraper.exe`.
5. Open `storage/results` to see downloaded files.
