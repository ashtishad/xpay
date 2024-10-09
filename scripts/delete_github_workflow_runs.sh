#!/bin/bash

# Replace with your GitHub details
GITHUB_TOKEN=""
REPO_OWNER=""
REPO_NAME=""

# Fetch all workflow runs
workflow_runs=$(curl -s -H "Authorization: token $GITHUB_TOKEN" \
  "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/actions/runs" | jq -r '.workflow_runs[].id')

# Loop through each workflow run and delete it
for run_id in $workflow_runs; do
  echo "Deleting workflow run $run_id"
  curl -s -X DELETE -H "Authorization: token $GITHUB_TOKEN" \
    "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/actions/runs/$run_id"
done

echo "All workflow runs deleted!"
