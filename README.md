# Backpack

Backs up your directories periodically with special handling like snapshoting sqlite file.

When running self-hosting applications, backup is important and we often find different
application has different backing up strategies. Simply copying stuff may not be good
for applications using database files (e.g. sqlite). Backpack offers a way to specify
different rules to snapshot different applications using directory rule and regex-based
file rules, and then it uses `librclone` to backup your files to storage platforms that
you wish to backup your files.

## Version

v0.1

## How to use?

### Run directly

Run the binary directly with the following command line arguments:

```txt
-config string
    config file path (required)
-try-first
    try backup before running cron
-try-only
    try backup only without starting cron
-version
    show version
```

For example `./backpack -config config.json -try-first`. The config file is a JSON file
structured as the following format:

```javascript
{
    // Config version (optional)
    "version": "0.1",
    // Backup rules
    "backups": [
        {
            // List directories to backup with their own rules
            "directories": [
                {
                    // Directory path
                    "path": "/data/run_files",
                    // Rules for file
                    "file_rules": [
                        {
                            // Regular expression for all files ending in `sqlite3`
                            "regex": ".*sqlite3",
                            // Snapshots these files using `sqlite3 .backup`
                            "command": "sqlite"
                        },
                        {
                            // The regular expression matches relative path in directory path
                            // The file rules are matched in order, the first matched rule
                            // will be executed
                            "regex": "files/\.unwanted/.*",
                            // Ignore these files
                            "command": "ignore"
                        }
                        // The rest of the files will be copied
                    ]
                }
                {
                    "path": "/data/photos"
                    // If no file rules, all directory will be copied,
                    // which is equivalent to:
                    // "file_rules": [
                    //     {
                    //         "regex": ".*",
                    //         "command": "copy"
                    //     }
                    // ]
                }
            ],
            // Rclone remote name in Rclone configuration (the name in [])
            "rclone_remote": "default",
            // Remote path (e.g. bucket name)
            "remote_path": "my_bucket/prefix",
            // The backup schedule
            "schedule": {
                // There are two options to schedule
                // 1. Specifying time HH:MM for every day
                "daily": [
                    "01:00",
                    "12:00",
                    "22:00" // 10:00pm
                ],
                // 2. Specify an interval which can be string like "30m", "1h", "300s", etc.
                "every": "30m" // backup every 30 minutes
            }
        }
        // More backup rules can be specified
    ]
}
```

> Note that comments are not supported in JSON file. They are used to explain the fields here.

The program can be exited using `Ctrl-C`. If you need to backup sqlite files, please make sure
`sqlite3` command line is in your `PATH`.

### In Docker

The backpack container can be used to backing up files from other container by
sharing their volume path. For example, if you want to back up your vaultwarden password
manager's vault files, you can use the following Docker compose file:

```yaml
version: '3'
services:
  vaultwarden:
    image: 'vaultwarden/server:latest'
    restart: unless-stopped
    volumes:
      - ${HOME}/vaultwarden:/data

  backpack:
    image: 'easyselfhost/backpack:latest'
    restart: unless-stopeed
    volumes:
      # Backpack config file
      - ${HOME}/backpack/config.json:/config/config.json
      # Rclone config file
      - ${HOME}/backpack/rclone.conf:/config/rclone.conf
      - ${HOME}/vaultwarden:/data/vaultwarden
```
