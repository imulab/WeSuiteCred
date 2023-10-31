# WeSuiteCred

Manage suite ticket and app authorizations from WeCom corporations by listening to messages from an MQTT broker published
by [WeTriage](https://github.com/imulab/WeTriage).

## Features

- [x] Listen for and store latest `suite_ticket`.
- [x] Listen for and store app authorizations, including `corp_id`, `corp_secret/permanent_code` and permissions.
- [x] Search and list app credentials and permissions.

## Usage

```bash
docker pull ghcr.io/imulab/wesuitecred:latest

# For a specific version, use the short commit SHA as the tag. For example:
#   docker pull ghcr.io/imulab/wesuitecred:117eb11f
#
# Note this is just an example, that's not the latest commit hash
```

## Listener

The `listener` command is the default command of the image. It listens for messages from the MQTT broker and interact
with the WeCom API to manage suite ticket and app authorizations.

The following flags are supported:

| Flag             | Description                        | Default | Env                |
|------------------|------------------------------------|---------|--------------------|
| `--debug`        | Enable debug mode                  | `false` | `WSC_DEBUG`        |
| `--mqtt-url`     | MQTT broker URL. See details below | -       | `WSC_MQTT_URL`     |
| `--suite-id`     | App template suite id              | -       | `WSC_SUITE_ID`     |
| `--suite-secret` | App template suite secret          | -       | `WSC_SUITE_SECRET` |

The SQLite database is written at `/var/WeSuiteCred` inside the container. You may want to mount a volume to this directory.

Below shows an example of using the image.

```bash
docker run -d \
    -v /var/WeSuiteCred:/var/WeSuiteCred:rw \
    -e WSC_MQTT_URL=tcp://localhost:1883 \
    -e WSC_SUITE_ID=your_suite_id \
    -e WSC_SUITE_SECRET=your_suite_secret \
    ghcr.io/imulab/wesuitecred:latest
```

## Show

The `show` command can be invoked by calling `WeSuiteCred show` in the image.

The following flags are supported:

| Flag            | Description                               |
|-----------------|-------------------------------------------|
| `--query`, `-q` | Query to match the corporation name or id |

Below shows an example:

```bash
docker run \
    -v /var/WeSuiteCred:/var/WeSuiteCred:ro \
    ghcr.io/imulab/wesuitecred:latest \
    WeSuiteCred show -q acme_corp
```

> Note that the database is mounted as read-only.

## Simulate `change_auth` event

For some reason, WeCom does not seem to push the `change_auth` event to the registered callback endpoint under some 
circumstances. As a result, WeTriage will not post a message to notify the change. As a workaround, this image provides
a utility to actively refresh app permissions for a corp authorization.

The following flags are supported:

| Flag         | Description                        |
|--------------|------------------------------------|
| `--mqtt-url` | MQTT broker URL. See details below |
| `--suite-id` | App template suite id              |
| `--corp-id`  | Authorized corporation id          |

Below shows an example:

```bash
docker run \
    ghcr.io/imulab/wesuitecred:latest \
    WeSuiteCred utils simulate-change-auth \
    --mqtt-url=tcp://localhost:1883 \
    --suite-id=your_suite_id \
    --corp-id=your_corp_id
```

This will trigger a standard `change_auth_info` message being published to the MQTT broker, and a running listener will
take care of refreshing the app permissions for the corporation.