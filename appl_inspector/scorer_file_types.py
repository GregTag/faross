EXECUTABLE_FILES = "executable_extensions.txt"
OBFUSCATED_FILES = "obfuscated_extensions.txt"
PROJECT_EXT = "file_types.txt"

with open(EXECUTABLE_FILES, 'r') as f:
    ex_ext = f.read().split()

with open(OBFUSCATED_FILES, 'r') as f:
    obf_ext = f.read().split()

with open(PROJECT_EXT, 'r') as f:
    project_ext = f.read().split()
    project_ext = [ext.lower().strip('\"') for ext in project_ext]

score=10
for executable in ex_ext:
    if executable in project_ext:
        score = 3
        break
for obfuscated in obf_ext:
    if obfuscated in project_ext:
        score = 1
        break

print(score)
