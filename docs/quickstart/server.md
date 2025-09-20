# Server

The Server configuration is managed via [empire/server/config.yaml](https://github.com/BC-SECURITY/Empire/blob/master/empire/client/config.yaml).

Once launched, Empire checks for user write permissions on paths specified in `config.yaml`. If the current user does not have write permissions on these paths, `~/.empire` will be set as fallback parent directory and the configuration file will be updated as well.
If `empire-priv.key` and `empire-chain.pem` are not found in ~/.local/share/empire directory, self-signed certs will be generated.

* **suppress-self-cert-warning** - Suppress the http warnings when launching an Empire instance that uses a self-signed cert.

* **api** - Configure the RESTful API.

ip - The IP address to bind the API and Starkiller to.
port - The port to bind the API and Starkiller to.
secure - Enable HTTPS for the API and Starkiller. Browsers will not work with self-signed certs. Uses .key and .pem file from empire/server/data

```yaml
api:
  ip: 0.0.0.0
  port: 1337
  secure: false
```

* **database** - Configure Empire's database. Empire utilizes MySQL by default for high performance database operations. It can be configured to use sqlite for more lightweight implementations if required For more info on the database, see the [Database](database/README.md) section.

MySQL supports customizing the default url, username, password, and database name. By default these are set to
```yaml
database:
  use: mysql
  mysql:
    url: localhost:3306
    username: empire_user
    password: empire_password
    database_name: empire
```


If using SQLite the database location is customizable with the default setting:

```yaml
database:
  use: sqlite
  sqlite:
    location: empire/server/data/empire.db
```




The defaults block defines the properties that are initially loaded into the database when it is first created. These include the staging key, default user and password and obfuscation settings.

```yaml
database:
  defaults:
    # staging key will first look at OS environment variables, then here.
    # If empty, will be prompted (like Empire <3.7).
    staging-key: RANDOM
    username: empireadmin
    password: password123
    # The default configuration for global obfuscation.
    obfuscation:
      - language: powershell
        enabled: false
        command: "Token\\All\\1"
        module: "invoke-obfuscation"
        preobfuscatable: true
      - language: csharp
        enabled: false
        command: ""
        module: "confuser"
        preobfuscatable: false
    ip_allow_list: []
    ip_deny_list: []
    keyword_obfuscation:
      - Invoke-Empire
      - Invoke-Mimikatz
```

* **empire_compiler** - Configure the Empire Compiler module. This block manages settings for the Empire Compiler, which is responsible for handling C# compilation tasks.

archive: The URL to the Empire Compiler archive. The {{platform}} variable will be replaced with the current platform/architecture. (e.g. linux-amd64, linux-arm64)

```yaml
empire_compiler:
  archive: https://github.com/BC-SECURITY/Empire-Compiler/releases/download/v0.3.2/EmpireCompiler-{{platform}}-v0.3.2.tgz
```


* **plugins** - Config related to plugins
auto_start - boolean, whether the plugin should start automatically. If this is not set, Empire will defer to the plugin's own configuration.
auto_execute - run an execute command on the plugin at startup. If this is not set, Empire will defer to the plugin's own configuration.

```yaml
plugins:
  # Auto-execute plugin with defined settings
  basic_reporting:
    auto_start: true
    auto_execute:
      enabled: true
      options:
        report: all
```
* **plugin_marketplace** - This points the server to where Empire should look for additional available plugins to install. This defaults to the BC Security plugin marketplace but can point to a private marketplace as well.
name - the display name for the marketplace in Empire
git_url - git project to pull plugins from


```yaml
plugin_marketplace:
  registries:
    - name: BC-SECURITY
      git_url: git@github.com:BC-SECURITY/Empire-Plugin-Registry-Sponsors.git
      ref: main
      file: registry.yaml
```
* **directories** - Control where Empire should read and write specific data.

```yaml
directories:
  downloads: downloads
```

* **logging** - See [Logging](../logging/logging.md) for more information on logging configuration.

* **submodules** - Control if submodules will be auto updated on startup.

```
submodules:
  auto_update: true
```
