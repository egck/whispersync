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

# Agent/Server configuration

## Principle

Agent and server configurations rely on an app_config.yaml file that stores default key/values. These values can be filled in this file or overwritten by exporting environment variables. Details of each field are provided in cmd/<agent|server>.

If you want more details about the app_config implementation mecanism, please refer to https://github.com/bxgstudio/goconfigloader.

You will find an example for the server configuration below.

## Server configuration example
Configure following variables in app_config.yaml file as following:
```yaml
api_endpoint: :8080
secret: ""
data_storage_path: "" # Here no path is provided, it will be filled by an environment variable
```

You can also use environment variables to override default app_config.yaml values, capitalizing app_config variables and prefixing them by a WS_:
```bash
export WS_DATA_STORAGE_PATH=/my/path/to/backup
```

# Run server locally

```bash
make run-server
```

# Run agent locally
```bash
make run-agent
```