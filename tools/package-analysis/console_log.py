import os
import json
import pandas as pd


class Checker:
    def __init__(self, db_filename="security_db.csv"):
        self.df = pd.read_csv(db_filename, index_col="Command")

    def score_check(self, command, environment):
        ''' Check is command is secure
            return: score (0-10)
        '''
        score = 10
        risk_to_score = {"High":3, "Medium":5, "Low":8}
        for cmd in command:
            if cmd in self.df.index:
                arg, env, risk = self.df.loc[cmd]
                score = min(score, risk_to_score[risk])
                break
        return score

    def run(self, start_dir="/result"):
        os.chdir(start_dir)
        for name in os.listdir():
            if os.path.isfile(name) and name[-5:] == ".json":
                with open(name, 'r') as f:
                    data = json.load(f)
                    new_data = dict()
                    new_data["package_info"] = data["Package"]
                    new_data["commands"] = data["Analysis"]["import"]["Commands"]
                    new_data["commands"] += data["Analysis"]["install"]["Commands"]
                    for entry in new_data["commands"]:
                        entry["Score"] = self.score_check(*entry.values())
                    json_string = json.dumps(new_data, indent=4)
                    print(json_string)


checker = Checker("security_db.csv")
checker.run("/results")