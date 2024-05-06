FROM gcr.io/openssf/scorecard:stable

ENV GITHUB_AUTH_TOKEN default_token
ENV Checks All
ENV REPO_URL repo_url

CMD ["--format=json", "--repo=${REPO_URL}", "--checks=${Checks}"]
