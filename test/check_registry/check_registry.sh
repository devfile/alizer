#!/bin/bash

# This script runs an E2E check for alizer project and a list
# of given registries. At the moment we check 2 different registries
# the devfile community registry and the redhat product registry.

# The steps of this script are:
#   1. Generates a JSON list for all registries and their entries.
#   In order to do that it fetches the list of all stacks from the
#   registry and tries to fetch more info for its starterprojects.
#   The format of the json list is:
#    [
#        {
#           "Devfile": "name of the stack",
#           "Repo": "url of the repo",
#           "Revision": "revision to clone",
#           "SubDir": "sub-directory inside repo",
#       },
#   ]

#   2. After the list of stacks is generated it loops all generated
#   enries.

#   2a. For each entry it clones the given repo. If a revision is given
#   it clones a --single-branch for the given revision:
#        git clone --single-branch <branch> <repoUrl> tmp/<devfile_name>

#   2b. It runs the alizer binary against the cloned project and checks
#   if the devfile name is inside the list of matched devfiles alizer
#   returns.
echo ":: Creating Registry Entries JSON..."
echo ""
REGISTRY_ENTRIES_OUTPUT=$(go run test/check_registry/check_registry.go)
ENTRIES_PASSED=0
ENTRIES_FAILED=0

for entry in $(echo $REGISTRY_ENTRIES_OUTPUT | jq -c '.[]')
do
    # Assign variables for this entry
    devfile=$(jq -r '.Devfile' <<< $entry)
    path="tmp/$devfile"
    found_matching="0"
    repo=$(jq -r '.Repo' <<< $entry)
    registry=$(jq -r '.Registry' <<< $entry)
    revision=$(jq -r '.Revision' <<< $entry)
    subdir=$(jq -r '.SubDir' <<< $entry)

    # Clone project according to data
    echo ":: Cloning project for entry <$devfile>"
    echo ""
    
    if [ "$revision" != "" ]; then
        echo "$devfile -> found revision $revision for repo $repo"
        git clone --single-branch --branch $revision $repo tmp/$devfile
        if [ "$devfile" == "java-wildfly-bootable-jar" ]; then
            echo "revision: $revision repo: $repo path: tmp/$devfile"
        fi
    else
        git clone $repo tmp/$devfile
        if [ "$subdir" != "" ]; then
        path="$path/$subdir"
        fi
    fi

    # Checking with alizer
    echo ""
    echo ":: Running alizer against path $path"
    echo ""
    alizer_output=$(./alizer devfile --registry $registry $path)
    for raw_selected_devfile_name in $(echo $alizer_output | jq -c '.[].Name')
    do
        selected_devfile_name=$(sed -e 's/^"//' -e 's/"$//' <<<"$raw_selected_devfile_name")
        # Loop through the list of proposed devfiles to find the correct one
        if [[ "$selected_devfile_name" ==  *"$devfile"* ]] && [ "$found_matching" == "0" ]; then
            # If devfile name is contained inside selected one success
            echo "------------------"
            echo "SUCCESS - Devfile Name: $devfile <> Matched Devfile Name: $selected_devfile_name"
            echo "------------------"
            echo ""
            let ENTRIES_PASSED++
            found_matching="1"
        fi
    done

    # If the correct devfile is not matched throw error
    if [ "$found_matching" == "0" ]; then
        # nodejs-mongodb similar to nodejs
        if [[ "$repo" == "https://github.com/che-samples/nodejs-mongodb-sample" && "$selected_devfile_name" == "nodejs" ]]; then
            echo "[WARNING] $repo matched with $selected_devfile_name instead of $devfile. Will not raise error."
        else
            let ENTRIES_FAILED++
            echo "[FAIL] Failed to match project $repo with expected $devfile. Command ./alizer devfile --registry $registry $path"
            echo "[FAIL] Output: $alizer_output"
            echo "(PASSED: $ENTRIES_PASSED / FAILED: $ENTRIES_FAILED)"
            rm -rf tmp/$devfile
            exit 1
        fi
    fi
    rm -rf tmp/$devfile
done
echo "[OK] PASSED: $ENTRIES_PASSED / FAILED: $ENTRIES_FAILED"