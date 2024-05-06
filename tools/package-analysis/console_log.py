import os
import json


### dont forget to change to /result
os.chdir("/tmp/results")
for name in os.listdir():
    if os.path.isfile(name) and name[-5:] == ".json":
        with open(name, 'r') as f:
            data = json.load(f)
            new_data = dict()
            new_data["package_info"] = data["Package"]
            new_data["import_commands"] = data["Analysis"]["import"]["Commands"]
            new_data["install_commands"] = data["Analysis"]["install"]["Commands"]
            json_string = json.dumps(new_data, indent=4)
            print(json_string)