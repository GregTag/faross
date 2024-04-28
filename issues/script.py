import argparse
import subprocess
import requests
import sys

def get_score(owner, repo):
	r = requests.get(f'https://api.github.com/repos/{args.owner}/{args.repo}')
	if r.status_code == 404:
		print(f"Issues tool could not find the specified repo: {args.owner}/{args.repo}", file=sys.stderr)
		return 0
	elif r.status_code != 200:
		print(f"Issues tool got non-200 response from GitHib: {r}", file=sys.stderr)
		return 0

	# Ok, repo exists
	try:
		last_month = int(subprocess.check_output(['./api.sh', args.owner, args.repo, "30"]).decode('utf-8').strip())
		last_3months = int(subprocess.check_output(['./api.sh', args.owner, args.repo, "90"]).decode('utf-8').strip())
		normal = min((last_3months - last_month) / 2, 1)
		increase = last_month / normal
		return max(min(12 - 2 * increase**1.5, 0), 10)
	except Exception as e:
		print("Unknown internal error in popularity tool:", e)
		return 0


TOOL_NAME = "Issues"
DESC = "Scores issue splash based on increase in issue number over the last month"

parser = argparse.ArgumentParser(
                    prog=TOOL_NAME + ' tool',
                    description=DESC)
parser.add_argument('owner')  # positional argument
parser.add_argument('repo')  # positional argument
args = parser.parse_args()

score = get_score(args.owner, args.repo)

result = {
	TOOL_NAME: {
		"score": score,
		"risk": "medium",
		"desc": DESC
	}
}

print(result)