def calc_max_penalty(risk_tags, penalties):
    ans = 0
    for tags, penalty in zip(risk_tags, penalties):
        ans += len(tags) * penalty
    return ans

def define_risk(project_tag, risk_tags):
    for risk, tags in enumerate(risk_tags):
        for tag in tags:
            if tag in project_tag:
                return risk
    return None

low_risk = [
    "Component.Executable.SharedLibraryLoading",
    "Miscellaneous.CodeHygiene.Comment",
    "OS.Environment",
    "WebApp.Communications",
]

medium_risk = [
    "CloudServices.DataStorage",
    "CloudServices.Finance.eCommerce",
    "Data.Zipfile",
    "OS.FileOperation",
    "OS.Process.ListRequest",
    "WebApp.Storage",
]

high_risk = [
    "OS.Network.Connection",
    "OS.ACL.Write.Unsafe",
    "OS.Process.DynamicExecution",
]

risk_tags = [low_risk, medium_risk, high_risk]
penalties = [1, 2, 3]
max_penalty = calc_max_penalty(risk_tags, penalties)


project_penalty = 0
with open('tags.txt', 'r') as f:
    for line in f:
        risk = define_risk(line, risk_tags)
        if risk != None:
            project_penalty += penalties[risk]

print((max_penalty - project_penalty) * 10 / 25)