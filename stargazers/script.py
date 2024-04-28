import numpy as np
import argparse
import subprocess
import sys


start = 10  # score = 1
stop = 50000  # score = 10
thresholds = np.geomspace(start, stop, 10).astype(int)

parser = argparse.ArgumentParser(
                    prog='Popularity tool',
                    description='Scores popularity based on number of stargazers')
parser.add_argument('owner')  # positional argument
parser.add_argument('repo')  # positional argument
args = parser.parse_args()

stars_str = subprocess.check_output(['./api.sh', args.owner, args.repo]).decode('utf-8').strip()
try:
	stars = int(stars_str)
	# thresholds[score-1] <= stars < thresholds[score]
	score = np.searchsorted(thresholds, stars, side="right")
except Exception as e:
	if stars_str == "null":
		print(f"Popularity tool could not find the specified repo: {args.owner}/{args.repo}", file=sys.stderr)
	else:
		print("Unknown internal error in popularity tool:", e)
	score = 0

result = {
	"Popularity": {
		"score": score,
		"risk": "low",
		"desc": "Scores popularity based on number of stargazers of repo on GitHub"
	}
}

print(result)