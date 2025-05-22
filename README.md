# whispersync
<img src="./assets/WS.png" alt="Texte alternatif" width="200"/>
<hr>
<b>Project still in WIP</b>
<hr>
Keep your critical files safe and always up to date — automatically.
A lightweight solution that monitors your folders in real time and syncs changes instantly to a secure remote server.
No complex setup. No maintenance. Just peace of mind.

# local install
```bash
make install-dev && echo 'export WS_SERVER_APP_CONFIG_PATH=</path/to_app_config.yaml>' > ignored/.env && echo 'export WS_AGENT_APP_CONFIG_PATH=</path/to/app_config.yaml>' >> ignored/.env
```