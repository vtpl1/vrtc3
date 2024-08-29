- By default vrtc3 will search config file `vrtc3.yaml` in current work directory
- vrtc3 support multiple config files:
  - `vrtc3 -c config1.yaml -c config2.yaml -c config3.yaml` 
- vrtc3 support inline config as multiple formats from command line:
  - **YAML**: `vrtc3 -c '{log: {format: text}}'`
  - **JSON**: `vrtc3 -c '{"log":{"format":"text"}}'`
  - **key=value**: `vrtc3 -c log.format=text`
- Every next config will overwrite previous (but only defined params)

```
vrtc3 -config "{log: {format: text}}" -config /config/vrtc3.yaml -config "{rtsp: {listen: ''}}" -config /usr/local/vrtc3/vrtc3.yaml
```

or simple version

```
vrtc3 -c log.format=text -c /config/vrtc3.yaml -c rtsp.listen='' -c /usr/local/vrtc3/vrtc3.yaml
```

## Environment variables

Also vrtc3 support templates for using environment variables in any part of config:

```yaml
streams:
  camera1: rtsp://rtsp:${CAMERA_PASSWORD}@192.168.1.123/av_stream/ch0

rtsp:
  username: ${RTSP_USER:admin}   # "admin" if env "RTSP_USER" not set
  password: ${RTSP_PASS:secret}  # "secret" if env "RTSP_PASS" not set
```

## JSON Schema

Editors like [GoLand](https://www.jetbrains.com/go/) and [VS Code](https://code.visualstudio.com/) supports autocomplete and syntax validation.

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/vtpl1/vrtc3/master/website/schema.json
```

## Defaults

- Default values may change in updates
- FFmpeg module has many presets, they are not listed here because they may also change in updates

```yaml
api:
  listen: ":1984"

ffmpeg:
  bin: "ffmpeg"

log:
  format: "color"
  level: "info"
  output: "stdout"
  time: "UNIXMS"

rtsp:
  listen: ":8554"
  default_query: "video&audio"

srtp:
  listen: ":8443"

webrtc:
  listen: ":8555/tcp"
  ice_servers:
    - urls: [ "stun:stun.l.google.com:19302" ]
```
