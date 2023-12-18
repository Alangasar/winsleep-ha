# Home Assistant MQTT Windows PC Sleep

## Install
1. Create switch in configuration.yaml

```
mqtt:
  - switch:
      command_topic: "pc/sleep"
```

2. Add wol switch

```
switch:
  - platform: wake_on_lan
    name: PC WOL
    mac: D8:BB:C1:70:2B:EF #You PC MAC adress
    host: 192.168.1.1 #You PC IP adress
    turn_off:
      service: mqtt.publish
      data:
        topic: "pc/sleep"
        payload: "SLEEP"

```