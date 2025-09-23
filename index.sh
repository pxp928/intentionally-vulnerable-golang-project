        # Get changed files using git
        echo "Detecting changed files..."
        
          # For PRs, compare against base branch
    base_sha="4efac6bb68a7737fa071265d2d4aae70d9268599"
    head_sha="3716b0cb6ecb1667305e5bec19fec506fb2e92c7"          
echo "PR detected: comparing $base_sha..$head_sha"
          changed_files=$(git diff --name-only "$base_sha..$head_sha" 2>/dev/null)
        
        # Filter for relevant file extensions and project files
        relevant_files=""
        if [[ -n "$changed_files" ]]; then
          echo "All changed files:"
          echo "$changed_files"
          
          # Filter for supported files
          relevant_files=$(echo "$changed_files" | grep -E '\.(go|ts|tsx|js|jsx|py|java|scala|kt)$|go\.mod$|package\.json$|setup\.py$|pyproject\.toml$|requirements\.txt$|Pipfile$|pom\.xml$|build\.gradle(\.kts)?$' 2>/dev/null || echo "")
        fi
        
        # Exit if no relevant changes
        if [[ -z "$relevant_files" ]]; then
          echo "No relevant files changed"
          exit 0
        fi
        
        echo "Relevant changed files:"
        echo "$relevant_files"
        
        # Function to find project root
        find_project_root() {
          local file_path="$1"
          local dir=$(dirname "$file_path")
          
          while [[ "$dir" != "." && "$dir" != "/" ]]; do
            if [[ -f "$dir/go.mod" || -f "$dir/package.json" || -f "$dir/setup.py" || \
                  -f "$dir/pyproject.toml" || -f "$dir/pom.xml" || \
                  -f "$dir/build.gradle" || -f "$dir/build.gradle.kts" ]]; then
              echo "$dir"
              return
            fi
            dir=$(dirname "$dir")
          done
          
          # Check root directory
          if [[ -f "go.mod" || -f "package.json" || -f "setup.py" || \
                -f "pyproject.toml" || -f "pom.xml" || \
                -f "build.gradle" || -f "build.gradle.kts" ]]; then
            echo "."
          fi
        }
        
        # Find unique project directories
        declare -A projects
        while IFS= read -r file; do
          [[ -z "$file" ]] && continue
          # Remove working directory prefix if present
          relative_file="${file#${{ inputs.working-directory }}/}"
          project_root=$(find_project_root "$relative_file")
          
          if [[ -n "$project_root" ]]; then
            projects["$project_root"]=1
          fi
        done <<< "$relevant_files"
        
        # Generate SCIP indexes for each project
        for project_dir in "${!projects[@]}"; do
          echo "Processing: $project_dir"
          
          if [[ ! -d "$project_dir" ]]; then
            echo "Directory not found: $project_dir"
            continue
          fi
          
          cd "$project_dir"
          
          # Go projects
          if [[ -f "go.mod" ]]; then
            echo "Generating Go SCIP index..."
            scip-go --output=index.scip
          fi
          
          # TypeScript/Node.js projects  
          if [[ -f "package.json" ]]; then
            echo "Generating TypeScript SCIP index..."
            # Install dependencies quietly
            if [[ -f "package-lock.json" ]]; then
              npm ci --silent 2>/dev/null || true
            elif [[ -f "yarn.lock" ]]; then
              yarn install --frozen-lockfile --silent 2>/dev/null || true
            else
              npm install --silent 2>/dev/null || true
            fi
            scip-typescript index --output=index.scip
          fi
          
          # Python projects
          if [[ -f "setup.py" || -f "pyproject.toml" || -f "requirements.txt" || -f "Pipfile" ]]; then
            echo "Generating Python SCIP index..."
            # Install dependencies quietly
            if [[ -f "requirements.txt" ]]; then
              pip3 install -q -r requirements.txt 2>/dev/null || true
            elif [[ -f "pyproject.toml" ]]; then
              pip3 install -q . 2>/dev/null || true
            fi
            scip-python index --output=index.scip
          fi
          
          # Java projects
          if [[ -f "pom.xml" || -f "build.gradle" || -f "build.gradle.kts" ]]; then
            echo "Generating Java SCIP index..."
            scip-java index --output=index.scip
          fi
          
          cd - >/dev/null
        done
        
        # Find and output all generated SCIP indexes
        scip_files=$(find . -name "index.scip" -type f | sort)
        if [[ -n "$scip_files" ]]; then
          echo "Generated SCIP indexes:"
          echo "$scip_files"
          # Format for GitHub output (newlines to spaces)
          echo "scip_indexes=$(echo "$scip_files" | tr '\n' ' ')" >> "$GITHUB_OUTPUT"
        else
          echo "No SCIP indexes generated"
          echo "scip_indexes=" >> "$GITHUB_OUTPUT"
        fi

