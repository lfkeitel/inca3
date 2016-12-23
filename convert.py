#!/usr/bin/env python3

import sys

filename = sys.argv[1]

profiles = {
    'cisco:ssh': 1,
    'cisco:telnet': 2,
    'cisco:telnetold': 4,
    'cisco:telnetold2': 5,
    'juniper:ssh': 3,
}

# INSERT INTO table VALUES(1,'test-device','Test Device','127.0.0.1',1);

def slug(str):
    str = str.lower()
    str = str.replace("_", "-")
    return str.replace(" ", "-")

with open(filename) as f:
    for line in f:
        line = line.strip()
        if line.startswith("#") or len(line) == 0:
            continue

        lineParts = line.split("::")

        deviceName = lineParts[0]
        address = lineParts[1]
        profile = lineParts[2]+":"+lineParts[3]

        stmt = "INSERT INTO \"device\" (name, slug, address, type) VALUES('{0}','{1}','{2}',{3});".format(
            deviceName,
            slug(deviceName),
            address,
            profiles[profile]
        )

        print(stmt)