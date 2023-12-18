# Home Assistant MQTT Windows PC Sleep

## Install
1. Add to config.json mqtt connection parameters
2. Add winsleep-ha.exe to windows startup
3. Create switch in configuration.yaml

```
mqtt:
  - switch:
      command_topic: "pc/sleep"
```

4. Add wol switch

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