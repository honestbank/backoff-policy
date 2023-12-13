#!/bin/zsh

RELEASE_WORKFLOW=".github/workflows/repository-release-prod.yaml"
CONFIG_FILE="config/config.go"
GO_MOD="go.mod"
META_DATA_FILE="catalog-info.meta.json"

if [ ! -f $META_DATA_FILE ]; then
  echo "missing meta file, generating example file for you"
  cat << EOF >> $META_DATA_FILE
{
  "squad_name": "example-squad",
  "design_document": "https://example.com",
  "runbook": "https://example.com",
  "manual_dependencies": [],
  "example-service-name": {
    "tags" : [
      "language:golang",
      "idempotent:false",
      "stateless:false"
    ]
  }
}
EOF
fi

# Define the output file name
OUTPUT_FILE="catalog-info.yaml"
: > $OUTPUT_FILE # Clear the output file before appending

typeset -A TEAM_MAP
TEAM_MAP[acquisition-squad]=backend-engineers
TEAM_MAP[acquisition]=backend-engineers
TEAM_MAP[data-engineering]=data-squad
TEAM_MAP[decisioning]=backend-engineers
TEAM_MAP[devops]=devops-engineers
TEAM_MAP[internal-infra]=devops-engineers
TEAM_MAP[self-service-squad]=backend-engineers
TEAM_MAP[self-service]=backend-engineers
TEAM_MAP[spend-squad]=backend-engineers
TEAM_MAP[spend]=backend-engineers

squad_exist_in_team() {
    local pattern=$1
    for key in "${(@k)TEAM_MAP}"; do
        if [[ $key == $pattern ]]; then
            return 0
        fi
    done
    return 1
}

REPO_NAME=$(basename "$(pwd)")
SERVICE_NAMES=(${(s: :)$(yq e '.jobs.repository-release-prod.with.helm_release_names' "$RELEASE_WORKFLOW")})
if [[ ${#SERVICE_NAMES[@]} == 0 || "$SERVICE_NAMES" == "null" ]]; then
  SERVICE_NAMES=($REPO_NAME)
fi

SQUAD_NAME=$(yq e '.jobs.repository-release-prod.with.argocd_state_repo' "$RELEASE_WORKFLOW")
SQUAD_NAME=$(echo "$SQUAD_NAME" | cut -c 14-50)
if [[ -z $SQUAD_NAME || "$SQUAD_NAME" == "null" ]]; then
  SQUAD_NAME=$(jq -r '.squad_name' $META_DATA_FILE)
fi
GH_TEAM=""
if squad_exist_in_team "$SQUAD_NAME"; then
  GH_TEAM=${TEAM_MAP[$SQUAD_NAME]}
else
  GH_TEAM="devops-engineers"
fi

if [[ "$GH_TEAM" == "null" ]]; then
  echo "couldn't find service owner"
  exit 1
fi

DESIGN_DOCUMENT=$(jq -r '.design_document' $META_DATA_FILE)
RUNBOOK=$(jq -r '.runbook' $META_DATA_FILE)

if [[ -z $DESIGN_DOCUMENT || -z $RUNBOOK ||  "$DESIGN_DOCUMENT" == "https://example.com" || "$RUNBOOK" == "https://example.com" ]]; then
  echo "couldn't find design document or runbook"
  exit 1
fi

# Loop through each subfolder in the charts directory
for SERVICE in $SERVICE_NAMES; do
    # Default dependencies
    DEPENDENCIES=(${(s: :)$(jq -r ".manual_dependencies[]" $META_DATA_FILE)})
    TOPICS=(${(s: :)$(grep Topic "config/config.go" | sed -n 's/.*default:"\([^"]*\)".*/\1/p')})
    for topic in $TOPICS; do
      DEPENDENCIES+=("resource:confluent-$topic")
    done
    BUCKETS=(${(s: :)$(grep Bucket "config/config.go" | sed -n 's/.*default:"\([^"]*\)".*/\1/p')})
    for bucket in $BUCKETS; do
      DEPENDENCIES+=("resource:$bucket")
    done
    LIBS=(${(s: :)$(grep honestbank go.mod | sed -n 's|.*honestbank/\(.*\) v.*|\1|p' | sed 's|/|-|g')})
    for lib in $LIBS; do
      DEPENDENCIES+=("component:$lib")
    done
    TAGS=(${(s: :)$(jq -r ".\"$SERVICE\".tags[]" $META_DATA_FILE)})
    if [[ ${#TAGS[@]} == 0 || "$TAGS" == "null" ]]; then
      TAGS=(
        "language:golang"
        "idempotent:false"
        "stateless:false"
      )
    fi
    TAGS=(${(o)TAGS})
    # Sorting by AESC
    DEPENDENCIES=(${(o)DEPENDENCIES})
    cat << EOF >> $OUTPUT_FILE
---
apiVersion: backstage.io/v1alpha1
kind: Component
metadata:
  name: $SERVICE
  description: The $SERVICE workload
  annotations:
    github.com/project-slug: honestbank/$REPO_NAME
    github.com/team-slug: honestbank/$GH_TEAM
    sonarqube.org/project-key: honestbank_$REPO_NAME
  tags:
$(for tag in "${TAGS[@]}"; do
    echo "    - $tag"
done)
  links:
    - url: $DESIGN_DOCUMENT
      title: Design Document
      icon: menubook
    - url: $RUNBOOK
      title: Runbook
      icon: help
spec:
  type: application
  lifecycle: production
  owner: group:$SQUAD_NAME-squad
  dependsOn:
$(for resource in "${DEPENDENCIES[@]}"; do
    echo "    - $resource"
done)
EOF
done

echo "File generated: $OUTPUT_FILE"
